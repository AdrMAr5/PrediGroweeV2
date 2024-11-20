package handlers

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"stats/internal/models"
	"stats/internal/storage"
	"strconv"
)

type GetUserStatsHandler struct {
	storage storage.Storage
	logger  *zap.Logger
}

func NewGetUserStatsHandler(storage storage.Storage, logger *zap.Logger) *GetUserStatsHandler {
	return &GetUserStatsHandler{storage: storage, logger: logger}
}

func (h *GetUserStatsHandler) Handle(rw http.ResponseWriter, r *http.Request) {
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
