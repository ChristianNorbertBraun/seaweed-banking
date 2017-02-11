package handler 

import (
	"net/http"

	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/database"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/model"
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
)

// GetAccount returns an account for the given bic and iban
func GetAccount(w http.ResponseWriter, r *http.Request) {
	bic := chi.URLParam(r, "bic")
	iban := chi.URLParam(r, "iban")

	account, err := database.ReadAccount(bic, iban)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, http.StatusText(http.StatusNotFound))
		return
	}
	render.JSON(w, r, account)
}

// GetAllAccounts retrieves all accounts from database
func GetAllAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := database.ReadAccounts()

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, http.StatusText(http.StatusInternalServerError))
		return
	}

	render.JSON(w, r, accounts)
}

// CreateAccount creates an account
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	account := model.Account{}

	if err := render.Bind(r.Body, &account); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, http.StatusText(http.StatusBadRequest))
		return
	}

	if err := database.CreateAccount(account); err != nil {
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, http.StatusText(http.StatusConflict))
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, account)
}
