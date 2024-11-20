package handlers

import (
	"admin/clients"
	"admin/internal/models"
	"admin/internal/storage"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type UsersHandler struct {
	storage     storage.Storage
	logger      *zap.Logger
	authClient  clients.AuthClient
	statsClient clients.StatsClient
}

func NewUsersHandler(storage storage.Storage, logger *zap.Logger, authClient clients.AuthClient, statsClient clients.StatsClient) *UsersHandler {
	return &UsersHandler{
		storage:     storage,
		logger:      logger,
		authClient:  authClient,
		statsClient: statsClient,
	}
}
func (u *UsersHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := u.authClient.GetUsers()
	if err != nil {
		u.logger.Error("failed to get users", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	usersJson, _ := json.Marshal(users)
	_, err = w.Write(usersJson)
	if err != nil {
		u.logger.Error("failed to write response", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
func (u *UsersHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user models.UserPayload
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		u.logger.Error("failed to decode request body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if user.ID == "" {
		user.ID = r.PathValue("id")
	}
	err = u.authClient.UpdateUser(user)
	if err != nil {
		u.logger.Error("failed to update user", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (u *UsersHandler) GetUserDetails(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	user, err := u.authClient.GetUser(userID)
	userStats, err := u.statsClient.GetUserStats(userID)
	if err != nil {
		u.logger.Error("failed to get user or user stats", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	var userDetails models.UserDetails
	userDetails.User = user
	userDetails.Stats = userStats
	w.Header().Set("Content-Type", "application/json")
	userDetailsJson, _ := json.Marshal(userDetails)
	_, err = w.Write(userDetailsJson)
	if err != nil {
		u.logger.Error("failed to write response", zap.Error(err))
		return
	}
}

func (u *UsersHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	err := u.authClient.DeleteUser(userID)
	if err != nil {
		u.logger.Error("failed to delete user", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
