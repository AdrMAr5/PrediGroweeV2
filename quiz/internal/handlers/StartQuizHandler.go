package handlers

import (
	"encoding/json"
	"fmt"
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
	groupID, err := h.storage.GetNextQuestionGroupID(0)
	order, err := h.storage.GetGroupQuestionsIDsRandomOrder(groupID)
	if err != nil {
		h.logger.Error("failed to start quiz", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}

	newQuizSession := models.QuizSession{
		Mode:              payload.Mode,
		UserID:            userID,
		Status:            models.QuizStatusNotStarted,
		ScreenSize:        fmt.Sprintf("%dx%d", payload.ScreenWidth, payload.ScreenHeight),
		CurrentQuestionID: order[0],
		CurrentGroup:      groupID,
		GroupOrder:        order,
	}
	session, err := h.storage.GetUserLastQuizSession(userID)
	if err == nil && session != nil {
		newQuizSession.CurrentQuestionID = session.CurrentQuestionID
		newQuizSession.CurrentGroup = session.CurrentGroup
		newQuizSession.GroupOrder = session.GroupOrder
		session.FinishedAt = session.UpdatedAt
		session.Status = models.QuizStatusFinished
		err = h.storage.UpdateQuizSession(*session)
	}

	sessionCreated, err := h.storage.CreateQuizSession(newQuizSession)
	if err != nil {
		h.logger.Error("failed to create quiz session in db", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}
	err = h.statsClient.SaveSession(sessionCreated)
	if err != nil {
		h.logger.Error("failed to save session in stats service", zap.Error(err))
	}
	timeLimit, err := h.storage.GetTimeLimit()
	if err != nil {
		h.logger.Error("failed to get time limit", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"session":    sessionCreated,
		"time_limit": timeLimit,
	}
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		h.logger.Error("failed to encode response", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}
}
