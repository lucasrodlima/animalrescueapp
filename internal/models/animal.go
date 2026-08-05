package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type AnimalStatus int

const (
	StatusAvailable AnimalStatus = iota
	StatusAdopted
	// StatusPending
	// StatusRejected
)

type Animal struct {
	ID        int          `json:"id"`
	Name      string       `json:"name"`
	Species   string       `json:"species"`
	Age       int          `json:"age"`
	Breed     string       `json:"breed"`
	Status    AnimalStatus `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt *time.Time   `json:"deleted_at"`
}

var (
	ErrInvalidID       = errors.New("id cannot be negative")
	ErrNameRequired    = errors.New("name is required")
	ErrSpeciesRequired = errors.New("species is required")
	ErrAgeInvalid      = errors.New("age must be between 0 and 50")
	ErrStatusInvalid   = errors.New("invalid animal status")
)

func (s AnimalStatus) IsValid() bool {
	switch s {
	case StatusAvailable, StatusAdopted:
		return true
	default:
		return false
	}
}

func (s AnimalStatus) String() string {
	switch s {
	case StatusAvailable:
		return "available"
	case StatusAdopted:
		return "adopted"
	default:
		return "unknown"
	}
}

func (a *Animal) Validate() error {
	var errs []string

	if a.ID < 0 {
		errs = append(errs, ErrInvalidID.Error())
	}
	if a.Age < 0 || a.Age > 50 {
		errs = append(errs, ErrAgeInvalid.Error())
	}
	if strings.TrimSpace(a.Name) == "" {
		errs = append(errs, ErrNameRequired.Error())
	}
	if strings.TrimSpace(a.Species) == "" {
		errs = append(errs, ErrSpeciesRequired.Error())
	}
	if !a.Status.IsValid() {
		errs = append(errs, ErrStatusInvalid.Error())
	}

	if len(errs) == 0 {
		return nil
	}

	return fmt.Errorf("animal validation failed: %v", strings.Join(errs, "; "))
}
