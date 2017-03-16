package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	weedharvester "github.com/ChristianNorbertBraun/Weedharvester"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/config"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/model"
)

// ErrEmptyID will be returned if there is no accountinfo for a given oldestTransaction
var ErrEmptyID = errors.New("Can't find accountinfo for empty ID.")

// GetAccountInfo returns a single accountinfo for the given bic and iban with the oldestTransaction as id
func GetAccountInfo(bic string, iban string, oldestTransaction string) (*model.AccountInfo, error) {
	if oldestTransaction == "" {
		return nil, ErrEmptyID
	}
	path := fmt.Sprintf("%s/%s/%s",
		config.Configuration.Seaweed.AccountFolder,
		bic,
		iban)
	reader, err := filer.Read(oldestTransaction, path)

	if err != nil {
		return nil, err
	}

	return parseAccountInfo(reader)
}

func GetAccountInfoFromTo(bic string, iban string, from time.Time, to time.Time) (*model.AccountInfo, error) {
	path := fmt.Sprintf("%s/%s/%s",
		config.Configuration.Seaweed.AccountFolder,
		bic,
		iban)

	var accountInfo *model.AccountInfo
	directory, err := filer.ReadDirectory(path, from.UTC().Format(time.RFC3339Nano))

	if err != nil {
		return nil, err
	}

	removeFilesAfterTo(directory, to)

	if len(directory.Files) == 0 {
		accountInfo, err = GetLatestAccountInfo(bic, iban)
		if err != nil {
			return nil, err
		}

		transactionsAfterFromAndBeforeTo := accountInfo.GetTransactionsAfterAndBefore(from, to)
		return model.NewAccountInfo(accountInfo.Name,
			bic,
			iban,
			accountInfo.Balance,
			transactionsAfterFromAndBeforeTo), nil
	}

	accountInfos, err := getAllAccountInfoFromDirectory(directory)

	if err != nil {
		return nil, err
	}

	predeccessorAccountInfo, err := GetAccountInfo(bic, iban, accountInfos[0].Predeccessor)
	if err != nil && err != ErrEmptyID {
		return nil, err
	} else if err == ErrEmptyID {
		accountInfo = createAccountInfoFromListOfAccountInfos(accountInfos, from, to)

		return accountInfo, nil
	}

	accountInfos = append([]*model.AccountInfo{predeccessorAccountInfo}, accountInfos...)
	accountInfo = createAccountInfoFromListOfAccountInfos(accountInfos, from, to)

	return accountInfo, nil
}

// GetAccountInfoFrom returns an accountinfo which holds all transactions from the given time
// TODO massive refactoring
func GetAccountInfoFrom(bic string, iban string, from time.Time) (*model.AccountInfo, error) {
	path := fmt.Sprintf("%s/%s/%s",
		config.Configuration.Seaweed.AccountFolder,
		bic,
		iban)

	var accountInfo *model.AccountInfo
	directory, err := filer.ReadDirectory(path, from.UTC().Format(time.RFC3339Nano))

	if err != nil {
		return nil, err
	}

	// No files found after given time
	if len(directory.Files) == 0 {
		accountInfo, err = GetLatestAccountInfo(bic, iban)
		if err != nil {
			return nil, err
		}

		transactionAfterFrom := accountInfo.GetTransactionsAfter(from)
		return model.NewAccountInfo(accountInfo.Name, bic, iban, accountInfo.Balance, transactionAfterFrom), nil
	}

	accountInfos, err := getAllAccountInfoFromDirectory(directory)

	if err != nil {
		return nil, err
	}

	predeccessorAccountInfo, err := GetAccountInfo(bic, iban, accountInfos[0].Predeccessor)
	if err != nil && err != ErrEmptyID {
		return nil, err
	} else if err == ErrEmptyID {
		accountInfo = createAccountInfoFromListOfAccountInfos(accountInfos, from, time.Time{})

		return accountInfo, nil
	}

	accountInfos = append([]*model.AccountInfo{predeccessorAccountInfo}, accountInfos...)
	accountInfo = createAccountInfoFromListOfAccountInfos(accountInfos, from, time.Time{})

	return accountInfo, nil
}

// GetLatestAccountInfo returns the latest account info for a given
func GetLatestAccountInfo(bic string, iban string) (*model.AccountInfo, error) {
	path := fmt.Sprintf("%s/%s/%s",
		config.Configuration.Seaweed.AccountFolder,
		bic,
		iban)
	directory, err := filer.ReadDirectory(path, "")

	if err != nil {
		return nil, err
	}

	if len(directory.Files) == 0 {
		return nil, fmt.Errorf("There is no latest accountInfo for %s", path)
	}

	latestAccountInfoFileName := directory.Files[len(directory.Files)-1].Name
	latestAccountInfoJSON, err := filer.Read(latestAccountInfoFileName, path)

	accountInfo, err := parseAccountInfo(latestAccountInfoJSON)

	if err != nil {

		return nil, err
	}

	return accountInfo, nil
}

func getAllAccountInfoFromDirectory(directory *weedharvester.Directory) ([]*model.AccountInfo, error) {
	accountInfos := []*model.AccountInfo{}

	for _, file := range directory.Files {
		data, err := filer.Read(file.Name, directory.Directory)
		if err != nil {
			return nil, err
		}

		accountInfo, err := parseAccountInfo(data)

		if err != nil {
			return nil, err
		}

		accountInfos = append(accountInfos, accountInfo)
	}

	return accountInfos, nil
}

func createAccountInfoFromListOfAccountInfos(accountInfos []*model.AccountInfo, after time.Time, before time.Time) *model.AccountInfo {
	transactions := []*model.Transaction{}
	lastAccountInfo := accountInfos[len(accountInfos)-1]

	if before.IsZero() {
		for _, accountInfo := range accountInfos {
			transactions = append(transactions, accountInfo.GetTransactionsAfter(after)...)
		}
	} else {
		for _, accountInfo := range accountInfos {
			transactions = append(transactions, accountInfo.GetTransactionsAfterAndBefore(after, before)...)
		}
	}

	return model.NewAccountInfo(lastAccountInfo.Name,
		lastAccountInfo.BIC,
		lastAccountInfo.IBAN,
		lastAccountInfo.Balance,
		transactions)
}

func parseAccountInfo(reader io.Reader) (*model.AccountInfo, error) {
	accountInfo := model.AccountInfo{}

	if err := json.NewDecoder(reader).Decode(&accountInfo); err != nil {
		return nil, err
	}

	return &accountInfo, nil
}

func removeFilesAfterTo(directory *weedharvester.Directory, to time.Time) {
	for i, file := range directory.Files {
		date, _ := time.Parse(time.RFC3339Nano, file.Name)
		if date.After(to) {
			directory.Files = directory.Files[:i]

			return
		}
	}
}
