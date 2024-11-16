package handlers

import (
	"admin/clients"
	"admin/internal/storage"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

type UsersHandler struct {
	storage    storage.Storage
	logger     *zap.Logger
	authClient *clients.AuthClient
}

func NewUsersHandler(storage storage.Storage, logger *zap.Logger, authClient *clients.AuthClient) *UsersHandler {
	return &UsersHandler{
		storage:    storage,
		logger:     logger,
		authClient: authClient,
	}
}
func (u *UsersHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := u.authClient.GetUsers()
	if err != nil {
		u.logger.Error("failed to get users", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Println(users)
	w.Header().Set("Content-Type", "application/json")
	usersJson, _ := json.Marshal(users)
	_, err = w.Write(usersJson)
	if err != nil {
		u.logger.Error("failed to write response", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
