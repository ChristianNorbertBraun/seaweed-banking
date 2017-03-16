package model

type NoBalanceAccount struct {
	Name string `json:"name"`
	BIC  string `json:"bic"`
	IBAN string `json:"iban"`
}

// IsVaild checks the validity of the given account
func (a *NoBalanceAccount) IsVaild() bool {
	return a.Name != "" && a.BIC != "" && a.IBAN != ""
}

func NewNoBalanceAccount(account Account) *NoBalanceAccount {
	return &NoBalanceAccount{Name: account.Name, IBAN: account.IBAN, BIC: account.BIC}
}

// Account represents a basic account with bic, iban and the balance
type Account struct {
	Name    string `json:"name"`
	BIC     string `json:"bic"`
	IBAN    string `json:"iban"`
	Balance int32  `json:"balance"`
}

// IsVaild checks the validity of the given account
func (a *Account) IsVaild() bool {
	return a.Name != "" && a.BIC != "" && a.IBAN != ""
}
