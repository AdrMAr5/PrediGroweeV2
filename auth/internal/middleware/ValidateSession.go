package middleware

import (
	"auth/internal/auth"
	"auth/internal/storage"
	"context"
	"log"
	"net/http"
	"time"
)

func ValidateSession(next http.HandlerFunc, storage storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := auth.ExtractSessionIDFromRequest(r)
		log.Println(sessionID)
		if err != nil {
			log.Println("Failed to extract session ID", err)
			http.Error(w, "Invalid session ID", http.StatusUnauthorized)
			return
		}
		session, err := storage.GetUserSessionBySessionID(sessionID)
		if err != nil {
			log.Println("Failed to get session from storage", err)
			http.Error(w, "Invalid session id", http.StatusUnauthorized)
			return
		}
		if session.Expiration.Before(time.Now()) {
			log.Println("Session expired")
			http.Error(w, "Session expired. Please log in", http.StatusUnauthorized)
			return
		}
		newCtx := context.WithValue(r.Context(), "user_id", session.UserID)
		next(w, r.WithContext(newCtx))
	}
}
