package storage

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"stats/internal/models"
)

type Storage interface {
	Ping() error
	Close() error
	SaveResponse(response *models.QuestionResponse) error
	GetSession(sessionID int) (*models.QuizSession, error)
	SaveSession(session *models.QuizSession) error
	GetUserStatsForMode(userID int, mode models.QuizMode) (correctCount int, wrongCount int, err error)
}

var ErrSessionNotFound = fmt.Errorf("session not found")
var ErrStatsNotFound = fmt.Errorf("stats not found")

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

func (p *PostgresStorage) GetSession(sessionID int) (*models.QuizSession, error) {
	var session models.QuizSession
	err := p.db.QueryRow(`SELECT user_id, finish_time, quiz_mode, session_id FROM quiz_sessions WHERE session_id = $1`, sessionID).Scan(&session.UserID, &session.FinishTime, &session.QuizMode, &session.SessionID)
	if err == sql.ErrNoRows {
		return nil, ErrSessionNotFound
	}
	return &session, nil
}

func (p *PostgresStorage) SaveSession(session *models.QuizSession) error {
	_, err := p.db.Exec(`INSERT INTO quiz_sessions (user_id, quiz_mode, session_id) values ($1, $2, $3)`, session.UserID, session.QuizMode, session.SessionID)
	if err != nil {
		return err
	}
	return nil
}
func (p *PostgresStorage) GetUserStatsForMode(userID int, mode models.QuizMode) (correctCount int, wrongCount int, err error) {
	rows, err := p.db.Query(`select correct, count(*) from answers a
    join quiz_sessions s on a.session_id = s.session_id
    where user_id=$1 and quiz_mode=$2
    group by quiz_mode, correct`, userID, mode)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, ErrStatsNotFound
		}
		return 0, 0, err
	}
	for rows.Next() {
		var isCorrect bool
		var count int
		err = rows.Scan(&isCorrect, &count)
		if isCorrect {
			correctCount = count
		} else {
			wrongCount = count
		}
	}
	return
}
