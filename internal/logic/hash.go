package logic

import (
	"context"
	"errors"

	"gitlab.com/pjrpc/pjrpc/v2"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/utils"
)

type Hash interface {
	Hash(ctx context.Context, data string) (string, error)
	IsHashEqual(ctx context.Context, data string, hashed string) (bool, error)
}

type hash struct {
	hashConfig configs.Hash
	logger     *zap.Logger
}

func NewHash(
	hashConfig configs.Hash,
	logger *zap.Logger,
) Hash {
	return &hash{
		hashConfig: hashConfig,
		logger:     logger,
	}
}

func (h hash) Hash(ctx context.Context, data string) (string, error) {
	logger := utils.LoggerWithContext(ctx, h.logger)

	hashed, err := bcrypt.GenerateFromPassword([]byte(data), h.hashConfig.Cost)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to generate hash")
		return "", pjrpc.JRPCErrInternalError()
	}

	return string(hashed), nil
}

func (h hash) IsHashEqual(ctx context.Context, data string, hashed string) (bool, error) {
	logger := utils.LoggerWithContext(ctx, h.logger)

	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(data)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}

		logger.With(zap.Error(err)).Error("failed to compare data with hash")
		return false, pjrpc.JRPCErrInternalError()
	}

	return true, nil
}
