package models

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
