package handlers

import (
	"go.uber.org/zap"
	"images/internal/storage"
	"net/http"
	"strconv"
)

type GetImageHandler struct {
	storage storage.Store
	logger  *zap.Logger
}

func NewGetImageHandler(store storage.Store, logger *zap.Logger) *GetImageHandler {
	return &GetImageHandler{
		storage: store,
		logger:  logger,
	}
}
func (h *GetImageHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(rw, "Invalid image id", http.StatusBadRequest)
		return
	}
	_, err = h.storage.GetImageById(id)
	if err != nil {
		http.Error(rw, "Image not found", http.StatusNotFound)
		return
	}
	rw.WriteHeader(http.StatusOK)
}
