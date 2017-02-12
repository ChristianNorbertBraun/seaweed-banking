package worker

import (
	"fmt"
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
			fmt.Println("Start Update at: ", t)
			runUpdate()
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
		if latestTransactionTime.
			After(update.LastTransaction) {
			log.Printf("AccountInfo already up to date")
			err := database.DeleteUpdate(update.BIC, update.IBAN)

			if err != nil {
				log.Println("Unable to delete update entry ", err)
			}
		}

		transactionsToUpdate, err := database.GetAllTransactionsForAccountAfter(
			update.BIC,
			update.IBAN,
			accountInfo.LatestTransaction)

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
