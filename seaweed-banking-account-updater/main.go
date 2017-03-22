package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"io/ioutil"

	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/config"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/database"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/handler"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/worker"
	_ "github.com/lib/pq"
	"github.com/pressly/chi"
	"github.com/pressly/chi/docgen"
	"github.com/pressly/chi/middleware"
	"github.com/pressly/chi/render"
)

var routes = flag.Bool("routes", false, "Generate router documentation")
var updateRate = flag.Duration("updateRate", 10*time.Second, "Time between two updates")
var configPath = flag.String("config", "./data/conf/config.json", "Path to json formated config")
var master = flag.Bool("master", false, "Declare update service as master")
var port = flag.String("port", "", "Declare port for updater")
var incomingConnections = flag.Bool("incomingConnections", false, "Enable incoming connections")

func init() {
	flag.Parse()

	err := config.Parse(*configPath)
	if err != nil {
		log.Fatalf("Unable to parse config from: %s because: %s",
			*configPath,
			err)
	}

	if *port != "" {
		config.Configuration.Server.Port = *port
	}

	if *incomingConnections {
		config.Configuration.Server.Host = ""
	}

	database.Configure()
	if *master {
		worker.SetUpUpdateWorker(*updateRate)
	} else {
		worker.SetUpSlavePing(time.Minute)
	}
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		address, err := ioutil.ReadAll(r.Body)

		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.PlainText(w, r, "No adress given")
			return
		}

		worker.Register(string(address))
	})

	r.Post("/do/update", handler.RunUpdates)

	r.Route("/updates", func(r chi.Router) {
		r.Get("/", handler.ReadAllUpdates)
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
