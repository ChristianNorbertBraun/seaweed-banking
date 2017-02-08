package database

import (
	"log"

	"github.com/ChristianNorbertBraun/seaweed-banking-backend/model"
)

// ReadAccount returns for a given bic and iban an account or an error if there
// is no matching account
func ReadAccount(bic string, iban string) (*model.Account, error) {
	account := model.Account{}

	if err := Connection.
		QueryRow("SELECT bic, iban, balance FROM accountbalance WHERE iban = $1 AND bic = $2", bic, iban).
		Scan(&account.BIC, &account.IBAN, &account.Balance); err != nil {
		log.Printf("Unable to read accounts with bic %s and iban %s: %s", bic, iban, err)
		return nil, err
	}

	return &account, nil
}

// ReadAccounts returns all accounts created so far with their balance
func ReadAccounts() ([]*model.Account, error) {
	rows, err := Connection.Query("SELECT bic, iban, balance FROM accountbalance")

	if err != nil {
		log.Printf("Unable to read all accounts: %s", err)
		return nil, err
	}

	accounts := []*model.Account{}

	for rows.Next() {
		current := model.Account{}
		err := rows.Scan(&current.BIC, &current.IBAN, &current.Balance)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &current)
	}

	return accounts, nil
}

// CreateAccount creates an account with the given data
func CreateAccount(account model.Account) error {
	if err := Connection.QueryRow("INSERT INTO accountbalance(bic, iban, balance) VALUES ($1, $2, $3)",
		account.BIC, account.IBAN, account.Balance).Scan(); err != nil {
		log.Printf("Unable to create account %s", err)
		return err
	}

	return nil
}
