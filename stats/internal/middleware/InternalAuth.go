package middleware

import (
	"go.uber.org/zap"
	"net/http"
)

// todo: implement InternalAuth middleware
func InternalAuth(next http.HandlerFunc, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("InternalAuth middleware")
		next(w, r)
	}
}
