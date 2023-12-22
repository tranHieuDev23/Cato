package db

import (
	"context"

	"gorm.io/gorm"
)

const (
	createProblemExampleListBatchSize = 5
)

type ProblemExample struct {
	gorm.Model
	OfProblemID uint64
	Problem     Problem `gorm:"foreignKey:OfProblemID"`
	Input       string  `gorm:"type:text"`
	Output      string  `gorm:"type:text"`
}

type ProblemExampleDataAccessor interface {
	CreateProblemExampleList(ctx context.Context, problemExampleList []*ProblemExample) error
	GetProblemExampleListOfProblem(ctx context.Context, problemID uint64) ([]*ProblemExample, error)
	DeleteProblemExampleOfProblem(ctx context.Context, problemID uint64) error
	WithDB(db *gorm.DB) ProblemExampleDataAccessor
}

type problemExampleDataAccessor struct {
	db *gorm.DB
}

func NewProblemExampleDataAccessor(
	db *gorm.DB,
) ProblemExampleDataAccessor {
	return &problemExampleDataAccessor{
		db: db,
	}
}

func (a problemExampleDataAccessor) CreateProblemExampleList(ctx context.Context, problemExampleList []*ProblemExample) error {
	db := a.db.WithContext(ctx)
	if err := db.CreateInBatches(problemExampleList, createProblemExampleListBatchSize).Error; err != nil {
		return err
	}

	return nil
}

func (a problemExampleDataAccessor) DeleteProblemExampleOfProblem(ctx context.Context, problemID uint64) error {
	db := a.db.WithContext(ctx)
	if err := db.Model(new(ProblemExample)).
		Where(&ProblemExample{
			OfProblemID: problemID,
		}).
		Error; err != nil {
		return err
	}

	return nil
}

func (a problemExampleDataAccessor) GetProblemExampleListOfProblem(ctx context.Context, problemID uint64) ([]*ProblemExample, error) {
	db := a.db.WithContext(ctx)
	problemExampleList := make([]*ProblemExample, 0)
	if err := db.Model(new(ProblemExample)).
		Where(&ProblemExample{
			OfProblemID: problemID,
		}).
		Find(problemExampleList).
		Error; err != nil {
		return make([]*ProblemExample, 0), err
	}

	return problemExampleList, nil
}

func (a problemExampleDataAccessor) WithDB(db *gorm.DB) ProblemExampleDataAccessor {
	return &problemExampleDataAccessor{
		db: db,
	}
}
