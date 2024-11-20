package models

import (
	"encoding/json"
	"io"
	"time"
)

type QuestionResponse struct {
	QuestionID int        `json:"question_id"`
	Answer     string     `json:"answer"`
	IsCorrect  bool       `json:"is_correct"`
	Time       *time.Time `json:"time,omitempty"`
	UserID     *int       `json:"user_id,omitempty"`
}

func (q *QuestionResponse) FromJSON(r io.Reader) error {
	return json.NewDecoder(r).Decode(q)
}
