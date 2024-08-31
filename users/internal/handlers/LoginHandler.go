package handlers

import (
	"PrediGroweeV2/users/internal/auth"
	"PrediGroweeV2/users/internal/models"
	"PrediGroweeV2/users/internal/storage"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type LoginHandler struct {
	store  storage.Store
	logger *zap.Logger
}

func NewLoginHandler(store storage.Store, logger *zap.Logger) *LoginHandler {
	return &LoginHandler{
		logger: logger,
		store:  store,
	}
}

func (h *LoginHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var userPayload models.LoginUserPayload
	if err := userPayload.FromJSON(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dbUser, err := h.store.GetUserByEmail(userPayload.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	//if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(userPayload.Password)); err != nil {
	//	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	//	return
	//}
	if dbUser.Password != userPayload.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	accessToken, err := auth.GenerateAccessToken(strconv.Itoa(dbUser.ID))
	if err != nil {
		h.logger.Error("Error generating access token", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := auth.GenerateRefreshToken(strconv.Itoa(dbUser.ID))
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
	err = json.NewEncoder(w).Encode(map[string]string{"token": accessToken, "message": "Login successful"})
	if err != nil {
		h.logger.Error("Error encoding response", zap.Error(err))
		return
	}
}
