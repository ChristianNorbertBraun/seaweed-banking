package database

import (
	"fmt"
	"time"

	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/config"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/model"
)

var updateCollection = "updates"

// InsertUpdate creates or updates an update object in the update
// collection
func InsertUpdate(update *model.Update) error {
	_, err := session.DB(config.Configuration.Db.DBName).
		C(updateCollection).
		UpsertId(update.ID, update)

	return err
}

// FindAllUpdates returns all updates currently in the collection
func FindAllUpdates() ([]*model.Update, error) {
	updates := []*model.Update{}
	err := session.DB(config.Configuration.Db.DBName).
		C(updateCollection).
		Find(nil).
		All(&updates)

	if err != nil {
		return nil, err
	}

	return updates, nil
}

// DeleteUpdate deletes the update for the given bic and iban
func DeleteUpdate(bic string, iban string) error {
	return session.DB(config.Configuration.Db.DBName).
		C(updateCollection).
		RemoveId(bic + iban)
}

// DeleteTransactionsFromBook deletes all transactions from bic until and including the given transaction
func DeleteTransactionsFromBook(update *model.Update) error {
	latestFileName := update.LastTransaction.Format(time.RFC3339Nano)
	path := fmt.Sprintf("%s/%s/%s", config.Configuration.Seaweed.BookFolder, update.BIC, update.IBAN)
	return filer.DeleteFilesUntil(path, latestFileName)
}
