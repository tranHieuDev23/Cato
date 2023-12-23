package db

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/tranHieuDev23/cato/internal/utils"
)

type AccountPassword struct {
	gorm.Model
	OfAccountID uint64
	Account     Account `gorm:"foreignKey:OfAccountID"`
	Hash        string
}

type AccountPasswordDataAccessor interface {
	CreateAccountPassword(ctx context.Context, accountPassword *AccountPassword) error
	UpdateAccountPassword(ctx context.Context, accountPassword *AccountPassword) error
	GetAccountPasswordOfAccountID(ctx context.Context, accountID uint64) (*AccountPassword, error)
	WithDB(db *gorm.DB) AccountPasswordDataAccessor
}

type accountPasswordDataAccessor struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewAccountPasswordDataAccessor(
	db *gorm.DB,
	logger *zap.Logger,
) AccountPasswordDataAccessor {
	return &accountPasswordDataAccessor{
		db:     db,
		logger: logger,
	}
}

func (a accountPasswordDataAccessor) CreateAccountPassword(ctx context.Context, accountPassword *AccountPassword) error {
	logger := utils.LoggerWithContext(ctx, a.logger)
	db := a.db.WithContext(ctx)
	if err := db.Create(accountPassword).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to create account password")
		return err
	}

	return nil
}

func (a accountPasswordDataAccessor) GetAccountPasswordOfAccountID(ctx context.Context, ofAccountID uint64) (*AccountPassword, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).
		With(zap.Uint64("of_account_id", ofAccountID))
	db := a.db.WithContext(ctx)
	accountPassword := new(AccountPassword)
	if err := db.Model(new(AccountPassword)).
		Where(&AccountPassword{
			OfAccountID: ofAccountID,
		}).
		First(accountPassword).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		logger.With(zap.Error(err)).Error("failed to get account password")
		return nil, err
	}

	return accountPassword, nil
}

func (a accountPasswordDataAccessor) UpdateAccountPassword(ctx context.Context, accountPassword *AccountPassword) error {
	logger := utils.LoggerWithContext(ctx, a.logger)
	db := a.db.WithContext(ctx)
	if err := db.Save(accountPassword).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to update account password")
		return err
	}

	return nil
}

func (a accountPasswordDataAccessor) WithDB(db *gorm.DB) AccountPasswordDataAccessor {
	return &accountPasswordDataAccessor{
		db:     db,
		logger: a.logger,
	}
}
