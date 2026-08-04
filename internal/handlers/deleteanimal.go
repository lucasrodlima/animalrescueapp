package handlers

import (
	"context"
	"database/sql"
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

		result, err := db.ExecContext(ctx, query, id)
		if err != nil {
			http.Error(w, "error deleting animal", http.StatusInternalServerError)
			return
		}

		rows, err := result.RowsAffected()
		if err != nil {
			http.Error(w, "error deleting animal", http.StatusInternalServerError)
			return
		}

		if rows == 0 {
			http.Error(w, "couldn't find animal", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

}
