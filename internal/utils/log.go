package utils

import "go.uber.org/zap"

func InitializeLogger() (*zap.Logger, func(), error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		logger.Sync()
	}

	return logger, cleanup, err
}
