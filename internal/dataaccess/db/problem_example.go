package db

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/tranHieuDev23/cato/internal/utils"
)

const (
	createProblemExampleListBatchSize = 5
)

type ProblemExample struct {
	gorm.Model
	OfProblemID uint64
	Input       string `gorm:"type:text"`
	Output      string `gorm:"type:text"`
}

type ProblemExampleDataAccessor interface {
	CreateProblemExampleList(ctx context.Context, problemExampleList []*ProblemExample) error
	GetProblemExampleListOfProblem(ctx context.Context, problemID uint64) ([]*ProblemExample, error)
	DeleteProblemExampleOfProblem(ctx context.Context, problemID uint64) error
	WithDB(db *gorm.DB) ProblemExampleDataAccessor
}

type problemExampleDataAccessor struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewProblemExampleDataAccessor(
	db *gorm.DB,
	logger *zap.Logger,
) ProblemExampleDataAccessor {
	return &problemExampleDataAccessor{
		db:     db,
		logger: logger,
	}
}

func (a problemExampleDataAccessor) CreateProblemExampleList(ctx context.Context, problemExampleList []*ProblemExample) error {
	logger := utils.LoggerWithContext(ctx, a.logger)
	db := a.db.WithContext(ctx)
	if err := db.CreateInBatches(problemExampleList, createProblemExampleListBatchSize).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to create problem example list")
		return err
	}

	return nil
}

func (a problemExampleDataAccessor) DeleteProblemExampleOfProblem(ctx context.Context, problemID uint64) error {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.Uint64("problem_id", problemID))
	db := a.db.WithContext(ctx)
	if err := db.
		Where(&ProblemExample{
			OfProblemID: problemID,
		}).
		Delete(new(ProblemExample)).
		Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to delete problem example list")
		return err
	}

	return nil
}

func (a problemExampleDataAccessor) GetProblemExampleListOfProblem(ctx context.Context, problemID uint64) ([]*ProblemExample, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.Uint64("problem_id", problemID))
	db := a.db.WithContext(ctx)
	problemExampleList := make([]*ProblemExample, 0)
	if err := db.Model(new(ProblemExample)).
		Where(&ProblemExample{
			OfProblemID: problemID,
		}).
		Find(&problemExampleList).
		Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to get problem example list")
		return make([]*ProblemExample, 0), err
	}

	return problemExampleList, nil
}

func (a problemExampleDataAccessor) WithDB(db *gorm.DB) ProblemExampleDataAccessor {
	return &problemExampleDataAccessor{
		db:     db,
		logger: a.logger,
	}
}
