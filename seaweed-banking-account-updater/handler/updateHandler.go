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
	updates, err := database.GetAllUpdates()

	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, err.Error())
		return
	}

	render.JSON(w, r, updates)
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

			log.Println("About to get all Transactions after", accountInfo.LatestTransaction)
			log.Println("Length of latestAccountInfo", len(accountInfo.Transactions))
			transactionsToUpdate, err := database.GetAllPendingTransactionsForAccount(
				current.BIC,
				current.IBAN)
			log.Println("Got all Transactions: ", len(transactionsToUpdate))

			if err != nil {
				log.Println("Error: ", err)

				return
			}

			if accountInfo, err = mergeTransactionsAndSaveAccountInfo(accountInfo, transactionsToUpdate); err != nil {
				log.Println("Can't create new account info for", current.BIC, current.IBAN)

				return
			}

			if err = database.CreateAccountInfo(accountInfo); err != nil {
				log.Println("Can't create new account info for", current.BIC, current.IBAN)

				return
			}

			writtenAccountInfo, _ := database.GetLatestAccountInfo(accountInfo.BIC, accountInfo.IBAN)

			latestTransaction, _ := time.Parse(time.RFC3339Nano, writtenAccountInfo.LatestTransaction)
			if latestTransaction.Equal(transactionsToUpdate.Last().BookingDate) {
				if err = database.DeletePendingTransactions(transactionsToUpdate); err != nil {
					log.Println("Can`t delete transactions for", accountInfo.BIC, accountInfo.IBAN, err)

					return
				}
				log.Println("Done with updating for ", current.BIC, current.IBAN)
			}
		}(update)
	}
}

func mergeTransactionsAndSaveAccountInfo(accountInfo *model.AccountInfo, transactions model.Transactions) (*model.AccountInfo, error) {
	var err error
	for _, transaction := range transactions {

		oldestTransaction, _ := time.Parse(time.RFC3339Nano, (*accountInfo).OldestTransaction)
		latestTransaction, _ := time.Parse(time.RFC3339Nano, (*accountInfo).LatestTransaction)
		if transaction.BookingDate.Before(oldestTransaction) ||
			transaction.BookingDate.Before(latestTransaction) {
			log.Println("Somthing is really really bad")

			return nil, nil
		}
		added, newAccountInfo := accountInfo.AddTransaction(transaction)

		if !added {
			err = database.CreateAccountInfo(accountInfo)

			accountInfo = newAccountInfo
		}
	}

	return accountInfo, err
}
