package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lucasrodlima/animalrescueapp/internal/models"
)

func TestGetAnimals(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	defer db.Close()

	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "name", "species", "age", "breed", "status", "created_at", "updated_at"})
	rows.AddRow(1, "Bob", "Dog", 2, "Labrador", models.StatusAvailable, now, now)
	rows.AddRow(2, "Sara", "Cat", 3, "Persian", models.StatusAvailable, now, now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, species, age, breed, status, created_at, updated_at FROM animals WHERE deleted_at IS NULL`)).WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/animals", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := GetAnimals(db)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "application/json")
	}

	var animals []models.Animal
	if err := json.NewDecoder(rec.Body).Decode(&animals); err != nil {
		t.Errorf("handler returned invalid JSON: %v", err)
	}

	if len(animals) != 2 {
		t.Errorf("handler returned wrong number of animals: got %v want %v", len(animals), 2)
	}

	if len(animals) > 0 && (animals[0].ID != 1 || animals[0].Name != "Bob" || animals[0].Breed != "Labrador") {
		t.Errorf("unexpected first animal: %+v", animals[0])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sql expectations: %v", err)
	}
}
