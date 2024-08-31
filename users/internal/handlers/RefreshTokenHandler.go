package handlers

import (
	"PrediGroweeV2/users/internal/auth"
	"PrediGroweeV2/users/internal/storage"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type RefreshTokenHandler struct {
	store  storage.Store
	logger *zap.Logger
}

func NewRefreshTokenHandler(store storage.Store, logger *zap.Logger) *RefreshTokenHandler {
	return &RefreshTokenHandler{
		store:  store,
		logger: logger,
	}
}

func (h *RefreshTokenHandler) Handle(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.ExtractRefreshTokenFromRequest(r)
	h.logger.Info("refresh token", zap.String("token", refreshToken))
	if err != nil {
		h.logger.Error("Error extracting refresh token", zap.Error(err))
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}
	token, err := auth.ValidateJWT(refreshToken)
	if err != nil || !token.Valid {
		h.logger.Error("Error validating refresh token", zap.Error(err))
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}
	userID, err := token.Claims.GetSubject()
	if err != nil {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}
	accessToken, err := auth.GenerateAccessToken(userID)
	if err != nil {
		h.logger.Error("Error generating access token", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := auth.GenerateRefreshToken(userID)
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
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})
	if err != nil {
		h.logger.Error("Error encoding response", zap.Error(err))
		return
	}
}
