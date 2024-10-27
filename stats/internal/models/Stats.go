package models

import (
	"encoding/json"
	"io"
)

type QuizMode = string

const (
	QuizModeEducational QuizMode = "educational"
	QuizModeClassic     QuizMode = "classic"
	QuizModeLimitedTime QuizMode = "time_limited"
)

type UserStats struct {
	TotalQuestions map[QuizMode]int
	CorrectAnswers map[QuizMode]int
	Accuracy       map[QuizMode]float64
}
type QuestionStats struct {
	QuestionID int
	Answer     string
	IsCorrect  bool
}

type QuizStats struct {
	Mode           QuizMode
	TotalQuestions int
	CorrectAnswers int
	Accuracy       float64
	Questions      []QuestionStats
}

func (s *QuizStats) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(s)
}

func (u *UserStats) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(u)
}
