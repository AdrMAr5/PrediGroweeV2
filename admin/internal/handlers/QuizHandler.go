package handlers

import (
	"admin/clients"
	"admin/internal/models"
	"admin/internal/storage"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type QuizHandler struct {
	storage     storage.Storage
	logger      *zap.Logger
	quizClient  clients.QuizClient
	statsClient clients.StatsClient
}

func NewQuizHandler(storage storage.Storage, logger *zap.Logger, quizClient clients.QuizClient, statsClient clients.StatsClient) *QuizHandler {
	return &QuizHandler{
		storage:     storage,
		logger:      logger,
		quizClient:  quizClient,
		statsClient: statsClient,
	}
}
func (h *QuizHandler) GetAllQuestions(w http.ResponseWriter, _ *http.Request) {
	questions, err := h.quizClient.GetAllQuestions()
	if err != nil {
		h.logger.Error("Failed to get questions", zap.Error(err))
		http.Error(w, "Failed to get questions", http.StatusInternalServerError)
		return
	}
	questionsJSON, err := json.Marshal(questions)
	if err != nil {
		h.logger.Error("Failed to marshal questions", zap.Error(err))
		http.Error(w, "Failed to marshal questions", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(questionsJSON)
	if err != nil {
		h.logger.Error("Failed to write response", zap.Error(err))
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (h *QuizHandler) GetAllParameters(w http.ResponseWriter, _ *http.Request) {
	parameters, err := h.quizClient.GetAllParameters()
	if err != nil {
		h.logger.Error("Failed to get parameters", zap.Error(err))
		http.Error(w, "Failed to get parameters", http.StatusInternalServerError)
		return
	}
	parametersJSON, err := json.Marshal(parameters)
	if err != nil {
		h.logger.Error("Failed to marshal parameters", zap.Error(err))
		http.Error(w, "Failed to marshal parameters", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(parametersJSON)
	if err != nil {
		h.logger.Error("Failed to write response", zap.Error(err))
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (h *QuizHandler) UpdateParameter(w http.ResponseWriter, r *http.Request) {
	paramId := r.PathValue("id")
	updatedParameter := models.Parameter{}
	err := json.NewDecoder(r.Body).Decode(&updatedParameter)
	if err != nil {
		h.logger.Error("Failed to decode request", zap.Error(err))
		http.Error(w, "Failed to decode request", http.StatusBadRequest)
		return
	}
	err = h.quizClient.UpdateParameter(paramId, updatedParameter)
}

func (h *QuizHandler) GetAllOptions(w http.ResponseWriter, _ *http.Request) {
	options, err := h.quizClient.GetAllOptions()
	if err != nil {
		h.logger.Error("Failed to get options", zap.Error(err))
		http.Error(w, "Failed to get options", http.StatusInternalServerError)
		return
	}
	optionsJSON, err := json.Marshal(options)
	if err != nil {
		h.logger.Error("Failed to marshal options", zap.Error(err))
		http.Error(w, "Failed to marshal options", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(optionsJSON)
}

func (h *QuizHandler) GetQuestion(w http.ResponseWriter, r *http.Request) {
	questionId := r.PathValue("id")
	question, err := h.quizClient.GetQuestion(questionId)
	if err != nil {
		h.logger.Error("Failed to get question", zap.Error(err))
		http.Error(w, "Failed to get question", http.StatusInternalServerError)
		return
	}
	questionJSON, err := json.Marshal(question)
	if err != nil {
		h.logger.Error("Failed to marshal question", zap.Error(err))
		http.Error(w, "Failed to marshal question", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(questionJSON)
	if err != nil {
		h.logger.Error("Failed to write response", zap.Error(err))
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (h *QuizHandler) UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	questionId := r.PathValue("id")
	updatedQuestion := models.Question{}
	err := json.NewDecoder(r.Body).Decode(&updatedQuestion)
	if err != nil {
		h.logger.Error("Failed to decode request", zap.Error(err))
		http.Error(w, "Failed to decode request", http.StatusBadRequest)
		return
	}
	err = h.quizClient.UpdateQuestion(questionId, updatedQuestion)
	if err != nil {
		h.logger.Error("Failed to update question", zap.Error(err))
		http.Error(w, "Failed to update question", http.StatusInternalServerError)
		return
	}
}

func (h *QuizHandler) CreateParameter(w http.ResponseWriter, r *http.Request) {
	newParameter := models.Parameter{}
	err := json.NewDecoder(r.Body).Decode(&newParameter)
	if err != nil {
		h.logger.Error("Failed to decode request", zap.Error(err))
		http.Error(w, "Failed to decode request", http.StatusBadRequest)
		return
	}
	param, err := h.quizClient.CreateParameter(newParameter)
	if err != nil {
		h.logger.Error("Failed to create parameter", zap.Error(err))
		http.Error(w, "Failed to create parameter", http.StatusInternalServerError)
		return
	}
	paramJSON, err := json.Marshal(param)
	if err != nil {
		h.logger.Error("Failed to marshal parameter", zap.Error(err))
		http.Error(w, "Failed to marshal parameter", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(paramJSON)
}

func (h *QuizHandler) UpdateOption(w http.ResponseWriter, r *http.Request) {
	optionId := r.PathValue("id")
	updatedOption := models.Option{}
	err := json.NewDecoder(r.Body).Decode(&updatedOption)
	if err != nil {
		h.logger.Error("Failed to decode request", zap.Error(err))
		http.Error(w, "Failed to decode request", http.StatusBadRequest)
		return
	}
	err = h.quizClient.UpdateOption(optionId, updatedOption)
	if err != nil {
		h.logger.Error("Failed to update option", zap.Error(err))
		http.Error(w, "Failed to update option", http.StatusInternalServerError)
		return
	}
}

func (h *QuizHandler) CreateOption(w http.ResponseWriter, r *http.Request) {
	newOption := models.Option{}
	err := json.NewDecoder(r.Body).Decode(&newOption)
	if err != nil {
		h.logger.Error("Failed to decode request", zap.Error(err))
		http.Error(w, "Failed to decode request", http.StatusBadRequest)
		return
	}
	option, err := h.quizClient.CreateOption(newOption)
	if err != nil {
		h.logger.Error("Failed to create option", zap.Error(err))
		http.Error(w, "Failed to create option", http.StatusInternalServerError)
		return
	}
	optionJSON, err := json.Marshal(option)
	if err != nil {
		h.logger.Error("Failed to marshal option", zap.Error(err))
		http.Error(w, "Failed to marshal option", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(optionJSON)
}

func (h *QuizHandler) DeleteOption(w http.ResponseWriter, r *http.Request) {
	optionId := r.PathValue("id")
	err := h.quizClient.DeleteOption(optionId)
	if err != nil {
		h.logger.Error("Failed to delete option", zap.Error(err))
		http.Error(w, "Failed to delete option", http.StatusInternalServerError)
		return
	}
}
