package models

type Question struct {
	ID          int
	Title       string
	Description string
	PatientID   string
	Gender      string
	Images      []Image
	Parameters  []Parameter
}

type Parameter struct {
	Name   string
	Unit   string
	Values []ParameterValue
}

type ParameterValue struct {
	Age   int
	Value float64
}
