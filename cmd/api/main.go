package main

import (
	"database/sql"
	"log"
	"net/http"
	"sync"

	_ "modernc.org/sqlite"

	"github.com/lucasrodlima/animalrescueapp/internal/handlers"
	// "github.com/lucasrodlima/animalrescueapp/internal/models"
)

type apiConfig struct {
	db   *sql.DB
	port string
	mu   *sync.RWMutex
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
		mu:   &sync.RWMutex{},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	animalHandler := handlers.NewAnimalHandler(cfg.db, cfg.mu)
	mux.HandleFunc("GET /animals", animalHandler.GetAnimals)
	mux.HandleFunc("POST /animals", animalHandler.CreateAnimal)

	log.Printf("Server started on %s", cfg.port)
	if err := http.ListenAndServe(cfg.port, mux); err != nil {
		panic(err)
	}
}
