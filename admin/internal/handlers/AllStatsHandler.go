package handlers

import (
	"admin/internal/storage"
	"go.uber.org/zap"
	"net/http"
)

type AllStatsHandler struct {
	storage storage.Storage
	logger  *zap.Logger
}

func NewAllStatsHandler(storage storage.Storage, logger *zap.Logger) *AllStatsHandler {
	return &AllStatsHandler{
		storage: storage,
		logger:  logger,
	}
}
func (h *AllStatsHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(200)
}
