package model

import "sync"
import "sort"
import "time"

// AccountInfo holds all information for an account from the oldest to
// the latest transaction
type AccountInfo struct {
	IBAN              string       `json:"iban"`
	BIC               string       `json:"bic"`
	Balance           int32        `json:"balance"`
	OldestTransaction string       `json:"oldestTransaction"`
	LatestTransaction string       `json:"latestTransaction"`
	Transactions      Transactions `json:"transactions"`
	mutex             sync.Mutex
}

// MaxTransactionsPerAccountInfo represents the maximum number of transaction stored
// within a single AccountInfo
const MaxTransactionsPerAccountInfo = 10

// NewAccountInfo creates a new accountInfo
func NewAccountInfo(bic string, iban string, balance int32) *AccountInfo {
	accountInfo := AccountInfo{BIC: bic, IBAN: iban, Balance: balance}
	accountInfo.Transactions = []*Transaction{}

	return &accountInfo
}

// AddTransaction adds a Transaction to the AccountInfo and updates the balance and
// the oldest and latest transaction date
func (ai *AccountInfo) AddTransaction(transaction *Transaction) (bool, *AccountInfo) {
	ai.mutex.Lock()
	defer ai.mutex.Unlock()
	if ai.Transactions.Len() < MaxTransactionsPerAccountInfo {
		ai.Transactions = append(ai.Transactions, transaction)
		ai.Balance += transaction.ValueInSmallestUnit
		sort.Sort(ai.Transactions)

		ai.OldestTransaction = ai.Transactions[0].
			BookingDate.
			Format(time.RFC3339Nano)
		ai.LatestTransaction = ai.Transactions[ai.Transactions.Len()-1].
			BookingDate.
			Format(time.RFC3339Nano)

		return true, nil
	}

	accountInfo := NewAccountInfo(ai.BIC, ai.IBAN, ai.Balance)
	accountInfo.AddTransaction(transaction)

	return false, accountInfo
}
