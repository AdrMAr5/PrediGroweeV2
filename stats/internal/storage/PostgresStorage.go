package storage

import (
	"database/sql"
	"go.uber.org/zap"
	"stats/internal/models"
)

type Storage interface {
	Ping() error
	Close() error
	SaveResponse(response *models.QuestionResponse) error
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
func (p *PostgresStorage) SaveResponse(response *models.QuestionResponse) error {
	if response.IsFirstAnswer {
		_, err := p.db.Exec(`INSERT INTO quiz_sessions (session_id, user_id) values ($1, $2)`, response.SessionID, response.UserID)
		if err != nil {
			return err
		}
	}
	if response.IsLastAnswer {
		_, err := p.db.Exec(`UPDATE quiz_sessions SET finish_time = now() WHERE session_id = $1`, response.SessionID)
		if err != nil {
			return err
		}
	}
	_, err := p.db.Exec(`INSERT INTO answers (session_id, question_id, answer) values ($1, $2, $3)`, response.SessionID, response.QuestionID, response.Answer)
	if err != nil {
		return err
	}
	return nil
}
