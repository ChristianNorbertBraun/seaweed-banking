package model

// NoBalanceAccount represents a basic account with bic, iban and the balance
type NoBalanceAccount struct {
	Name string `json:"name"`
	BIC  string `json:"bic"`
	IBAN string `json:"iban"`
}
