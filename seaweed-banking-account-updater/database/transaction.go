package database

import (
	"time"

	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/model"
)

func GetAllPendingTransactionsForAccount(bic string, iban string) (model.Transactions, error) {
	rows, err := Connection.Query(
		`SELECT recipientName, recipientBic, recipientIban,
		senderName, senderBic, senderIban,
		valueInSmallestUnit, currency, bookingDate, intendedUse
		FROM latestTransaction WHERE recipientBic = $1 AND recipientIban = $2 ORDER BY bookingDate`,
		bic,
		iban)

	if err != nil {
		return nil, err
	}

	transactions := model.Transactions{}

	for rows.Next() {
		current := model.Transaction{}
		accountRecipient := model.NoBalanceAccount{}
		accountSender := model.NoBalanceAccount{}

		err := rows.Scan(&accountRecipient.Name, &accountRecipient.BIC, &accountRecipient.IBAN,
			&accountSender.Name, &accountSender.BIC, &accountSender.IBAN,
			&current.ValueInSmallestUnit, &current.Currency, &current.BookingDate, &current.IntendedUse)

		if err != nil {
			return nil, err
		}

		current.Recipient = accountRecipient
		current.Sender = accountSender

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, &current)
	}

	return transactions, nil
}

func DeletePendingTransactions(transactions model.Transactions) error {
	for _, transaction := range transactions {
		err := DeletePendingTransactionForAccount(transaction.Recipient.BIC, transaction.Recipient.IBAN, transaction.BookingDate)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeletePendingTransactionForAccount(bic string, iban string, bookingDate time.Time) error {
	_, err := Connection.Exec(
		`DELETE FROM latestTransaction WHERE
		recipientBic = $1 AND recipientIban = $2 AND
		bookingDate = $3`,
		bic,
		iban,
		bookingDate)

	return err
}
