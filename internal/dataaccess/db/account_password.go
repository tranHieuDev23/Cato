package db

import (
	"context"
	"errors"

	"gorm.io/gorm"
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
	db *gorm.DB
}

func NewAccountPasswordDataAccessor(
	db *gorm.DB,
) AccountPasswordDataAccessor {
	return &accountPasswordDataAccessor{
		db: db,
	}
}

func (a accountPasswordDataAccessor) CreateAccountPassword(ctx context.Context, accountPassword *AccountPassword) error {
	db := a.db.WithContext(ctx)
	if err := db.Create(accountPassword).Error; err != nil {
		return err
	}

	return nil
}

func (a accountPasswordDataAccessor) GetAccountPasswordOfAccountID(ctx context.Context, ofAccountID uint64) (*AccountPassword, error) {
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

		return nil, err
	}

	return accountPassword, nil
}

func (a accountPasswordDataAccessor) UpdateAccountPassword(ctx context.Context, accountPassword *AccountPassword) error {
	db := a.db.WithContext(ctx)
	if err := db.Save(accountPassword).Error; err != nil {
		return err
	}

	return nil
}

func (a accountPasswordDataAccessor) WithDB(db *gorm.DB) AccountPasswordDataAccessor {
	return &accountPasswordDataAccessor{
		db: db,
	}
}
