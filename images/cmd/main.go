package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"images/internal/api"
	"images/internal/clients"
	"images/internal/storage"
	"log"
	"time"
)

func main() {
	// Initialize logger
	var err error
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Fatalf("Failed to sync logger: %v", err)
		}
	}(logger)

	// Initialize database
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", "localhost", "5432", "postgres", "postgres", "quiz_images"))
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Fatal("Failed to close database connection: %v", zap.Error(err))
		}
	}(db)

	// Set up database connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify database connection
	if err = db.Ping(); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}
	postgresStorage := storage.NewPostgresStorage(db, logger)
	authClient := clients.NewAuthClient("http://localhost:8080", logger)
	apiServer := api.NewApiServer(":8081", postgresStorage, logger, authClient)
	apiServer.Run()
}
