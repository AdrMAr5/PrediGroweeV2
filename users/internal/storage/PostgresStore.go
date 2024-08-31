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
	p.logger.Info("creating user")
	return &models.User{
		ID:        0,
		FirstName: "mock",
		LastName:  "user",
		Email:     "email",
		Password:  "hash",
	}, nil
}

func (p *PostgresStorage) GetUserById(id int) (*models.User, error) {
	return nil, nil
}

func (p *PostgresStorage) GetUserByEmail(email string) (*models.User, error) {
	p.logger.Info("get user")
	return &models.User{
		ID:        0,
		FirstName: "mock",
		LastName:  "user",
		Email:     email,
		Password:  "password123",
	}, nil
}
