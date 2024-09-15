package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type VerifyTokenHandler struct {
}

func NewVerifyTokenHandler() *VerifyTokenHandler {
	return &VerifyTokenHandler{}
}
func (h *VerifyTokenHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	userRole := r.Context().Value("user_role")
	resp, err := json.Marshal(map[string]interface{}{
		"user_id": userID,
		"role":    userRole,
	})
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = rw.Write(bytes.NewBuffer(resp).Bytes())
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}
