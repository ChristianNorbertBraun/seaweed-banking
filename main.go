package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ChristianNorbertBraun/seaweed-banking-backend/config"
	"github.com/ChristianNorbertBraun/seaweed-banking-backend/database"
	_ "github.com/mattes/migrate/driver/postgres"
	"github.com/mattes/migrate/migrate"
	"github.com/pressly/chi"
	"github.com/pressly/chi/docgen"
	"github.com/pressly/chi/middleware"
)

var routes = flag.Bool("routes", false, "Generate router documentation")
var configPath = flag.String("config", "./data/conf/config.json", "Path to json formated config")

func init() {
	flag.Parse()

	err := config.Parse(*configPath)
	if err != nil {
		log.Fatalf("Unable to parse config from: %s because: %s",
			*configPath,
			err)
	}

	database.Configure()
}

func main() {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	allErrors, ok := migrate.UpSync(config.Configuration.Db.URL,
		"./data/migration")
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
				r.Post("/", CreateTransaction)
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

	serverURL := config.Configuration.Server.Host +
		":" +
		config.Configuration.Server.Port
	http.ListenAndServe(serverURL, r)
}
