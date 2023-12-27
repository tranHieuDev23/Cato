package db

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/tranHieuDev23/cato/internal/utils"
)

type TokenPublicKey struct {
	gorm.Model
	PublicKey []byte
}

type TokenPublicKeyDataAccessor interface {
	CreatePublicKey(ctx context.Context, tokenPublicKey *TokenPublicKey) error
	GetPublicKey(ctx context.Context, id uint64) (*TokenPublicKey, error)
	WithDB(db *gorm.DB) TokenPublicKeyDataAccessor
}

type tokenPublicKeyDataAccessor struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewTokenPublicKeyDataAccessor(
	db *gorm.DB,
	logger *zap.Logger,
) TokenPublicKeyDataAccessor {
	return &tokenPublicKeyDataAccessor{
		db:     db,
		logger: logger,
	}
}

func (a tokenPublicKeyDataAccessor) CreatePublicKey(ctx context.Context, tokenPublicKey *TokenPublicKey) error {
	logger := utils.LoggerWithContext(ctx, a.logger)
	db := a.db.WithContext(ctx)
	if err := db.Create(tokenPublicKey).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to create token public key")
		return err
	}

	return nil
}

func (a tokenPublicKeyDataAccessor) GetPublicKey(ctx context.Context, id uint64) (*TokenPublicKey, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.Uint64("id", id))
	db := a.db.WithContext(ctx)
	tokenPublicKey := new(TokenPublicKey)
	if err := db.First(tokenPublicKey, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		logger.With(zap.Error(err)).Error("failed to get token public key")
		return nil, err
	}

	return tokenPublicKey, nil
}

func (a tokenPublicKeyDataAccessor) WithDB(db *gorm.DB) TokenPublicKeyDataAccessor {
	return &tokenPublicKeyDataAccessor{
		db:     db,
		logger: a.logger,
	}
}
