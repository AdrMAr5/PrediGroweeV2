package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
	"os"
	"quiz/internal/api"
	"quiz/internal/clients"
	"quiz/internal/storage"
	"time"
)

const PingDbAttempts = 3

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
	logger.Info("Connected to database")

	// Set up database connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify database connection
	for i := 1; i <= PingDbAttempts; i++ {
		err = db.Ping()
		if err == nil {
			break
		} else {
			logger.Error(fmt.Sprintf("Failed to Ping the database (attempt: %d/%d)", i, PingDbAttempts), zap.Error(err))
		}
		time.Sleep(2 * time.Second)
	}
	if err = db.Ping(); err != nil {
		logger.Fatal("Failed to ping database, exiting", zap.Error(err))
	}
	postgresStorage := storage.NewPostgresStorage(db, logger)
	authClient := clients.NewAuthClient("http://auth:8080/auth", logger)
	logger.Info("Connected to auth service")
	statsClient := clients.NewStatsClient("http://stats:8080/stats", logger)
	logger.Info("Connected to stats service")
	apiServer := api.NewApiServer(":8080", postgresStorage, logger, authClient, statsClient)
	apiServer.Run()
}
func connectToPostgres() (*sql.DB, error) {
	//env := os.Getenv("ENV")
	//sslMode := "require"
	//if env == "local" {
	//	sslMode = "disable"
	//}
	sslMode := "disable"
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbName, sslMode)
	return sql.Open("postgres", connString)
}
