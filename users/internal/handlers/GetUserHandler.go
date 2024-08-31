package handlers

import (
	"PrediGroweeV2/users/internal/storage"
	"go.uber.org/zap"
	"net/http"
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

}
