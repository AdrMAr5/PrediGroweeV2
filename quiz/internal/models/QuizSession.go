package models

type QuizSession struct {
	ID                string
	Mode              QuizMode
	UserId            int
	CurrentQuestionId int
	State             string
}
