package models

import (
	"encoding/json"
	"io"
)

type Question struct {
	ID          int
	Question    string
	PatientCode string
	Gender      string
	Ages        PatientAges
	Images      map[int]string
	Parameters  map[int][]Parameter
}
type PatientAges struct {
	Age1          int
	Age2          int
	PredictionAge int
}

type Parameter struct {
	Name  string
	Value float64
}

func (q *Question) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(q)
}
