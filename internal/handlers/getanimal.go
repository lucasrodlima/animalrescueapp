package handlers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/lucasrodlima/animalrescueapp/internal/models"
)

func GetAnimal(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		query := `SELECT * FROM animals
				  WHERE id = ?`

		id := r.PathValue("id")

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

			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Scan error: %v", err)
			return
		}

		if err := row.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Row error: %v", err)
			return
		}

		respondJson(w, http.StatusOK, animal)
	})
}
