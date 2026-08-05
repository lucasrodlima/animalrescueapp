package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

type DB struct {
	Client *sql.DB
}

func NewDatabase() *DB {
	new_db, err := sql.Open("sqlite", "app.db")
	if err != nil {
		log.Fatal(err)
		return &DB{}
	}

	return &DB{
		Client: new_db,
	}
}
