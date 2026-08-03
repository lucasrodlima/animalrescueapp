package models

import (
	"database/sql"
	"time"
)

type AnimalStatus int

const (
	StatusAvailable AnimalStatus = iota
	StatusAdopted
	StatusPending
	StatusRejected
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
	DeletedAt sql.NullTime `json:"deleted_at"`
}
