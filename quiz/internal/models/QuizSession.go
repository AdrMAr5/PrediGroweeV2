package models

import (
	"encoding/json"
	"io"
	"time"
)

type QuizSession struct {
	ID                int
	UserID            int
	Mode              QuizMode
	Status            QuizStatus
	CurrentQuestionID int
	CreatedAt         *time.Time
	UpdatedAt         *time.Time
	FinishedAt        *time.Time
}

func (qs *QuizSession) ToJSON(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(qs)
}
