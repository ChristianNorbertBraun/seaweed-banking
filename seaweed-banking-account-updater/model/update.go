package model

import "time"

// Update represents when the account with the given BIC and IBAN
// got it's last transaction
type Update struct {
	ID              string    `json:"id" bson:"_id"`
	BIC             string    `json:"bic" bson:"bic"`
	IBAN            string    `json:"iban" bson:"iban"`
	LastTransaction time.Time `json:"lastTransaction" bson:"lastTransaction"`
}

// NewUpdate creates a new update
func NewUpdate(transaction Transaction) *Update {

	return &Update{
		ID:              transaction.Recipient.BIC + transaction.Recipient.IBAN,
		BIC:             transaction.Recipient.BIC,
		IBAN:            transaction.Recipient.IBAN,
		LastTransaction: transaction.BookingDate}
}
