package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"time"

	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/config"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/database"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/handler"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/worker"
	"github.com/pressly/chi"
	"github.com/pressly/chi/docgen"
	"github.com/pressly/chi/middleware"
	"github.com/pressly/chi/render"
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
	worker.SetUpUpdateWorker(10 * time.Second)
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, "Hallo Welt")
	})
	r.Route("/updates", func(r chi.Router) {
		r.Get("/", handler.ReadAllUpdates)
		r.Post("/", handler.CreateUpdate)
	})

	if *routes {
		fmt.Println(docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{
			ProjectPath: "github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater",
			Intro:       "Welcome to the seaweed-banking-account-updater generated docs.",
		}))
		return
	}

	serverURL := config.Configuration.Server.Host +
		":" +
		config.Configuration.Server.Port

	http.ListenAndServe(serverURL, r)
}
