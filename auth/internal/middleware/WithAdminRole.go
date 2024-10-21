package middleware

import "net/http"

func WithAdminRole(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value("user_role").(string)
		if !ok {
			http.Error(w, "User role not found", http.StatusUnauthorized)
			return
		}
		if role != "admin" {
			http.Error(w, "User is not an admin", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}
