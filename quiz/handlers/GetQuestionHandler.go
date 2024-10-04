package handlers

import (
	"go.uber.org/zap"
	"net/http"
	"quiz/internal/storage"
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
