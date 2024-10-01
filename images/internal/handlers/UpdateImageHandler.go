package handlers

import (
	"go.uber.org/zap"
	"images/internal/storage"
	"net/http"
)

type UpdateImageHandler struct {
	storage storage.Store
	logger  *zap.Logger
}

func NewUpdateImageHandler(store storage.Store, logger *zap.Logger) *UpdateImageHandler {
	return &UpdateImageHandler{
		storage: store,
		logger:  logger,
	}
}
func (h *UpdateImageHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotImplemented)
}
