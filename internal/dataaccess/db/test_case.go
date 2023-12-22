package db

import (
	"context"

	"gorm.io/gorm"
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
	IsHidden    bool
}

type TestCaseDataAccessor interface {
	CreateTestCase(ctx context.Context, testCase TestCase) error
	CreateTestCaseList(ctx context.Context, testCaseList []*TestCase) error
	GetTestCaseListOfProblem(ctx context.Context, problemID uint64, offset uint64, limit uint64) ([]*TestCase, error)
	GetTestCaseCountOfProblem(ctx context.Context, problemID uint64) (uint64, error)
	UpdateTestCase(ctx context.Context, testCase TestCase) error
	DeleteTestCase(ctx context.Context, id uint64) error
	WithDB(db *gorm.DB) TestCaseDataAccessor
}

type testCaseDataAccessor struct {
	db *gorm.DB
}

func NewTestCaseDataAccessor(
	db *gorm.DB,
) TestCaseDataAccessor {
	return &testCaseDataAccessor{
		db: db,
	}
}

func (a testCaseDataAccessor) CreateTestCase(ctx context.Context, testCase TestCase) error {
	db := a.db.WithContext(ctx)
	if err := db.Create(testCase).Error; err != nil {
		return err
	}

	return nil
}

func (a testCaseDataAccessor) UpdateTestCase(ctx context.Context, testCase TestCase) error {
	db := a.db.WithContext(ctx)
	if err := db.Save(testCase).Error; err != nil {
		return err
	}

	return nil
}

func (a testCaseDataAccessor) CreateTestCaseList(ctx context.Context, testCaseList []*TestCase) error {
	db := a.db.WithContext(ctx)
	if err := db.CreateInBatches(testCaseList, createTestCaseListBatchSize).Error; err != nil {
		return err
	}

	return nil
}

func (a testCaseDataAccessor) DeleteTestCase(ctx context.Context, id uint64) error {
	db := a.db.WithContext(ctx)
	if err := db.Model(new(TestCase)).
		Delete(id).
		Error; err != nil {
		return err
	}

	return nil
}

func (a testCaseDataAccessor) GetTestCaseListOfProblem(ctx context.Context, problemID uint64, offset uint64, limit uint64) ([]*TestCase, error) {
	db := a.db.WithContext(ctx)
	testCaseList := make([]*TestCase, 0)
	if err := db.Model(new(TestCase)).
		Where(&TestCase{
			OfProblemID: problemID,
		}).
		Find(testCaseList).
		Error; err != nil {
		return make([]*TestCase, 0), err
	}

	return testCaseList, nil
}

func (a testCaseDataAccessor) GetTestCaseCountOfProblem(ctx context.Context, problemID uint64) (uint64, error) {
	db := a.db.WithContext(ctx)
	count := int64(0)
	if err := db.Model(new(TestCase)).
		Where(&TestCase{
			OfProblemID: problemID,
		}).
		Count(&count).
		Error; err != nil {
		return 0, err
	}

	return uint64(count), nil
}

func (a testCaseDataAccessor) WithDB(db *gorm.DB) TestCaseDataAccessor {
	return &testCaseDataAccessor{
		db: db,
	}
}
