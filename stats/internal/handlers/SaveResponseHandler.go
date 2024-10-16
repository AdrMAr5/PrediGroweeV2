package handlers

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"stats/internal/models"
	"stats/internal/storage"
)

type SaveResponseHandler struct {
	storage storage.Storage
	logger  *zap.Logger
}

func NewSaveResponseHandler(storage storage.Storage, logger *zap.Logger) *SaveResponseHandler {
	return &SaveResponseHandler{
		storage: storage,
		logger:  logger,
	}
}

func (h *SaveResponseHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	h.logger.Info("SaveResponseHandler.Handle")
	var response models.QuestionResponse
	err := response.FromJSON(r.Body)
	h.logger.Info(fmt.Sprintf("response to save: %+v", response))
	if err != nil {
		h.logger.Error("failed to decode response", zap.Error(err))
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = h.storage.GetSession(response.SessionID)
	if err == storage.ErrSessionNotFound {
		err = h.storage.SaveSession(&models.QuizSession{
			SessionID: response.SessionID,
			UserID:    response.UserID,
			QuizMode:  "educational",
		})
		if err != nil {
			h.logger.Error("failed to save session", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	err = h.storage.SaveResponse(&response)
	rw.WriteHeader(http.StatusOK)
}
