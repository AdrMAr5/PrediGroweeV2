package models

import "time"

type QuizSession struct {
	SessionID  int
	UserID     int
	FinishTime time.Time
	QuizMode   string
}
