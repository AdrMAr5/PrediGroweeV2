package handlers

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"quiz/internal/models"
	"quiz/internal/storage"
	"strconv"
)

type GetQuestionHandler struct {
	Store  storage.Store
	logger *zap.Logger
}

func NewGetQuestionHandler(store storage.Store, logger *zap.Logger) *GetQuestionHandler {
	return &GetQuestionHandler{
		Store:  store,
		logger: logger,
	}
}
func (h *GetQuestionHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	sessionId := r.PathValue("quizSessionId")
	if sessionId == "" {
		h.logger.Info("no session id provided")
		http.Error(rw, "invalid session id", http.StatusBadRequest)
		return
	}
	sessionID, err := strconv.Atoi(sessionId)
	if err != nil {
		h.logger.Info("invalid session id")
		http.Error(rw, "invalid session id", http.StatusBadRequest)
		return
	}
	session, err := h.Store.GetQuizSessionByID(sessionID)
	if err != nil {
		h.logger.Error("failed to get session from db", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}
	userID := r.Context().Value("user_id").(int)
	if session.UserID != userID {
		h.logger.Info("user is not allowed to access this session")
		http.Error(rw, fmt.Sprintf("no session with id: %s", sessionId), http.StatusForbidden)
		return
	}
	if session.Status == models.QuizStatusFinished {
		h.logger.Info("session is already finished")
		http.Error(rw, "quiz already finished", http.StatusForbidden)
		return
	}
	questionID := r.PathValue("id")
	if questionID == "" {
		h.logger.Info("no question id provided")
		http.Error(rw, "invalid question id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(questionID)
	if err != nil {
		h.logger.Info("invalid question id")
		http.Error(rw, "invalid question id", http.StatusBadRequest)
		return
	}
	// returning the same question for now due to lack of proper data in db
	question, err := h.Store.GetQuestionByID(1)
	if err != nil {
		h.logger.Error("failed to get question from db", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}
	session.CurrentQuestionID = id
	err = h.Store.UpdateQuizSession(session)

	// todo: query image urls from images service

	rw.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"question": question,
	}
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}
}
