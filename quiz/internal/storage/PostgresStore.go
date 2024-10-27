package storage

import (
	"database/sql"
	"go.uber.org/zap"
	"quiz/internal/models"
)

type Store interface {
	Ping() error
	Close() error

	// sessions
	CreateQuizSession(session models.QuizSession) (models.QuizSession, error)
	GetQuizSessionByID(id int) (models.QuizSession, error)
	UpdateQuizSession(session models.QuizSession) error
	GetUserActiveQuizSessions(userID int) ([]models.QuizSession, error)

	// questions
	GetQuestionByID(id int) (models.Question, error)
	GetAllQuestions() ([]models.Question, error)
	CreateQuestion(newCase models.QuestionPayload) (models.QuestionPayload, error)
	UpdateQuestionByID(questionID int, updatedCase models.QuestionPayload) (models.QuestionPayload, error)
	DeleteQuestionByID(id int) error
	CountQuestions() (int, error)
	GetQuestionOptions(id int) ([]string, error)
	GetQuestionCorrectOption(id int) (string, error)

	// cases
	CreateCase(newCase models.Case) (models.Case, error)
	UpdateCase(updatedCase models.Case) (models.Case, error)
	DeleteCaseWithParameters(id int) error
	GetAllCases() ([]models.Case, error)
	GetCaseByID(id int) (models.Case, error)
	CreateCaseParameter(caseID int, parameter models.ParameterValue) (models.ParameterValue, error)

	// parameters
	CreateParameter(parameter models.Parameter) (models.Parameter, error)
	UpdateParameter(parameter models.Parameter) error
	DeleteParameter(id int) error
	GetAllParameters() ([]models.Parameter, error)
	GetParameterByID(id int) (models.Parameter, error)
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

func (s *PostgresStorage) Ping() error {
	return s.db.Ping()
}

func (s *PostgresStorage) Close() error {
	return s.db.Close()
}

// Quiz Sessions
func (s *PostgresStorage) CreateQuizSession(session models.QuizSession) (models.QuizSession, error) {
	query := `
        INSERT INTO quiz_sessions (user_id, status, mode, current_question, created_at, updated_at)
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        RETURNING id, created_at, updated_at`

	err := s.db.QueryRow(
		query,
		session.UserID,
		session.Status,
		session.Mode,
		session.CurrentQuestionID,
	).Scan(&session.ID, &session.CreatedAt, &session.UpdatedAt)

	return session, err
}

func (s *PostgresStorage) GetQuizSessionByID(id int) (models.QuizSession, error) {
	var session models.QuizSession
	query := `
        SELECT id, user_id, status, mode, current_question, created_at, updated_at, finished_at
        FROM quiz_sessions
        WHERE id = $1`

	err := s.db.QueryRow(query, id).Scan(
		&session.ID,
		&session.UserID,
		&session.Status,
		&session.Mode,
		&session.CurrentQuestionID,
		&session.CreatedAt,
		&session.UpdatedAt,
		&session.FinishedAt,
	)

	return session, err
}

func (s *PostgresStorage) UpdateQuizSession(session models.QuizSession) error {
	query := `
        UPDATE quiz_sessions
        SET status = $1, mode = $2, current_question = $3, updated_at = NOW(), finished_at = $4
        WHERE id = $5`

	_, err := s.db.Exec(
		query,
		session.Status,
		session.Mode,
		session.CurrentQuestionID,
		session.FinishedAt,
		session.ID,
	)

	return err
}

func (s *PostgresStorage) GetUserActiveQuizSessions(userID int) ([]models.QuizSession, error) {
	query := `
        SELECT id, user_id, status, mode, current_question, created_at, updated_at, finished_at
        FROM quiz_sessions
        WHERE user_id = $1 and status != 'finished'
        ORDER BY created_at DESC`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.QuizSession
	for rows.Next() {
		var session models.QuizSession
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.Status,
			&session.Mode,
			&session.CurrentQuestionID,
			&session.CreatedAt,
			&session.UpdatedAt,
			&session.FinishedAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, rows.Err()
}

// Questions
func (s *PostgresStorage) GetQuestionByID(id int) (models.Question, error) {
	query := `
        SELECT q.id, q.question, q.prediction_age,
               c.id, c.code, c.patient_gender, c.age1, c.age2
        FROM questions q
        JOIN cases c ON q.case_id = c.id
        WHERE q.id = $1`

	var question models.Question
	err := s.db.QueryRow(query, id).Scan(
		&question.ID,
		&question.Question,
		&question.PredictionAge,
		&question.Case.ID,
		&question.Case.Code,
		&question.Case.Gender,
		&question.Case.Age1,
		&question.Case.Age2,
	)
	if err != nil {
		return question, err
	}

	question.Options, err = s.GetQuestionOptions(id)
	if err != nil {
		return question, err
	}

	question.Case.Parameters, question.Case.ParameterValues, err = s.getCaseParameters(question.Case.ID)
	if err != nil {
		return question, err
	}

	return question, nil
}

func (s *PostgresStorage) GetQuestionOptions(id int) ([]string, error) {
	query := `
		SELECT o.option from options o
			JOIN question_options qo on o.id = qo.option_id
			WHERE qo.question_id = $1`

	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	// todo: add answers
	question.Answers = make([]string, 0)

	// Get case parameters
	question.Case.Parameters, question.Case.ParameterValues, err = s.getCaseParameters(question.Case.ID)
	if err != nil {
		return question, err
	}

	return question, nil
}
func (s *PostgresStorage) GetAllQuestions() ([]models.Question, error) {
	query := `
				SELECT q.id, q.question, q.prediction_age,
               c.id, c.code, c.patient_gender, c.age1, c.age2
        FROM questions q
        JOIN cases c ON q.case_id = c.id`

	var options []string
	for rows.Next() {
		var option string
		err := rows.Scan(&option)
		if err != nil {
			return nil, err
		}
		options = append(options, option)
	}

	return options, rows.Err()
}
func (s *PostgresStorage) GetQuestionCorrectOption(id int) (string, error) {
	query := `
		SELECT o.option from options o
			JOIN question_options qo on o.id = qo.option_id
			WHERE qo.question_id = $1 and qo.is_correct = true`

	var option string
	err := s.db.QueryRow(query, id).Scan(&option)
	return option, err
}
func (s *PostgresStorage) GetAllQuestions() ([]models.Question, error) {
	query := `
				SELECT q.id, q.question, q.prediction_age,
               c.id, c.code, c.patient_gender, c.age1, c.age2
        FROM questions q
        JOIN cases c ON q.case_id = c.id`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	var questions []models.Question
	for rows.Next() {
		var question models.Question
		err = rows.Scan(
			&question.ID,
			&question.Question,
			&question.PredictionAge,
			&question.Case.ID,
			&question.Case.Code,
			&question.Case.Gender,
			&question.Case.Age1,
			&question.Case.Age2)
		if err != nil {
			return nil, err
		}
		question.Case.Parameters, question.Case.ParameterValues, err = s.getCaseParameters(question.Case.ID)
		if err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}
	defer rows.Close()
	return questions, nil
}
func (s *PostgresStorage) CreateQuestion(payload models.QuestionPayload) (models.QuestionPayload, error) {
	query := `
        INSERT INTO questions (question, prediction_age, case_id)
        VALUES ($1, $2, $3)
        RETURNING id`

	err := s.db.QueryRow(
		query,
		payload.Question,
		payload.PredictionAge,
		payload.CaseID,
	).Scan(&payload.ID)

	return payload, err
}

func (s *PostgresStorage) UpdateQuestionByID(questionID int, payload models.QuestionPayload) (models.QuestionPayload, error) {
	query := `
        UPDATE questions
        SET question = $1, prediction_age = $2, case_id = $3
        WHERE id = $4`

	_, err := s.db.Exec(
		query,
		payload.Question,
		payload.PredictionAge,
		payload.CaseID,
		questionID,
	)

	payload.ID = questionID
	return payload, err
}

func (s *PostgresStorage) DeleteQuestionByID(id int) error {
	query := "DELETE FROM questions WHERE id = $1"
	_, err := s.db.Exec(query, id)
	return err
}

func (s *PostgresStorage) CountQuestions() (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM questions").Scan(&count)
	return count, err
}

// Cases
func (s *PostgresStorage) CreateCase(newCase models.Case) (models.Case, error) {
	query := `
        INSERT INTO cases (code, patient_gender, age1, age2)
        VALUES ($1, $2, $3, $4)
        RETURNING id`

	err := s.db.QueryRow(
		query,
		newCase.Code,
		newCase.Gender,
		newCase.Age1,
		newCase.Age2,
	).Scan(&newCase.ID)

	return newCase, err
}

func (s *PostgresStorage) UpdateCase(updatedCase models.Case) (models.Case, error) {
	query := `
        UPDATE cases
        SET code = $1, patient_gender = $2, age1 = $3, age2 = $4
        WHERE id = $5`

	_, err := s.db.Exec(
		query,
		updatedCase.Code,
		updatedCase.Gender,
		updatedCase.Age1,
		updatedCase.Age2,
		updatedCase.ID,
	)

	return updatedCase, err
}

func (s *PostgresStorage) DeleteCaseWithParameters(id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM case_parameters WHERE case_id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM cases WHERE id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *PostgresStorage) GetAllCases() ([]models.Case, error) {
	query := `
        SELECT id, code, patient_gender, age1, age2
        FROM cases
        ORDER BY id`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cases []models.Case
	for rows.Next() {
		var c models.Case
		err = rows.Scan(
			&c.ID,
			&c.Code,
			&c.Gender,
			&c.Age1,
			&c.Age2,
		)
		if err != nil {
			return nil, err
		}

		// Get parameters for each case
		c.Parameters, c.ParameterValues, err = s.getCaseParameters(c.ID)
		if err != nil {
			return nil, err
		}

		cases = append(cases, c)
	}

	return cases, rows.Err()
}

func (s *PostgresStorage) GetCaseByID(id int) (models.Case, error) {
	query := `
        SELECT id, code, patient_gender, age1, age2
        FROM cases
        WHERE id=$1`

	var c models.Case
	err := s.db.QueryRow(query, id).Scan(
		&c.ID,
		&c.Code,
		&c.Gender,
		&c.Age1,
		&c.Age2)
	if err != nil {
		return c, err
	}
	c.Parameters, c.ParameterValues, err = s.getCaseParameters(c.ID)
	if err != nil {
		return c, err
	}
	return c, nil
}
func (s *PostgresStorage) CreateCaseParameter(caseID int, parameter models.ParameterValue) (models.ParameterValue, error) {
	query := `
		INSERT INTO case_parameters (case_id, parameter_id, value_1, value_2)
		VALUES ($1, $2, $3, $4)
		RETURNING parameter_id, value_1, value_2`

	err := s.db.QueryRow(
		query,
		caseID,
		parameter.ParameterID,
		parameter.Value1,
		parameter.Value2,
	).Scan(&parameter.ParameterID, &parameter.Value1, &parameter.Value2)

	return parameter, err
}

// Parameters
func (s *PostgresStorage) CreateParameter(parameter models.Parameter) (models.Parameter, error) {
	query := `
        INSERT INTO parameters (name, description, reference_value)
        VALUES ($1, $2, $3)
        RETURNING id`

	err := s.db.QueryRow(
		query,
		parameter.Name,
		parameter.Description,
		parameter.ReferenceValues,
	).Scan(&parameter.ID)

	return parameter, err
}

func (s *PostgresStorage) UpdateParameter(parameter models.Parameter) error {
	query := `
        UPDATE parameters
        SET name = $1, description = $2, reference_value = $3
        WHERE id = $4`

	_, err := s.db.Exec(
		query,
		parameter.Name,
		parameter.Description,
		parameter.ReferenceValues,
		parameter.ID,
	)

	return err
}

func (s *PostgresStorage) DeleteParameter(id int) error {
	query := "DELETE FROM parameters WHERE id = $1"
	_, err := s.db.Exec(query, id)
	return err
}

func (s *PostgresStorage) GetParameterByID(id int) (models.Parameter, error) {
	query := `
		SELECT id, name, description, reference_value
		FROM parameters
		WHERE id = $1`

	var p models.Parameter
	err := s.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.ReferenceValues,
	)
	return p, err
}

func (s *PostgresStorage) GetAllParameters() ([]models.Parameter, error) {
	query := `
        SELECT id, name, description, reference_value
        FROM parameters
        ORDER BY id`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parameters []models.Parameter
	for rows.Next() {
		var p models.Parameter
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.ReferenceValues,
		)
		if err != nil {
			return nil, err
		}
		parameters = append(parameters, p)
	}

	return parameters, rows.Err()
}

// Helper functions
func (s *PostgresStorage) getCaseParameters(caseID int) ([]models.Parameter, []models.ParameterValue, error) {
	query := `select cp.parameter_id, cp.value_1, cp.value_2, p.description, p.name, p.reference_value from cases c
		join case_parameters cp on c.id = cp.case_id
		join parameters p on cp.parameter_id = p.id
		where c.id=$1;`
	rows, err := s.db.Query(query, caseID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	parameters := make([]models.Parameter, 0)
	parameterValues := make([]models.ParameterValue, 0)
	for rows.Next() {
		var p models.Parameter
		var pv models.ParameterValue
		err := rows.Scan(&p.ID, &pv.Value1, &pv.Value2, &p.Description, &p.Name, &p.ReferenceValues)
		if err != nil {
			return nil, nil, err
		}
		pv.ParameterID = p.ID
		parameters = append(parameters, p)
		parameterValues = append(parameterValues, pv)
	}
	return parameters, parameterValues, nil
}
