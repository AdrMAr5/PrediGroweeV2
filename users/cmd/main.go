package main

import (
	"PrediGroweeV2/users/internal/api"
	"PrediGroweeV2/users/internal/storage"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type JWTClaim struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

var (
	db     *sql.DB
	logger *zap.Logger
	jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
)

func main() {
	// Initialize logger
	var err error
	logger, err = zap.NewProduction()
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
	db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", "localhost", "5432", "postgres", "postgres", "api-db"))
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
	apiServer := api.NewApiServer(":8080", postgresStorage, logger)
	apiServer.Run()
}

//func instrumentHandler(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		start := time.Now()
//		rec := statusRecorder{w, http.StatusOK}
//		next.ServeHTTP(&rec, r)
//		duration := time.Since(start).Seconds()
//		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, fmt.Sprintf("%d", rec.status)).Inc()
//		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
//	})
//}

//type statusRecorder struct {
//	http.ResponseWriter
//	status int
//}
//
//func (rec *statusRecorder) WriteHeader(code int) {
//	rec.status = code
//	rec.ResponseWriter.WriteHeader(code)
//}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("X-User-Email")
	var user User
	err := db.QueryRow("SELECT id, first_name,last_name, email FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		logger.Error("Error fetching user", zap.Error(err))
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func VerifyTokenHandler(w http.ResponseWriter, r *http.Request) {
	var tokenRequest struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&tokenRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := jwt.ParseWithClaims(tokenRequest.Token, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":    true,
		"username": claims.Username,
		"email":    claims.Email,
	})
}
