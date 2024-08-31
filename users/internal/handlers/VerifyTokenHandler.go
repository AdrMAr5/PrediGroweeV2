package handlers

import (
	"PrediGroweeV2/users/internal/auth"
	"PrediGroweeV2/users/internal/storage"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
)

type VerifyTokenHandler struct {
	store  storage.Store
	logger *zap.Logger
}

func NewVerifyTokenHandler(store storage.Store, logger *zap.Logger) *VerifyTokenHandler {
	return &VerifyTokenHandler{
		logger: logger,
		store:  store,
	}
}

func (h *VerifyTokenHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(rw, "Missing authorization token", http.StatusUnauthorized)
		return
	}
	token, err := auth.ValidateJWT(tokenString)
	if err != nil || !token.Valid {
		http.Error(rw, "Invalid or expired token", http.StatusUnauthorized)
		return
	}
	tokenClaims := token.Claims.(jwt.MapClaims)
	userID, err := strconv.Atoi(tokenClaims["sub"].(string))
	if err != nil {
		log.Println("Error parsing user_id", err)
		http.Error(rw, "Permission denied", http.StatusForbidden)
		return
	}
	_, err = h.store.GetUserById(userID)
	if err != nil {
		log.Println("Error getting user", err)
		http.Error(rw, "Permission denied", http.StatusForbidden)
		return
	}
	rw.WriteHeader(http.StatusOK)

}
