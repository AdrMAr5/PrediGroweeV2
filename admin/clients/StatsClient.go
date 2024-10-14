package clients

import (
	"go.uber.org/zap"
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
