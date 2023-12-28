package cache

import (
	"context"

	"github.com/tranHieuDev23/cato/internal/utils"
	"go.uber.org/zap"
)

const (
	cacheKeySetting = "setting"
)

type Setting interface {
	Get(ctx context.Context) ([]byte, error)
	Set(ctx context.Context, data []byte) error
}
type setting struct {
	client Client
	logger *zap.Logger
}

func NewSetting(
	client Client,
	logger *zap.Logger,
) Setting {
	return &setting{
		client: client,
		logger: logger,
	}
}

func (c setting) Get(ctx context.Context) ([]byte, error) {
	logger := utils.LoggerWithContext(ctx, c.logger)

	cacheEntry, err := c.client.Get(ctx, cacheKeySetting)
	if err != nil {
		return nil, err
	}

	if cacheEntry == nil {
		return nil, nil
	}

	publicKey, ok := cacheEntry.([]byte)
	if !ok {
		logger.Error("cache entry is not of type bytes")
		return nil, nil
	}

	return publicKey, nil
}

func (c setting) Set(ctx context.Context, data []byte) error {
	return c.client.Set(ctx, cacheKeySetting, data)
}
