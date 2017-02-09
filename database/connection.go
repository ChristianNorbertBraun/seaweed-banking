package database

import (
	"database/sql"
	"log"

	"github.com/ChristianNorbertBraun/seaweed-banking-backend/config"
)

// Connection contains the open db connection
var Connection *sql.DB

// Configure sets the db url
func Configure() {
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
