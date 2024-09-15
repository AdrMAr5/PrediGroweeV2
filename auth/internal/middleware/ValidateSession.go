package middleware

import (
	"PrediGroweeV2/auth/internal/auth"
	"PrediGroweeV2/auth/internal/storage"
	"context"
	"net/http"
	"time"
)

func ValidateSession(next http.HandlerFunc, storage storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := auth.ExtractSessionIDFromRequest(r)
		if err != nil {
			http.Error(w, "Invalid session ID", http.StatusUnauthorized)
			return
		}
		session, err := storage.GetUserSessionBySessionID(sessionID)
		if err != nil {
			http.Error(w, "Invalid session id", http.StatusUnauthorized)
			return
		}
		if session.Expiration.Before(time.Now()) {
			http.Error(w, "Session expired. Please log in", http.StatusUnauthorized)
			return
		}
		newCtx := context.WithValue(r.Context(), "user_id", session.UserID)
		next(w, r.WithContext(newCtx))
	}
}
