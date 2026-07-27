package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/lucasrodlima/animalrescueapp/internal/handlers"
	"github.com/lucasrodlima/animalrescueapp/internal/models"
)

type apiConfig struct {
	db   []models.Animal
	port string
	mu   *sync.RWMutex
}

func main() {
	cfg := apiConfig{
		db: []models.Animal{
			{ID: 1, Name: "Bob", Age: 2, Species: "Dog", Breed: "Labrador", Status: models.StatusAvailable},
			{ID: 2, Name: "Sara", Age: 3, Species: "Cat", Breed: "Persian", Status: models.StatusAvailable},
			{ID: 3, Name: "Craig", Age: 4, Species: "Dog", Breed: "Pinscher", Status: models.StatusAvailable},
		},
		port: "8080",
		mu:   &sync.RWMutex{},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	animalHandler := handlers.NewAnimalHandler(cfg.db, cfg.mu)
	mux.HandleFunc("GET /animals", animalHandler.GetAnimals)

	log.Printf("Server started on %s", cfg.port)
	if err := http.ListenAndServe(cfg.port, mux); err != nil {
		panic(err)
	}
}
