package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"images/internal/api"
	"images/internal/clients"
	"log"
	"os"
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

	authClient := clients.NewAuthClient("http://auth:8080/auth", logger)
	apiServer := api.NewApiServer(":8080", logger, authClient)
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
