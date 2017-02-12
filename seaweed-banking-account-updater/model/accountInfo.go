package model

// AccountInfo holds all information for an account from the oldest to
// the latest transaction
type AccountInfo struct {
	IBAN              string        `json:"iban"`
	BIC               string        `json:"bic"`
	Balance           int32         `json:"balance"`
	OldestTransaction string        `json:"oldestTransaction"`
	LatestTransaction string        `json:"latestTransaction"`
	Transactions      []Transaction `json:"transactions"`
}
