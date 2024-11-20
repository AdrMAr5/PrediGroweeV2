package handlers

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"stats/internal/storage"
	"strconv"
)

type GetAllStatsHandler struct {
	storage storage.Storage
	logger  *zap.Logger
}

func NewGetAllStatsHandler(storage storage.Storage, logger *zap.Logger) *GetAllStatsHandler {
	return &GetAllStatsHandler{
		storage: storage,
		logger:  logger,
	}
}

func (h *GetAllStatsHandler) GetResponses(w http.ResponseWriter, _ *http.Request) {
	stats, err := h.storage.GetAllResponses()
	if err != nil {
		h.logger.Error("failed to get stats", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	statsJson, _ := json.Marshal(stats)
	_, err = w.Write(statsJson)
	if err != nil {
		h.logger.Error("failed to write response", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *GetAllStatsHandler) GetStatsForQuestion(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetStatsForQuestion")
	questionId := r.PathValue("id")
	if questionId == "-" {
		stats, err := h.storage.GetStatsForAllQuestions()
		if err != nil {
			h.logger.Error("failed to get stats", zap.Error(err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		statsJson, _ := json.Marshal(stats)
		_, err = w.Write(statsJson)
		return
	}
	questionID, err := strconv.Atoi(questionId)
	if err != nil {
		h.logger.Error("failed to parse question id", zap.Error(err))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	stats, err := h.storage.GetStatsForQuestion(questionID)
	if err != nil {
		h.logger.Error("failed to get stats", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	statsJson, _ := json.Marshal(stats)
	_, err = w.Write(statsJson)
	if err != nil {
		h.logger.Error("failed to write response", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
