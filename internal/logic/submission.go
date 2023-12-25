package logic

import (
	"context"

	"gitlab.com/pjrpc/pjrpc/v2"
	"go.uber.org/zap"

	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc"
	"github.com/tranHieuDev23/cato/internal/utils"
)

type Submission interface {
	CreateSubmission(ctx context.Context, in *rpc.CreateSubmissionRequest, token string) (*rpc.CreateSubmissionResponse, error)
	GetSubmissionSnippetList(ctx context.Context, in *rpc.GetSubmissionSnippetListRequest, token string) (*rpc.GetSubmissionSnippetListResponse, error)
	GetSubmission(ctx context.Context, in *rpc.GetSubmissionRequest, token string) (*rpc.GetSubmissionResponse, error)
	DeleteSubmission(ctx context.Context, in *rpc.DeleteSubmissionRequest, token string) error
	GetAccountSubmissionSnippetList(ctx context.Context, in *rpc.GetAccountSubmissionSnippetListRequest, token string) (*rpc.GetAccountSubmissionSnippetListResponse, error)
	GetProblemSubmissionSnippetList(ctx context.Context, in *rpc.GetProblemSubmissionSnippetListRequest, token string) (*rpc.GetProblemSubmissionSnippetListResponse, error)
	GetAccountProblemSubmissionSnippetList(ctx context.Context, in *rpc.GetAccountProblemSubmissionSnippetListRequest, token string) (*rpc.GetAccountProblemSubmissionSnippetListResponse, error)
	ScheduleSubmittedExecutingSubmissionToJudge(ctx context.Context) error
}

type submission struct {
	token                                  Token
	role                                   Role
	judge                                  Judge
	accountDataAccessor                    db.AccountDataAccessor
	problemDataAccessor                    db.ProblemDataAccessor
	submissionDataAccessor                 db.SubmissionDataAccessor
	logger                                 *zap.Logger
	shouldScheduleCreatedSubmissionToJudge bool
}

func NewSubmission(
	token Token,
	role Role,
	judge Judge,
	accountDataAccessor db.AccountDataAccessor,
	problemDataAccessor db.ProblemDataAccessor,
	submissionDataAccessor db.SubmissionDataAccessor,
	logger *zap.Logger,
	shouldScheduleCreatedSubmissionToJudge bool,
) Submission {
	return &submission{
		token:                                  token,
		role:                                   role,
		judge:                                  judge,
		accountDataAccessor:                    accountDataAccessor,
		problemDataAccessor:                    problemDataAccessor,
		submissionDataAccessor:                 submissionDataAccessor,
		logger:                                 logger,
		shouldScheduleCreatedSubmissionToJudge: shouldScheduleCreatedSubmissionToJudge,
	}
}

func (s submission) dbSubmissionToRPCSubmission(
	submission *db.Submission,
	problem *db.Problem,
	author *db.Account,
) rpc.Submission {
	return rpc.Submission{
		ID: uint64(submission.ID),
		Problem: rpc.SubmissionProblemSnippet{
			ID:          submission.OfProblemID,
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
			ID:          submission.OfProblemID,
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

func (s submission) CreateSubmission(ctx context.Context, in *rpc.CreateSubmissionRequest, token string) (*rpc.CreateSubmissionResponse, error) {
	logger := utils.LoggerWithContext(ctx, s.logger)

	account, err := s.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	if hasPermission, err := s.role.AccountHasPermission(
		ctx,
		string(account.Role),
		PermissionSubmissionsSelfWrite,
		PermissionSubmissionsAllWrite,
	); err != nil {
		return nil, err
	} else if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	problem, err := s.problemDataAccessor.GetProblem(ctx, in.ProblemID)
	if err != nil {
		return nil, err
	}

	if problem == nil {
		logger.With(zap.Uint64("problem_id", in.ProblemID)).Error("problem not found")
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
	}

	submission := &db.Submission{
		OfProblemID:     in.ProblemID,
		AuthorAccountID: uint64(account.ID),
		Content:         in.Content,
		Language:        in.Language,
		Status:          db.SubmissionStatusSubmitted,
	}
	if err := s.submissionDataAccessor.CreateSubmission(ctx, submission); err != nil {
		return nil, err
	}

	if s.shouldScheduleCreatedSubmissionToJudge {
		s.judge.ScheduleSubmissionToJudge(uint64(submission.ID))
	}

	return &rpc.CreateSubmissionResponse{
		SubmissionSnippet: s.dbSubmissionToRPCSubmissionSnippet(submission, problem, account),
	}, nil
}

func (s submission) DeleteSubmission(ctx context.Context, in *rpc.DeleteSubmissionRequest, token string) error {
	account, err := s.token.GetAccount(ctx, token)
	if err != nil {
		return err
	}

	if hasPermission, err := s.role.AccountHasPermission(
		ctx,
		string(account.Role),
		PermissionSubmissionsSelfWrite,
		PermissionSubmissionsAllWrite,
	); err != nil {
		return err
	} else if !hasPermission {
		return pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	submission, err := s.submissionDataAccessor.GetSubmission(ctx, in.ID)
	if err != nil {
		return err
	}

	if submission == nil {
		return pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
	}

	if submission.AuthorAccountID != uint64(account.ID) {
		if hasPermission, err := s.role.AccountHasPermission(
			ctx,
			string(account.Role),
			PermissionSubmissionsAllWrite,
		); err != nil {
			return err
		} else if !hasPermission {
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}
	}

	return s.submissionDataAccessor.DeleteSubmission(ctx, in.ID)
}

func (s submission) GetAccountSubmissionSnippetList(ctx context.Context, in *rpc.GetAccountSubmissionSnippetListRequest, token string) (*rpc.GetAccountSubmissionSnippetListResponse, error) {
	logger := utils.LoggerWithContext(ctx, s.logger)

	account, err := s.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	if account.ID == uint(in.AccountID) {
		if hasPermission, err := s.role.AccountHasPermission(
			ctx,
			string(account.Role),
			PermissionSubmissionsSelfRead,
		); err != nil {
			return nil, err
		} else if !hasPermission {
			return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}
	} else {
		if hasPermission, err := s.role.AccountHasPermission(
			ctx,
			string(account.Role),
			PermissionSubmissionsAllRead,
		); err != nil {
			return nil, err
		} else if !hasPermission {
			return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}
	}

	author, err := s.accountDataAccessor.GetAccount(ctx, in.AccountID)
	if err != nil {
		return nil, err
	}

	if author == nil {
		logger.With(zap.Uint64("account_id", in.AccountID)).Error("account not found")
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
	}

	if uint64(author.ID) != uint64(account.ID) {
		if hasPermission, err := s.role.AccountHasPermission(
			ctx,
			string(account.Role),
			PermissionSubmissionsAllRead,
		); err != nil {
			return nil, err
		} else if !hasPermission {
			return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}
	}

	totalSubmissionCount, err := s.submissionDataAccessor.GetAccountSubmissionCount(ctx, in.AccountID)
	if err != nil {
		return nil, err
	}

	submissionList, err := s.submissionDataAccessor.GetAccountSubmissionList(ctx, in.AccountID, in.Offset, in.Limit)
	if err != nil {
		return nil, err
	}

	submissionSnippetList := make([]rpc.SubmissionSnippet, 0, len(submissionList))
	for i := range submissionList {
		problem, err := s.problemDataAccessor.GetProblem(ctx, submissionList[i].OfProblemID)
		if err != nil {
			return nil, err
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

func (s submission) GetProblemSubmissionSnippetList(ctx context.Context, in *rpc.GetProblemSubmissionSnippetListRequest, token string) (*rpc.GetProblemSubmissionSnippetListResponse, error) {
	logger := utils.LoggerWithContext(ctx, s.logger)

	account, err := s.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	if hasPermission, err := s.role.AccountHasPermission(
		ctx,
		string(account.Role),
		PermissionSubmissionsSelfRead,
		PermissionSubmissionsAllRead,
	); err != nil {
		return nil, err
	} else if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	problem, err := s.problemDataAccessor.GetProblem(ctx, in.ProblemID)
	if err != nil {
		return nil, err
	}

	if problem == nil {
		logger.With(zap.Uint64("problem_id", in.ProblemID)).Error("problem not found")
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
	}

	if problem.AuthorAccountID == uint64(account.ID) {
		if hasPermission, err := s.role.AccountHasPermission(
			ctx,
			string(account.Role),
			PermissionSubmissionsAllRead,
		); err != nil {
			return nil, err
		} else if !hasPermission {
			return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}
	}

	totalSubmissionCount, err := s.submissionDataAccessor.GetProblemSubmissionCount(ctx, in.ProblemID)
	if err != nil {
		return nil, err
	}

	submissionList, err := s.submissionDataAccessor.GetProblemSubmissionList(ctx, in.ProblemID, in.Offset, in.Limit)
	if err != nil {
		return nil, err
	}

	submissionSnippetList := make([]rpc.SubmissionSnippet, 0, len(submissionList))
	for i := range submissionList {
		author, err := s.accountDataAccessor.GetAccount(ctx, submissionList[i].AuthorAccountID)
		if err != nil {
			return nil, err
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

	if hasPermission, err := s.role.AccountHasPermission(
		ctx,
		string(account.Role),
		PermissionSubmissionsSelfRead,
		PermissionSubmissionsAllRead,
	); err != nil {
		return nil, err
	} else if !hasPermission {
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

	if uint64(author.ID) != uint64(account.ID) {
		if hasPermission, err := s.role.AccountHasPermission(
			ctx,
			string(account.Role),
			PermissionSubmissionsAllRead,
		); err != nil {
			return nil, err
		} else if !hasPermission {
			return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}
	}

	problem, err := s.problemDataAccessor.GetProblem(ctx, in.ProblemID)
	if err != nil {
		return nil, err
	}

	if problem == nil {
		logger.With(zap.Uint64("problem_id", in.ProblemID)).Error("problem not found")
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
	}

	if problem.AuthorAccountID == uint64(account.ID) {
		if hasPermission, err := s.role.AccountHasPermission(
			ctx,
			string(account.Role),
			PermissionSubmissionsAllRead,
		); err != nil {
			return nil, err
		} else if !hasPermission {
			return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}
	}

	totalSubmissionCount, err := s.submissionDataAccessor.GetAccountProblemSubmissionCount(ctx, in.AccountID, in.ProblemID)
	if err != nil {
		return nil, err
	}

	submissionList, err := s.submissionDataAccessor.GetAccountProblemSubmissionList(ctx, in.AccountID, in.ProblemID, in.Offset, in.Limit)
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

func (s submission) GetSubmission(ctx context.Context, in *rpc.GetSubmissionRequest, token string) (*rpc.GetSubmissionResponse, error) {
	account, err := s.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	if hasPermission, err := s.role.AccountHasPermission(
		ctx,
		string(account.Role),
		PermissionSubmissionsSelfRead,
		PermissionSubmissionsAllRead,
	); err != nil {
		return nil, err
	} else if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	submission, err := s.submissionDataAccessor.GetSubmission(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	if submission == nil {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
	}

	if submission.AuthorAccountID != uint64(account.ID) {
		if hasPermission, err := s.role.AccountHasPermission(
			ctx,
			string(account.Role),
			PermissionSubmissionsAllRead,
		); err != nil {
			return nil, err
		} else if !hasPermission {
			return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}
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

func (s submission) GetSubmissionSnippetList(ctx context.Context, in *rpc.GetSubmissionSnippetListRequest, token string) (*rpc.GetSubmissionSnippetListResponse, error) {
	account, err := s.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	if hasPermission, err := s.role.AccountHasPermission(
		ctx,
		string(account.Role),
		PermissionSubmissionsAllRead,
	); err != nil {
		return nil, err
	} else if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	totalSubmissionCount, err := s.submissionDataAccessor.GetSubmissionCount(ctx)
	if err != nil {
		return nil, err
	}

	submissionList, err := s.submissionDataAccessor.GetSubmissionList(ctx, in.Offset, in.Limit)
	if err != nil {
		return nil, err
	}

	submissionSnippetList := make([]rpc.SubmissionSnippet, 0, len(submissionList))
	for i := range submissionList {
		problem, err := s.problemDataAccessor.GetProblem(ctx, submissionList[i].OfProblemID)
		if err != nil {
			return nil, err
		}

		author, err := s.accountDataAccessor.GetAccount(ctx, submissionList[i].AuthorAccountID)
		if err != nil {
			return nil, err
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

func (s submission) ScheduleSubmittedExecutingSubmissionToJudge(ctx context.Context) error {
	submittedSubmissionIDList, err := s.submissionDataAccessor.GetSubmissionIDListWithStatus(ctx, db.SubmissionStatusSubmitted)
	if err != nil {
		return err
	}

	for _, id := range submittedSubmissionIDList {
		s.judge.ScheduleSubmissionToJudge(id)
	}

	executingSubmissionIDList, err := s.submissionDataAccessor.GetSubmissionIDListWithStatus(ctx, db.SubmissionStatusExecuting)
	if err != nil {
		return err
	}

	for _, id := range executingSubmissionIDList {
		submission, err := s.submissionDataAccessor.GetSubmission(ctx, id)
		if err != nil {
			return err
		}

		submission.Status = db.SubmissionStatusSubmitted
		if err := s.submissionDataAccessor.UpdateSubmission(ctx, submission); err != nil {
			return err
		}

		s.judge.ScheduleSubmissionToJudge(id)
	}

	return nil
}

type LocalSubmission Submission

func NewLocalSubmission(
	token Token,
	role Role,
	judge LocalJudge,
	accountDataAccessor db.AccountDataAccessor,
	problemDataAccessor db.ProblemDataAccessor,
	submissionDataAccessor db.SubmissionDataAccessor,
	logger *zap.Logger,
) LocalSubmission {
	return NewSubmission(
		token, role, judge, accountDataAccessor, problemDataAccessor, submissionDataAccessor, logger, true)
}

type DistributedSubmission Submission

func NewDistributedSubmission(
	token Token,
	role Role,
	judge DistributedJudge,
	accountDataAccessor db.AccountDataAccessor,
	problemDataAccessor db.ProblemDataAccessor,
	submissionDataAccessor db.SubmissionDataAccessor,
	logger *zap.Logger,
) DistributedSubmission {
	return NewSubmission(
		token, role, judge, accountDataAccessor, problemDataAccessor, submissionDataAccessor, logger, false)
}
