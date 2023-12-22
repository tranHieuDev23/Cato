package utils

import (
	"context"

	"gitlab.com/pjrpc/pjrpc/v2"
	"go.uber.org/zap"
)

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

func LoggerWithContext(ctx context.Context, logger *zap.Logger) *zap.Logger {
	contextData, ok := pjrpc.ContextGetData(ctx)
	if !ok {
		return logger
	}

	requestID := contextData.JRPCRequest.GetID()
	return logger.With(zap.String("request_id", requestID))
}
