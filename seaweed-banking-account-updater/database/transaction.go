package database

import (
	"fmt"
	"io"

	"encoding/json"

	weedharvester "github.com/ChristianNorbertBraun/Weedharvester"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/config"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/model"
)

//GetAllTransactionsForAccountAfter fetches all transaction from seaweed which occured after
// the given time as a string. The time has to be formated in time.RFC3339Nano
func GetAllTransactionsForAccountAfter(bic string, iban string, time string) (model.Transactions, error) {
	path := fmt.Sprintf("%s/%s/%s",
		config.Configuration.Seaweed.BookFolder,
		bic,
		iban)

	transactionsInDirectory, err := filer.ReadDirectory(path, time)
	if err != nil {
		return nil, err
	}

	return readCompleteDirectory(transactionsInDirectory)
}

func readCompleteDirectory(directory *weedharvester.Directory) (model.Transactions, error) {
	transactions := model.Transactions{}

	for _, file := range directory.Files {
		reader, err := filer.Read(file.Name, directory.Directory)
		if err != nil {
			return nil, err
		}
		transaction, err := parseTransaction(reader)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func parseTransaction(reader io.Reader) (*model.Transaction, error) {
	transaction := model.Transaction{}

	if err := json.NewDecoder(reader).Decode(&transaction); err != nil {
		return nil, err
	}

return &transaction, nil
}
