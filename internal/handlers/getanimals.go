package handlers

import "net/http"

func (a *AnimalHandler) GetAnimals(w http.ResponseWriter, r *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()

	animals := a.mockDB

	respondJson(w, http.StatusOK, animals)
}
