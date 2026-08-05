package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/lucasrodlima/animalrescueapp/internal/database"
	"github.com/lucasrodlima/animalrescueapp/internal/models"
)

func respondJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func GetAnimals(repo *database.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		animals, err := repo.ReadAnimals()
		if err != nil {
			http.Error(w, "couldn't list animals", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		respondJson(w, http.StatusOK, animals)
	})
}

func GetAnimal(repo *database.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		animal, err := repo.ReadAnimalByID(r.PathValue("id"))
		if err != nil {
			log.Println(err)
			http.Error(w, "couldn't find animal", http.StatusNotFound)
			return
		}

		respondJson(w, http.StatusOK, animal)
	})
}

func CreateAnimal(db *database.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var animal models.Animal
		if err := json.NewDecoder(r.Body).Decode(&animal); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		err := db.CreateAnimal(animal)
		if err != nil {
			http.Error(w, "error creating animal", http.StatusInternalServerError)
			return
		}

		respondJson(w, http.StatusCreated, animal)
	})
}

func UpdateAnimal(db *database.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var animal models.Animal
		if err := json.NewDecoder(r.Body).Decode(&animal); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		err = db.UpdateAnimalByID(id, animal)
		if err != nil {
			http.Error(w, "error updating animal", http.StatusInternalServerError)
			return
		}

		respondJson(w, http.StatusOK, animal)
	})
}

func DeleteAnimal(db *database.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		err = db.DeleteAnimalByID(id)
		if err != nil {
			http.Error(w, "couldn't delete animal", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
