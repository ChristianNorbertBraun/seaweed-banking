package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/database"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/model"
	"github.com/pressly/chi/render"
)

// ReadAllUpdates returns all updates
func ReadAllUpdates(w http.ResponseWriter, r *http.Request) {
	updates, err := database.FindAllUpdates()

	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, err.Error())
		return
	}

	render.JSON(w, r, updates)
}

// CreateUpdate transforms the given transaction into an update
// and saves it
func CreateUpdate(w http.ResponseWriter, r *http.Request) {
	var transaction model.Transaction

	if err := render.Bind(r.Body, &transaction); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, err.Error())
		return
	}

	log.Println("TransactionDate: ", transaction)

	update := model.NewUpdate(transaction)
	if err := database.InsertUpdate(update); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, err.Error())
		return
	}

	render.JSON(w, r, update)
}

// RunUpdates executes the creation of the account info for the given updates
func RunUpdates(w http.ResponseWriter, r *http.Request) {
	var updates []*model.Update

	if err := render.Bind(r.Body, &updates); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, err.Error())
		return
	}
	log.Printf("Got %d updates to do!", len(updates))
	UpdateAccountInfo(updates)
}

// UpdateAccountInfo is a helper function which will gather all transactions for the given updates
// and merges them into account infos
func UpdateAccountInfo(updates []*model.Update) {
	for _, update := range updates {
		go func(current *model.Update) {
			accountInfo, err := database.GetLatestAccountInfo(current.BIC, current.IBAN)

			log.Printf("Working on update for bic %s iban %s", current.BIC, current.IBAN)

			if err != nil {
				log.Println("Unable to get latest account info", err)

				return
			}

			if err = deleteUpdateWhenAlreadyUpToDate(accountInfo, current); err != nil {
				log.Println("Unable to delete stale update:", err)
				return
			}

			log.Println("About to get all Transactions after", accountInfo.LatestTransaction)
			transactionsToUpdate, err := database.GetAllTransactionsForAccountAfter(
				current.BIC,
				current.IBAN,
				accountInfo.LatestTransaction)
			log.Println("Got all Transactions: ", len(transactionsToUpdate))

			if err != nil {
				log.Println("Unable to fetch transactions after: ", accountInfo.LatestTransaction)
				log.Println("Error was: ", err)

				return
			}

			if err = mergeTransactionsAndSaveAccountInfo(accountInfo, transactionsToUpdate); err != nil {
				log.Println("Can't create new account info for", current.BIC, current.IBAN)

				return
			}
			if err = database.CreateAccountInfo(accountInfo); err != nil {
				log.Println("Can't create new account info for", current.BIC, current.IBAN)

				return
			}
			log.Println("Done with updating for ", current.BIC, current.IBAN)
		}(update)
	}
}

func deleteUpdateWhenAlreadyUpToDate(accountInfo *model.AccountInfo, update *model.Update) error {
	latestTransactionTime, _ := time.Parse(time.RFC3339Nano, accountInfo.LatestTransaction)
	log.Println("Accountinfo time: ", accountInfo.LatestTransaction)
	log.Println("Update last transaction: ", update.LastTransaction.Format(time.RFC3339Nano))
	if latestTransactionTime.
		After(update.LastTransaction) {
		log.Printf("AccountInfo already up to date")
		if err := database.DeleteUpdate(update.BIC, update.IBAN); err != nil {
			return err
		}
	}
	return nil
}

func mergeTransactionsAndSaveAccountInfo(accountInfo *model.AccountInfo, transactions model.Transactions) error {
	var err error
	for _, transaction := range transactions {
		added, newAccountInfo := accountInfo.AddTransaction(transaction)

		if !added {
			err = database.CreateAccountInfo(accountInfo)
			accountInfo = newAccountInfo
		}
	}

	return err
}
