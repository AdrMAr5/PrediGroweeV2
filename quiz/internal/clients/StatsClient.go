package clients

import (
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"quiz/internal/models"
)

type StatsClient struct {
	addr   string
	logger *zap.Logger
}

func NewStatsClient(addr string, logger *zap.Logger) *StatsClient {
	return &StatsClient{
		addr:   addr,
		logger: logger,
	}
}

func (c *StatsClient) SaveResponse(answer models.QuestionAnswer) error {
	jsonPayload, err := json.Marshal(answer)
	if err != nil {
		c.logger.Error("failed to marshal request body", zap.Error(err))
		return err
	}
	req, err := http.NewRequest("POST", c.addr+"/saveResponse", bytes.NewBuffer(jsonPayload))
	if err != nil {
		c.logger.Error("failed to create request", zap.Error(err))
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.logger.Error("failed to send request", zap.Error(err))
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		c.logger.Error("unexpected status code", zap.Int("status_code", resp.StatusCode))
		return err
	}
	return nil
}
