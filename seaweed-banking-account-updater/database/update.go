package database

import "github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/model"

func GetAllUpdates() ([]*model.Update, error) {
	rows, err := Connection.Query(
		`Select recipientBic, recipientIban, max(bookingDate)
		FROM latestTransaction
		GROUP BY recipientBic, recipientIban
	`)

	if err != nil {
		return nil, err
	}

	updates := []*model.Update{}

	for rows.Next() {
		current := model.Update{}

		if err := rows.Scan(&current.BIC, &current.IBAN, &current.LastTransaction); err != nil {
			return nil, err
		}

		updates = append(updates, &current)
	}

	return updates, nil
}
