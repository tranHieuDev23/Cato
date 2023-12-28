package cache

import (
	"context"
	"errors"

	"github.com/bluele/gcache"
	"go.uber.org/zap"

	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/utils"
)

type Client interface {
	Set(ctx context.Context, key string, value any) error
	Get(ctx context.Context, key string) (any, error)
}

type client struct {
	cache  gcache.Cache
	logger *zap.Logger
}

func NewClient(cacheConfig configs.Cache, logger *zap.Logger) Client {
	return &client{
		cache:  gcache.New(int(cacheConfig.Size)).LRU().Build(),
		logger: logger,
	}
}

func (c client) Get(ctx context.Context, key string) (any, error) {
	logger := utils.LoggerWithContext(ctx, c.logger)
	value, err := c.cache.GetIFPresent(key)
	if err != nil {
		if errors.Is(err, gcache.KeyNotFoundError) {
			return nil, nil
		}

		logger.With(zap.String("key", key)).With(zap.Error(err)).Error("failed to set cache entry")
		return nil, err
	}

	return value, nil
}

func (c client) Set(ctx context.Context, key string, value any) error {
	logger := utils.LoggerWithContext(ctx, c.logger)
	if err := c.cache.Set(key, value); err != nil {
		logger.With(zap.String("key", key)).With(zap.Error(err)).Error("failed to set cache entry")
		return err
	}

	return nil
}
