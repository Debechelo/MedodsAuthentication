package pkg

import (
	"database/sql"
	"log"
)

var db *sql.DB

func InitializeDB() {
	var err error
	db, err = sql.Open("postgres", "postgres://user:password@localhost/authdb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}
