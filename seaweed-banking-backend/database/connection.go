package database

import (
	"database/sql"
	"log"

	weedharvester "github.com/ChristianNorbertBraun/Weedharvester"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/config"
)

// Connection contains the open db connection
var Connection *sql.DB
var filer weedharvester.Filer

// Configure sets the db url
func Configure() {
	configureDb()
	configureSeaweedFiler()
}

func configureDb() {
	con, err := sql.Open("postgres", config.Configuration.Db.URL)
	if err != nil {
		log.Fatal("Could not open DB: ", err)
	}

	err = con.Ping()
	if err != nil {
		log.Fatal("Could not open DB: ", err)
	}

	log.Println("DB initialized")
	Connection = con
}

func configureSeaweedFiler() {
	fil := weedharvester.NewFiler(config.Configuration.Seaweed.FilerURL)
	err := fil.Ping()

	if err != nil {
		log.Fatal("Could not connect to seaweed filer: ", err)
	}

	log.Println("Created filer for: ", config.Configuration.Seaweed.FilerURL)
	filer = fil
}
