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
	"github.com/lucasrodlima/animalrescueapp/internal/database"
	"github.com/lucasrodlima/animalrescueapp/internal/models"
)

func newMockRepo(t *testing.T) (*database.DB, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return &database.DB{Client: db}, mock, cleanup
}

func TestGetAnimals(t *testing.T) {
	repo, mock, cleanup := newMockRepo(t)
	defer cleanup()

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
	handler := GetAnimals(repo)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rec.Code, http.StatusOK)
	}

	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v", rec.Header().Get("Content-Type"), "application/json")
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

func TestGetAnimals_AllSoftDeleted(t *testing.T) {
	repo, mock, cleanup := newMockRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name", "species", "age", "breed", "status", "created_at", "updated_at"})
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, species, age, breed, status, created_at, updated_at FROM animals WHERE deleted_at IS NULL`)).WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/animals", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := GetAnimals(repo)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code when all animals are soft-deleted: got %v want %v", rec.Code, http.StatusOK)
	}

	var animals []models.Animal
	if err := json.NewDecoder(rec.Body).Decode(&animals); err != nil {
		t.Errorf("handler returned invalid JSON: %v", err)
	}

	if len(animals) != 0 {
		t.Errorf("handler returned soft-deleted animals: got %v want 0", len(animals))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sql expectations: %v", err)
	}
}

func TestGetAnimal(t *testing.T) {
	repo, mock, cleanup := newMockRepo(t)
	defer cleanup()

	now := time.Date(2026, time.July, 31, 12, 0, 0, 0, time.UTC)
	rows := sqlmock.NewRows([]string{"id", "name", "species", "age", "breed", "status", "created_at", "updated_at"})
	rows.AddRow(1, "Bob", "Dog", 2, "Labrador", models.StatusAvailable, now, now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, species, age, breed, status, created_at, updated_at FROM animals WHERE id = ? AND deleted_at IS NULL`)).
		WithArgs("1").
		WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/animals/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "1")

	rec := httptest.NewRecorder()
	handler := GetAnimal(repo)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rec.Code, http.StatusOK)
	}

	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v", rec.Header().Get("Content-Type"), "application/json")
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
	repo, mock, cleanup := newMockRepo(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, species, age, breed, status, created_at, updated_at FROM animals WHERE id = ? AND deleted_at IS NULL`)).
		WithArgs("99").
		WillReturnError(sql.ErrNoRows)

	req, err := http.NewRequest("GET", "/animals/99", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "99")

	rec := httptest.NewRecorder()
	handler := GetAnimal(repo)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", rec.Code, http.StatusNotFound)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sql expectations: %v", err)
	}
}

func TestGetAnimal_SoftDeleted(t *testing.T) {
	repo, mock, cleanup := newMockRepo(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, species, age, breed, status, created_at, updated_at FROM animals WHERE id = ? AND deleted_at IS NULL`)).
		WithArgs("7").
		WillReturnError(sql.ErrNoRows)

	req, err := http.NewRequest("GET", "/animals/7", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "7")

	rec := httptest.NewRecorder()
	handler := GetAnimal(repo)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for soft-deleted animal: got %v want %v", rec.Code, http.StatusNotFound)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sql expectations: %v", err)
	}
}

func TestGetAnimal_BadID(t *testing.T) {
	repo, mock, cleanup := newMockRepo(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, species, age, breed, status, created_at, updated_at FROM animals WHERE id = ? AND deleted_at IS NULL`)).
		WithArgs("abc").
		WillReturnError(sql.ErrNoRows)

	req, err := http.NewRequest("GET", "/animals/abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "abc")

	rec := httptest.NewRecorder()
	handler := GetAnimal(repo)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", rec.Code, http.StatusNotFound)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sql expectations: %v", err)
	}
}

func TestCreateAnimal(t *testing.T) {
	repo, mock, cleanup := newMockRepo(t)
	defer cleanup()

	animal := models.Animal{
		Name:    "Bob",
		Species: "Dog",
		Age:     2,
		Breed:   "Labrador",
		Status:  models.StatusAvailable,
	}

	body, err := json.Marshal(animal)
	if err != nil {
		t.Fatalf("failed to marshal animal: %v", err)
	}

	mock.ExpectExec(`INSERT INTO animals`).
		WithArgs(animal.Name, animal.Species, animal.Age, animal.Breed, animal.Status).
		WillReturnResult(sqlmock.NewResult(1, 1))

	req, err := http.NewRequest("POST", "/animals", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := CreateAnimal(repo)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", rec.Code, http.StatusCreated)
	}

	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v", rec.Header().Get("Content-Type"), "application/json")
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
	repo, mock, cleanup := newMockRepo(t)
	defer cleanup()

	req, err := http.NewRequest("POST", "/animals", bytes.NewBufferString(`{"name":`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := CreateAnimal(repo)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", rec.Code, http.StatusBadRequest)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unexpected sql interaction: %v", err)
	}
}

func TestUpdateAnimal(t *testing.T) {
	repo, mock, cleanup := newMockRepo(t)
	defer cleanup()

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

	mock.ExpectExec(`UPDATE animals`).
		WithArgs(animal.Name, animal.Species, animal.Age, animal.Breed, animal.Status, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req, err := http.NewRequest("PUT", "/animals/1", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := UpdateAnimal(repo)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rec.Code, http.StatusOK)
	}

	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v", rec.Header().Get("Content-Type"), "application/json")
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

func TestUpdateAnimal_BadJSON(t *testing.T) {
	repo, mock, cleanup := newMockRepo(t)
	defer cleanup()

	req, err := http.NewRequest("PUT", "/animals/1", bytes.NewBufferString(`{"name":`))
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := UpdateAnimal(repo)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", rec.Code, http.StatusBadRequest)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unexpected sql interaction: %v", err)
	}
}

func TestUpdateAnimal_NotFound(t *testing.T) {
	repo, mock, cleanup := newMockRepo(t)
	defer cleanup()

	mock.ExpectExec(`UPDATE animals`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 99).
		WillReturnResult(sqlmock.NewResult(0, 0))

	req, err := http.NewRequest("PUT", "/animals/99", bytes.NewBufferString(`{"name":"Ghost","species":"Dog","age":4,"breed":"Mixed","status":0}`))
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "99")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := UpdateAnimal(repo)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", rec.Code, http.StatusInternalServerError)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sql expectations: %v", err)
	}
}

func TestUpdateAnimal_SoftDeleted(t *testing.T) {
	repo, mock, cleanup := newMockRepo(t)
	defer cleanup()

	mock.ExpectExec(`UPDATE animals`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 7).
		WillReturnResult(sqlmock.NewResult(0, 0))

	req, err := http.NewRequest("PUT", "/animals/7", bytes.NewBufferString(`{"name":"Ghost","species":"Dog","age":4,"breed":"Mixed","status":0}`))
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "7")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := UpdateAnimal(repo)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code for soft-deleted animal: got %v want %v", rec.Code, http.StatusInternalServerError)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sql expectations: %v", err)
	}
}

func TestUpdateAnimal_BadID(t *testing.T) {
	repo, mock, cleanup := newMockRepo(t)
	defer cleanup()

	req, err := http.NewRequest("PUT", "/animals/abc", bytes.NewBufferString(`{"name":"Bob","species":"Dog","age":3,"breed":"Labrador","status":1}`))
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "abc")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := UpdateAnimal(repo)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", rec.Code, http.StatusBadRequest)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unexpected sql interaction: %v", err)
	}
}

func TestDeleteAnimal(t *testing.T) {
	repo, mock, cleanup := newMockRepo(t)
	defer cleanup()

	mock.ExpectExec(`UPDATE animals`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req, err := http.NewRequest("DELETE", "/animals/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "1")

	rec := httptest.NewRecorder()
	handler := DeleteAnimal(repo)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", rec.Code, http.StatusNoContent)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sql expectations: %v", err)
	}
}

func TestDeleteAnimal_NotFound(t *testing.T) {
	repo, mock, cleanup := newMockRepo(t)
	defer cleanup()

	mock.ExpectExec(`UPDATE animals`).
		WithArgs(99).
		WillReturnResult(sqlmock.NewResult(0, 0))

	req, err := http.NewRequest("DELETE", "/animals/99", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "99")

	rec := httptest.NewRecorder()
	handler := DeleteAnimal(repo)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", rec.Code, http.StatusInternalServerError)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sql expectations: %v", err)
	}
}

func TestDeleteAnimal_BadID(t *testing.T) {
	repo, mock, cleanup := newMockRepo(t)
	defer cleanup()

	req, err := http.NewRequest("DELETE", "/animals/abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "abc")

	rec := httptest.NewRecorder()
	handler := DeleteAnimal(repo)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", rec.Code, http.StatusBadRequest)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unexpected sql interaction: %v", err)
	}
}
