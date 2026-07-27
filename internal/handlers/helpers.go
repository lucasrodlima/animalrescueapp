package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/lucasrodlima/animalrescueapp/internal/models"
)

type AnimalHandler struct {
	mockDB []models.Animal
}

func NewAnimalHandler(db []models.Animal) *AnimalHandler {
	return &AnimalHandler{
		mockDB: db,
	}
}

func respondJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
