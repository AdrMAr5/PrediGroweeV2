package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"quiz/internal/clients"
	"quiz/internal/models"
	"quiz/internal/storage"
)

type StartQuizHandler struct {
	storage     storage.Store
	logger      *zap.Logger
	statsClient *clients.StatsClient
}

func NewStartQuizHandler(store storage.Store, logger *zap.Logger, client *clients.StatsClient) *StartQuizHandler {
	return &StartQuizHandler{
		storage:     store,
		logger:      logger,
		statsClient: client,
	}
}

func (h *StartQuizHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	h.logger.Info("Starting quiz session")
	userID := r.Context().Value("user_id").(int)
	var payload models.StartQuizPayload
	if err := payload.FromJSON(r.Body); err != nil {
		http.Error(rw, "invalid request payload", http.StatusBadRequest)
		return
	}
	if err := payload.Validate(); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	quizSession := models.QuizSession{
		Mode:              payload.Mode,
		UserID:            userID,
		Status:            models.QuizStatusNotStarted,
		CurrentQuestionID: 3,
	}
	sessionCreated, err := h.storage.CreateQuizSession(quizSession)
	if err != nil {
		h.logger.Error("failed to create quiz session in db", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}
	err = h.statsClient.SaveSession(sessionCreated)
	if err != nil {
		h.logger.Error("failed to save session in stats service", zap.Error(err))
	}
	rw.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"session": sessionCreated,
	}
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		h.logger.Error("failed to encode response", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}
}
