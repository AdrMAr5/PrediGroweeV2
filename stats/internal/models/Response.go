package models

import (
	"encoding/json"
	"io"
)

type QuestionResponse struct {
	QuestionID    int
	UserID        int
	SessionID     int
	Answer        string
	IsFirstAnswer bool
	IsLastAnswer  bool
}

func (q *QuestionResponse) FromJSON(r io.Reader) error {
	return json.NewDecoder(r).Decode(q)
}
