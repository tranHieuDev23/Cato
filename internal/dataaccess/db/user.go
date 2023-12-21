package db

import (
	"context"

	"gorm.io/gorm"
)

type UserRole string

const (
	UserRoleAdmin         UserRole = "admin"
	UserRoleProblemSetter UserRole = "problem_setter"
	UserRoleContestant    UserRole = "contestant"
)

type User struct {
	gorm.Model
	Username    string
	DisplayName string
	Role        UserRole
}

type UserDataAccessor interface {
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id uint64) (*User, error)
	GetUserList(ctx context.Context, offset, limit uint64) ([]*User, error)
	GetUserCount(ctx context.Context) (uint64, error)
	WithDB(db *gorm.DB) UserDataAccessor
}

type userDataAccessor struct {
	db *gorm.DB
}

func NewUserDataAccessor(
	db *gorm.DB,
) UserDataAccessor {
	return &userDataAccessor{
		db: db,
	}
}

func (a userDataAccessor) CreateUser(ctx context.Context, user *User) error {
	db := a.db.WithContext(ctx)
	if err := db.Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (a userDataAccessor) GetUser(ctx context.Context, id uint64) (*User, error) {
	db := a.db.WithContext(ctx)
	user := new(User)
	if err := db.First(user, id).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (a userDataAccessor) GetUserCount(ctx context.Context) (uint64, error) {
	db := a.db.WithContext(ctx)
	count := int64(0)
	if err := db.Model(new(User)).Count(&count).Error; err != nil {
		return 0, err
	}

	return uint64(count), nil
}

func (a userDataAccessor) GetUserList(ctx context.Context, offset uint64, limit uint64) ([]*User, error) {
	db := a.db.WithContext(ctx)
	userList := make([]*User, 0)
	if err := db.Model(new(User)).Limit(int(limit)).Offset(int(offset)).Find(&userList).Error; err != nil {
		return make([]*User, 0), err
	}

	return userList, nil
}

func (a userDataAccessor) UpdateUser(ctx context.Context, user *User) error {
	db := a.db.WithContext(ctx)
	if err := db.Save(user).Error; err != nil {
		return err
	}

	return nil
}

func (a userDataAccessor) WithDB(db *gorm.DB) UserDataAccessor {
	return &userDataAccessor{
		db: db,
	}
}
