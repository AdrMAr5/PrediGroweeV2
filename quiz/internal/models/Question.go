package models

type Question struct {
	ID       int
	Question string
	Answer   string
	Images   []Image
}
