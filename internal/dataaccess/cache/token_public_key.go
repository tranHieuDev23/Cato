package cache

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/tranHieuDev23/cato/internal/utils"
)

type TokenPublicKey interface {
	Get(ctx context.Context, id uint64) ([]byte, error)
	Set(ctx context.Context, id uint64, bytes []byte) error
}

type tokenPublicKey struct {
	client Client
	logger *zap.Logger
}

func NewTokenPublicKey(
	client Client,
	logger *zap.Logger,
) TokenPublicKey {
	return &tokenPublicKey{
		client: client,
		logger: logger,
	}
}

func (c tokenPublicKey) getTokenPublicKeyCacheKey(id uint64) string {
	return fmt.Sprintf("token_public_key:%d", id)
}

func (c tokenPublicKey) Get(ctx context.Context, id uint64) ([]byte, error) {
	logger := utils.LoggerWithContext(ctx, c.logger).With(zap.Uint64("id", id))

	cacheKey := c.getTokenPublicKeyCacheKey(id)
	cacheEntry, err := c.client.Get(ctx, cacheKey)
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

func (c tokenPublicKey) Set(ctx context.Context, id uint64, bytes []byte) error {
	cacheKey := c.getTokenPublicKeyCacheKey(id)
	return c.client.Set(ctx, cacheKey, bytes)
}
