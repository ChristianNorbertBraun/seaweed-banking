package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ChristianNorbertBraun/seaweed-banking-backend/database"
	_ "github.com/mattes/migrate/driver/postgres"
	"github.com/mattes/migrate/migrate"
	"github.com/pressly/chi"
	"github.com/pressly/chi/docgen"
	"github.com/pressly/chi/middleware"
)

var routes = flag.Bool("routes", false, "Generate router documentation")

func init() {
	flag.Parse()
	database.Configure()
}

func main() {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	allErrors, ok := migrate.UpSync("postgres://go:go@docker/go?sslmode=disable", "./data/migration")
	if !ok {
		log.Println("Unable to do migration for reasons:")
		for _, err := range allErrors {
			log.Println(err)
		}
	}

	r.Route("/accounts", func(r chi.Router) {
		r.Get("/", GetAllAccounts)
		r.Post("/", CreateAccount)

		r.Route("/:bic/:iban", func(r chi.Router) {
			r.Get("/", GetAccount)

			r.Route("/transactions", func(r chi.Router) {
				r.Get("/", GetTransaction)
				r.Post("/", CreateTransactionAndUpdateBalance)
			})
		})
	})

	if *routes {
		// fmt.Println(docgen.JSONRoutesDoc(r))
		fmt.Println(docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{
			ProjectPath: "github.com/ChristianNorbertBraun/seaweed-banking-backend",
			Intro:       "Welcome to the seaweed-banking-backend generated docs.",
		}))
		return
	}

	http.ListenAndServe("localhost:3333", r)
}
