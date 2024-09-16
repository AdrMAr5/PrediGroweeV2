package models

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"io"
)

type QuizMode = string

const (
	QuizModeEducational QuizMode = "educational"
	QuizModeClassic     QuizMode = "classic"
	QuizModeLimitedTime QuizMode = "limited_time"
)

type StartQuizPayload struct {
	Mode QuizMode `json:"mode" ,validate:"required,oneof=educational classic limited_time"`
}

func (p *StartQuizPayload) Validate() error {
	return validator.New().Struct(p)
}
func (p *StartQuizPayload) FromJSON(ioReader io.Reader) error {
	return json.NewDecoder(ioReader).Decode(p)
}
