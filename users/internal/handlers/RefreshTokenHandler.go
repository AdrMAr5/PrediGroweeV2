package handlers

import (
	"PrediGroweeV2/users/internal/auth"
	"PrediGroweeV2/users/internal/storage"
	"go.uber.org/zap"
	"net/http"
	"strconv"
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
	userID := r.Context().Value("user_id").(int)
	accessToken, err := auth.GenerateAccessToken(strconv.Itoa(userID))
	if err != nil {
		h.logger.Error("Error generating access token", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Path:     "/",
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   false, // Set to true if using HTTPS
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(15 * time.Minute),
	})
}
