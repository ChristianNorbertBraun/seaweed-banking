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
	Recipient           NoBalanceAccount `json:"receipient"`
	Sender              NoBalanceAccount `json:"sender"`
	BookingDate         time.Time        `json:"bookingDate"`
	Currency            Currency         `json:"currency"`
	ValueInSmallestUnit int32            `json:"valueInSmallestUnit"`
	IntendedUse         string           `json:"intendedUse"`
}

// Transactions is an array of Transaction
type Transactions []*Transaction

func (slice Transactions) Len() int {
	return len(slice)
}

func (slice Transactions) Less(i, j int) bool {
	return slice[i].BookingDate.Before(slice[j].BookingDate)
}

func (slice Transactions) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
