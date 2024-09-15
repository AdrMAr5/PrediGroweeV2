package middleware

import (
	"PrediGroweeV2/quiz/internal/clients"
	"context"
	"log"
	"net/http"
)

func VerifyToken(next http.HandlerFunc, authClient *clients.AuthClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			log.Println("No token cookie provided")
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		userData, err := authClient.VerifyAuthToken(cookie.Value)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "user_id", userData.UserID))
		r = r.WithContext(context.WithValue(r.Context(), "user_role", userData.Role))
		next(w, r)
	}
}
