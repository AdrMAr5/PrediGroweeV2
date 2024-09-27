package models

type QuizSession struct {
	ID                int
	Mode              QuizMode
	UserId            int
	CurrentQuestionId int
	State             string
}
