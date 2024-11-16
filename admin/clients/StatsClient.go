package clients

import (
	"go.uber.org/zap"
)

type StatsClient struct {
	addr   string
	apiKey string
	logger *zap.Logger
}

func NewStatsClient(addr string, apiKey string, logger *zap.Logger) *StatsClient {
	return &StatsClient{
		addr:   addr,
		apiKey: apiKey,
		logger: logger,
	}
}
