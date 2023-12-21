package db

import (
	"context"

	"gorm.io/gorm"
)

type Problem struct {
	gorm.Model
	DisplayName            string
	AuthorUserID           uint64
	Author                 User   `gorm:"foreignKey:AuthorUserID"`
	Description            string `gorm:"type:text"`
	TimeLimitInMillisecond uint64
	MemoryLimitInByte      uint64
}

type ProblemDataAccessor interface {
	CreateProblem(ctx context.Context, Problem *Problem) error
	UpdateProblem(ctx context.Context, Problem *Problem) error
	GetProblem(ctx context.Context, id uint64) (*Problem, error)
	GetProblemList(ctx context.Context, offset, limit uint64) ([]*Problem, error)
	GetUserProblemList(ctx context.Context, userID uint64, offset, limit uint64) ([]*Problem, error)
	GetProblemCount(ctx context.Context) (uint64, error)
	GetUserProblemCount(ctx context.Context, userID uint64) (uint64, error)
	WithDB(db *gorm.DB) ProblemDataAccessor
}

type problemDataAccessor struct {
	db *gorm.DB
}

func NewProblemDataAccessor(db *gorm.DB) ProblemDataAccessor {
	return &problemDataAccessor{
		db: db,
	}
}

func (*problemDataAccessor) CreateProblem(ctx context.Context, Problem *Problem) error {
	panic("unimplemented")
}

func (*problemDataAccessor) GetProblem(ctx context.Context, id uint64) (*Problem, error) {
	panic("unimplemented")
}

func (*problemDataAccessor) GetProblemCount(ctx context.Context) (uint64, error) {
	panic("unimplemented")
}

func (*problemDataAccessor) GetUserProblemCount(ctx context.Context, userID uint64) (uint64, error) {
	panic("unimplemented")
}

func (*problemDataAccessor) GetProblemList(ctx context.Context, offset uint64, limit uint64) ([]*Problem, error) {
	panic("unimplemented")
}

func (*problemDataAccessor) GetUserProblemList(ctx context.Context, userID uint64, offset, limit uint64) ([]*Problem, error) {
	panic("unimplemented")
}

func (*problemDataAccessor) UpdateProblem(ctx context.Context, Problem *Problem) error {
	panic("unimplemented")
}

func (a problemDataAccessor) WithDB(db *gorm.DB) ProblemDataAccessor {
	return &problemDataAccessor{
		db: db,
	}
}
