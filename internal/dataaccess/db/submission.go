package db

import (
	"context"

	"gorm.io/gorm"
)

type SubmissionStatus uint8
type SubmissionResult uint8

const (
	SubmissionStatusSubmitted SubmissionStatus = 1
	SubmissionStatusExecuting SubmissionStatus = 2
	SubmissionStatusFinished  SubmissionStatus = 3

	SubmissionResultOK                SubmissionResult = 1
	SubmissionResultCompileError      SubmissionResult = 2
	SubmissionResultRuntimeError      SubmissionResult = 3
	SubmissionResultTimeLimitExceeded SubmissionResult = 4
	SubmissionResultMemoryLimitExceed SubmissionResult = 5
	SubmissionResultWrongAnswer       SubmissionResult = 6
)

type Submission struct {
	gorm.Model
	OfProblemID  uint64
	Problem      Problem `gorm:"foreignKey:OfProblemID"`
	AuthorUserID uint64
	Author       User   `gorm:"foreignKey:AuthorUserID"`
	Content      string `gorm:"type:text"`
	Language     string
	Status       SubmissionStatus
	Result       SubmissionResult
}

type SubmissionDataAccessor interface {
	CreateSubmission(ctx context.Context, Submission *Submission) error
	UpdateSubmission(ctx context.Context, Submission *Submission) error
	GetSubmission(ctx context.Context, id uint64) (*Submission, error)
	GetSubmissionList(ctx context.Context, offset, limit uint64) ([]*Submission, error)
	GetUserSubmissionList(ctx context.Context, userID uint64, offset, limit uint64) ([]*Submission, error)
	GetProblemSubmissionList(ctx context.Context, problemID uint64, offset, limit uint64) ([]*Submission, error)
	GetSubmissionCount(ctx context.Context) (uint64, error)
	GetUserSubmissionCount(ctx context.Context, userID uint64) (uint64, error)
	GetProblemSubmissionCount(ctx context.Context, problemID uint64) (uint64, error)
	WithDB(db *gorm.DB) SubmissionDataAccessor
}

type submissionDataAccessor struct {
	db *gorm.DB
}

func NewSubmissionDataAccessor(db *gorm.DB) SubmissionDataAccessor {
	return &submissionDataAccessor{
		db: db,
	}
}

func (a submissionDataAccessor) CreateSubmission(ctx context.Context, Submission *Submission) error {
	panic("unimplemented")
}

func (a submissionDataAccessor) GetSubmission(ctx context.Context, id uint64) (*Submission, error) {
	panic("unimplemented")
}

func (a submissionDataAccessor) UpdateSubmission(ctx context.Context, Submission *Submission) error {
	panic("unimplemented")
}

func (a submissionDataAccessor) GetSubmissionCount(ctx context.Context) (uint64, error) {
	panic("unimplemented")
}

func (a submissionDataAccessor) GetUserSubmissionCount(ctx context.Context, userID uint64) (uint64, error) {
	panic("unimplemented")
}

func (a submissionDataAccessor) GetProblemSubmissionCount(ctx context.Context, problemID uint64) (uint64, error) {
	panic("unimplemented")
}

func (a submissionDataAccessor) GetSubmissionList(ctx context.Context, offset uint64, limit uint64) ([]*Submission, error) {
	panic("unimplemented")
}

func (a submissionDataAccessor) GetUserSubmissionList(ctx context.Context, userID uint64, offset, limit uint64) ([]*Submission, error) {
	panic("unimplemented")
}

func (a submissionDataAccessor) GetProblemSubmissionList(ctx context.Context, problemID uint64, offset, limit uint64) ([]*Submission, error) {
	panic("unimplemented")
}

func (a submissionDataAccessor) WithDB(db *gorm.DB) SubmissionDataAccessor {
	return &submissionDataAccessor{
		db: db,
	}
}
