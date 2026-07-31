package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/lucasrodlima/animalrescueapp/internal/handlers"
	_ "modernc.org/sqlite"
)

type apiConfig struct {
	db   *sql.DB
	port string
}

func main() {

	db, err := sql.Open("sqlite3", "app.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	cfg := apiConfig{
		db:   db,
		port: ":8080",
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	mux.HandleFunc("GET /animals", handlers.GetAnimals(cfg.db))
	mux.HandleFunc("GET /animals", handlers.CreateAnimal(cfg.db))

	log.Printf("Server started on %s", cfg.port)
	if err := http.ListenAndServe(cfg.port, mux); err != nil {
		panic(err)
	}
}
