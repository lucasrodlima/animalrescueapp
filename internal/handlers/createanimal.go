package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/lucasrodlima/animalrescueapp/internal/models"
)

func CreateAnimal(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		var animal models.Animal
		if err := json.NewDecoder(r.Body).Decode(&animal); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		row := db.QueryRowContext(ctx, `
			INSERT INTO animals (name, species, age, breed, status)
			VALUES (?, ?, ?, ?, ?)
			RETURNING id, created_at, updated_at
			`,
			animal.Name, animal.Species, animal.Age, animal.Breed, animal.Status)

		if err := row.Scan(&animal.ID, &animal.CreatedAt, &animal.UpdatedAt); err != nil {
			http.Error(w, "failed to create animal", http.StatusInternalServerError)
			log.Printf("CreateAnimal scan error: %v", err)
			return
		}

		respondJson(w, http.StatusCreated, animal)
	})
}
