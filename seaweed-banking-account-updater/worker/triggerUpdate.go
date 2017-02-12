package worker

import (
	"log"
	"time"

	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/database"
)

var ticker *time.Ticker

// SetUpUpdateWorker starts a worker which checks the update table for new entries
// and updates the matching accountInfos
func SetUpUpdateWorker(duration time.Duration) {
	ticker = time.NewTicker(duration)
	go func() {
		for t := range ticker.C {
			log.Println("-------------------------------------------------------------")
			log.Println("Start Update at: ", t)

			runUpdate()
			log.Println("-------------------------------------------------------------")
		}
	}()
}

// StopUpdateWorker stops the update worker
func StopUpdateWorker() {
	if ticker != nil {
		ticker.Stop()
	}
}

func runUpdate() {
	updates, err := database.FindAllUpdates()
	if err != nil {
		log.Println("Unable to get updates", err)

		return
	}

	for _, update := range updates {
		accountInfo, err := database.GetLatestAccountInfo(update.BIC, update.IBAN)

		log.Printf("Working on update for bic %s iban %s", update.BIC, update.IBAN)

		if err != nil {
			log.Println("Unable to get latest account info", err)

			return
		}

		latestTransactionTime, _ := time.Parse(time.RFC3339Nano, accountInfo.LatestTransaction)
		log.Println("Accountinfo time: ", accountInfo.LatestTransaction)
		log.Println("Update last transaction: ", update.LastTransaction.Format(time.RFC3339Nano))
		if latestTransactionTime.
			After(update.LastTransaction) {
			log.Printf("AccountInfo already up to date")
			err := database.DeleteUpdate(update.BIC, update.IBAN)

			if err != nil {
				log.Println("Unable to delete update entry ", err)
			}
		}

		log.Println("About to get all Transactions after")
		transactionsToUpdate, err := database.GetAllTransactionsForAccountAfter(
			update.BIC,
			update.IBAN,
			accountInfo.LatestTransaction)
		log.Println("Got all Transactions: ", len(transactionsToUpdate))

		if err != nil {
			log.Println("Unable to fetch transactions after: ", accountInfo.LatestTransaction)
			log.Println("Error was: ", err)
			return
		}

		for _, transaction := range transactionsToUpdate {
			added, newAccountInfo := accountInfo.AddTransaction(transaction)

			if !added {
				database.CreateAccountInfo(accountInfo)
				accountInfo = newAccountInfo
			}
		}

		database.CreateAccountInfo(accountInfo)
	}

	log.Println("Done with updating")
}
