package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"quiz/internal/storage"
)

type GetUserSessionsHandler struct {
	storage storage.Store
	logger  *zap.Logger
}

func NewGetUserSessionsHandler(store storage.Store, logger *zap.Logger) *GetUserSessionsHandler {
	return &GetUserSessionsHandler{
		storage: store,
		logger:  logger,
	}
}

func (h *GetUserSessionsHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	sessions, err := h.storage.GetUserQuizSessions(userID)
	if err != nil {
		h.logger.Error("failed to get user sessions from db", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	data := map[string]interface{}{
		"sessions": sessions,
	}
	if err := json.NewEncoder(rw).Encode(data); err != nil {
		h.logger.Error("failed to encode response", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}
	//rw.WriteHeader(http.StatusOK)
}
