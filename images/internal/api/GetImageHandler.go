package api

import (
	"go.uber.org/zap"

	"net/http"
)

type GetImageHandler struct {
	logger *zap.Logger
}

func NewGetImageHandler(logger *zap.Logger) *GetImageHandler {
	return &GetImageHandler{
		logger: logger,
	}
}
func (h *GetImageHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	h.logger.Info("Getting image with id: " + id)
	http.ServeFile(rw, r, "/app/images/"+"xray1"+".jpg")
}
