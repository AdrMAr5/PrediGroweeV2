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
	SaveResponse(sessionID int, response *models.QuestionResponse) error
	SaveSession(session *models.QuizSession) error
	GetUserStatsForMode(userID int, mode models.QuizMode) (correctCount int, wrongCount int, err error)
	GetQuizSessionByID(quizSessionID int) (*models.QuizSession, error)
	GetQuizQuestionsStats(quizSessionID int) ([]models.QuestionStats, error)
	GetUserQuizStats(quizSessionID int) (*models.QuizStats, error)
	FinishQuizSession(quizSessionID int) error
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
func (p *PostgresStorage) SaveResponse(sessionID int, response *models.QuestionResponse) error {
	_, err := p.db.Exec(`INSERT INTO answers (session_id, question_id, answer, correct) values ($1, $2, $3, $4)`, sessionID, response.QuestionID, response.Answer, response.IsCorrect)
	if err != nil {
		return err
	}
	return nil
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
func (p *PostgresStorage) GetQuizSessionByID(quizSessionID int) (*models.QuizSession, error) {
	var session models.QuizSession
	err := p.db.QueryRow(`SELECT user_id, session_id, finish_time, quiz_mode FROM quiz_sessions WHERE session_id = $1`, quizSessionID).Scan(&session.UserID, &session.SessionID, &session.FinishTime, &session.QuizMode)
	if err == sql.ErrNoRows {
		return nil, ErrSessionNotFound
	}
	return &session, nil
}
func (p *PostgresStorage) GetQuizQuestionsStats(quizSessionID int) ([]models.QuestionStats, error) {
	query := `select a.question_id, a.answer, a.correct from answers a
where a.session_id = $1`
	rows, err := p.db.Query(query, quizSessionID)
	if err != nil {
		return nil, err
	}
	var questionsStats []models.QuestionStats
	for rows.Next() {
		var qs models.QuestionStats
		err = rows.Scan(&qs.QuestionID, &qs.Answer, &qs.IsCorrect)
		if err != nil {
			return nil, err
		}
		questionsStats = append(questionsStats, qs)
	}
	return questionsStats, nil
}
func (p *PostgresStorage) GetUserQuizStats(quizSessionID int) (*models.QuizStats, error) {
	var quizStats models.QuizStats
	err := p.db.QueryRow(`select quiz_mode, count(*) as total_questions, sum(CASE WHEN correct THEN 1 END) as correct_answers from answers a
join quiz_sessions s on a.session_id = s.session_id
where a.session_id = $1
group by quiz_mode`, quizSessionID).Scan(&quizStats.Mode, &quizStats.TotalQuestions, &quizStats.CorrectAnswers)
	if err != nil {
		return nil, err
	}
	if quizStats.TotalQuestions != 0 {
		quizStats.Accuracy = float64(quizStats.CorrectAnswers) / float64(quizStats.TotalQuestions)
	} else {
		quizStats.Accuracy = 0
	}
	quizStats.Questions, err = p.GetQuizQuestionsStats(quizSessionID)
	if err != nil {
		return nil, err
	}
	return &quizStats, nil
}
func (p *PostgresStorage) FinishQuizSession(quizSessionID int) error {
	_, err := p.db.Exec(`UPDATE quiz_sessions SET finish_time = now() WHERE session_id = $1`, quizSessionID)
	if err != nil {
		return err
	}
	return nil
}
