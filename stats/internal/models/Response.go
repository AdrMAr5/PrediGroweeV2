package models

import (
	"encoding/json"
	"io"
)

type QuestionResponse struct {
	QuestionID int    `json:"question_id"`
	Answer     string `json:"answer"`
	IsCorrect  bool   `json:"is_correct"`
}

func (q *QuestionResponse) FromJSON(r io.Reader) error {
	return json.NewDecoder(r).Decode(q)
}
