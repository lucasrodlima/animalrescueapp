package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/lucasrodlima/animalrescueapp/internal/models"
)

func UpdateAnimal(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		var animal models.Animal
		if err := json.NewDecoder(r.Body).Decode(&animal); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		query := `UPDATE animals
				  SET name = ?, species = ?, age = ?, breed = ?, status = ?, updated_at = CURRENT_TIMESTAMP
				  WHERE id = ?
				  RETURNING id, created_at, updated_at`

		row := db.QueryRowContext(ctx, query, animal.Name,
			animal.Species,
			animal.Age,
			animal.Breed,
			animal.Status,
			id,
		)

		if err := row.Scan(&animal.ID,
			&animal.CreatedAt,
			&animal.UpdatedAt); err != nil {

			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "animal not found", http.StatusNotFound)
				return
			}

			http.Error(w, "failed to update animal", http.StatusInternalServerError)
			log.Printf("UpdateAnimal scan error: %v", err)
			return
		}

		respondJson(w, http.StatusOK, animal)
	})
}
