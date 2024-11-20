package models

import "time"

type QuestionResponse struct {
	QuestionID int        `json:"question_id"`
	Answer     string     `json:"answer"`
	IsCorrect  bool       `json:"is_correct"`
	Time       *time.Time `json:"time,omitempty"`
	UserID     *int       `json:"user_id,omitempty"`
}

type QuestionStats struct {
	QuestionID int `json:"question_id"`
	Total      int `json:"total"`
	Correct    int `json:"correct"`
}

type ActivityStats struct {
	Date    time.Time `json:"date"`
	Total   int       `json:"total"`
	Correct int       `json:"correct"`
}
