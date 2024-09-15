package handlers

import (
	"PrediGroweeV2/quiz/internal/storage"
	"go.uber.org/zap"
	"net/http"
)

type FinishQuizHandler struct {
	storage storage.Store
	logger  *zap.Logger
}

func NewFinishQuizHandler(store storage.Store, logger *zap.Logger) *FinishQuizHandler {
	return &FinishQuizHandler{
		storage: store,
		logger:  logger,
	}
}
func (h *FinishQuizHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
}
