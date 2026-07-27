package handlers

import "net/http"

func (a *AnimalHandler) GetAnimals(w http.ResponseWriter, r *http.Request) {
	animals := a.mockDB

	respondJson(w, http.StatusOK, animals)
}
