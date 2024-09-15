package handlers

import (
	"PrediGroweeV2/quiz/internal/storage"
	"go.uber.org/zap"
	"net/http"
)

type StartQuizHandler struct {
	storage storage.Store
	logger  *zap.Logger
}

func NewStartQuizHandler(store storage.Store, logger *zap.Logger) *StartQuizHandler {
	return &StartQuizHandler{
		storage: store,
		logger:  logger,
	}
}

func (h *StartQuizHandler) Handle(rw http.ResponseWriter, r *http.Request) {

}
