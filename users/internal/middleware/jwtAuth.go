package middleware

import (
	"PrediGroweeV2/users/internal/auth"
	"PrediGroweeV2/users/internal/storage"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strconv"
)

func WithJWTAuth(next http.HandlerFunc, storage storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}
		token, err := auth.ValidateJWT(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		tokenClaims := token.Claims.(jwt.MapClaims)
		userID, err := strconv.Atoi(tokenClaims["sub"].(string))
		if err != nil {
			log.Println("Error parsing user_id", err)
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		}
		// check if user exists
		_, err = storage.GetUserById(userID)
		if err != nil {
			log.Println("Error getting user", err)
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		}
		r.Header.Set("X-User-ID", tokenClaims["sub"].(string))

		next(w, r)
	}
}
