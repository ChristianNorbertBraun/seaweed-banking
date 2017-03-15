package handler

import (
	"fmt"
	"log"
	"net/http"

	"time"

	"errors"

	"bytes"

	"encoding/json"

	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/config"
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

	transaction.BookingDate = time.Now().UTC()

	if !transaction.IsValid() {
		log.Println("Transaction is not valid: ", transaction)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, http.StatusText(http.StatusBadRequest))
	} else {

		if err := database.UpdateAccountBalance(transaction, createTransactionInDFS); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, http.StatusText(http.StatusBadRequest))
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, transaction)
	}
}

func createTransactionInDFS(transaction model.Transaction) error {

	if err := sendTransactionToUpdater(transaction); err != nil {
		return err
	}
	if err := database.CreateTransaction(transaction); err != nil {
		return err
	}

	return nil
}

func sendTransactionToUpdater(transaction model.Transaction) error {
	buffer := bytes.Buffer{}

	url := fmt.Sprintf("%s:%s/updates",
		config.Configuration.Updater.Host,
		config.Configuration.Updater.Port)

	if err := json.NewEncoder(&buffer).Encode(transaction); err != nil {
		return err
	}
	response, err := http.Post(url, "application/json", &buffer)

	if err != nil {
		return err
	}

	if response.StatusCode >= 300 {
		return errors.New("Bad Statuscode while sending transaction update")
	}

	return nil
}
