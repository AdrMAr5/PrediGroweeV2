package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"quiz/internal/storage"
)

type OptionsHandler struct {
	storage storage.Store
	logger  *zap.Logger
}

func NewOptionsHandler(store storage.Store, logger *zap.Logger) *OptionsHandler {
	return &OptionsHandler{
		storage: store,
		logger:  logger,
	}
}

func (h *OptionsHandler) GetAllOptions(w http.ResponseWriter, _ *http.Request) {
	options, err := h.storage.GetAllOptions()
	if err != nil {
		h.logger.Error("Failed to get options", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(options)
	if err != nil {
		h.logger.Error("Failed to encode response", zap.Error(err))
	}
}
