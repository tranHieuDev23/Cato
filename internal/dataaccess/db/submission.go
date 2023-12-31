package db

import (
	"context"
	"time"

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

	SubmissionResultOK                  SubmissionResult = 1
	SubmissionResultCompileError        SubmissionResult = 2
	SubmissionResultRuntimeError        SubmissionResult = 3
	SubmissionResultTimeLimitExceeded   SubmissionResult = 4
	SubmissionResultMemoryLimitExceeded SubmissionResult = 5
	SubmissionResultWrongAnswer         SubmissionResult = 6
	SubmissionResultUnsupportedLanguage SubmissionResult = 7
)

type Submission struct {
	gorm.Model
	OfProblemID     uint64           `gorm:"index"`
	AuthorAccountID uint64           `gorm:"index"`
	Content         string           `gorm:"type:text"`
	Language        string           `gorm:"type:varchar(16)"`
	Status          SubmissionStatus `gorm:"index"`
	Result          SubmissionResult
}

type SubmissionListFilterParams struct {
	OfProblemID     *uint64
	AuthorAccountID *uint64
	Status          SubmissionStatus
}

type SubmissionDataAccessor interface {
	CreateSubmission(ctx context.Context, submission *Submission) error
	UpdateSubmission(ctx context.Context, submission *Submission) error
	UpdateExecutingSubmissionUpdatedBeforeThresholdToSubmitted(ctx context.Context, threshold time.Time) error
	DeleteSubmission(ctx context.Context, id uint64) error
	DeleteSubmissionOfProblem(ctx context.Context, problemID uint64) error
	GetSubmission(ctx context.Context, id uint64) (*Submission, error)
	GetSubmissionList(
		ctx context.Context,
		filterParams SubmissionListFilterParams,
		offset,
		limit uint64,
	) ([]*Submission, error)
	GetSubmissionIDList(ctx context.Context, filterParams SubmissionListFilterParams) ([]uint64, error)
	GetSubmissionCount(ctx context.Context, filterParams SubmissionListFilterParams) (uint64, error)
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
	logger := utils.LoggerWithContext(ctx, a.logger)
	db := a.db.WithContext(ctx)
	if err := db.Create(submission).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to create submission")
		return pjrpc.JRPCErrInternalError()
	}

	return nil
}

func (a submissionDataAccessor) GetSubmission(ctx context.Context, id uint64) (*Submission, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.Uint64("id", id))
	db := a.db.WithContext(ctx)
	submission := &Submission{}
	if err := db.First(submission, id).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to get submission")
		return nil, pjrpc.JRPCErrInternalError()
	}

	return submission, nil
}

func (a submissionDataAccessor) UpdateSubmission(ctx context.Context, submission *Submission) error {
	logger := utils.LoggerWithContext(ctx, a.logger)
	db := a.db.WithContext(ctx)
	if err := db.Save(submission).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to update submission")
		return pjrpc.JRPCErrInternalError()
	}

	return nil
}

func (a submissionDataAccessor) UpdateExecutingSubmissionUpdatedBeforeThresholdToSubmitted(
	ctx context.Context,
	threshold time.Time,
) error {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.Time("threshold", threshold))
	db := a.db.WithContext(ctx)
	if err := db.Model(new(Submission)).
		Where("updated_at < ? and status = ?", threshold, SubmissionStatusExecuting).
		UpdateColumn("status", SubmissionStatusSubmitted).
		Error; err != nil {
		logger.
			With(zap.Error(err)).
			Error("failed to update executing submission updated before threshold to submitted")
		return err
	}

	return nil
}

func (a submissionDataAccessor) DeleteSubmission(ctx context.Context, id uint64) error {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.Uint64("id", id))
	db := a.db.WithContext(ctx)
	if err := db.Delete(&Submission{Model: gorm.Model{ID: uint(id)}}).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to delete submission")
		return pjrpc.JRPCErrInternalError()
	}

	return nil
}

func (a submissionDataAccessor) getGormDBFromFilterParams(
	ctx context.Context,
	filterParams SubmissionListFilterParams,
) *gorm.DB {
	whereClause := &Submission{}
	if filterParams.OfProblemID != nil {
		whereClause.OfProblemID = *filterParams.OfProblemID
	}

	if filterParams.AuthorAccountID != nil {
		whereClause.AuthorAccountID = *filterParams.AuthorAccountID
	}

	if filterParams.Status != 0 {
		whereClause.Status = filterParams.Status
	}

	return a.db.WithContext(ctx).Model(new(Submission)).Where(whereClause)
}

func (a submissionDataAccessor) GetSubmissionCount(
	ctx context.Context,
	filterParams SubmissionListFilterParams,
) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.Any("filter_params", filterParams))
	db := a.getGormDBFromFilterParams(ctx, filterParams)
	count := int64(0)
	if err := db.Model(new(Submission)).Count(&count).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to get submission count")
		return 0, err
	}

	return uint64(count), nil
}

func (a submissionDataAccessor) GetSubmissionList(
	ctx context.Context,
	filterParams SubmissionListFilterParams,
	offset uint64,
	limit uint64,
) ([]*Submission, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).
		With(zap.Any("filter_params", filterParams)).
		With(zap.Uint64("limit", limit)).
		With(zap.Uint64("offset", offset))
	db := a.getGormDBFromFilterParams(ctx, filterParams)
	submissionList := make([]*Submission, 0)
	if err := db.Model(new(Submission)).
		Order("id desc").
		Limit(int(limit)).
		Offset(int(offset)).
		Find(&submissionList).Error; err != nil {
		logger.
			With(zap.Error(err)).
			Error("failed to get submission list")
		return make([]*Submission, 0), err
	}

	return submissionList, nil
}

func (a submissionDataAccessor) GetSubmissionIDList(
	ctx context.Context,
	filterParams SubmissionListFilterParams,
) ([]uint64, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.Any("filter_params", filterParams))
	db := a.getGormDBFromFilterParams(ctx, filterParams)
	idList := make([]uint64, 0)
	if err := db.Pluck("id", &idList).Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to get submitted submission id list")
		return make([]uint64, 0), err
	}

	return idList, nil
}

func (a submissionDataAccessor) DeleteSubmissionOfProblem(ctx context.Context, problemID uint64) error {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.Uint64("problem_id", problemID))
	db := a.db.WithContext(ctx)
	if err := db.
		Where(&Submission{
			OfProblemID: problemID,
		}).
		Delete(new(Submission)).
		Error; err != nil {
		logger.With(zap.Error(err)).Error("failed to delete submission of problem")
		return err
	}

	return nil
}

func (a submissionDataAccessor) WithDB(db *gorm.DB) SubmissionDataAccessor {
	return &submissionDataAccessor{
		db:     db,
		logger: a.logger,
	}
}
