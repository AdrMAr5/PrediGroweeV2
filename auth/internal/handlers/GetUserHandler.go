package handlers

import (
	"PrediGroweeV2/auth/internal/storage"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type GetUserHandler struct {
	store  storage.Store
	logger *zap.Logger
}

func NewGetUserHandler(store storage.Store, logger *zap.Logger) *GetUserHandler {
	return &GetUserHandler{
		logger: logger,
		store:  store,
	}
}
func (h *GetUserHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	sessionUserID := r.Context().Value("user_id").(int)
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(rw, "Invalid user id", http.StatusBadRequest)
		return
	}
	if sessionUserID != id {
		http.Error(rw, "Permission denied", http.StatusForbidden)
		return
	}
	user, err := h.store.GetUserById(id)
	if err != nil {
		http.Error(rw, "User not found", http.StatusNotFound)
		return
	}
	rw.WriteHeader(http.StatusOK)
	err = user.ToJSON(rw)
	if err != nil {
		h.logger.Error("Error marshalling user", zap.Error(err))
	}
}
