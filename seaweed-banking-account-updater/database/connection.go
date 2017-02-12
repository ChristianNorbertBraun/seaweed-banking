package database

import (
	"log"

	weedharvester "github.com/ChristianNorbertBraun/Weedharvester"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/config"
	mgo "gopkg.in/mgo.v2"
)

var session *mgo.Session
var filer weedharvester.Filer

// Configure establish connection to mongodb
func Configure() {
	configureMongodb()
	configureSeaweedFiler()
}

func configureMongodb() {
	s, err := mgo.Dial(config.Configuration.Db.URL)

	if err != nil {
		log.Fatal("Could not connect to mongodb: ", err)
	}
	session = s
	log.Print("Connected to mongodb: ", config.Configuration.Db.URL)
}

func configureSeaweedFiler() {
	fil := weedharvester.NewFiler(config.Configuration.Seaweed.FilerURL)

	if err := fil.Ping(); err != nil {
		log.Fatal("Could not connect to filer: ", config.Configuration.Seaweed.FilerURL)
	}

	filer = fil
	log.Print("Connected to seaweed filer at: ", config.Configuration.Seaweed.FilerURL)
}
