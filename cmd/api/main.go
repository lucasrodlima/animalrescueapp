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

	mux.HandleFunc("GET /animals", handlers.GetAnimals(cfg.db))
	mux.HandleFunc("GET /animals/{id}", handlers.GetAnimal(cfg.db))
	mux.HandleFunc("POST /animals", handlers.CreateAnimal(cfg.db))
	mux.HandleFunc("PUT /animals/{id}", handlers.UpdateAnimal(cfg.db))
	mux.HandleFunc("DELETE /animals/{id}", handlers.DeleteAnimal(cfg.db))

	log.Printf("Server started on %s", cfg.port)
	if err := http.ListenAndServe(cfg.port, mux); err != nil {
		panic(err)
	}
}
