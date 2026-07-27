package main

import (
	"log"
	"net/http"

	"github.com/lucasrodlima/animalrescueapp/internal/handlers"
	"github.com/lucasrodlima/animalrescueapp/internal/models"
)

func main() {
	mockDB := []models.Animal{
		{ID: 1, Name: "Bob", Age: 2, Species: "Dog", Breed: "Labrador", Status: models.StatusAvailable},
		{ID: 2, Name: "Sara", Age: 3, Species: "Cat", Breed: "Persian", Status: models.StatusAvailable},
		{ID: 3, Name: "Craig", Age: 4, Species: "Dog", Breed: "Pinscher", Status: models.StatusAvailable},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	animalHandler := handlers.NewAnimalHandler(mockDB)
	mux.HandleFunc("GET /animals", animalHandler.GetAnimals)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
