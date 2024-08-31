package handlers

import (
	"PrediGroweeV2/users/internal/auth"
	"PrediGroweeV2/users/internal/models"
	"PrediGroweeV2/users/internal/storage"
	"encoding/json"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
)

type RegisterHandler struct {
	store  storage.Store
	logger *zap.Logger
}

func NewRegisterHandler(store storage.Store, logger *zap.Logger) *RegisterHandler {
	return &RegisterHandler{
		logger: logger,
		store:  store,
	}
}

func (h *RegisterHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := user.FromJSON(r.Body)
	if err != nil {
		h.logger.Error("Error unmarshalling json", zap.Error(err))
		http.Error(w, "Invalid user", http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Error("Error hashing password", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)
	userCreated, err := h.store.CreateUser(&user)
	if err != nil {
		h.logger.Error("Error creating user", zap.Error(err))
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	accessToken, err := auth.GenerateAccessToken(strconv.Itoa(userCreated.ID))
	if err != nil {
		h.logger.Error("Error generating access token", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := auth.GenerateRefreshToken(strconv.Itoa(userCreated.ID))
	if err != nil {
		h.logger.Error("Error generating refresh token", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		HttpOnly: true,
		Secure:   false, // Set to true if using HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(3 * 24 * time.Hour),
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   false, // Set to true if using HTTPS
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(15 * time.Minute),
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"user":  userCreated,
		"token": accessToken,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Error writing response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
