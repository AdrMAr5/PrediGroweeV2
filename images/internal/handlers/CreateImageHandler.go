package handlers

import (
	"PrediGroweeV2/images/internal/storage"
	"go.uber.org/zap"
	"net/http"
)

type CreateImageHandler struct {
	storage storage.Store
	logger  *zap.Logger
}

func NewCreateImageHandler(store storage.Store, logger *zap.Logger) *CreateImageHandler {
	return &CreateImageHandler{
		storage: store,
		logger:  logger,
	}
}
func (h *CreateImageHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotImplemented)
}
