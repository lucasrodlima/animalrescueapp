package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/lucasrodlima/animalrescueapp/internal/models"
)

func (db *DB) ReadAnimals() ([]models.Animal, error) {
	query := `SELECT id, name, species, age, breed, status, created_at, updated_at FROM animals WHERE deleted_at IS NULL`

	rows, err := db.Client.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var animals []models.Animal
	for rows.Next() {
		var animal models.Animal
		if err := rows.Scan(&animal.ID, &animal.Name, &animal.Species, &animal.Age, &animal.Breed, &animal.Status, &animal.CreatedAt, &animal.UpdatedAt); err != nil {
			log.Println(err)
			return nil, err
		}
		animals = append(animals, animal)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return animals, nil
}

func (db *DB) ReadAnimalByID(id string) (models.Animal, error) {
	query := `SELECT id, name, species, age, breed, status, created_at, updated_at FROM animals WHERE id = ? AND deleted_at IS NULL`

	row := db.Client.QueryRow(query, id)

	var animal models.Animal

	if err := row.Scan(&animal.ID,
		&animal.Name,
		&animal.Species,
		&animal.Age,
		&animal.Breed,
		&animal.Status,
		&animal.CreatedAt,
		&animal.UpdatedAt); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return models.Animal{}, fmt.Errorf("animal not found: %w", err)
		}

		return models.Animal{}, fmt.Errorf("error during query: %w", err)
	}

	return animal, nil
}

func (db *DB) CreateAnimal(animal models.Animal) error {
	result, err := db.Client.Exec(`
				INSERT INTO animals (name, species, age, breed, status)
				VALUES (?, ?, ?, ?, ?)
				`,
		animal.Name, animal.Species, animal.Age, animal.Breed, animal.Status)
	if err != nil {
		return fmt.Errorf("Failed to create animal: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("Failed to create animal: %w", err)
	}

	return nil
}

func (db *DB) UpdateAnimalByID(id int, animal models.Animal) error {
	query := `UPDATE animals
					  SET name = ?, species = ?, age = ?, breed = ?, status = ?, updated_at = CURRENT_TIMESTAMP
					  WHERE id = ? AND deleted_at IS NULL
					  RETURNING id, created_at, updated_at`

	result, err := db.Client.Exec(query, animal.Name,
		animal.Species,
		animal.Age,
		animal.Breed,
		animal.Status,
		id,
	)
	if err != nil {
		return fmt.Errorf("Failed to update animal: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("Failed to update animal: %w", err)
	}

	return nil
}

func (db *DB) DeleteAnimalByID(id int) error {
	query := `UPDATE animals
					          SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
							  WHERE id = ? AND deleted_at IS NULL`

	result, err := db.Client.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting animal: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("error deleting animal: %w", err)
	}

	return nil
}
