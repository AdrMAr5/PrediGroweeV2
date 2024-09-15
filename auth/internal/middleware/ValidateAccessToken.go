package middleware

import (
	"PrediGroweeV2/auth/internal/auth"
	"PrediGroweeV2/auth/internal/storage"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
)

func ValidateAccessToken(next http.HandlerFunc, storage storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := auth.ExtractAccessTokenFromRequest(r)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		token, err := auth.ValidateJWT(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		tokenClaims := token.Claims.(jwt.MapClaims)
		tokenUserID, err := strconv.Atoi(tokenClaims["sub"].(string))
		if err != nil {
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		}
		sessionUserID := r.Context().Value("user_id")
		if sessionUserID != tokenUserID {
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		}
		_, err = storage.GetUserById(tokenUserID)
		if err != nil {
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}
