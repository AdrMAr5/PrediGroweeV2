package storage

import (
	"PrediGroweeV2/users/internal/models"
	"database/sql"
	"go.uber.org/zap"
)

type Store interface {
	Ping() error
	CreateUser(user *models.User) (*models.User, error)
	GetUserById(id int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
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

func (p *PostgresStorage) CreateUser(user *models.User) (*models.User, error) {
	var userCreated models.User
	err := p.db.QueryRow("INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING id, email, first_name, last_name", user.FirstName, user.LastName, user.Email, user.Password).Scan(&userCreated.ID, &userCreated.Email, &userCreated.FirstName, &userCreated.LastName)
	if err != nil {
		return nil, err
	}
	return &userCreated, nil
}

func (p *PostgresStorage) GetUserById(id int) (*models.User, error) {
	var user models.User
	err := p.db.QueryRow("SELECT id, first_name, last_name, email FROM users WHERE id = $1", id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *PostgresStorage) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := p.db.QueryRow("SELECT id, first_name, last_name, email FROM users WHERE email = $1", email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
