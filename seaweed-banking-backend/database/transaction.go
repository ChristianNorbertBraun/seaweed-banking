package database

import (
	"bytes"
	"encoding/json"
	"time"

	"log"

	"fmt"

	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/config"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/model"
)

//CreateTransaction creates a transaction within the given bic and iban
func CreateTransaction(transaction model.Transaction) error {
	buffer := bytes.Buffer{}
	filename := transaction.BookingDate.Format(time.RFC3339Nano)
	err := json.NewEncoder(&buffer).Encode(transaction)

	if err != nil {
		log.Println("Error while decoding transaction")

		return err
	}

	err = filer.Create(&buffer,
		filename,
		fmt.Sprintf("%s/%s/%s",
			config.Configuration.Seaweed.BookFolder,
			transaction.BIC,
			transaction.IBAN))

	if err != nil {
		return err
	}

	return nil
}
