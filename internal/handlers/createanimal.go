package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/lucasrodlima/animalrescueapp/internal/models"
)

func CreateAnimal(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var animal models.Animal
		if err := json.NewDecoder(r.Body).Decode(&animal); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if _, err := db.Exec(`
			INSERT INTO animals (name, species, age, breed, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			`,
			animal.Name, animal.Species, animal.Age, animal.Breed, animal.Status, animal.CreatedAt, animal.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error executing query: %v", err)
			return
		}

		respondJson(w, http.StatusCreated, animal)
	})
}
