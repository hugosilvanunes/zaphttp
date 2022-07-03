package zaphttp

import "go.uber.org/zap"

func New() *Client {
	logger := newLogger()

	return newClient(logger)
}

func NewWithLogger(logger *zap.Logger) *Client {
	return newClient(logger)
}
