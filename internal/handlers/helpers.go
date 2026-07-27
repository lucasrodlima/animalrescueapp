package handlers

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/lucasrodlima/animalrescueapp/internal/models"
)

type AnimalHandler struct {
	mockDB []models.Animal
	mu     *sync.RWMutex
}

func NewAnimalHandler(db []models.Animal, mu *sync.RWMutex) *AnimalHandler {
	return &AnimalHandler{
		mockDB: db,
		mu:     mu,
	}
}

func respondJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
