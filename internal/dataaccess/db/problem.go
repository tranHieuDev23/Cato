package db

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/tranHieuDev23/cato/internal/utils"
)

type Problem struct {
	gorm.Model
	DisplayName            string
	AuthorAccountID        uint64
	Author                 Account `gorm:"foreignKey:AuthorAccountID"`
	Description            string  `gorm:"type:text"`
	TimeLimitInMillisecond uint64
	MemoryLimitInByte      uint64
}

type ProblemDataAccessor interface {
	CreateProblem(ctx context.Context, problem *Problem) error
	UpdateProblem(ctx context.Context, problem *Problem) error
	DeleteProblem(ctx context.Context, id uint64) error
	GetProblem(ctx context.Context, id uint64) (*Problem, error)
	GetProblemList(ctx context.Context, offset, limit uint64) ([]*Problem, error)
	GetAccountProblemList(ctx context.Context, accountID uint64, offset, limit uint64) ([]*Problem, error)
	GetProblemCount(ctx context.Context) (uint64, error)
	GetAccountProblemCount(ctx context.Context, accountID uint64) (uint64, error)
	WithDB(db *gorm.DB) ProblemDataAccessor
}

type problemDataAccessor struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewProblemDataAccessor(
	db *gorm.DB,
	logger *zap.Logger,
) ProblemDataAccessor {
	return &problemDataAccessor{
		db:     db,
		logger: logger,
	}
}

func (a problemDataAccessor) CreateProblem(ctx context.Context, problem *Problem) error {
	logger := utils.LoggerWithContext(ctx, a.logger)
	db := a.db.WithContext(ctx)
	if err := db.Create(problem).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to create problem")
		return err
	}

	return nil
}

func (a problemDataAccessor) GetProblem(ctx context.Context, id uint64) (*Problem, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).
		With(zap.Uint64("id", id))
	db := a.db.WithContext(ctx)
	problem := new(Problem)
	if err := db.First(problem, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		logger.With(zap.Error(err)).Error("failed to get problem")
		return nil, err
	}

	return problem, nil
}

func (a problemDataAccessor) GetProblemCount(ctx context.Context) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, a.logger)
	db := a.db.WithContext(ctx)
	count := int64(0)
	if err := db.Model(new(Problem)).Count(&count).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to get problem account")
		return 0, err
	}

	return uint64(count), nil
}

func (a problemDataAccessor) GetAccountProblemCount(ctx context.Context, accountID uint64) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).
		With(zap.Uint64("account_id", accountID))
	db := a.db.WithContext(ctx)
	count := int64(0)
	if err := db.Model(new(Problem)).
		Where(&Problem{
			AuthorAccountID: accountID,
		}).
		Count(&count).
		Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to get account problem account")
		return 0, err
	}

	return uint64(count), nil
}

func (a problemDataAccessor) GetProblemList(ctx context.Context, offset uint64, limit uint64) ([]*Problem, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).
		With(zap.Uint64("limit", limit)).
		With(zap.Uint64("offset", offset))
	db := a.db.WithContext(ctx)
	problemList := make([]*Problem, 0)
	if err := db.Model(new(Problem)).Limit(int(limit)).Offset(int(offset)).Find(&problemList).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to get problem list")
		return make([]*Problem, 0), err
	}

	return problemList, nil
}

func (a problemDataAccessor) GetAccountProblemList(ctx context.Context, accountID uint64, offset, limit uint64) ([]*Problem, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).
		With(zap.Uint64("account_id", accountID)).
		With(zap.Uint64("limit", limit)).
		With(zap.Uint64("offset", offset))
	db := a.db.WithContext(ctx)
	problemList := make([]*Problem, 0)
	if err := db.Model(new(Problem)).
		Limit(int(limit)).
		Offset(int(offset)).
		Where(&Problem{
			AuthorAccountID: accountID,
		}).
		Find(&problemList).
		Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to get account problem list")
		return make([]*Problem, 0), err
	}

	return problemList, nil
}

func (a problemDataAccessor) UpdateProblem(ctx context.Context, problem *Problem) error {
	logger := utils.LoggerWithContext(ctx, a.logger)
	db := a.db.WithContext(ctx)
	if err := db.Save(problem).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to update problem")
		return err
	}

	return nil
}

func (a problemDataAccessor) DeleteProblem(ctx context.Context, id uint64) error {
	logger := utils.LoggerWithContext(ctx, a.logger).
		With(zap.Uint64("id", id))
	db := a.db.WithContext(ctx)
	if err := db.Delete(&Problem{
		Model: gorm.Model{ID: uint(id)},
	}).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to delete problem")
		return err
	}

	return nil
}

func (a problemDataAccessor) WithDB(db *gorm.DB) ProblemDataAccessor {
	return &problemDataAccessor{
		db:     db,
		logger: a.logger,
	}
}
