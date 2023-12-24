package db

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/tranHieuDev23/cato/internal/utils"
)

type ProblemTestCaseHash struct {
	gorm.Model
	OfProblemID uint64
	Problem     Problem `gorm:"foreignKey:OfProblemID"`
	Hash        string
}

type ProblemTestCaseHashDataAccessor interface {
	UpsertProblemTestCaseHash(ctx context.Context, problemTestCaseHash *ProblemTestCaseHash) error
	GetProblemTestCaseHashOfProblem(ctx context.Context, problemID uint64) (*ProblemTestCaseHash, error)
	WithDB(db *gorm.DB) ProblemTestCaseHashDataAccessor
}

type problemTestCaseHashDataAccessor struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewProblemTestCaseHashDataAccessor(
	db *gorm.DB,
	logger *zap.Logger,
) ProblemTestCaseHashDataAccessor {
	return &problemTestCaseHashDataAccessor{
		db:     db,
		logger: logger,
	}
}

func (a problemTestCaseHashDataAccessor) GetProblemTestCaseHashOfProblem(ctx context.Context, problemID uint64) (*ProblemTestCaseHash, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).
		With(zap.Uint64("problem_id", problemID))
	db := a.db.WithContext(ctx)
	hash := new(ProblemTestCaseHash)
	if err := db.Model(new(TestCase)).
		Where(&TestCase{
			OfProblemID: problemID,
		}).
		First(hash).
		Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to get test case hash list of problem")
		return nil, err
	}

	return hash, nil
}

func (a problemTestCaseHashDataAccessor) UpsertProblemTestCaseHash(ctx context.Context, problemTestCaseHash *ProblemTestCaseHash) error {
	logger := utils.LoggerWithContext(ctx, a.logger)
	db := a.db.WithContext(ctx)
	if err := db.Save(problemTestCaseHash).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to update problem test case hash")
		return err
	}

	return nil
}

func (a problemTestCaseHashDataAccessor) WithDB(db *gorm.DB) ProblemTestCaseHashDataAccessor {
	return &problemTestCaseHashDataAccessor{
		db:     db,
		logger: a.logger,
	}
}
