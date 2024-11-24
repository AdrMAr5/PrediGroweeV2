package models

import "time"

type QuestionResponse struct {
	QuestionID int        `json:"question_id"`
	Answer     string     `json:"answer"`
	IsCorrect  bool       `json:"is_correct"`
	Time       *time.Time `json:"time,omitempty"`
	UserID     *int       `json:"user_id,omitempty"`
	ScreenSize string     `json:"screen_size"`
	TimeSpent  int        `json:"time_spent"`
	CaseCode   string     `json:"case_code"`
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

type SurveyGroupedStats struct {
	Group    string  `json:"group"`
	Value    string  `json:"value"`
	Total    int     `json:"total"`
	Correct  int     `json:"correct"`
	Accuracy float64 `json:"accuracy"`
}
