package handlers

import (
	"PrediGroweeV2/quiz/internal/storage"
	"go.uber.org/zap"
	"net/http"
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
	rw.WriteHeader(http.StatusOK)
}
