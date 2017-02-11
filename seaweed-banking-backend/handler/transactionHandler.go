package handler

import (
	"log"
	"net/http"

	"time"

	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/database"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/model"
	"github.com/pressly/chi/render"
)

// GetTransaction returns a demo transaction for testing purposes
func GetTransaction(w http.ResponseWriter, r *http.Request) {
	transaction := model.Transaction{
		BIC:                 "BIC",
		IBAN:                "IBAN",
		BookingDate:         time.Now(),
		Currency:            "EUR",
		ValueInSmallestUnit: 100,
		IntendedUse:         "Nothing"}

	render.JSON(w, r, transaction)
}

// CreateTransactionAndUpdateBalance creates the in the body of the request defined posting
// TODO Currently only updating the account balance!
func CreateTransactionAndUpdateBalance(w http.ResponseWriter, r *http.Request) {
	transaction := model.Transaction{}
	if err := render.Bind(r.Body, &transaction); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, err.Error())
		return
	}

	if !transaction.IsValid() {
		log.Println("Transaction is not valid: ", transaction)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, http.StatusText(http.StatusBadRequest))
	} else {
		if err := database.UpdateAccountBalance(transaction); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, http.StatusText(http.StatusBadRequest))
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, transaction)
	}
}

// CreateTransaction checks the transaction
func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	transaction := model.Transaction{}
	if err := render.Bind(r.Body, &transaction); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, err.Error())
		return
	}

	if err := database.CreateTransaction(transaction); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, http.StatusText(http.StatusBadRequest))
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, transaction)
}