package db

import (
	"context"

	"gorm.io/gorm"
)

type UserPassword struct {
	gorm.Model
	OfUserID uint64
	User     User `gorm:"foreignKey:OfUserID"`
	Hash     string
}

type UserPasswordDataAccessor interface {
	CreateUserPassword(ctx context.Context, userPassword *UserPassword) error
	UpdateUserPassword(ctx context.Context, userPassword *UserPassword) error
	GetUserPassword(ctx context.Context, id uint64) (*UserPassword, error)
	WithDB(db *gorm.DB) UserPasswordDataAccessor
}

type userPasswordDataAccessor struct {
	db *gorm.DB
}

func NewUserPasswordDataAccessor(
	db *gorm.DB,
) UserPasswordDataAccessor {
	return &userPasswordDataAccessor{
		db: db,
	}
}

func (a userPasswordDataAccessor) CreateUserPassword(ctx context.Context, userPassword *UserPassword) error {
	db := a.db.WithContext(ctx)
	if err := db.Create(userPassword).Error; err != nil {
		return err
	}

	return nil
}

func (a userPasswordDataAccessor) GetUserPassword(ctx context.Context, ofUserID uint64) (*UserPassword, error) {
	db := a.db.WithContext(ctx)
	userPassword := new(UserPassword)
	if err := db.Model(new(UserPassword)).
		Where(&UserPassword{
			OfUserID: ofUserID,
		}).
		First(userPassword).
		Error; err != nil {
		return nil, err
	}

	return userPassword, nil
}

func (a userPasswordDataAccessor) UpdateUserPassword(ctx context.Context, userPassword *UserPassword) error {
	db := a.db.WithContext(ctx)
	if err := db.Save(userPassword).Error; err != nil {
		return err
	}

	return nil
}

func (a userPasswordDataAccessor) WithDB(db *gorm.DB) UserPasswordDataAccessor {
	return &userPasswordDataAccessor{
		db: db,
	}
}
