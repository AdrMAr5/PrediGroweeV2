package models

import (
	"encoding/json"
	"io"
)

type Question struct {
	ID            int      `json:"id"`
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	PredictionAge int      `json:"prediction_age"`
	Case          Case     `json:"case"`
}

func (q *Question) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(q)
}
func (q *Question) FromJSON(r io.Reader) error {
	return json.NewDecoder(r).Decode(q)
}

type QuestionPayload struct {
	ID            int      `json:"id,omitempty"`
	Question      string   `json:"question"`
	Answers       []string `json:"answers"`
	PredictionAge int      `json:"prediction_age"`
	CaseID        int      `json:"case_id"`
}

func (q *QuestionPayload) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(q)
}
func (q *QuestionPayload) FromJSON(r io.Reader) error {
	return json.NewDecoder(r).Decode(q)
}

type QuestionAnswer struct {
	QuestionID int    `json:"question_id"`
	Answer     string `json:"answer"`
	IsCorrect  bool   `json:"is_correct"`
}
