package database

import (
	"database/sql"
	"log"
)

// Connection contains the open db connection
var Connection *sql.DB

// Configure sets the db url
func Configure() {
	con, err := sql.Open("postgres", "postgres://go:go@docker/go?sslmode=disable")
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
