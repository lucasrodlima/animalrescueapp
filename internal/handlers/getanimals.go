package handlers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/lucasrodlima/animalrescueapp/internal/models"
)

func GetAnimals(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		query := `SELECT id, name, species, age, breed, status, created_at, updated_at FROM animals WHERE deleted_at IS NULL`

		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			http.Error(w, "failed to list animals", http.StatusInternalServerError)
			log.Printf("GetAnimals query error: %v", err)
			return
		}
		defer rows.Close()

		var animals []models.Animal
		for rows.Next() {
			var animal models.Animal
			if err := rows.Scan(&animal.ID, &animal.Name, &animal.Species, &animal.Age, &animal.Breed, &animal.Status, &animal.CreatedAt, &animal.UpdatedAt); err != nil {
				http.Error(w, "failed to list animals", http.StatusInternalServerError)
				log.Printf("GetAnimals scan error: %v", err)
				return
			}
			animals = append(animals, animal)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "failed to list animals", http.StatusInternalServerError)
			log.Printf("GetAnimals rows error: %v", err)
			return
		}

		respondJson(w, http.StatusOK, animals)
	})
}
