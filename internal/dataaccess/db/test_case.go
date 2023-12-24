package db

import (
	"context"
	"encoding/base64"
	"errors"
	"hash/fnv"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/tranHieuDev23/cato/internal/utils"
)

const (
	createTestCaseListBatchSize = 10
)

type TestCase struct {
	gorm.Model
	OfProblemID uint64
	Problem     Problem `gorm:"foreignKey:OfProblemID"`
	Input       string  `gorm:"type:text"`
	Output      string  `gorm:"type:text"`
	Hash        string
	IsHidden    bool
}

func (t *TestCase) BeforeCreate() {
	fnvHash := fnv.New64a()
	fnvHash.Write([]byte(t.Input))
	fnvHash.Write([]byte{0})
	fnvHash.Write([]byte(t.Output))
	t.Hash = base64.StdEncoding.EncodeToString(fnvHash.Sum(nil))
}

type TestCaseDataAccessor interface {
	CreateTestCase(ctx context.Context, testCase *TestCase) error
	CreateTestCaseList(ctx context.Context, testCaseList []*TestCase) error
	GetTestCase(ctx context.Context, id uint64) (*TestCase, error)
	GetTestCaseListOfProblem(ctx context.Context, problemID uint64, offset uint64, limit uint64) ([]*TestCase, error)
	GetTestCaseCountOfProblem(ctx context.Context, problemID uint64) (uint64, error)
	UpdateTestCase(ctx context.Context, testCase *TestCase) error
	DeleteTestCase(ctx context.Context, id uint64) error
	WithDB(db *gorm.DB) TestCaseDataAccessor
}

type testCaseDataAccessor struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewTestCaseDataAccessor(
	db *gorm.DB,
	logger *zap.Logger,
) TestCaseDataAccessor {
	return &testCaseDataAccessor{
		db:     db,
		logger: logger,
	}
}

func (a testCaseDataAccessor) CreateTestCase(ctx context.Context, testCase *TestCase) error {
	logger := utils.LoggerWithContext(ctx, a.logger)
	db := a.db.WithContext(ctx)
	if err := db.Create(testCase).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to create test case")
		return err
	}

	return nil
}

func (a testCaseDataAccessor) UpdateTestCase(ctx context.Context, testCase *TestCase) error {
	logger := utils.LoggerWithContext(ctx, a.logger)
	db := a.db.WithContext(ctx)
	if err := db.Save(testCase).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to update test case")
		return err
	}

	return nil
}

func (a testCaseDataAccessor) GetTestCase(ctx context.Context, id uint64) (*TestCase, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).
		With(zap.Uint64("id", id))
	db := a.db.WithContext(ctx)
	testCase := new(TestCase)
	if err := db.First(testCase).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		logger.With(zap.Error(err)).Error("failed to get test case")
		return nil, err
	}

	return testCase, nil
}

func (a testCaseDataAccessor) CreateTestCaseList(ctx context.Context, testCaseList []*TestCase) error {
	db := a.db.WithContext(ctx)
	if err := db.CreateInBatches(testCaseList, createTestCaseListBatchSize).Error; err != nil {
		return err
	}

	return nil
}

func (a testCaseDataAccessor) DeleteTestCase(ctx context.Context, id uint64) error {
	logger := utils.LoggerWithContext(ctx, a.logger).
		With(zap.Uint64("id", id))
	db := a.db.WithContext(ctx)
	if err := db.Model(new(TestCase)).
		Delete(id).
		Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to update test case")
		return err
	}

	return nil
}

func (a testCaseDataAccessor) GetTestCaseListOfProblem(ctx context.Context, problemID uint64, offset uint64, limit uint64) ([]*TestCase, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).
		With(zap.Uint64("problem_id", problemID))
	db := a.db.WithContext(ctx)
	testCaseList := make([]*TestCase, 0)
	if err := db.Model(new(TestCase)).
		Where(&TestCase{
			OfProblemID: problemID,
		}).
		Find(testCaseList).
		Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to get test case list of problem")
		return make([]*TestCase, 0), err
	}

	return testCaseList, nil
}

func (a testCaseDataAccessor) GetTestCaseCountOfProblem(ctx context.Context, problemID uint64) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).
		With(zap.Uint64("problem_id", problemID))
	db := a.db.WithContext(ctx)
	count := int64(0)
	if err := db.Model(new(TestCase)).
		Where(&TestCase{
			OfProblemID: problemID,
		}).
		Count(&count).
		Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to get test case count of problem")
		return 0, err
	}

	return uint64(count), nil
}

func (a testCaseDataAccessor) WithDB(db *gorm.DB) TestCaseDataAccessor {
	return &testCaseDataAccessor{
		db:     db,
		logger: a.logger,
	}
}
