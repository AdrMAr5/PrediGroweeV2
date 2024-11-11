package handlers

import (
	"go.uber.org/zap"
	"net/http"
	"quiz/internal/storage"
	"strconv"
)

type GetNextQuestionHandler struct {
	storage storage.Store
	logger  *zap.Logger
}

func NewGetNextQuestionHandler(storage storage.Store, logger *zap.Logger) *GetNextQuestionHandler {
	return &GetNextQuestionHandler{
		storage: storage,
		logger:  logger,
	}
}

func (h *GetNextQuestionHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	h.logger.Info("get next question handler")
	sessionId := r.PathValue("quizSessionId")
	if sessionId == "" {
		http.Error(rw, "missing session id", http.StatusBadRequest)
		return
	}
	sessionID, err := strconv.Atoi(sessionId)
	if err != nil {
		http.Error(rw, "invalid session id", http.StatusBadRequest)
		return
	}
	userID := r.Context().Value("user_id").(int)
	session, err := h.storage.GetQuizSessionByID(sessionID)
	if err != nil {
		h.logger.Error("failed to get session", zap.Error(err))
		http.Error(rw, "failed to get session", http.StatusNotFound)
		return
	}
	if session.UserID != userID {
		http.Error(rw, "failed to get session", http.StatusNotFound)
		return
	}
	if session.Status == "finished" {
		http.Error(rw, "quiz is finished", http.StatusNotFound)
		return
	}
	question, err := h.storage.GetQuestionByID(session.CurrentQuestionID)
	if err != nil {
		http.Error(rw, "failed to get question", http.StatusNotFound)
		return
	}
	// reset the value3 field to 0 not to expose the parameters of correct answer
	for _, pv := range question.Case.ParameterValues {
		pv.Value3 = 0
	}
	err = question.ToJSON(rw)
	if err != nil {
		http.Error(rw, "failed to get question", http.StatusInternalServerError)
		return
	}
}
