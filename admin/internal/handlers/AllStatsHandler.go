package handlers

import (
	"admin/clients"
	"admin/internal/storage"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type AllStatsHandler struct {
	storage     storage.Storage
	logger      *zap.Logger
	statsClient clients.StatsClient
}

func NewAllStatsHandler(storage storage.Storage, logger *zap.Logger, statsClient clients.StatsClient) *AllStatsHandler {
	return &AllStatsHandler{
		storage:     storage,
		logger:      logger,
		statsClient: statsClient,
	}
}
func (h *AllStatsHandler) GetAllResponses(w http.ResponseWriter, _ *http.Request) {
	stats, err := h.statsClient.GetAllResponses()
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

func (h *AllStatsHandler) GetStatsForQuestion(w http.ResponseWriter, r *http.Request) {
	questionId := r.PathValue("questionId")
	stats, err := h.statsClient.GetStatsForQuestion(questionId)
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

func (h *AllStatsHandler) GetStatsForAllQuestions(w http.ResponseWriter, r *http.Request) {
	stats, err := h.statsClient.GetStatsForAllQuestions()
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
