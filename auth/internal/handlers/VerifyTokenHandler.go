package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type VerifyTokenHandler struct {
}

func NewVerifyTokenHandler() *VerifyTokenHandler {
	return &VerifyTokenHandler{}
}
func (h *VerifyTokenHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	log.Println("VerifyTokenHandler called")

	userID := r.Context().Value("user_id")
	userRole := r.Context().Value("user_role")

	log.Printf("User ID from context: %v", userID)
	log.Printf("User Role from context: %v", userRole)

	resp := map[string]interface{}{
		"user_id": userID,
		"role":    userRole,
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	err := json.NewEncoder(rw).Encode(resp)
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(rw, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("VerifyTokenHandler completed successfully")
}
