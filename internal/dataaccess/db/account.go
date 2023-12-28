package db

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/tranHieuDev23/cato/internal/utils"
)

type AccountRole string

const (
	AccountRoleAdmin         AccountRole = "admin"
	AccountRoleProblemSetter AccountRole = "problem_setter"
	AccountRoleContestant    AccountRole = "contestant"
	AccountRoleWorker        AccountRole = "worker"
)

type Account struct {
	gorm.Model
	AccountName string      `gorm:"type:varchar(32);uniqueIndex"`
	DisplayName string      `gorm:"type:varchar(32)"`
	Role        AccountRole `gorm:"type:varchar(32)"`
}

type AccountDataAccessor interface {
	CreateAccount(ctx context.Context, account *Account) error
	UpdateAccount(ctx context.Context, account *Account) error
	GetAccount(ctx context.Context, id uint64) (*Account, error)
	GetAccountByAccountName(ctx context.Context, accountName string) (*Account, error)
	GetAccountList(ctx context.Context, offset, limit uint64) ([]*Account, error)
	GetAccountCount(ctx context.Context) (uint64, error)
	WithDB(db *gorm.DB) AccountDataAccessor
}

type accountDataAccessor struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewAccountDataAccessor(
	db *gorm.DB,
	logger *zap.Logger,
) AccountDataAccessor {
	return &accountDataAccessor{
		db:     db,
		logger: logger,
	}
}

func (a accountDataAccessor) CreateAccount(ctx context.Context, account *Account) error {
	logger := utils.LoggerWithContext(ctx, a.logger)
	db := a.db.WithContext(ctx)
	if err := db.Create(account).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to create account")
		return err
	}

	return nil
}

func (a accountDataAccessor) GetAccount(ctx context.Context, id uint64) (*Account, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.Uint64("id", id))
	db := a.db.WithContext(ctx)
	account := new(Account)
	if err := db.First(account, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		logger.With(zap.Error(err)).Error("failed to get account")
		return nil, err
	}

	return account, nil
}

func (a accountDataAccessor) GetAccountByAccountName(ctx context.Context, accountName string) (*Account, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).
		With(zap.String("account_name", accountName))
	db := a.db.WithContext(ctx)
	account := new(Account)
	if err := db.Model(new(Account)).Where(&Account{
		AccountName: accountName,
	}).First(account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		logger.With(zap.Error(err)).Error("failed to get account by account name")
		return nil, err
	}

	return account, nil
}

func (a accountDataAccessor) GetAccountCount(ctx context.Context) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, a.logger)
	db := a.db.WithContext(ctx)
	count := int64(0)
	if err := db.Model(new(Account)).Count(&count).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to get account count")
		return 0, err
	}

	return uint64(count), nil
}

func (a accountDataAccessor) GetAccountList(ctx context.Context, offset uint64, limit uint64) ([]*Account, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).
		With(zap.Uint64("limit", limit)).
		With(zap.Uint64("offset", offset))
	db := a.db.WithContext(ctx)
	accountList := make([]*Account, 0)
	if err := db.Model(new(Account)).Limit(int(limit)).Offset(int(offset)).Find(&accountList).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to get account list")
		return make([]*Account, 0), err
	}

	return accountList, nil
}

func (a accountDataAccessor) UpdateAccount(ctx context.Context, account *Account) error {
	logger := utils.LoggerWithContext(ctx, a.logger)
	db := a.db.WithContext(ctx)
	if err := db.Save(account).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to update account")
		return err
	}

	return nil
}

func (a accountDataAccessor) WithDB(db *gorm.DB) AccountDataAccessor {
	return &accountDataAccessor{
		db:     db,
		logger: a.logger,
	}
}
