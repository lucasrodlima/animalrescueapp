package handlers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/lucasrodlima/animalrescueapp/internal/models"
)

func GetAnimal(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		query := `SELECT id, name, species, age, breed, status, created_at, updated_at FROM animals WHERE id = ? AND deleted_at IS NULL`

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		row := db.QueryRowContext(ctx, query, id)

		var animal models.Animal

		if err := row.Scan(&animal.ID,
			&animal.Name,
			&animal.Species,
			&animal.Age,
			&animal.Breed,
			&animal.Status,
			&animal.CreatedAt,
			&animal.UpdatedAt); err != nil {

			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "animal not found", http.StatusNotFound)
				return
			}

			http.Error(w, "failed to fetch animal", http.StatusInternalServerError)
			log.Printf("GetAnimal scan error: %v", err)
			return
		}

		respondJson(w, http.StatusOK, animal)
	})
}
