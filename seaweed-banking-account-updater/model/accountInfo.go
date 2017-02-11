package model

// AccountInfo holds all information for an account from the oldest to
// the latest transaction
type AccountInfo struct {
	IBAN              string        `json:"iban"`
	BIC               string        `json:"bic"`
	OldestTransaction string        `json:"oldestTransaction"`
	LatestTransaction string        `json:"latestTransaction"`
	Transactions      []Transaction `json:"transactions"`
}
