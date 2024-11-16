package clients

import (
	"admin/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

type AuthClient struct {
	addr   string
	apiKey string
	logger *zap.Logger
}

func NewAuthClient(addr string, apiKey string, logger *zap.Logger) *AuthClient {
	return &AuthClient{
		addr:   addr,
		apiKey: apiKey,
		logger: logger,
	}
}

func (c *AuthClient) VerifyAuthToken(token string) (models.UserAuthData, error) {
	body := struct {
		AuthToken string `json:"token"`
	}{
		AuthToken: token,
	}

	jsonPayload, err := json.Marshal(body)
	if err != nil {
		return models.UserAuthData{}, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", c.addr+"/verify", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return models.UserAuthData{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.UserAuthData{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("unexpected status code", zap.Error(err), zap.Int("status_code", resp.StatusCode))
		return models.UserAuthData{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var userDataResponse models.UserAuthData
	err = json.NewDecoder(resp.Body).Decode(&userDataResponse)
	if err != nil {
		return models.UserAuthData{}, fmt.Errorf("failed to decode response: %w", err)
	}

	c.logger.Info("response", zap.Any("response", userDataResponse))

	return userDataResponse, nil
}

func (c *AuthClient) GetUsers() ([]models.User, error) {
	req, err := http.NewRequest("GET", c.addr+"/users", nil)
	if err != nil {
		c.logger.Error("failed to create request", zap.Error(err))
		return nil, err
	}

	client := &http.Client{}
	req.Header.Set("X-Api-Key", c.apiKey)
	resp, err := client.Do(req)
	if err != nil {
		c.logger.Error("failed to send request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("unexpected status code", zap.Error(err), zap.Int("status_code", resp.StatusCode))
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var users []models.User
	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		c.logger.Error("failed to decode response", zap.Error(err))
		return nil, err
	}

	c.logger.Info("response", zap.Any("response", users))
	return users, nil
}
