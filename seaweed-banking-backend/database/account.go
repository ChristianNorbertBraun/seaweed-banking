package database

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"

	"bytes"

	"time"

	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/config"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/model"
)

// ReadAccount returns for a given bic and iban an account or an error if there
// is no matching account
func ReadAccount(bic string, iban string) (*model.Account, error) {
	account := model.Account{}

	if err := Connection.
		QueryRow("SELECT name, bic, iban, balance FROM accountbalance WHERE bic = $1 AND iban = $2",
			bic,
			iban).Scan(&account.Name, &account.BIC, &account.IBAN, &account.Balance); err != nil {
		log.Printf("Unable to read accounts with bic %s and iban %s: %s", bic, iban, err)
		return nil, err
	}

	return &account, nil
}

// ReadAccounts returns all accounts created so far with their balance
func ReadAccounts() ([]*model.Account, error) {
	rows, err := Connection.Query("SELECT name, bic, iban, balance FROM accountbalance")

	if err != nil {
		log.Printf("Unable to read all accounts: %s", err)
		return nil, err
	}

	accounts := []*model.Account{}

	for rows.Next() {
		current := model.Account{}
		err := rows.Scan(&current.Name, &current.BIC, &current.IBAN, &current.Balance)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &current)
	}

	return accounts, nil
}

// UpdateAccountBalance takes a transaction and applies the
// transaction value to the given account
//
// If the the transaction value would make the account balance go below zero
// there will be returned an error an the transaction will be canceled
func UpdateAccountBalance(transaction model.Transaction) error {
	account := model.Account{}
	tx, err := Connection.Begin()
	if err != nil {
		return err
	}
	defer rollback(err, tx)

	row := tx.QueryRow("SELECT bic, iban, balance FROM accountbalance WHERE bic = $1 AND IBAN = $2",
		transaction.Recipient.BIC,
		transaction.Recipient.IBAN)
	if err = row.Scan(&account.BIC, &account.IBAN, &account.Balance); err != nil {
		return err
	}

	if (account.Balance + transaction.ValueInSmallestUnit) < 0 {
		err = fmt.Errorf("Tried to withdraw %d from account bic: %s iban: %s with balance: %d",
			transaction.ValueInSmallestUnit,
			transaction.Recipient.BIC,
			transaction.Recipient.IBAN,
			account.Balance)

		return err
	}

	if _, err = tx.Exec("UPDATE accountbalance SET balance = $1 where bic = $2 AND iban = $3",
		(account.Balance + transaction.ValueInSmallestUnit),
		transaction.Recipient.BIC,
		transaction.Recipient.IBAN); err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO 
		latestTransaction(
			recipientName,
			recipientBic,
			recipientIban,
			senderName,
			senderBic,
			senderIban,
			valueInSmallestUnit,
			currency,
			bookingDate,
			intendedUse) VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		transaction.Recipient.Name,
		transaction.Recipient.BIC,
		transaction.Recipient.IBAN,
		transaction.Sender.Name,
		transaction.Sender.BIC,
		transaction.Sender.IBAN,
		transaction.ValueInSmallestUnit,
		string(transaction.Currency),
		transaction.BookingDate,
		transaction.IntendedUse)

	return err
}

// CreateAccount creates an account with the given data
func CreateAccount(account model.Account) error {
	tx, err := Connection.Begin()
	if err != nil {
		return err
	}
	defer rollback(err, tx)

	if _, err := Connection.Exec("INSERT INTO accountbalance(name,bic, iban, balance) VALUES ($1, $2, $3, $4)",
		account.Name,
		account.BIC,
		account.IBAN,
		account.Balance); err != nil {
		log.Printf("Unable to create account %s", err)
		return err
	}

	if err := createAccountInfo(account); err != nil {
		log.Printf("Unable to create account info for bic %s, iban %s: %s",
			account.BIC,
			account.IBAN,
			err.Error())

		return err
	}

	return nil
}

func createAccountInfo(account model.Account) error {
	buffer := bytes.Buffer{}
	fileName := time.Now().UTC().Format(time.RFC3339Nano)
	path := fmt.Sprintf("%s/%s/%s",
		config.Configuration.Seaweed.AccountFolder,
		account.BIC,
		account.IBAN)

	if err := json.NewEncoder(&buffer).Encode(account); err != nil {
		return err
	}

	return filer.Create(&buffer, fileName, path)
}

func rollback(err error, tx driver.Tx) {
	if err != nil {
		tx.Rollback()

		return
	}
	err = tx.Commit()
}
