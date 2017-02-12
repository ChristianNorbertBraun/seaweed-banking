package handler

import (
	"log"
	"net/http"

	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/database"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/model"
	"github.com/pressly/chi/render"
)

// ReadAllUpdates returns all updates
func ReadAllUpdates(w http.ResponseWriter, r *http.Request) {
	updates, err := database.FindAllUpdates()

	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, err.Error())
		return
	}

	render.JSON(w, r, updates)
}

// CreateUpdate transforms the given transaction into an update
// and saves it
func CreateUpdate(w http.ResponseWriter, r *http.Request) {
	var transaction model.Transaction

	if err := render.Bind(r.Body, &transaction); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, err.Error())
		return
	}

	log.Println("TransactionDate: ", transaction)

	update := model.NewUpdate(transaction)
	if err := database.InsertUpdate(update); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, err.Error())
		return
	}

	render.JSON(w, r, update)
}
