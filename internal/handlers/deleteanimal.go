package handlers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

func DeleteAnimal(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		query := `UPDATE animals
		          SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
				  WHERE id = ? AND deleted_at IS NULL`

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		// change to exec and see the error with result
		row := db.QueryRowContext(ctx, query, id)
		if row.Err() != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "couldn't find animal", http.StatusNotFound)
				return
			}

			http.Error(w, "query error", http.StatusInternalServerError)
			log.Printf("Error during query: %v", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

}
