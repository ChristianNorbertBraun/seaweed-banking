package model

// Account represents a basic account with bic, iban and the balance
type Account struct {
	BIC     string `json:"bic"`
	IBAN    string `json:"iban"`
	Balance int32  `json:"balance"`
}
