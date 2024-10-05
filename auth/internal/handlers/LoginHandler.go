package handlers

import (
	"auth/internal/auth"
	"auth/internal/models"
	"auth/internal/storage"
	"encoding/json"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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
		h.logger.Error("Error getting user by email", zap.Error(err))
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	// todo: uncomment this block after testing phase
	//userSession, err := h.store.GetUserSession(dbUser.ID)
	//if err == nil {
	//	if userSession.Expiration.After(time.Now()) {
	//		http.Error(w, "User already logged in", http.StatusConflict)
	//		return
	//	}
	//}
	if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(userPayload.Password)); err != nil {
		h.logger.Error("Error comparing passwords", zap.Error(err))
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	sessionId, err := auth.GenerateSessionID(64)
	if err != nil {
		h.logger.Error("Error generating session id", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = h.store.SaveUserSession(models.UserSession{
		UserID:     dbUser.ID,
		SessionID:  sessionId,
		Expiration: time.Now().Add(7 * 24 * time.Hour),
	})
	accessToken, err := auth.GenerateAccessToken(strconv.Itoa(dbUser.ID))
	if err != nil {
		h.logger.Error("Error generating access token", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Path:     "/",
		Name:     "session_id",
		Value:    sessionId,
		HttpOnly: true,
		Secure:   false, // Set to true if using HTTPS
	})
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken, "message": "Login successful"})
	if err != nil {
		h.logger.Error("Error encoding response", zap.Error(err))
		return
	}
}
