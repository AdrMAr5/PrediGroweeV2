package handlers

import (
	"PrediGroweeV2/users/internal/storage"
	"go.uber.org/zap"
	"net/http"
)

type VerifyTokenHandler struct {
	store  storage.Store
	logger *zap.Logger
}

func NewVerifyTokenHandler(store storage.Store, logger *zap.Logger) *VerifyTokenHandler {
	return &VerifyTokenHandler{
		logger: logger,
		store:  store,
	}
}

func (h *VerifyTokenHandler) Handle(rw http.ResponseWriter, r *http.Request) {

}
