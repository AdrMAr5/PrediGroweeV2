package handlers

import (
	"PrediGroweeV2/quiz/internal/models"
	"PrediGroweeV2/quiz/internal/storage"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type StartQuizHandler struct {
	storage storage.Store
	logger  *zap.Logger
}

func NewStartQuizHandler(store storage.Store, logger *zap.Logger) *StartQuizHandler {
	return &StartQuizHandler{
		storage: store,
		logger:  logger,
	}
}

func (h *StartQuizHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.Context().Value("userID").(string))
	if err != nil {
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}
	var payload models.StartQuizPayload
	if err := payload.FromJSON(r.Body); err != nil {
		http.Error(rw, "invalid request payload", http.StatusBadRequest)
		return
	}
	if err := payload.Validate(); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	quizSessionID, err := generateSessionID(32)
	if err != nil {
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}
	quizSession := models.QuizSession{
		ID:     quizSessionID,
		Mode:   payload.Mode,
		UserId: userID,
	}
	sessionCreated, err := h.storage.CreateQuizSession(quizSession)
	if err != nil {
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"session": sessionCreated,
	}
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}
func generateSessionID(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
