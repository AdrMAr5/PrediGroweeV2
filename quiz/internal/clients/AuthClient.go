package clients

import (
	"PrediGroweeV2/quiz/internal/models"
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type AuthClient struct {
	addr   string
	logger *zap.Logger
}

func NewAuthClient(addr string, logger *zap.Logger) *AuthClient {
	return &AuthClient{
		addr:   addr,
		logger: logger,
	}
}

func (c *AuthClient) VerifyAuthToken(token string) (models.UserData, error) {
	body := struct {
		AuthToken string `json:"token"`
	}{
		AuthToken: token,
	}
	jsonPayload, err := json.Marshal(body)
	payload := bytes.NewBuffer(jsonPayload)
	resp, err := http.Post(c.addr+"/verify", "application/json", payload)
	if err != nil {
		return models.UserData{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return models.UserData{}, err
	}
	var userDataResponse models.UserData
	err = userDataResponse.FromJSON(resp.Body)
	if err != nil {
		return models.UserData{}, err
	}
	return userDataResponse, nil

}
