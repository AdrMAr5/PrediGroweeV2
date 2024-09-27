package handlers

import (
	"PrediGroweeV2/quiz/internal/storage"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
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
	question, err := h.Store.GetQuestionById(id)
	if err != nil {
		h.logger.Error("failed to get question from db", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}

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
