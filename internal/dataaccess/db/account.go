package db

import (
	"context"

	"gorm.io/gorm"
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
	AccountName string
	DisplayName string
	Role        AccountRole
}

type AccountDataAccessor interface {
	CreateAccount(ctx context.Context, account *Account) error
	UpdateAccount(ctx context.Context, account *Account) error
	GetAccount(ctx context.Context, id uint64) (*Account, error)
	GetAccountList(ctx context.Context, offset, limit uint64) ([]*Account, error)
	GetAccountCount(ctx context.Context) (uint64, error)
	WithDB(db *gorm.DB) AccountDataAccessor
}

type accountDataAccessor struct {
	db *gorm.DB
}

func NewAccountDataAccessor(
	db *gorm.DB,
) AccountDataAccessor {
	return &accountDataAccessor{
		db: db,
	}
}

func (a accountDataAccessor) CreateAccount(ctx context.Context, account *Account) error {
	db := a.db.WithContext(ctx)
	if err := db.Create(account).Error; err != nil {
		return err
	}

	return nil
}

func (a accountDataAccessor) GetAccount(ctx context.Context, id uint64) (*Account, error) {
	db := a.db.WithContext(ctx)
	account := new(Account)
	if err := db.First(account, id).Error; err != nil {
		return nil, err
	}

	return account, nil
}

func (a accountDataAccessor) GetAccountCount(ctx context.Context) (uint64, error) {
	db := a.db.WithContext(ctx)
	count := int64(0)
	if err := db.Model(new(Account)).Count(&count).Error; err != nil {
		return 0, err
	}

	return uint64(count), nil
}

func (a accountDataAccessor) GetAccountList(ctx context.Context, offset uint64, limit uint64) ([]*Account, error) {
	db := a.db.WithContext(ctx)
	accountList := make([]*Account, 0)
	if err := db.Model(new(Account)).Limit(int(limit)).Offset(int(offset)).Find(&accountList).Error; err != nil {
		return make([]*Account, 0), err
	}

	return accountList, nil
}

func (a accountDataAccessor) UpdateAccount(ctx context.Context, account *Account) error {
	db := a.db.WithContext(ctx)
	if err := db.Save(account).Error; err != nil {
		return err
	}

	return nil
}

func (a accountDataAccessor) WithDB(db *gorm.DB) AccountDataAccessor {
	return &accountDataAccessor{
		db: db,
	}
}
