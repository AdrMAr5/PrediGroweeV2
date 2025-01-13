package handlers

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"quiz/internal/clients"
	"quiz/internal/models"
	"quiz/internal/storage"
	"strconv"
	"strings"
	"time"
)

type SubmitAnswerHandler struct {
	storage     storage.Store
	logger      *zap.Logger
	statsClient *clients.StatsClient
}

func NewSubmitAnswerHandler(store storage.Store, logger *zap.Logger, statsClient *clients.StatsClient) *SubmitAnswerHandler {
	return &SubmitAnswerHandler{
		storage:     store,
		logger:      logger,
		statsClient: statsClient,
	}
}
func (h *SubmitAnswerHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	quizSessionIdString := r.PathValue("quizSessionId")
	if quizSessionIdString == "" {
		h.logger.Info("no quiz session id provided")
		http.Error(rw, "invalid quiz session id", http.StatusBadRequest)
		return
	}
	quizSessionID, err := strconv.Atoi(quizSessionIdString)
	if err != nil {
		h.logger.Info("invalid quiz session id")
		http.Error(rw, "invalid quiz session id", http.StatusBadRequest)
		return
	}

	session, err := h.storage.GetQuizSessionByID(quizSessionID)
	if err != nil {
		h.logger.Error("failed to get quiz session from db", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}
	timeSpend := time.Now().Sub(session.QuestionRequestedTime)

	userID := r.Context().Value("user_id").(int)
	if userID != session.UserID {
		http.Error(rw, "internal server error", http.StatusInternalServerError)
	}
	data := map[string]interface{}{}

	correct, err := h.storage.GetQuestionCorrectOption(session.CurrentQuestionID)
	if err != nil {
		h.logger.Error("failed to get question correct option", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}
	if session.Mode == models.QuizModeEducational {
		h.logger.Info("educational mode")
		data["correct"] = correct
		h.logger.Info("educational mode, returning correct answer")
	}
	h.logger.Info("submitting answer")

	session.Status = models.QuizStatusInProgress
	err = h.storage.UpdateQuizSession(session)
	if err != nil {
		h.logger.Error("failed to update quiz session", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}
	var answer models.QuestionAnswer
	err = json.NewDecoder(r.Body).Decode(&answer)
	if err != nil {
		h.logger.Error("failed to decode answer", zap.Error(err))
		http.Error(rw, "invalid answer", http.StatusBadRequest)
		return
	}
	question, err := h.storage.GetQuestionByID(session.CurrentQuestionID)
	if err != nil {
		h.logger.Error("failed to get question by id", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Println("question", session.CurrentQuestionID, "answer", answer.Answer, "correct", correct)
	err = h.statsClient.SaveResponse(session.ID, models.QuestionAnswer{
		QuestionID: session.CurrentQuestionID,
		Answer:     answer.Answer,
		IsCorrect:  strings.EqualFold(strings.TrimSpace(answer.Answer), strings.TrimSpace(correct)),
		ScreenSize: answer.ScreenSize,
		TimeSpent:  int(timeSpend.Seconds()),
		CaseCode:   question.Case.Code,
	})
	if err != nil {
		h.logger.Error("failed to save response", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}
	err = h.SetNextQuestionID(&session)
	if err != nil {
		h.logger.Error("failed to set next question id", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(data); err != nil {
		h.logger.Error("failed to encode response", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *SubmitAnswerHandler) SetNextQuestionID(qs *models.QuizSession) error {
	for i, q := range qs.GroupOrder {
		if q == qs.CurrentQuestionID {
			if i+1 < len(qs.GroupOrder) {
				qs.CurrentQuestionID = qs.GroupOrder[i+1]
				return nil
			} else {
				nextGroup, err := h.storage.GetNextQuestionGroupID(qs.CurrentGroup)
				if err != nil {
					return err
				}
				qs.CurrentGroup = nextGroup
				qs.GroupOrder, err = h.storage.GetGroupQuestionsIDsRandomOrder(qs.CurrentGroup)
				if err != nil {
					return err
				}
				qs.CurrentQuestionID = qs.GroupOrder[0]
				return nil
			}
		}
	}
	return fmt.Errorf("question not found in group order")
}
