package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"quiz/internal/storage"
)

type SettingsHandler struct {
	storage storage.Store
	logger  *zap.Logger
}
type Settings struct {
	Name  string
	Value string
}

func NewSettingsHandler(store storage.Store, logger *zap.Logger) *SettingsHandler {
	return &SettingsHandler{
		storage: store,
		logger:  logger,
	}
}

func (h *SettingsHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	var setting []Settings
	err := json.NewDecoder(r.Body).Decode(&setting)
	if err != nil {
		h.logger.Error("failed to decode request", zap.Error(err))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	for _, s := range setting {
		if s.Name == "" {
			h.logger.Error("invalid setting name")
			http.Error(w, "invalid setting name", http.StatusBadRequest)
			return
		}
		err = h.storage.SaveSettings(s.Name, s.Value)
		if err != nil {
			h.logger.Error("failed to save settings", zap.Error(err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (h *SettingsHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.storage.GetSettings()
	if err != nil {
		h.logger.Error("failed to get settings", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(settings)
	if err != nil {
		h.logger.Error("failed to encode response", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
