package model

// NoBalanceAccount represents an account without balance
type NoBalanceAccount struct {
	Name string `json:"name"`
	BIC  string `json:"bic"`
	IBAN string `json:"iban"`
}

// IsVaild checks the validity of the given account
func (a *NoBalanceAccount) IsVaild() bool {
	return a.Name != "" && a.BIC != "" && a.IBAN != ""
}

// Account is a NoBalanceAccount with balance
type Account struct {
	NoBalanceAccount
	Balance int32 `json:"balance"`
}
