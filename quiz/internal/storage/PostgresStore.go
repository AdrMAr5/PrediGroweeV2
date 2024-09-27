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

	// Query for basic question info and patient info
	err := p.db.QueryRow(`
        SELECT q.id, q.title, q.description, pi.patient_id, pi.gender 
        FROM questions q
        JOIN patient_info pi ON q.patient_id = pi.id
        WHERE q.id = $1`, id).Scan(
		&question.ID, &question.Title, &question.Description, &question.PatientID, &question.Gender)
	if err != nil {
		return models.Question{}, err
	}

	// Query for images
	rows, err := p.db.Query("SELECT image_url, age FROM images WHERE question_id = $1", id)
	if err != nil {
		return models.Question{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var img models.Image
		if err := rows.Scan(&img.URL, &img.Age); err != nil {
			return models.Question{}, err
		}
		question.Images = append(question.Images, img)
	}

	// Query for parameters
	rows, err = p.db.Query(`
        SELECT p.name, p.unit, qp.age, qp.value 
        FROM question_parameters qp
        JOIN parameters p ON qp.parameter_id = p.id
        WHERE qp.question_id = $1
        ORDER BY p.name, qp.age`, id)
	if err != nil {
		return models.Question{}, err
	}
	defer rows.Close()

	paramMap := make(map[string][]models.ParameterValue)
	for rows.Next() {
		var name, unit string
		var pv models.ParameterValue
		if err := rows.Scan(&name, &unit, &pv.Age, &pv.Value); err != nil {
			return models.Question{}, err
		}
		paramMap[name] = append(paramMap[name], pv)
	}

	for name, values := range paramMap {
		question.Parameters = append(question.Parameters, models.Parameter{
			Name:   name,
			Unit:   "unitMock",
			Values: values,
		})
	}

	return question, nil
}

func (p *PostgresStorage) CreateQuizSession(session models.QuizSession) (models.QuizSession, error) {
	err := p.db.QueryRow("INSERT INTO quiz_sessions (mode, user_id) VALUES ($1, $2) RETURNING session_id, mode, user_id, current_question_id, state", session.Mode, session.UserId).Scan(&session.ID, &session.Mode, &session.UserId, &session.CurrentQuestionId, &session.State)
	if err != nil {
		return models.QuizSession{}, err
	}
	return session, nil
}
