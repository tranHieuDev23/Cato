package logic

import (
	"context"
	"strings"

	"gitlab.com/pjrpc/pjrpc/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/mikespook/gorbac"
	"github.com/samber/lo"

	"github.com/microcosm-cc/bluemonday"

	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc"
	"github.com/tranHieuDev23/cato/internal/utils"
)

type Problem interface {
	CreateProblem(ctx context.Context, in *rpc.CreateProblemRequest, token string) (*rpc.CreateProblemResponse, error)
	GetProblemSnippetList(
		ctx context.Context,
		in *rpc.GetProblemSnippetListRequest,
		token string,
	) (*rpc.GetProblemSnippetListResponse, error)
	GetProblem(ctx context.Context, in *rpc.GetProblemRequest, token string) (*rpc.GetProblemResponse, error)
	UpdateProblem(ctx context.Context, in *rpc.UpdateProblemRequest, token string) (*rpc.UpdateProblemResponse, error)
	DeleteProblem(ctx context.Context, in *rpc.DeleteProblemRequest, token string) error
	GetAccountProblemSnippetList(
		ctx context.Context,
		in *rpc.GetAccountProblemSnippetListRequest,
		token string,
	) (*rpc.GetAccountProblemSnippetListResponse, error)
	WithDB(db *gorm.DB) Problem
}

type problem struct {
	token                           Token
	role                            Role
	testCase                        TestCase
	accountDataAccessor             db.AccountDataAccessor
	problemDataAccessor             db.ProblemDataAccessor
	problemExampleDataAccessor      db.ProblemExampleDataAccessor
	problemTestCaseHashDataAccessor db.ProblemTestCaseHashDataAccessor
	testCaseDataAccessor            db.TestCaseDataAccessor
	submissionDataAccessor          db.SubmissionDataAccessor
	logger                          *zap.Logger
	db                              *gorm.DB
	displayNameSanitizePolicy       *bluemonday.Policy
	descriptionSanitizePolicy       *bluemonday.Policy
}

func NewProblem(
	token Token,
	role Role,
	testCase TestCase,
	accountDataAccessor db.AccountDataAccessor,
	problemDataAccessor db.ProblemDataAccessor,
	problemExampleDataAccessor db.ProblemExampleDataAccessor,
	problemTestCaseHashDataAccessor db.ProblemTestCaseHashDataAccessor,
	testCaseDataAccessor db.TestCaseDataAccessor,
	submissionDataAccessor db.SubmissionDataAccessor,
	logger *zap.Logger,
	db *gorm.DB,
) Problem {
	return &problem{
		token:                           token,
		role:                            role,
		testCase:                        testCase,
		accountDataAccessor:             accountDataAccessor,
		problemDataAccessor:             problemDataAccessor,
		problemExampleDataAccessor:      problemExampleDataAccessor,
		problemTestCaseHashDataAccessor: problemTestCaseHashDataAccessor,
		testCaseDataAccessor:            testCaseDataAccessor,
		submissionDataAccessor:          submissionDataAccessor,
		logger:                          logger,
		db:                              db,
		displayNameSanitizePolicy:       bluemonday.StrictPolicy(),
		descriptionSanitizePolicy:       bluemonday.UGCPolicy(),
	}
}

func (p problem) cleanupDisplayName(displayName string) string {
	displayName = strings.Trim(displayName, " ")
	displayName = p.displayNameSanitizePolicy.Sanitize(displayName)
	return displayName
}

func (p problem) isValidDisplayName(displayName string) bool {
	return displayName != ""
}

func (p problem) cleanupDescription(description string) string {
	description = strings.Trim(description, " ")
	description = p.descriptionSanitizePolicy.Sanitize(description)
	return description
}

func (p problem) dbProblemExampleToRPCProblemExample(problemExample *db.ProblemExample) rpc.ProblemExample {
	return rpc.ProblemExample{
		Input:  problemExample.Input,
		Output: problemExample.Output,
	}
}

func (p problem) dbProblemToRPCProblem(
	problem *db.Problem,
	author *db.Account,
	problemExampleList []*db.ProblemExample,
) rpc.Problem {
	return rpc.Problem{
		ID:          uint64(problem.ID),
		DisplayName: problem.DisplayName,
		Author: rpc.Account{
			ID:          uint64(author.ID),
			AccountName: author.AccountName,
			DisplayName: author.DisplayName,
		},
		Description:            problem.Description,
		TimeLimitInMillisecond: problem.TimeLimitInMillisecond,
		MemoryLimitInByte:      problem.MemoryLimitInByte,
		ExampleList: lo.Map[*db.ProblemExample, rpc.ProblemExample](
			problemExampleList,
			func(item *db.ProblemExample, _ int) rpc.ProblemExample {
				return p.dbProblemExampleToRPCProblemExample(item)
			},
		),
		CreatedTime: uint64(problem.CreatedAt.UnixMilli()),
		UpdatedTime: uint64(problem.UpdatedAt.UnixMilli()),
	}
}

func (p problem) dbProblemToRPCProblemSnippet(problem *db.Problem, author *db.Account) rpc.ProblemSnippet {
	return rpc.ProblemSnippet{
		ID:          uint64(problem.ID),
		DisplayName: problem.DisplayName,
		Author: rpc.Account{
			ID:          uint64(author.ID),
			AccountName: author.AccountName,
			DisplayName: author.DisplayName,
		},
		TimeLimitInMillisecond: problem.TimeLimitInMillisecond,
		MemoryLimitInByte:      problem.MemoryLimitInByte,
		CreatedTime:            uint64(problem.CreatedAt.UnixMilli()),
		UpdatedTime:            uint64(problem.UpdatedAt.UnixMilli()),
	}
}

func (p problem) CreateProblem(
	ctx context.Context,
	in *rpc.CreateProblemRequest,
	token string,
) (*rpc.CreateProblemResponse, error) {
	account, err := p.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	hasPermission, err := p.role.AccountHasPermission(
		ctx,
		string(account.Role),
		PermissionProblemsSelfWrite,
		PermissionProblemsAllWrite,
	)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	cleanedDisplayName := p.cleanupDisplayName(in.DisplayName)
	if !p.isValidDisplayName(cleanedDisplayName) {
		return nil, pjrpc.JRPCErrInvalidParams()
	}

	cleanedDescription := p.cleanupDescription(in.Description)

	response := &rpc.CreateProblemResponse{}
	if txErr := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		problem := &db.Problem{
			DisplayName:            cleanedDisplayName,
			AuthorAccountID:        uint64(account.ID),
			Description:            cleanedDescription,
			TimeLimitInMillisecond: in.TimeLimitInMillisecond,
			MemoryLimitInByte:      in.MemoryLimitInByte,
		}

		var problemExampleList []*db.ProblemExample

		err = utils.ExecuteUntilFirstError(
			func() error {
				return p.problemDataAccessor.WithDB(tx).CreateProblem(ctx, problem)
			},
			func() error {
				problemExampleList = lo.Map(in.ExampleList, func(item rpc.ProblemExample, _ int) *db.ProblemExample {
					return &db.ProblemExample{
						OfProblemID: uint64(problem.ID),
						Input:       utils.TrimSpaceRight(item.Input),
						Output:      utils.TrimSpaceRight(item.Output),
					}
				})
				return p.problemExampleDataAccessor.WithDB(tx).CreateProblemExampleList(ctx, problemExampleList)
			},
			func() error {
				return p.testCase.WithDB(tx).UpsertProblemTestCaseHash(ctx, uint64(problem.ID))
			},
		)
		if err != nil {
			return err
		}

		response.Problem = p.dbProblemToRPCProblem(problem, account, problemExampleList)

		return nil
	}); txErr != nil {
		return nil, txErr
	}

	return response, nil
}

func (p problem) DeleteProblem(ctx context.Context, in *rpc.DeleteProblemRequest, token string) error {
	logger := utils.LoggerWithContext(ctx, p.logger)

	account, err := p.token.GetAccount(ctx, token)
	if err != nil {
		return err
	}

	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		problem, problemErr := p.problemDataAccessor.WithDB(tx).GetProblem(ctx, in.ID)
		if problemErr != nil {
			return problemErr
		}

		if problem == nil {
			logger.With(zap.Uint64("id", in.ID)).Error("cannot find problem")
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
		}

		requiredPermissionList := []gorbac.Permission{PermissionProblemsAllWrite}
		if problem.AuthorAccountID == uint64(account.ID) {
			requiredPermissionList = append(requiredPermissionList, PermissionProblemsSelfWrite)
		}

		hasPermission, problemErr := p.role.AccountHasPermission(ctx, string(account.Role), requiredPermissionList...)
		if problemErr != nil {
			return problemErr
		}
		if !hasPermission {
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}

		return utils.ExecuteUntilFirstError(
			func() error {
				return p.problemExampleDataAccessor.WithDB(tx).DeleteProblemExampleOfProblem(ctx, uint64(problem.ID))
			},
			func() error {
				return p.testCaseDataAccessor.WithDB(tx).DeleteTestCaseOfProblem(ctx, uint64(problem.ID))
			},
			func() error {
				return p.submissionDataAccessor.WithDB(tx).DeleteSubmissionOfProblem(ctx, uint64(problem.ID))
			},
			func() error {
				return p.problemTestCaseHashDataAccessor.WithDB(tx).DeleteProblemTestCaseHashOfProblem(ctx, uint64(problem.ID))
			},
			func() error {
				return p.problemDataAccessor.WithDB(tx).DeleteProblem(ctx, uint64(problem.ID))
			},
		)
	})
}

func (p problem) GetAccountProblemSnippetList(
	ctx context.Context,
	in *rpc.GetAccountProblemSnippetListRequest,
	token string,
) (*rpc.GetAccountProblemSnippetListResponse, error) {
	account, err := p.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	requiredPermissionList := []gorbac.Permission{PermissionProblemsAllRead}
	if in.AccountID == uint64(account.ID) {
		requiredPermissionList = append(requiredPermissionList, PermissionProblemsSelfRead)
	}

	hasPermission, err := p.role.AccountHasPermission(ctx, string(account.Role), requiredPermissionList...)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	totalProblemCount, err := p.problemDataAccessor.GetAccountProblemCount(ctx, in.AccountID)
	if err != nil {
		return nil, err
	}

	problemList, err := p.problemDataAccessor.GetAccountProblemList(ctx, in.AccountID, in.Offset, in.Limit)
	if err != nil {
		return nil, err
	}

	problemSnippetList := make([]rpc.ProblemSnippet, 0, len(problemList))
	for i := range problemList {
		author, authorErr := p.accountDataAccessor.GetAccount(ctx, problemList[i].AuthorAccountID)
		if authorErr != nil {
			return nil, authorErr
		}

		if author == nil {
			return nil, pjrpc.JRPCErrInternalError()
		}

		problemSnippetList = append(problemSnippetList, p.dbProblemToRPCProblemSnippet(problemList[i], author))
	}

	return &rpc.GetAccountProblemSnippetListResponse{
		TotalProblemCount:  totalProblemCount,
		ProblemSnippetList: problemSnippetList,
	}, nil
}

func (p problem) GetProblem(
	ctx context.Context,
	in *rpc.GetProblemRequest,
	token string,
) (*rpc.GetProblemResponse, error) {
	logger := utils.LoggerWithContext(ctx, p.logger)

	account, err := p.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	problem, err := p.problemDataAccessor.GetProblem(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	if problem == nil {
		logger.With(zap.Uint64("id", in.ID)).Error("cannot find problem")
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
	}

	requiredPermissionList := []gorbac.Permission{PermissionProblemsAllRead}
	if problem.AuthorAccountID == uint64(account.ID) {
		requiredPermissionList = append(requiredPermissionList, PermissionProblemsSelfRead)
	}

	hasPermission, err := p.role.AccountHasPermission(ctx, string(account.Role), requiredPermissionList...)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	author, err := p.accountDataAccessor.GetAccount(ctx, problem.AuthorAccountID)
	if err != nil {
		return nil, err
	}

	if author == nil {
		return nil, pjrpc.JRPCErrInternalError()
	}

	problemExampleList, err := p.problemExampleDataAccessor.GetProblemExampleListOfProblem(ctx, uint64(problem.ID))
	if err != nil {
		return nil, err
	}

	return &rpc.GetProblemResponse{
		Problem: p.dbProblemToRPCProblem(problem, author, problemExampleList),
	}, nil
}

func (p problem) GetProblemSnippetList(
	ctx context.Context,
	in *rpc.GetProblemSnippetListRequest,
	token string,
) (*rpc.GetProblemSnippetListResponse, error) {
	account, err := p.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	hasPermission, err := p.role.AccountHasPermission(ctx, string(account.Role), PermissionProblemsAllRead)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	totalProblemCount, err := p.problemDataAccessor.GetProblemCount(ctx)
	if err != nil {
		return nil, err
	}

	problemList, err := p.problemDataAccessor.GetProblemList(ctx, in.Offset, in.Limit)
	if err != nil {
		return nil, err
	}

	problemSnippetList := make([]rpc.ProblemSnippet, 0, len(problemList))
	for i := range problemList {
		author, authorErr := p.accountDataAccessor.GetAccount(ctx, problemList[i].AuthorAccountID)
		if authorErr != nil {
			return nil, authorErr
		}

		if author == nil {
			return nil, pjrpc.JRPCErrInternalError()
		}

		problemSnippetList = append(problemSnippetList, p.dbProblemToRPCProblemSnippet(problemList[i], author))
	}

	return &rpc.GetProblemSnippetListResponse{
		TotalProblemCount:  totalProblemCount,
		ProblemSnippetList: problemSnippetList,
	}, nil
}

func (p problem) applyUpdateProblem(in *rpc.UpdateProblemRequest, problem *db.Problem) error {
	if in.DisplayName != nil {
		cleanedDisplayName := p.cleanupDisplayName(*in.DisplayName)
		if !p.isValidDisplayName(cleanedDisplayName) {
			return pjrpc.JRPCErrInvalidParams()
		}

		problem.DisplayName = cleanedDisplayName
	}

	if in.Description != nil {
		problem.Description = p.cleanupDescription(*in.Description)
	}

	if in.TimeLimitInMillisecond != nil {
		problem.TimeLimitInMillisecond = *in.TimeLimitInMillisecond
	}

	if in.MemoryLimitInByte != nil {
		problem.MemoryLimitInByte = *in.MemoryLimitInByte
	}

	return nil
}

func (p problem) updateProblemExampleList(
	ctx context.Context,
	in *rpc.UpdateProblemRequest,
	tx *gorm.DB,
) ([]*db.ProblemExample, error) {
	if in.ExampleList == nil {
		return p.problemExampleDataAccessor.WithDB(tx).GetProblemExampleListOfProblem(ctx, in.ID)
	}

	if err := p.problemExampleDataAccessor.DeleteProblemExampleOfProblem(ctx, in.ID); err != nil {
		return make([]*db.ProblemExample, 0), err
	}

	problemExampleList := lo.Map(*in.ExampleList, func(item rpc.ProblemExample, _ int) *db.ProblemExample {
		return &db.ProblemExample{OfProblemID: in.ID, Input: item.Input, Output: item.Output}
	})
	if err := p.problemExampleDataAccessor.WithDB(tx).CreateProblemExampleList(ctx, problemExampleList); err != nil {
		return make([]*db.ProblemExample, 0), err
	}

	return problemExampleList, nil
}

func (p problem) UpdateProblem(
	ctx context.Context,
	in *rpc.UpdateProblemRequest,
	token string,
) (*rpc.UpdateProblemResponse, error) {
	logger := utils.LoggerWithContext(ctx, p.logger)

	account, err := p.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	response := &rpc.UpdateProblemResponse{}
	if txErr := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		problem, problemErr := p.problemDataAccessor.WithDB(tx).GetProblem(ctx, in.ID)
		if problemErr != nil {
			return problemErr
		}

		if problem == nil {
			logger.With(zap.Uint64("id", in.ID)).Error("cannot find problem")
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
		}

		requiredPermissionList := []gorbac.Permission{PermissionProblemsAllWrite}
		if problem.AuthorAccountID == uint64(account.ID) {
			requiredPermissionList = append(requiredPermissionList, PermissionProblemsSelfWrite)
		}

		hasPermission, problemErr := p.role.AccountHasPermission(ctx, string(account.Role), requiredPermissionList...)
		if problemErr != nil {
			return problemErr
		}
		if !hasPermission {
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}

		problemErr = utils.ExecuteUntilFirstError(
			func() error { return p.applyUpdateProblem(in, problem) },
			func() error { return p.problemDataAccessor.WithDB(tx).UpdateProblem(ctx, problem) },
		)
		if problemErr != nil {
			return err
		}

		author, problemErr := p.accountDataAccessor.GetAccount(ctx, problem.AuthorAccountID)
		if problemErr != nil {
			return problemErr
		}

		problemExampleList, problemErr := p.updateProblemExampleList(ctx, in, tx)
		if problemErr != nil {
			return problemErr
		}

		response.Problem = p.dbProblemToRPCProblem(problem, author, problemExampleList)

		return nil
	}); txErr != nil {
		return nil, err
	}

	return response, nil
}

func (p problem) WithDB(db *gorm.DB) Problem {
	return &problem{
		token:                           p.token.WithDB(db),
		role:                            p.role,
		accountDataAccessor:             p.accountDataAccessor.WithDB(db),
		problemDataAccessor:             p.problemDataAccessor.WithDB(db),
		problemExampleDataAccessor:      p.problemExampleDataAccessor.WithDB(db),
		problemTestCaseHashDataAccessor: p.problemTestCaseHashDataAccessor.WithDB(db),
		testCase:                        p.testCase.WithDB(db),
		testCaseDataAccessor:            p.testCaseDataAccessor.WithDB(db),
		submissionDataAccessor:          p.submissionDataAccessor.WithDB(db),
		logger:                          p.logger,
		db:                              db,
		displayNameSanitizePolicy:       p.displayNameSanitizePolicy,
		descriptionSanitizePolicy:       p.descriptionSanitizePolicy,
	}
}
