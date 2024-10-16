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

func (u *UserStats) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(u)
}
