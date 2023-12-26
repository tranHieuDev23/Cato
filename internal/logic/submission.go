package logic

import (
	"context"
	"time"

	"gitlab.com/pjrpc/pjrpc/v2"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"

	"github.com/mikespook/gorbac"

	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc"
	"github.com/tranHieuDev23/cato/internal/utils"
)

type Submission interface {
	CreateSubmission(
		ctx context.Context,
		in *rpc.CreateSubmissionRequest,
		token string,
	) (*rpc.CreateSubmissionResponse, error)
	UpdateSubmission(
		ctx context.Context,
		in *rpc.UpdateSubmissionRequest,
		token string,
	) (*rpc.UpdateSubmissionResponse, error)
	GetSubmissionSnippetList(
		ctx context.Context,
		in *rpc.GetSubmissionSnippetListRequest,
		token string,
	) (*rpc.GetSubmissionSnippetListResponse, error)
	GetSubmission(ctx context.Context, in *rpc.GetSubmissionRequest, token string) (*rpc.GetSubmissionResponse, error)
	DeleteSubmission(ctx context.Context, in *rpc.DeleteSubmissionRequest, token string) error
	GetAccountSubmissionSnippetList(
		ctx context.Context,
		in *rpc.GetAccountSubmissionSnippetListRequest,
		token string,
	) (*rpc.GetAccountSubmissionSnippetListResponse, error)
	GetProblemSubmissionSnippetList(
		ctx context.Context,
		in *rpc.GetProblemSubmissionSnippetListRequest,
		token string,
	) (*rpc.GetProblemSubmissionSnippetListResponse, error)
	GetAccountProblemSubmissionSnippetList(
		ctx context.Context,
		in *rpc.GetAccountProblemSubmissionSnippetListRequest,
		token string,
	) (*rpc.GetAccountProblemSubmissionSnippetListResponse, error)
	GetAndUpdateFirstSubmittedSubmissionToExecuting(
		ctx context.Context,
		token string,
	) (*rpc.GetAndUpdateFirstSubmittedSubmissionToExecutingResponse, error)
	ScheduleSubmittedExecutingSubmissionToJudge(ctx context.Context) error
	UpdateExecutingSubmissionUpdatedBeforeThresholdToSubmitted(ctx context.Context) error
}

type submission struct {
	token                                       Token
	role                                        Role
	judge                                       Judge
	accountDataAccessor                         db.AccountDataAccessor
	problemDataAccessor                         db.ProblemDataAccessor
	submissionDataAccessor                      db.SubmissionDataAccessor
	db                                          *gorm.DB
	logger                                      *zap.Logger
	appArguments                                utils.Arguments
	revertExecutingSubmissionsThresholdDuration time.Duration
}

func NewSubmission(
	token Token,
	role Role,
	judge Judge,
	accountDataAccessor db.AccountDataAccessor,
	problemDataAccessor db.ProblemDataAccessor,
	submissionDataAccessor db.SubmissionDataAccessor,
	db *gorm.DB,
	logger *zap.Logger,
	appArguments utils.Arguments,
	logicConfig configs.Logic,
) (Submission, error) {
	revertExecutingSubmissionsThresholdDuration, err := logicConfig.RevertExecutingSubmissions.GetThresholdDuration()
	if err != nil {
		logger.
			With(zap.String("revert_executing_submission_threshold", logicConfig.RevertExecutingSubmissions.Threshold)).
			With(zap.Error(err)).
			Error("failed to parse revert execution submission threshold")
		return nil, err
	}

	return &submission{
		token:                  token,
		role:                   role,
		judge:                  judge,
		accountDataAccessor:    accountDataAccessor,
		problemDataAccessor:    problemDataAccessor,
		submissionDataAccessor: submissionDataAccessor,
		db:                     db,
		logger:                 logger,
		appArguments:           appArguments,
		revertExecutingSubmissionsThresholdDuration: revertExecutingSubmissionsThresholdDuration,
	}, nil
}

func (s submission) dbSubmissionToRPCSubmission(
	submission *db.Submission,
	problem *db.Problem,
	author *db.Account,
) rpc.Submission {
	return rpc.Submission{
		ID: uint64(submission.ID),
		Problem: rpc.SubmissionProblemSnippet{
			UUID:        problem.UUID,
			DisplayName: problem.DisplayName,
		},
		Author: rpc.Account{
			ID:          uint64(author.ID),
			AccountName: author.AccountName,
			DisplayName: author.DisplayName,
			Role:        string(author.Role),
		},
		Language:    submission.Language,
		Content:     submission.Content,
		Status:      uint8(submission.Status),
		Result:      uint8(submission.Result),
		CreatedTime: uint64(submission.CreatedAt.UnixMilli()),
	}
}

func (s submission) dbSubmissionToRPCSubmissionSnippet(
	submission *db.Submission,
	problem *db.Problem,
	author *db.Account,
) rpc.SubmissionSnippet {
	return rpc.SubmissionSnippet{
		ID: uint64(submission.ID),
		Problem: rpc.SubmissionProblemSnippet{
			UUID:        problem.UUID,
			DisplayName: problem.DisplayName,
		},
		Author: rpc.Account{
			ID:          uint64(author.ID),
			AccountName: author.AccountName,
			DisplayName: author.DisplayName,
			Role:        string(author.Role),
		},
		Language:    submission.Language,
		Status:      uint8(submission.Status),
		Result:      uint8(submission.Result),
		CreatedTime: uint64(submission.CreatedAt.UnixMilli()),
	}
}

func (s submission) CreateSubmission(
	ctx context.Context,
	in *rpc.CreateSubmissionRequest,
	token string,
) (*rpc.CreateSubmissionResponse, error) {
	logger := utils.LoggerWithContext(ctx, s.logger)

	account, err := s.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	hasPermission, err := s.role.AccountHasPermission(
		ctx,
		string(account.Role),
		PermissionSubmissionsSelfWrite,
		PermissionSubmissionsAllWrite,
	)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	problem, err := s.problemDataAccessor.GetProblemByUUID(ctx, in.ProblemUUID)
	if err != nil {
		return nil, err
	}

	if problem == nil {
		logger.With(zap.String("problem_uuid", in.ProblemUUID)).Error("problem not found")
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
	}

	submission := &db.Submission{
		OfProblemID:     uint64(problem.ID),
		AuthorAccountID: uint64(account.ID),
		Content:         in.Content,
		Language:        in.Language,
		Status:          db.SubmissionStatusSubmitted,
	}
	err = s.submissionDataAccessor.CreateSubmission(ctx, submission)
	if err != nil {
		return nil, err
	}

	if !s.appArguments.Distributed {
		s.judge.ScheduleJudgeLocalSubmission(uint64(submission.ID))
	}

	return &rpc.CreateSubmissionResponse{
		SubmissionSnippet: s.dbSubmissionToRPCSubmissionSnippet(submission, problem, account),
	}, nil
}

func (s submission) isValidSubmissionStatusTransition(
	oldStatus db.SubmissionStatus,
	newStatus rpc.SubmissionStatus,
) bool {
	if newStatus == rpc.SubmissionStatusSubmitted {
		return oldStatus == db.SubmissionStatus(rpc.SubmissionStatusExecuting)
	}

	if newStatus == rpc.SubmissionStatusExecuting {
		return oldStatus == db.SubmissionStatus(rpc.SubmissionStatusSubmitted)
	}

	if newStatus == rpc.SubmissionStatusFinished {
		return oldStatus == db.SubmissionStatus(rpc.SubmissionStatusExecuting)
	}

	return false
}

func (s submission) applyUpdateSubmission(
	ctx context.Context,
	in *rpc.UpdateSubmissionRequest,
	submission *db.Submission,
) error {
	logger := utils.LoggerWithContext(ctx, s.logger)

	if in.Status != 0 {
		if !s.isValidSubmissionStatusTransition(submission.Status, rpc.SubmissionStatus(in.Status)) {
			logger.
				With(zap.Uint8("old_status", uint8(submission.Status))).
				With(zap.Uint8("new_status", in.Status)).
				Error("invalid submission status transition")
			return pjrpc.JRPCErrInvalidParams()
		}

		submission.Status = db.SubmissionStatus(in.Status)

		if in.Status == uint8(rpc.SubmissionStatusFinished) {
			if in.Result == 0 {
				return pjrpc.JRPCErrInvalidParams()
			}

			submission.Result = db.SubmissionResult(in.Result)
		}
	}

	return nil
}

func (s submission) UpdateSubmission(
	ctx context.Context,
	in *rpc.UpdateSubmissionRequest,
	token string,
) (*rpc.UpdateSubmissionResponse, error) {
	account, err := s.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	hasPermission, err := s.role.AccountHasPermission(ctx, string(account.Role), PermissionSubmissionsAllWrite)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	response := &rpc.UpdateSubmissionResponse{}
	if txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		submission, submissionErr := s.submissionDataAccessor.WithDB(tx).GetSubmission(ctx, in.ID)
		if submissionErr != nil {
			return submissionErr
		}

		if submission == nil {
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
		}

		submissionErr = utils.ExecuteUntilFirstError(
			func() error { return s.applyUpdateSubmission(ctx, in, submission) },
			func() error { return s.submissionDataAccessor.WithDB(tx).UpdateSubmission(ctx, submission) },
		)
		if submissionErr != nil {
			return submissionErr
		}

		problem, submissionErr := s.problemDataAccessor.WithDB(tx).GetProblem(ctx, submission.OfProblemID)
		if err != nil {
			return submissionErr
		}

		author, submissionErr := s.accountDataAccessor.WithDB(tx).GetAccount(ctx, submission.AuthorAccountID)
		if err != nil {
			return submissionErr
		}

		response.SubmissionSnippet = s.dbSubmissionToRPCSubmissionSnippet(submission, problem, author)

		return nil
	}); txErr != nil {
		return nil, txErr
	}

	return response, nil
}

func (s submission) DeleteSubmission(ctx context.Context, in *rpc.DeleteSubmissionRequest, token string) error {
	account, err := s.token.GetAccount(ctx, token)
	if err != nil {
		return err
	}

	submission, err := s.submissionDataAccessor.GetSubmission(ctx, in.ID)
	if err != nil {
		return err
	}

	if submission == nil {
		return pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
	}

	requiredPermissionList := []gorbac.Permission{PermissionSubmissionsAllWrite}
	if submission.AuthorAccountID == uint64(account.ID) {
		requiredPermissionList = append(requiredPermissionList, PermissionSubmissionsSelfWrite)
	}

	hasPermission, err := s.role.AccountHasPermission(ctx, string(account.Role), requiredPermissionList...)
	if err != nil {
		return err
	}
	if !hasPermission {
		return pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	return s.submissionDataAccessor.DeleteSubmission(ctx, in.ID)
}

func (s submission) GetAccountSubmissionSnippetList(
	ctx context.Context,
	in *rpc.GetAccountSubmissionSnippetListRequest,
	token string,
) (*rpc.GetAccountSubmissionSnippetListResponse, error) {
	logger := utils.LoggerWithContext(ctx, s.logger)

	account, err := s.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	requiredPermissionList := []gorbac.Permission{PermissionSubmissionsAllRead}
	if account.ID == uint(in.AccountID) {
		requiredPermissionList = append(requiredPermissionList, PermissionSubmissionsSelfRead)
	}

	hasPermission, err := s.role.AccountHasPermission(ctx, string(account.Role), requiredPermissionList...)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	author, err := s.accountDataAccessor.GetAccount(ctx, in.AccountID)
	if err != nil {
		return nil, err
	}

	if author == nil {
		logger.With(zap.Uint64("account_id", in.AccountID)).Error("account not found")
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
	}

	totalSubmissionCount, err := s.submissionDataAccessor.GetSubmissionCount(ctx, db.SubmissionListFilterParams{
		AuthorAccountID: &in.AccountID,
	})
	if err != nil {
		return nil, err
	}

	submissionList, err := s.submissionDataAccessor.GetSubmissionList(ctx, db.SubmissionListFilterParams{
		AuthorAccountID: &in.AccountID,
	}, in.Offset, in.Limit)
	if err != nil {
		return nil, err
	}

	submissionSnippetList := make([]rpc.SubmissionSnippet, 0, len(submissionList))
	for i := range submissionList {
		problem, getProblemErr := s.problemDataAccessor.GetProblem(ctx, submissionList[i].OfProblemID)
		if getProblemErr != nil {
			return nil, getProblemErr
		}

		submissionSnippetList = append(
			submissionSnippetList,
			s.dbSubmissionToRPCSubmissionSnippet(submissionList[i], problem, author),
		)
	}

	return &rpc.GetAccountSubmissionSnippetListResponse{
		TotalSubmissionCount:  totalSubmissionCount,
		SubmissionSnippetList: submissionSnippetList,
	}, nil
}

func (s submission) GetProblemSubmissionSnippetList(
	ctx context.Context,
	in *rpc.GetProblemSubmissionSnippetListRequest,
	token string,
) (*rpc.GetProblemSubmissionSnippetListResponse, error) {
	logger := utils.LoggerWithContext(ctx, s.logger)

	account, err := s.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	problem, err := s.problemDataAccessor.GetProblemByUUID(ctx, in.ProblemUUID)
	if err != nil {
		return nil, err
	}

	if problem == nil {
		logger.With(zap.String("problem_uuid", in.ProblemUUID)).Error("problem not found")
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
	}

	requiredPermissionList := []gorbac.Permission{PermissionSubmissionsAllRead}
	if problem.AuthorAccountID == uint64(account.ID) {
		requiredPermissionList = append(requiredPermissionList, PermissionSubmissionsSelfRead)
	}

	hasPermission, err := s.role.AccountHasPermission(ctx, string(account.Role), requiredPermissionList...)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	totalSubmissionCount, err := s.submissionDataAccessor.GetSubmissionCount(ctx, db.SubmissionListFilterParams{
		OfProblemID: proto.Uint64(uint64(problem.ID)),
	})
	if err != nil {
		return nil, err
	}

	submissionList, err := s.submissionDataAccessor.GetSubmissionList(ctx, db.SubmissionListFilterParams{
		OfProblemID: proto.Uint64(uint64(problem.ID)),
	}, in.Offset, in.Limit)
	if err != nil {
		return nil, err
	}

	submissionSnippetList := make([]rpc.SubmissionSnippet, 0, len(submissionList))
	for i := range submissionList {
		author, authorErr := s.accountDataAccessor.GetAccount(ctx, submissionList[i].AuthorAccountID)
		if authorErr != nil {
			return nil, authorErr
		}

		submissionSnippetList = append(
			submissionSnippetList,
			s.dbSubmissionToRPCSubmissionSnippet(submissionList[i], problem, author),
		)
	}

	return &rpc.GetProblemSubmissionSnippetListResponse{
		TotalSubmissionCount:  totalSubmissionCount,
		SubmissionSnippetList: submissionSnippetList,
	}, nil
}

func (s submission) GetAccountProblemSubmissionSnippetList(
	ctx context.Context,
	in *rpc.GetAccountProblemSubmissionSnippetListRequest,
	token string,
) (*rpc.GetAccountProblemSubmissionSnippetListResponse, error) {
	logger := utils.LoggerWithContext(ctx, s.logger)

	account, err := s.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	author, err := s.accountDataAccessor.GetAccount(ctx, in.AccountID)
	if err != nil {
		return nil, err
	}

	if author == nil {
		logger.With(zap.Uint64("account_id", in.AccountID)).Error("account not found")
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
	}

	problem, err := s.problemDataAccessor.GetProblemByUUID(ctx, in.ProblemUUID)
	if err != nil {
		return nil, err
	}

	if problem == nil {
		logger.With(zap.String("problem_uuid", in.ProblemUUID)).Error("problem not found")
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
	}

	requiredPermissionList := []gorbac.Permission{PermissionSubmissionsAllRead}
	if uint64(author.ID) == uint64(account.ID) && problem.AuthorAccountID == uint64(account.ID) {
		requiredPermissionList = append(requiredPermissionList, PermissionSubmissionsSelfRead)
	}

	hasPermission, err := s.role.AccountHasPermission(ctx, string(account.Role), requiredPermissionList...)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	totalSubmissionCount, err := s.submissionDataAccessor.GetSubmissionCount(ctx, db.SubmissionListFilterParams{
		OfProblemID:     proto.Uint64(uint64(problem.ID)),
		AuthorAccountID: &in.AccountID,
	})
	if err != nil {
		return nil, err
	}

	submissionList, err := s.submissionDataAccessor.GetSubmissionList(ctx, db.SubmissionListFilterParams{
		OfProblemID:     proto.Uint64(uint64(problem.ID)),
		AuthorAccountID: &in.AccountID,
	}, in.Offset, in.Limit)
	if err != nil {
		return nil, err
	}

	submissionSnippetList := make([]rpc.SubmissionSnippet, 0, len(submissionList))
	for i := range submissionList {
		submissionSnippetList = append(
			submissionSnippetList,
			s.dbSubmissionToRPCSubmissionSnippet(submissionList[i], problem, author),
		)
	}

	return &rpc.GetAccountProblemSubmissionSnippetListResponse{
		TotalSubmissionCount:  totalSubmissionCount,
		SubmissionSnippetList: submissionSnippetList,
	}, nil
}

func (s submission) GetSubmission(
	ctx context.Context,
	in *rpc.GetSubmissionRequest,
	token string,
) (*rpc.GetSubmissionResponse, error) {
	account, err := s.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	submission, err := s.submissionDataAccessor.GetSubmission(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	if submission == nil {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
	}

	requiredPermissionList := []gorbac.Permission{PermissionSubmissionsAllRead}
	if submission.AuthorAccountID == uint64(account.ID) {
		requiredPermissionList = append(requiredPermissionList, PermissionSubmissionsSelfRead)
	}

	hasPermission, err := s.role.AccountHasPermission(ctx, string(account.Role), requiredPermissionList...)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	problem, err := s.problemDataAccessor.GetProblem(ctx, submission.OfProblemID)
	if err != nil {
		return nil, err
	}

	author, err := s.accountDataAccessor.GetAccount(ctx, submission.AuthorAccountID)
	if err != nil {
		return nil, err
	}

	return &rpc.GetSubmissionResponse{
		Submission: s.dbSubmissionToRPCSubmission(submission, problem, author),
	}, nil
}

func (s submission) GetSubmissionSnippetList(
	ctx context.Context,
	in *rpc.GetSubmissionSnippetListRequest,
	token string,
) (*rpc.GetSubmissionSnippetListResponse, error) {
	account, err := s.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	hasPermission, err := s.role.AccountHasPermission(ctx, string(account.Role), PermissionSubmissionsAllRead)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	totalSubmissionCount, err := s.submissionDataAccessor.GetSubmissionCount(ctx, db.SubmissionListFilterParams{})
	if err != nil {
		return nil, err
	}

	submissionList, err := s.submissionDataAccessor.
		GetSubmissionList(ctx, db.SubmissionListFilterParams{}, in.Offset, in.Limit)
	if err != nil {
		return nil, err
	}

	submissionSnippetList := make([]rpc.SubmissionSnippet, 0, len(submissionList))
	for i := range submissionList {
		problem, problemErr := s.problemDataAccessor.GetProblem(ctx, submissionList[i].OfProblemID)
		if problemErr != nil {
			return nil, problemErr
		}

		author, authorErr := s.accountDataAccessor.GetAccount(ctx, submissionList[i].AuthorAccountID)
		if authorErr != nil {
			return nil, authorErr
		}

		submissionSnippetList = append(
			submissionSnippetList,
			s.dbSubmissionToRPCSubmissionSnippet(submissionList[i], problem, author),
		)
	}

	return &rpc.GetSubmissionSnippetListResponse{
		TotalSubmissionCount:  totalSubmissionCount,
		SubmissionSnippetList: submissionSnippetList,
	}, nil
}

func (s submission) GetAndUpdateFirstSubmittedSubmissionToExecuting(
	ctx context.Context,
	token string,
) (*rpc.GetAndUpdateFirstSubmittedSubmissionToExecutingResponse, error) {
	account, err := s.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	hasPermission, err := s.role.AccountHasPermission(ctx, string(account.Role), PermissionSubmissionsAllWrite)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	response := &rpc.GetAndUpdateFirstSubmittedSubmissionToExecutingResponse{}
	if txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		submissionList, submissionErr := s.submissionDataAccessor.
			WithDB(tx).GetSubmissionList(ctx, db.SubmissionListFilterParams{Status: db.SubmissionStatusSubmitted}, 0, 1)
		if submissionErr != nil {
			return submissionErr
		}

		if len(submissionList) == 0 {
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
		}

		submission := submissionList[0]
		submission.Status = db.SubmissionStatusExecuting

		submissionErr = s.submissionDataAccessor.WithDB(tx).UpdateSubmission(ctx, submission)
		if submissionErr != nil {
			return submissionErr
		}

		problem, submissionErr := s.problemDataAccessor.WithDB(tx).GetProblem(ctx, submission.OfProblemID)
		if err != nil {
			return submissionErr
		}

		author, submissionErr := s.accountDataAccessor.WithDB(tx).GetAccount(ctx, submission.AuthorAccountID)
		if err != nil {
			return submissionErr
		}

		response.Submission = s.dbSubmissionToRPCSubmission(submission, problem, author)

		return nil
	}); txErr != nil {
		return nil, txErr
	}

	return response, nil
}

func (s submission) ScheduleSubmittedExecutingSubmissionToJudge(ctx context.Context) error {
	submittedSubmissionIDList, err := s.submissionDataAccessor.GetSubmissionIDList(ctx, db.SubmissionListFilterParams{
		Status: db.SubmissionStatusSubmitted,
	})
	if err != nil {
		return err
	}

	for _, id := range submittedSubmissionIDList {
		s.judge.ScheduleJudgeLocalSubmission(id)
	}

	executingSubmissionIDList, err := s.submissionDataAccessor.GetSubmissionIDList(ctx, db.SubmissionListFilterParams{
		Status: db.SubmissionStatusExecuting,
	})
	if err != nil {
		return err
	}

	for _, id := range executingSubmissionIDList {
		submission, submissionErr := s.submissionDataAccessor.GetSubmission(ctx, id)
		if submissionErr != nil {
			return submissionErr
		}

		submission.Status = db.SubmissionStatusSubmitted
		submissionErr = s.submissionDataAccessor.UpdateSubmission(ctx, submission)
		if submissionErr != nil {
			return submissionErr
		}

		s.judge.ScheduleJudgeLocalSubmission(id)
	}

	return nil
}

func (s submission) UpdateExecutingSubmissionUpdatedBeforeThresholdToSubmitted(ctx context.Context) error {
	logger := utils.LoggerWithContext(ctx, s.logger)

	threshold := time.Now().Add(-s.revertExecutingSubmissionsThresholdDuration)
	logger.
		With(zap.Time("threshold", threshold)).
		Info("reverting executing submissions with update time before threshold to submitted")

	return s.submissionDataAccessor.UpdateExecutingSubmissionUpdatedBeforeThresholdToSubmitted(ctx, threshold)
}
