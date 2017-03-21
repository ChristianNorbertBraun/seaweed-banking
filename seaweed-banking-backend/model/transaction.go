package model

import "time"

// Currency is the type of all given currencies
type Currency string

// Represents all possible currencies
const (
	EUR Currency = "EUR"
	USD          = "USD"
)

// Transaction represents a complete transaction from one account
type Transaction struct {
	Recipient           NoBalanceAccount `json:"recipient"`
	Sender              NoBalanceAccount `json:"sender"`
	BookingDate         time.Time        `json:"bookingDate"`
	Currency            Currency         `json:"currency"`
	ValueInSmallestUnit int32            `json:"valueInSmallestUnit"`
	IntendedUse         string           `json:"intendedUse"`
}

// IsValid checks if all fields of the transaction are set
func (t *Transaction) IsValid() bool {
	return t.Recipient.IsVaild() && t.Sender.IsVaild() && !t.BookingDate.IsZero() &&
		t.Currency != "" && t.ValueInSmallestUnit != 0 && t.IntendedUse != ""
}
