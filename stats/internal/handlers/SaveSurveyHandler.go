package handlers

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"stats/internal/models"
	"stats/internal/storage"
)

type SaveSurveyHandler struct {
	storage storage.Storage
	logger  *zap.Logger
}

func NewSaveSurveyHandler(storage storage.Storage, logger *zap.Logger) *SaveSurveyHandler {
	return &SaveSurveyHandler{
		storage: storage,
		logger:  logger,
	}
}

func (h *SaveSurveyHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var surveyResponse models.SurveyResponse
	err := surveyResponse.FromJSON(r.Body)
	if err != nil {
		h.logger.Error("failed to parse survey response", zap.Error(err))
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	fmt.Println(surveyResponse)
	userID := r.Context().Value("user_id").(int)
	_, err = h.storage.GetSurveyResponseForUser(userID)
	if err == nil {
		http.Error(w, "survey response already exists", http.StatusConflict)
		return
	}
	surveyResponse.UserID = userID
	err = h.storage.SaveSurveyResponse(&surveyResponse)
	if err != nil {
		h.logger.Error("failed to save survey response", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
