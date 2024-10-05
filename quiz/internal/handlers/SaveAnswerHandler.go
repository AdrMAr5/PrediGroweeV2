package handlers

import (
	"go.uber.org/zap"
	"net/http"
	"quiz/internal/models"
	"quiz/internal/storage"
	"strconv"
)

type SubmitAnswerHandler struct {
	storage storage.Store
	logger  *zap.Logger
}

func NewSubmitAnswerHandler(store storage.Store, logger *zap.Logger) *SubmitAnswerHandler {
	return &SubmitAnswerHandler{
		storage: store,
		logger:  logger,
	}
}
func (h *SubmitAnswerHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	quizSessionIdString := r.PathValue("quizSessionId")
	if quizSessionIdString == "" {
		h.logger.Info("no quiz session id provided")
		http.Error(rw, "invalid quiz session id", http.StatusBadRequest)
		return
	}
	quizSessionID, err := strconv.Atoi(quizSessionIdString)
	if err != nil {
		h.logger.Info("invalid quiz session id")
		http.Error(rw, "invalid quiz session id", http.StatusBadRequest)
		return
	}
	session, err := h.storage.GetQuizSessionByID(quizSessionID)
	if err != nil {
		h.logger.Error("failed to get quiz session from db", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}
	if session.Mode == models.QuizModeEducational {
		// todo: handle return correct answer
		rw.WriteHeader(http.StatusOK)
	}
	// todo: implement storing answers
	rw.WriteHeader(http.StatusOK)
}
