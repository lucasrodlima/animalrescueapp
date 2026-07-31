package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lucasrodlima/animalrescueapp/internal/models"
)

func TestCreateAnimal(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	now := time.Date(2026, time.July, 31, 12, 0, 0, 0, time.UTC)
	animal := models.Animal{
		Name:      "Bob",
		Species:   "Dog",
		Age:       2,
		Breed:     "Labrador",
		Status:    models.StatusAvailable,
		CreatedAt: now,
		UpdatedAt: now,
	}

	body, err := json.Marshal(animal)
	if err != nil {
		t.Fatalf("failed to marshal animal: %v", err)
	}

	mock.ExpectExec(`INSERT INTO animals`).
		WithArgs(animal.Name, animal.Species, animal.Age, animal.Breed, animal.Status, animal.CreatedAt, animal.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	req, err := http.NewRequest("POST", "/animals", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := CreateAnimal(db)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "application/json")
	}

	var got models.Animal
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Errorf("handler returned invalid JSON: %v", err)
	}

	if got.Name != animal.Name || got.Species != animal.Species || got.Age != animal.Age || got.Breed != animal.Breed || got.Status != animal.Status {
		t.Errorf("handler returned unexpected animal: got %+v want %+v", got, animal)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sql expectations: %v", err)
	}
}

func TestCreateAnimal_BadJSON(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	req, err := http.NewRequest("POST", "/animals", bytes.NewBufferString(`{"name":`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := CreateAnimal(db)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unexpected sql interaction: %v", err)
	}
}
