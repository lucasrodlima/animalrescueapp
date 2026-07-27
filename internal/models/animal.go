package models

type AnimalStatus int

const (
	StatusAvailable AnimalStatus = iota
	StatusAdopted
	StatusPending
	StatusRejected
)

type Animal struct {
	ID      int          `json:"id"`
	Name    string       `json:"name"`
	Age     int          `json:"age"`
	Species string       `json:"species"`
	Breed   string       `json:"breed"`
	Status  AnimalStatus `json:"status"`
}
