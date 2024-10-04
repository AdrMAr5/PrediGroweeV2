package storage

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"quiz/internal/models"
	"time"
)

type Store interface {
	Ping() error
	Close() error
	GetQuestionByID(id int) (models.Question, error)
	CreateQuizSession(session models.QuizSession) (models.QuizSession, error)
	GetQuizSessionByID(id int) (models.QuizSession, error)
	UpdateQuizSession(session models.QuizSession) error
	GetUserQuizSessions(userID int) ([]models.QuizSession, error)
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
func (p *PostgresStorage) GetQuestionByID(id int) (models.Question, error) {
	var patientID int
	var q models.Question
	err := p.db.QueryRow(`
		SELECT q.id, q.question, p.code, p.gender, p.age1, p.age2, p.prediction_age, p.id
		FROM questions q
		JOIN patients p ON q.patient_id = p.id
		WHERE q.id = $1`, id).Scan(
		&q.ID, &q.Question, &q.PatientCode, &q.Gender,
		&q.Ages.Age1, &q.Ages.Age2, &q.Ages.PredictionAge, &patientID)
	if err != nil {
		return models.Question{}, fmt.Errorf("error reading question: %w", err)
	}

	rows, err := p.db.Query(`
		SELECT pp.age, p.name, pp.value
		FROM patient_parameters pp
		JOIN parameters p ON pp.parameter_id = p.id
		JOIN patients pat ON pp.patient_id = pat.id
		WHERE pat.id = $1 AND (pp.age = $2 OR pp.age = $3)`,
		patientID, q.Ages.Age1, q.Ages.Age2)
	if err != nil {
		return models.Question{}, fmt.Errorf("error reading parameters: %w", err)
	}
	defer rows.Close()

	q.Parameters = make(map[int][]models.Parameter)
	for rows.Next() {
		var age int
		var param models.Parameter
		err := rows.Scan(&age, &param.Name, &param.Value)
		if err != nil {
			return models.Question{}, fmt.Errorf("error scanning parameter: %w", err)
		}
		q.Parameters[age] = append(q.Parameters[age], param)
	}
	return q, nil
}

func (p *PostgresStorage) GetQuizSessionByID(id int) (models.QuizSession, error) {
	var session models.QuizSession
	err := p.db.QueryRow(`
		SELECT id, user_id, status, mode, current_question, created_at, updated_at, finished_at 
		FROM quiz_sessions WHERE id = $1`, id).
		Scan(&session.ID, &session.UserID, &session.Status, &session.Mode,
			&session.CurrentQuestionID, &session.CreatedAt, &session.UpdatedAt, &session.FinishedAt)
	if err != nil {
		return models.QuizSession{}, err
	}
	return session, nil
}
func (p *PostgresStorage) GetUserQuizSessions(userID int) ([]models.QuizSession, error) {
	var sessions []models.QuizSession
	rows, err := p.db.Query(`SELECT id, user_id, status, mode, current_question, created_at, updated_at, finished_at
		FROM quiz_sessions WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var session models.QuizSession
		err := rows.Scan(&session.ID, &session.UserID, &session.Status, &session.Mode,
			&session.CurrentQuestionID, &session.CreatedAt, &session.UpdatedAt, &session.FinishedAt)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (p *PostgresStorage) CreateQuizSession(session models.QuizSession) (models.QuizSession, error) {
	err := p.db.QueryRow(`
		INSERT INTO quiz_sessions (user_id, status, mode, current_question) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, user_id, status, mode, current_question, created_at`,
		session.UserID, session.Status, session.Mode, 1).
		Scan(&session.ID, &session.UserID, &session.Status, &session.Mode,
			&session.CurrentQuestionID, &session.CreatedAt)
	if err != nil {
		return models.QuizSession{}, err
	}
	return session, nil
}

func (p *PostgresStorage) UpdateQuizSession(session models.QuizSession) error {
	_, err := p.db.Exec(`
		UPDATE quiz_sessions 
		SET status = $1, current_question = $2, updated_at = $3, finished_at = $4 
		WHERE id = $5`,
		session.Status, session.CurrentQuestionID, time.Now(), session.FinishedAt, session.ID)
	return err
}
