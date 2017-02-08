package main

import (
	"log"
	"net/http"

	"time"

	"github.com/ChristianNorbertBraun/seaweed-banking-backend/model"
	"github.com/pressly/chi/render"
)

// GetTransaction returns a demo transaction for testing purposes
func GetTransaction(w http.ResponseWriter, r *http.Request) {
	transaction := model.Transaction{BIC: "BIC", IBAN: "IBAN", BookingDate: time.Now(), Currency: "EUR", ValueInSmallestUnit: 100, IntendedUse: "Nothing"}

	render.JSON(w, r, transaction)
}

// CreateTransactionAndUpdateBalance creates the in the body of the request defined posting
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
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, transaction)
	}
}
