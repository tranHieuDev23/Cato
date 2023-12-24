package db

import (
	"context"

	"gitlab.com/pjrpc/pjrpc/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/tranHieuDev23/cato/internal/utils"
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
	OfProblemID     uint64
	Problem         Problem `gorm:"foreignKey:OfProblemID"`
	AuthorAccountID uint64
	Author          Account `gorm:"foreignKey:AuthorAccountID"`
	Content         string  `gorm:"type:text"`
	Language        string
	Status          SubmissionStatus
	Result          SubmissionResult
}

type SubmissionDataAccessor interface {
	CreateSubmission(ctx context.Context, submission *Submission) error
	UpdateSubmission(ctx context.Context, submission *Submission) error
	DeleteSubmission(ctx context.Context, id uint64) error
	GetSubmission(ctx context.Context, id uint64) (*Submission, error)
	GetSubmissionList(ctx context.Context, offset, limit uint64) ([]*Submission, error)
	GetAccountSubmissionList(ctx context.Context, accountID uint64, offset, limit uint64) ([]*Submission, error)
	GetProblemSubmissionList(ctx context.Context, problemID uint64, offset, limit uint64) ([]*Submission, error)
	GetAccountProblemSubmissionList(ctx context.Context, accountID, problemID, offset, limit uint64) ([]*Submission, error)
	GetSubmissionCount(ctx context.Context) (uint64, error)
	GetAccountSubmissionCount(ctx context.Context, accountID uint64) (uint64, error)
	GetProblemSubmissionCount(ctx context.Context, problemID uint64) (uint64, error)
	GetAccountProblemSubmissionCount(ctx context.Context, accountID, problemID uint64) (uint64, error)
	WithDB(db *gorm.DB) SubmissionDataAccessor
}

type submissionDataAccessor struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewSubmissionDataAccessor(
	db *gorm.DB,
	logger *zap.Logger,
) SubmissionDataAccessor {
	return &submissionDataAccessor{
		db:     db,
		logger: logger,
	}
}

func (a submissionDataAccessor) CreateSubmission(ctx context.Context, submission *Submission) error {
	logger := utils.LoggerWithContext(ctx, a.logger).
		With(zap.Any("submission", submission))
	db := a.db.WithContext(ctx)
	if err := db.Create(submission).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to create submission")
		return pjrpc.JRPCErrInternalError()
	}

	return nil
}

func (a submissionDataAccessor) GetSubmission(ctx context.Context, id uint64) (*Submission, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).
		With(zap.Any("id", id))
	db := a.db.WithContext(ctx)
	submission := &Submission{}
	if err := db.First(submission, id).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to get submission")
		return nil, pjrpc.JRPCErrInternalError()
	}

	return submission, nil
}

func (a submissionDataAccessor) UpdateSubmission(ctx context.Context, submission *Submission) error {
	logger := utils.LoggerWithContext(ctx, a.logger).
		With(zap.Any("submission", submission))
	db := a.db.WithContext(ctx)
	if err := db.Save(submission).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to update submission")
		return pjrpc.JRPCErrInternalError()
	}

	return nil
}

func (a submissionDataAccessor) DeleteSubmission(ctx context.Context, id uint64) error {
	logger := utils.LoggerWithContext(ctx, a.logger).
		With(zap.Any("id", id))
	db := a.db.WithContext(ctx)
	if err := db.Delete(new(Submission), id).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to delete submission")
		return pjrpc.JRPCErrInternalError()
	}

	return nil
}

func (a submissionDataAccessor) GetSubmissionCount(ctx context.Context) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, a.logger)
	db := a.db.WithContext(ctx)
	count := int64(0)
	if err := db.Model(new(Submission)).Count(&count).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to get submission count")
		return 0, err
	}

	return uint64(count), nil
}

func (a submissionDataAccessor) GetAccountSubmissionCount(ctx context.Context, accountID uint64) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).
		With(zap.Uint64("account_id", accountID))
	db := a.db.WithContext(ctx)
	count := int64(0)
	if err := db.Model(new(Submission)).
		Where(&Submission{
			AuthorAccountID: accountID,
		}).
		Count(&count).
		Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to get submission count of account")
		return 0, err
	}

	return uint64(count), nil
}

func (a submissionDataAccessor) GetProblemSubmissionCount(ctx context.Context, problemID uint64) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).
		With(zap.Uint64("problem_id", problemID))
	db := a.db.WithContext(ctx)
	count := int64(0)
	if err := db.Model(new(Submission)).
		Where(&Submission{
			OfProblemID: problemID,
		}).
		Count(&count).
		Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to get submission count of problem")
		return 0, err
	}

	return uint64(count), nil
}

func (a submissionDataAccessor) GetSubmissionList(ctx context.Context, offset uint64, limit uint64) ([]*Submission, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).
		With(zap.Uint64("limit", limit)).
		With(zap.Uint64("offset", offset))
	db := a.db.WithContext(ctx)
	submissionList := make([]*Submission, 0)
	if err := db.Model(new(Submission)).Limit(int(limit)).Offset(int(offset)).Find(&submissionList).Error; err != nil {
		logger.
			With(zap.Error(err)).
			Error("failed to get submission list")
		return make([]*Submission, 0), err
	}

	return submissionList, nil
}

func (a submissionDataAccessor) GetAccountSubmissionList(ctx context.Context, accountID uint64, offset, limit uint64) ([]*Submission, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).
		With(zap.Uint64("account_id", accountID)).
		With(zap.Uint64("limit", limit)).
		With(zap.Uint64("offset", offset))
	db := a.db.WithContext(ctx)
	submissionList := make([]*Submission, 0)
	if err := db.Model(new(Submission)).
		Where(&Submission{
			AuthorAccountID: accountID,
		}).
		Limit(int(limit)).
		Offset(int(offset)).
		Find(&submissionList).
		Error; err != nil {
		logger.
			With(zap.Error(err)).
			Error("failed to get submission list of account")
		return make([]*Submission, 0), err
	}

	return submissionList, nil
}

func (a submissionDataAccessor) GetProblemSubmissionList(ctx context.Context, problemID uint64, offset, limit uint64) ([]*Submission, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).
		With(zap.Uint64("problem_id", problemID)).
		With(zap.Uint64("limit", limit)).
		With(zap.Uint64("offset", offset))
	db := a.db.WithContext(ctx)
	submissionList := make([]*Submission, 0)
	if err := db.Model(new(Submission)).
		Where(&Submission{
			OfProblemID: problemID,
		}).
		Limit(int(limit)).
		Offset(int(offset)).
		Find(&submissionList).
		Error; err != nil {
		logger.
			With(zap.Error(err)).
			Error("failed to get submission list of problem")
		return make([]*Submission, 0), err
	}

	return submissionList, nil
}

func (a submissionDataAccessor) GetAccountProblemSubmissionCount(ctx context.Context, accountID, problemID uint64) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).
		With(zap.Uint64("account_id", accountID)).
		With(zap.Uint64("problem_id", problemID))
	db := a.db.WithContext(ctx)
	count := int64(0)
	if err := db.Model(new(Submission)).
		Where(&Submission{
			AuthorAccountID: accountID,
			OfProblemID:     problemID,
		}).
		Count(&count).
		Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to get submission count of problem")
		return 0, err
	}

	return uint64(count), nil
}

func (a submissionDataAccessor) GetAccountProblemSubmissionList(
	ctx context.Context,
	accountID,
	problemID,
	offset,
	limit uint64,
) ([]*Submission, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).
		With(zap.Uint64("account_id", accountID)).
		With(zap.Uint64("problem_id", problemID)).
		With(zap.Uint64("limit", limit)).
		With(zap.Uint64("offset", offset))
	db := a.db.WithContext(ctx)
	submissionList := make([]*Submission, 0)
	if err := db.Model(new(Submission)).
		Where(&Submission{
			AuthorAccountID: accountID,
			OfProblemID:     problemID,
		}).
		Limit(int(limit)).
		Offset(int(offset)).
		Find(&submissionList).
		Error; err != nil {
		logger.
			With(zap.Error(err)).
			Error("failed to get submission list of account")
		return make([]*Submission, 0), err
	}

	return submissionList, nil
}

func (a submissionDataAccessor) WithDB(db *gorm.DB) SubmissionDataAccessor {
	return &submissionDataAccessor{
		db:     db,
		logger: a.logger,
	}
}
