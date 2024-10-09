package storage

import (
	"auth/internal/models"
	"database/sql"
	"go.uber.org/zap"
)

type Store interface {
	Ping() error
	Close() error
	CreateUser(user *models.User) (*models.User, error)
	GetUserById(id int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	SaveUserSession(token models.UserSession) error
	GetUserSession(userID int) (models.UserSession, error)
	UpdateUserSession(token models.UserSession) error
	GetUserSessionBySessionID(sessionID string) (models.UserSession, error)
}
type FirestoreStorage struct {
	config string
}
type PostgresStorage struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewPostgresStorage(db *sql.DB, logger *zap.Logger) *PostgresStorage {
	return &PostgresStorage{db, logger}
}
func (p *PostgresStorage) Ping() error {
	return p.db.Ping()
}

func (p *PostgresStorage) Close() error {
	return p.db.Close()
}

func (p *PostgresStorage) CreateUser(user *models.User) (*models.User, error) {
	var userCreated models.User
	err := p.db.QueryRow("INSERT INTO users (first_name, last_name, email, pwd) VALUES ($1, $2, $3, $4) RETURNING id, email, first_name, last_name", user.FirstName, user.LastName, user.Email, user.Password).Scan(&userCreated.ID, &userCreated.Email, &userCreated.FirstName, &userCreated.LastName)
	if err != nil {
		return nil, err
	}
	return &userCreated, nil
}

func (p *PostgresStorage) GetUserById(id int) (*models.User, error) {
	var user models.User
	err := p.db.QueryRow("SELECT id, first_name, last_name, email, role FROM users WHERE id = $1", id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *PostgresStorage) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := p.db.QueryRow("SELECT id, first_name, last_name, email, pwd, role FROM users WHERE email = $1", email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *PostgresStorage) SaveUserSession(session models.UserSession) error {
	_, err := p.db.Exec("INSERT INTO users_sessions (user_id, session_id, expiration) VALUES ($1, $2, $3)", session.UserID, session.SessionID, session.Expiration)
	return err
}
func (p *PostgresStorage) UpdateUserSession(session models.UserSession) error {
	_, err := p.db.Exec("UPDATE users_sessions SET session_id = $1, expiration = $2 WHERE user_id = $3", session.SessionID, session.Expiration, session.UserID)
	return err
}
func (p *PostgresStorage) GetUserSession(userID int) (models.UserSession, error) {
	var session models.UserSession
	err := p.db.QueryRow("SELECT user_id, session_id, expiration FROM users_sessions WHERE user_id = $1", userID).Scan(&session.UserID, &session.SessionID, &session.Expiration)
	return session, err
}
func (p *PostgresStorage) GetUserSessionBySessionID(sessionID string) (models.UserSession, error) {
	var session models.UserSession
	err := p.db.QueryRow("SELECT user_id, session_id, expiration FROM users_sessions WHERE session_id = $1", sessionID).Scan(&session.UserID, &session.SessionID, &session.Expiration)
	return session, err
}
