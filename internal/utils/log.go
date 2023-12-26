package utils

import (
	"context"

	"gitlab.com/pjrpc/pjrpc/v2"
	"go.uber.org/zap"

	"github.com/tranHieuDev23/cato/internal/configs"
)

func getZapLoggerLevel(level string) zap.AtomicLevel {
	switch level {
	case "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "panic":
		return zap.NewAtomicLevelAt(zap.PanicLevel)
	default:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
}

func InitializeLogger(logConfig configs.Log) (*zap.Logger, func(), error) {
	zapLoggerConfig := zap.NewProductionConfig()
	zapLoggerConfig.Level = getZapLoggerLevel(logConfig.Level)

	logger, err := zapLoggerConfig.Build()
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		// deliberately ignore the returned error here
		_ = logger.Sync()
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
