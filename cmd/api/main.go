package main

import (
	"log"
	"net/http"

	"github.com/lucasrodlima/animalrescueapp/internal/database"
	"github.com/lucasrodlima/animalrescueapp/internal/handlers"
)

type apiConfig struct {
	db   *database.DB
	port string
}

func main() {
	newDB := database.NewDatabase()
	defer newDB.Client.Close()

	cfg := apiConfig{
		db:   newDB,
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
