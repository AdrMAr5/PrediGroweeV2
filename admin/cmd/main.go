package main

import (
	"admin/clients"
	"admin/internal/api"
	"admin/internal/storage"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
	"os"
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
	db, err := connectToPostgres()
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
	authClient := clients.NewAuthClient("http://auth:8080/auth", logger)
	statsClient := clients.NewStatsClient("http://stats:8080/stats", logger)
	apiServer := api.NewApiServer(":8080", postgresStorage, logger, authClient, statsClient)
	apiServer.Run()
}
func connectToPostgres() (*sql.DB, error) {
	env := os.Getenv("ENV")
	sslMode := "require"
	if env == "local" {
		sslMode = "disable"
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbName, sslMode)
	return sql.Open("postgres", connString)
}
