package handler

import (
	"log"
	"net/http"

	"time"

	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/database"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/model"
	"github.com/pressly/chi/render"
)

// CreateTransactionAndUpdateBalance creates the in the body of the request defined posting
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

		return
	}

	if err := database.UpdateAccountBalance(transaction); err != nil {
		render.Status(r, http.StatusBadRequest)
		log.Println("Error  while creating transaction: ", err.Error())
		render.JSON(w, r, err.Error())

		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, transaction)

}
