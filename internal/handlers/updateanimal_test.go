package handlers

import (
	"bytes"
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

func TestUpdateAnimal(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	now := time.Date(2026, time.August, 2, 18, 0, 0, 0, time.UTC)
	animal := models.Animal{
		Name:    "Bob",
		Species: "Dog",
		Age:     3,
		Breed:   "Labrador",
		Status:  models.StatusAdopted,
	}

	body, err := json.Marshal(animal)
	if err != nil {
		t.Fatalf("failed to marshal animal: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(1, now.Add(-24*time.Hour), now)

	mock.ExpectQuery(regexp.QuoteMeta(`UPDATE animals
				  SET name = ?, species = ?, age = ?, breed = ?, status = ?, updated_at = CURRENT_TIMESTAMP
				  WHERE id = ? AND deleted_at IS NULL
				  RETURNING id, created_at, updated_at`)).
		WithArgs(animal.Name, animal.Species, animal.Age, animal.Breed, animal.Status, 1).
		WillReturnRows(rows)

	req, err := http.NewRequest("PUT", "/animals/1", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := UpdateAnimal(db)
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

	if got.ID != 1 || got.Name != animal.Name || got.Species != animal.Species || got.Age != animal.Age || got.Breed != animal.Breed || got.Status != animal.Status {
		t.Errorf("handler returned unexpected animal: got %+v want %+v", got, animal)
	}

	if got.CreatedAt.IsZero() || got.UpdatedAt.IsZero() {
		t.Errorf("handler returned empty timestamps: got %+v", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sql expectations: %v", err)
	}
}

func TestUpdateAnimal_BadJSON(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	req, err := http.NewRequest("PUT", "/animals/1", bytes.NewBufferString(`{"name":`))
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := UpdateAnimal(db)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unexpected sql interaction: %v", err)
	}
}

func TestUpdateAnimal_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(`UPDATE animals
					  SET name = ?, species = ?, age = ?, breed = ?, status = ?, updated_at = CURRENT_TIMESTAMP
					  WHERE id = ? AND deleted_at IS NULL
					  RETURNING id, created_at, updated_at`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 99).
		WillReturnError(sql.ErrNoRows)

	req, err := http.NewRequest("PUT", "/animals/99", bytes.NewBufferString(`{"name":"Ghost","species":"Dog","age":4,"breed":"Mixed","status":0}`))
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "99")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := UpdateAnimal(db)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sql expectations: %v", err)
	}
}

func TestUpdateAnimal_SoftDeleted(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(`UPDATE animals
					  SET name = ?, species = ?, age = ?, breed = ?, status = ?, updated_at = CURRENT_TIMESTAMP
					  WHERE id = ? AND deleted_at IS NULL
					  RETURNING id, created_at, updated_at`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 7).
		WillReturnError(sql.ErrNoRows)

	req, err := http.NewRequest("PUT", "/animals/7", bytes.NewBufferString(`{"name":"Ghost","species":"Dog","age":4,"breed":"Mixed","status":0}`))
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "7")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := UpdateAnimal(db)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for soft-deleted animal: got %v want %v", status, http.StatusNotFound)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sql expectations: %v", err)
	}
}

func TestUpdateAnimal_BadID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	req, err := http.NewRequest("PUT", "/animals/abc", bytes.NewBufferString(`{"name":"Bob","species":"Dog","age":3,"breed":"Labrador","status":1}`))
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "abc")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := UpdateAnimal(db)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unexpected sql interaction: %v", err)
	}
}
