package clients

import (
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

func (c *AuthClient) VerifyAuthToken(token string) error {
	body := struct {
		AuthToken string `json:"token"`
	}{
		AuthToken: token,
	}
	jsonData, err := json.Marshal(body)
	data := bytes.NewBuffer(jsonData)
	resp, err := http.Post(c.addr+"/verify", "application/json", data)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return err
	}
	return nil
}
