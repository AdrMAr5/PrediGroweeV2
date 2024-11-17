package clients

import (
	"admin/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

type QuizClient interface {
	GetAllQuestions() ([]models.Question, error)
	GetAllParameters() ([]models.Parameter, error)
	UpdateParameter(id string, parameter models.Parameter) error
	GetAllOptions() ([]models.Option, error)
}

type QuizRestClient struct {
	addr   string
	apiKey string
	logger *zap.Logger
}

func NewQuizRestClient(addr string, apiKey string, logger *zap.Logger) *QuizRestClient {
	return &QuizRestClient{
		addr:   addr,
		apiKey: apiKey,
		logger: logger,
	}
}
func (c *QuizRestClient) NewRequestWithAuth(method, path string, body interface{}) (*http.Request, error) {
	jsonPayload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(method, c.addr+path, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", c.apiKey)

	return req, nil
}

func (c *QuizRestClient) GetAllQuestions() ([]models.Question, error) {
	req, err := c.NewRequestWithAuth("GET", "/questions", nil)
	if err != nil {
		return []models.Question{}, fmt.Errorf("failed to create request: %w", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []models.Question{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return []models.Question{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var questions []models.Question
	err = json.NewDecoder(resp.Body).Decode(&questions)
	return questions, err
}

func (c *QuizRestClient) GetAllParameters() ([]models.Parameter, error) {
	req, err := c.NewRequestWithAuth("GET", "/parameters", nil)
	if err != nil {
		return []models.Parameter{}, fmt.Errorf("failed to create request: %w", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []models.Parameter{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return []models.Parameter{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var parameters []models.Parameter
	err = json.NewDecoder(resp.Body).Decode(&parameters)
	return parameters, err
}

func (c *QuizRestClient) UpdateParameter(id string, parameter models.Parameter) error {
	req, err := c.NewRequestWithAuth("PATCH", fmt.Sprintf("/parameters/%s", id), parameter)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func (c *QuizRestClient) GetAllOptions() ([]models.Option, error) {
	req, err := c.NewRequestWithAuth("GET", "/options", nil)
	if err != nil {
		return []models.Option{}, fmt.Errorf("failed to create request: %w", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []models.Option{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return []models.Option{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var options []models.Option
	err = json.NewDecoder(resp.Body).Decode(&options)
	return options, err
}
