package storage

import (
	"database/sql"
	"go.uber.org/zap"
	"images/internal/models"
)

type Store interface {
	Ping() error
	Close() error
	GetImageById(id int) (models.Image, error)
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

// todo: implement db methods
func (p *PostgresStorage) GetImageById(id int) (models.Image, error) {
	return models.Image{ID: 1}, nil
}
