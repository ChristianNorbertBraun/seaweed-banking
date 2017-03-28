package model

import "time"

// AccountInfo holds all information for an account from the oldest to
// the latest transaction
type AccountInfo struct {
	Name              string         `json:"name"`
	BIC               string         `json:"bic"`
	IBAN              string         `json:"iban"`
	Balance           int32          `json:"balance"`
	Predecessor       string         `json:"predecessor,omitempty"`
	OldestTransaction string         `json:"oldestTransaction,omitempty"`
	LatestTransaction string         `json:"latestTransaction,omitempty"`
	Transactions      []*Transaction `json:"transactions"`
}

// GetTransactionsAfter returns all transaction after the given time
func (ai *AccountInfo) GetTransactionsAfter(after time.Time) []*Transaction {
	length := len(ai.Transactions)

	for i := 0; i < length; i++ {
		if ai.Transactions[i].BookingDate.After(after) ||
			ai.Transactions[i].BookingDate.Equal(after) {
			return ai.Transactions[i:]
		}
	}

	return []*Transaction{}
}

// GetTransactionsAfterAndBefore returns all transactions after and before the given time
func (ai *AccountInfo) GetTransactionsAfterAndBefore(after time.Time, before time.Time) []*Transaction {
	if before.Before(after) {
		return []*Transaction{}
	}

	length := len(ai.Transactions)

	for i := 0; i < length; i++ {
		if ai.Transactions[i].BookingDate.After(after) ||
			ai.Transactions[i].BookingDate.Equal(after) {
			for k := i; k < length; k++ {
				if ai.Transactions[k].BookingDate.After(before) {

					return ai.Transactions[i:k]
				}
			}
			return ai.Transactions[i:]
		}
	}

	return []*Transaction{}
}

// NewAccountInfo creates a new accountInf
func NewAccountInfo(name string, bic string, iban string, balance int32, transactions []*Transaction) *AccountInfo {
	return &AccountInfo{Name: name, BIC: bic, IBAN: iban, Balance: balance, Transactions: transactions}
}
