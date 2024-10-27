package middleware

import (
	"auth/internal/models"
	"net/http"
)

func WithAdminRole(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value("user_role").(models.UserRole)
		if !ok {
			http.Error(w, "User role not found", http.StatusUnauthorized)
			return
		}
		if role != models.RoleAdmin {
			http.Error(w, "User is not an admin", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}
