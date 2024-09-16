package storage

import (
	"PrediGroweeV2/quiz/internal/models"
	"database/sql"
	"go.uber.org/zap"
)

type Store interface {
	Ping() error
	Close() error
	GetQuestionById(id int) (models.Question, error)
	CreateQuizSession(session models.QuizSession) (models.QuizSession, error)
}

type PostgresStorage struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewPostgresStorage(db *sql.DB, logger *zap.Logger) *PostgresStorage {
	return &PostgresStorage{
		db:     db,
		logger: logger,
	}
}

func (p *PostgresStorage) Ping() error {
	return p.db.Ping()
}

func (p *PostgresStorage) Close() error {
	return p.db.Close()
}

func (p *PostgresStorage) GetQuestionById(id int) (models.Question, error) {
	var question models.Question
	images := make([]models.Image, 2)
	err := p.db.QueryRow("SELECT id, question, image_id_1, image_id_2 FROM questions WHERE id = $1", id).Scan(&question.ID, &question, &question.Answer, &images[0].ID, &images[1].ID)
	if err != nil {
		return models.Question{}, err
	}
	question.Images = images
	return question, nil
}

func (p *PostgresStorage) CreateQuizSession(session models.QuizSession) (models.QuizSession, error) {
	err := p.db.QueryRow("INSERT INTO quiz_sessions (session_id, mode, user_id) VALUES ($1, $2, $3) RETURNING session_id, mode, user_id, current_question_id, state", session.ID, session.Mode, session.UserId).Scan(&session.ID, &session.Mode, &session.UserId, &session.CurrentQuestionId, &session.State)
	if err != nil {
		return models.QuizSession{}, err
	}
	return session, nil
}
