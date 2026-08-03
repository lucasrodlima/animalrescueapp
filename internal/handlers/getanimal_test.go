package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lucasrodlima/animalrescueapp/internal/models"
)

func TestGetAnimal(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	now := time.Date(2026, time.July, 31, 12, 0, 0, 0, time.UTC)

	rows := sqlmock.NewRows([]string{"id", "name", "species", "age", "breed", "status", "created_at", "updated_at"})
	rows.AddRow(1, "Bob", "Dog", 2, "Labrador", models.StatusAvailable, now, now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, species, age, breed, status, created_at, updated_at FROM animals WHERE id = ? AND deleted_at IS NULL`)).
		WithArgs(1).
		WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/animals/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "1")

	rec := httptest.NewRecorder()
	handler := GetAnimal(db)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "application/json")
	}

	var got models.Animal
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Errorf("handler returned invalid JSON: %v", err)
	}

	if got.ID != 1 || got.Name != "Bob" || got.Species != "Dog" || got.Age != 2 || got.Breed != "Labrador" {
		t.Errorf("handler returned unexpected animal: got %+v", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sql expectations: %v", err)
	}
}

func TestGetAnimal_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, species, age, breed, status, created_at, updated_at FROM animals WHERE id = ? AND deleted_at IS NULL`)).
		WithArgs(99).
		WillReturnError(sql.ErrNoRows)

	req, err := http.NewRequest("GET", "/animals/99", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "99")

	rec := httptest.NewRecorder()
	handler := GetAnimal(db)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sql expectations: %v", err)
	}
}

func TestGetAnimal_SoftDeleted(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, species, age, breed, status, created_at, updated_at FROM animals WHERE id = ? AND deleted_at IS NULL`)).
		WithArgs(7).
		WillReturnError(sql.ErrNoRows)

	req, err := http.NewRequest("GET", "/animals/7", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "7")

	rec := httptest.NewRecorder()
	handler := GetAnimal(db)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for soft-deleted animal: got %v want %v", status, http.StatusNotFound)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sql expectations: %v", err)
	}
}

func TestGetAnimal_BadID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	req, err := http.NewRequest("GET", "/animals/abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "abc")

	rec := httptest.NewRecorder()
	handler := GetAnimal(db)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unexpected sql interaction: %v", err)
	}
}
