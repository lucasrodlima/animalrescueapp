package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/lucasrodlima/animalrescueapp/internal/models"
)

func (a *AnimalHandler) CreateAnimal(w http.ResponseWriter, r *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()

	var animal models.Animal

	if err := json.NewDecoder(r.Body).Decode(&animal); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	a.mockDB = append(a.mockDB, animal)

	respondJson(w, http.StatusCreated, nil)
}
