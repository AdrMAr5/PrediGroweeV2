package storage

import (
	"database/sql"
	"go.uber.org/zap"
)

type Storage interface {
	Ping() error
	Close() error
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
