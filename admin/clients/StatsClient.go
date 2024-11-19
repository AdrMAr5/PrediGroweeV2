package clients

import (
	"admin/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

type StatsClient interface {
	GetUserStats(userID string) (models.UserStats, error)
}

type StatsRestClient struct {
	addr   string
	apiKey string
	logger *zap.Logger
}

func NewStatsRestClient(addr string, apiKey string, logger *zap.Logger) *StatsRestClient {
	return &StatsRestClient{
		addr:   addr,
		apiKey: apiKey,
		logger: logger,
	}
}
func (c *StatsRestClient) NewRequestWithAuth(method, path string, body interface{}) (*http.Request, error) {
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

func (c *StatsRestClient) GetUserStats(userID string) (models.UserStats, error) {
	req, err := c.NewRequestWithAuth("GET", fmt.Sprintf("/users/%s", userID), nil)
	if err != nil {
		return models.UserStats{}, fmt.Errorf("failed to create request: %w", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.UserStats{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return models.UserStats{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var userStats models.UserStats
	err = json.NewDecoder(resp.Body).Decode(&userStats)
	if err != nil {
		return models.UserStats{}, fmt.Errorf("failed to decode response body: %w", err)
	}
	return userStats, nil
}
