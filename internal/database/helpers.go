package database

import (
	"database/sql"
	"log"
)

type DB struct {
	Client *sql.DB
}

func NewDatabase() *DB {
	new_db, err := sql.Open("sqlite3", "app.db")
	if err != nil {
		log.Fatal(err)
		return &DB{}
	}

	return &DB{
		Client: new_db,
	}
}
