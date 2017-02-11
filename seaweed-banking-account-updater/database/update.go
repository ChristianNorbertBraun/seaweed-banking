package database

import (
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/config"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/model"
)

var collectionName = "updates"

// InsertUpdate creates or updates an update object in the update
// collection
func InsertUpdate(update *model.Update) error {
	_, err := session.DB(config.Configuration.Db.DBName).
		C(collectionName).
		UpsertId(update.ID, update)

	return err
}

// FindAllUpdates returns all updates currently in the collection
func FindAllUpdates() ([]*model.Update, error) {
	updates := []*model.Update{}
	err := session.DB(config.Configuration.Db.DBName).
		C(collectionName).
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
		C(collectionName).
		RemoveId(bic + iban)
}
