package models

import "time"

type QuizSession struct {
	ID                int
	Mode              QuizMode
	UserID            int
	CurrentQuestionID int
	Status            QuizStatus
	CreatedAt         *time.Time
	UpdatedAt         *time.Time
	FinishedAt        *time.Time
}
