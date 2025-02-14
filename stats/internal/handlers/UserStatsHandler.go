package handlers

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"stats/internal/models"
	"stats/internal/storage"
	"strconv"
)

type UserStatsHandler struct {
	storage storage.Storage
	logger  *zap.Logger
}

func NewUserStatsHandler(storage storage.Storage, logger *zap.Logger) *UserStatsHandler {
	return &UserStatsHandler{storage: storage, logger: logger}
}

func (h *UserStatsHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	var userID int
	userId := r.PathValue("id")
	if userId == "" {
		userID = r.Context().Value("user_id").(int)
	} else {
		var err error
		userID, err = strconv.Atoi(userId)
		if err != nil {
			http.Error(rw, "invalid user id", http.StatusBadRequest)
		}
	}

	stats := models.UserStats{
		TotalQuestions: make(map[models.QuizMode]int),
		CorrectAnswers: make(map[models.QuizMode]int),
		Accuracy:       make(map[models.QuizMode]float64),
	}
	for _, mode := range []string{models.QuizModeEducational, models.QuizModeClassic, models.QuizModeLimitedTime} {
		correct, wrong, err := h.storage.GetUserStatsForMode(userID, mode)
		if err == storage.ErrStatsNotFound {
			continue
		}
		if err != nil {
			h.logger.Error(fmt.Sprintf("failed to get stats for userID: %s, quizMode: %s", userID, mode))
			http.Error(rw, "failed to get statistics", http.StatusInternalServerError)
		}
		stats.TotalQuestions[mode] = correct + wrong
		stats.CorrectAnswers[mode] = correct
		if stats.TotalQuestions[mode] != 0 {
			stats.Accuracy[mode] = float64(correct) / float64(correct+wrong)
		} else {
			stats.Accuracy[mode] = 0
		}

	}
	fmt.Println(stats)
	err := stats.ToJSON(rw)
	if err != nil {
		h.logger.Error("failed to encode stats")
		http.Error(rw, "internal server error", http.StatusInternalServerError)
	}
}

func (h *UserStatsHandler) GetUserSessions(rw http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	stats, err := h.storage.GetUserQuizSessionsStats(userID)
	if err != nil {
		h.logger.Error("failed to get user sessions", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(stats)
	if err != nil {
		h.logger.Error("failed to encode response", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
	}
}

func (h *UserStatsHandler) DeleteUserResponses(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("deleting user responses")
	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	err = h.storage.DeleteUserResponses(userID)
	if err != nil {
		h.logger.Error("failed to delete user responses", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserStatsHandler) GetAllUsersStats(w http.ResponseWriter, _ *http.Request) {
	stats, err := h.storage.GetAllUsersStats()
	if err != nil {
		h.logger.Error("failed to get user stats", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}
