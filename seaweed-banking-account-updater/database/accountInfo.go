package database

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/config"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/model"
)

// GetLatestAccountInfo returns the latest account info for a given
func GetLatestAccountInfo(bic string, iban string) (*model.AccountInfo, error) {
	path := fmt.Sprintf("%s/%s/%s",
		config.Configuration.Seaweed.AccountFolder,
		bic,
		iban)
	accountInfo := model.AccountInfo{}
	directory, err := filer.ReadDirectory(path, "")

	if err != nil {
		return nil, err
	}

	if len(directory.Files) == 0 {
		return nil, fmt.Errorf("There is no latest accountInfo for %s", path)
	}

	latestAccountInfoFileName := directory.Files[len(directory.Files)-1].Name
	latestAccountInfoJSON, err := filer.Read(latestAccountInfoFileName, path)

	err = json.NewDecoder(latestAccountInfoJSON).Decode(&accountInfo)

	if err != nil {
		return nil, err
	}

	return &accountInfo, nil
}

// CreateAccountInfo creates a new accountinfo with the oldest Transaction
// as its name
func CreateAccountInfo(accountInfo model.AccountInfo) error {
	buffer := bytes.Buffer{}
	path := fmt.Sprintf("%s/%s/%s",
		config.Configuration.Seaweed.AccountFolder,
		accountInfo.BIC,
		accountInfo.IBAN)

	if err := json.NewEncoder(&buffer).Encode(accountInfo); err != nil {
		return err
	}

	err := filer.Create(&buffer,
		accountInfo.OldestTransaction,
		path)

	if err != nil {
		return err
	}

	return nil
}
