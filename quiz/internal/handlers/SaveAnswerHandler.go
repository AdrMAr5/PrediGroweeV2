package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"quiz/internal/clients"
	"quiz/internal/models"
	"quiz/internal/storage"
	"strconv"
)

type SubmitAnswerHandler struct {
	storage     storage.Store
	logger      *zap.Logger
	statsClient *clients.StatsClient
}

func NewSubmitAnswerHandler(store storage.Store, logger *zap.Logger, statsClient *clients.StatsClient) *SubmitAnswerHandler {
	return &SubmitAnswerHandler{
		storage:     store,
		logger:      logger,
		statsClient: statsClient,
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
	userID := r.Context().Value("user_id").(int)
	if userID != session.UserID {
		http.Error(rw, "internal server error", http.StatusInternalServerError)
	}
	if session.Mode == models.QuizModeEducational {
		// todo: handle return correct answer
		h.logger.Info("educational mode, returning correct answer")
	}
	session.CurrentQuestionID = session.CurrentQuestionID + 1
	err = h.storage.UpdateQuizSession(session)
	if err != nil {
		h.logger.Error("failed to update quiz session", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}
	var answer models.QuestionAnswer
	err = json.NewDecoder(r.Body).Decode(&answer)
	if err != nil {
		h.logger.Error("failed to decode answer", zap.Error(err))
		http.Error(rw, "invalid answer", http.StatusBadRequest)
		return
	}
	totalQuestions, err := h.storage.CountQuestions()
	if err != nil {
		h.logger.Error("failed to count questions", zap.Error(err))
		totalQuestions = 50
	}

	err = h.statsClient.SaveResponse(models.QuestionAnswer{
		QuestionID:    session.CurrentQuestionID,
		UserID:        session.UserID,
		SessionID:     session.ID,
		Answer:        answer.Answer,
		IsFirstAnswer: session.CurrentQuestionID == 1,
		IsLastAnswer:  session.CurrentQuestionID == totalQuestions,
	})
	if err != nil {
		h.logger.Error("failed to save response", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}
