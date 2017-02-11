package database

import (
	"log"

	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/config"
	mgo "gopkg.in/mgo.v2"
)

var session *mgo.Session

// Configure establish connection to mongodb
func Configure() {
	s, err := mgo.Dial(config.Configuration.Db.URL)

	if err != nil {
		log.Fatal("Could not connect to mongodb: ", err)
	}
	session = s
	log.Print("Connected to mongodb: ", config.Configuration.Db.URL)
}
