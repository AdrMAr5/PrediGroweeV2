package handlers

import (
	"PrediGroweeV2/quiz/internal/storage"
	"go.uber.org/zap"
	"net/http"
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
	rw.WriteHeader(http.StatusNotImplemented)
}
