package model

import (
	"time"
)

// Update represents when the account with the given BIC and IBAN
// got it's last transaction
type Update struct {
	BIC             string    `json:"bic"`
	IBAN            string    `json:"iban"`
	LastTransaction time.Time `json:"lastTransaction"`
}

// NewUpdate creates a new update
func NewUpdate(transaction Transaction) *Update {

	return &Update{
		BIC:             transaction.Recipient.BIC,
		IBAN:            transaction.Recipient.IBAN,
		LastTransaction: transaction.BookingDate}
}
