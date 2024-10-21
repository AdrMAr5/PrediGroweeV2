package handlers

import (
	"auth/internal/models"
	"auth/internal/storage"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// UpdateUserHandler updates user details
type UpdateUserHandler struct {
	storage storage.Store
	logger  *zap.Logger
}

func NewUpdateUserHandler(store storage.Store, logger *zap.Logger) *UpdateUserHandler {
	return &UpdateUserHandler{
		storage: store,
		logger:  logger,
	}
}

func (h *UpdateUserHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(rw, "invalid user id", http.StatusBadRequest)
		return
	}

	var updatedUser models.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(rw, "invalid request payload", http.StatusBadRequest)
		return
	}

	updatedUser.ID = userID
	if err := h.storage.UpdateUser(updatedUser); err != nil {
		h.logger.Error("failed to update user", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
