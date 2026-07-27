package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/lucasrodlima/animalrescueapp/internal/models"
)

func TestGetAnimals(t *testing.T) {
	req, err := http.NewRequest("GET", "/animals", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mockDB := []models.Animal{
		{ID: 1, Name: "Bob", Age: 2, Species: "Dog", Breed: "Labrador", Status: models.StatusAvailable},
		{ID: 2, Name: "Sara", Age: 3, Species: "Cat", Breed: "Persian", Status: models.StatusAvailable},
	}

	animalHandler := NewAnimalHandler(mockDB, &sync.RWMutex{})
	handler := http.HandlerFunc(animalHandler.GetAnimals)
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

	if animals[0].ID == 0 || animals[0].Name == "" {
		t.Errorf("handler returned empty animal ID or name: %v", animals[0])
	}
}
