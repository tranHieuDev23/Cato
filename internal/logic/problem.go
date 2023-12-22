package logic

import (
	"context"
	"strings"

	"gitlab.com/pjrpc/pjrpc/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/samber/lo"

	"github.com/microcosm-cc/bluemonday"

	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc"
	"github.com/tranHieuDev23/cato/internal/utils"
)

type Problem interface {
	CreateProblem(ctx context.Context, in *rpc.CreateProblemRequest, token string) (*rpc.CreateProblemResponse, error)
	GetProblemSnippetList(ctx context.Context, in *rpc.GetProblemSnippetListRequest, token string) (*rpc.GetProblemSnippetListResponse, error)
	GetProblem(ctx context.Context, in *rpc.GetProblemRequest, token string) (*rpc.GetProblemResponse, error)
	UpdateProblem(ctx context.Context, in *rpc.UpdateProblemRequest, token string) (*rpc.UpdateProblemResponse, error)
	DeleteProblem(ctx context.Context, in *rpc.DeleteProblemRequest, token string) error
	GetAccountProblemSnippetList(ctx context.Context, in *rpc.GetAccountProblemSnippetListRequest, token string) (*rpc.GetAccountProblemSnippetListResponse, error)
	WithDB(db *gorm.DB) Problem
}

type problem struct {
	token                      Token
	role                       Role
	accountDataAccessor        db.AccountDataAccessor
	problemDataAccessor        db.ProblemDataAccessor
	problemExampleDataAccessor db.ProblemExampleDataAccessor
	logger                     *zap.Logger
	db                         *gorm.DB
	displayNameSanitizePolicy  *bluemonday.Policy
	descriptionSanitizePolicy  *bluemonday.Policy
}

func NewProblem(
	token Token,
	role Role,
	accountDataAccessor db.AccountDataAccessor,
	problemDataAccessor db.ProblemDataAccessor,
	problemExampleDataAccessor db.ProblemExampleDataAccessor,
	logger *zap.Logger,
	db *gorm.DB,
) Problem {
	return &problem{
		token:                      token,
		role:                       role,
		accountDataAccessor:        accountDataAccessor,
		problemDataAccessor:        problemDataAccessor,
		problemExampleDataAccessor: problemExampleDataAccessor,
		logger:                     logger,
		db:                         db,
		displayNameSanitizePolicy:  bluemonday.StrictPolicy(),
		descriptionSanitizePolicy:  bluemonday.UGCPolicy(),
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
	description = p.displayNameSanitizePolicy.Sanitize(description)
	return description
}

func (p problem) dbProblemExampleToRPCProblemExample(problemExample *db.ProblemExample) rpc.ProblemExample {
	return rpc.ProblemExample{
		Input:  problemExample.Input,
		Output: problemExample.Output,
	}
}

func (p problem) dbProblemToRPCProblem(problem *db.Problem, author *db.Account, problemExampleList []*db.ProblemExample) rpc.Problem {
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
		CreatedTime: uint64(problem.CreatedAt.Unix()),
		UpdatedTime: uint64(problem.UpdatedAt.Unix()),
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
		CreatedTime:            uint64(problem.CreatedAt.Unix()),
		UpdatedTime:            uint64(problem.UpdatedAt.Unix()),
	}
}

func (p problem) CreateProblem(ctx context.Context, in *rpc.CreateProblemRequest, token string) (*rpc.CreateProblemResponse, error) {
	account, err := p.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	if hasPermission, err := p.role.AccountHasPermission(ctx, string(account.Role), PermissionProblemsWrite); err != nil {
		return nil, err
	} else if !hasPermission {
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
		if err := p.problemDataAccessor.WithDB(tx).CreateProblem(ctx, problem); err != nil {
			return err
		}

		problemExampleList := lo.Map(in.ExampleList, func(item rpc.ProblemExample, _ int) *db.ProblemExample {
			return &db.ProblemExample{
				OfProblemID: uint64(problem.ID),
				Input:       item.Input,
				Output:      item.Output,
			}
		})
		if err := p.problemExampleDataAccessor.WithDB(tx).CreateProblemExampleList(ctx, problemExampleList); err != nil {
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

	if hasPermission, err := p.role.AccountHasPermission(ctx, string(account.Role), PermissionProblemsWrite); err != nil {
		return err
	} else if !hasPermission {
		return pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		problem, err := p.problemDataAccessor.WithDB(tx).GetProblem(ctx, in.ID)
		if err != nil {
			return err
		}

		if problem == nil {
			logger.With(zap.Uint64("id", in.ID)).Error("cannot find problem")
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
		}

		if hasAccess, err := p.role.AccountCanAccessResource(
			ctx,
			uint64(account.ID),
			string(account.Role),
			problem.AuthorAccountID,
		); err != nil {
			return err
		} else if !hasAccess {
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}

		return p.problemDataAccessor.WithDB(tx).DeleteProblem(ctx, uint64(problem.ID))
	})
}

func (p problem) GetAccountProblemSnippetList(ctx context.Context, in *rpc.GetAccountProblemSnippetListRequest, token string) (*rpc.GetAccountProblemSnippetListResponse, error) {
	account, err := p.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	if hasPermission, err := p.role.AccountHasPermission(ctx, string(account.Role), PermissionProblemsRead); err != nil {
		return nil, err
	} else if !hasPermission {
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
		author, err := p.accountDataAccessor.GetAccount(ctx, problemList[i].AuthorAccountID)
		if err != nil {
			return nil, err
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

func (p problem) GetProblem(ctx context.Context, in *rpc.GetProblemRequest, token string) (*rpc.GetProblemResponse, error) {
	logger := utils.LoggerWithContext(ctx, p.logger)

	account, err := p.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	if hasPermission, err := p.role.AccountHasPermission(ctx, string(account.Role), PermissionProblemsRead); err != nil {
		return nil, err
	} else if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	problem, err := p.problemDataAccessor.GetProblem(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	if problem == nil {
		logger.With(zap.Uint64("id", in.ID)).Error("cannot find problem")
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
	}

	if hasAccess, err := p.role.AccountCanAccessResource(
		ctx,
		uint64(account.ID),
		string(account.Role),
		problem.AuthorAccountID,
	); err != nil {
		return nil, err
	} else if !hasAccess {
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

func (p problem) GetProblemSnippetList(ctx context.Context, in *rpc.GetProblemSnippetListRequest, token string) (*rpc.GetProblemSnippetListResponse, error) {
	account, err := p.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	if hasPermission, err := p.role.AccountHasPermission(ctx, string(account.Role), PermissionProblemsRead); err != nil {
		return nil, err
	} else if !hasPermission {
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
		author, err := p.accountDataAccessor.GetAccount(ctx, problemList[i].AuthorAccountID)
		if err != nil {
			return nil, err
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

func (p problem) UpdateProblem(ctx context.Context, in *rpc.UpdateProblemRequest, token string) (*rpc.UpdateProblemResponse, error) {
	logger := utils.LoggerWithContext(ctx, p.logger)

	account, err := p.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	if hasPermission, err := p.role.AccountHasPermission(ctx, string(account.Role), PermissionProblemsWrite); err != nil {
		return nil, err
	} else if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	response := &rpc.UpdateProblemResponse{}
	if txErr := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		problem, err := p.problemDataAccessor.WithDB(tx).GetProblem(ctx, in.ID)
		if err != nil {
			return err
		}

		if problem == nil {
			logger.With(zap.Uint64("id", in.ID)).Error("cannot find problem")
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
		}

		if hasAccess, err := p.role.AccountCanAccessResource(
			ctx,
			uint64(account.ID),
			string(account.Role),
			problem.AuthorAccountID,
		); err != nil {
			return err
		} else if !hasAccess {
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}

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

		if err := p.problemDataAccessor.WithDB(tx).UpdateProblem(ctx, problem); err != nil {
			return err
		}

		author, err := p.accountDataAccessor.GetAccount(ctx, problem.AuthorAccountID)
		if err != nil {
			return err
		}

		if author == nil {
			return pjrpc.JRPCErrInternalError()
		}

		var problemExampleList []*db.ProblemExample
		if in.ExampleList == nil {
			problemExampleList, err = p.problemExampleDataAccessor.WithDB(tx).GetProblemExampleListOfProblem(ctx, uint64(problem.ID))
			if err != nil {
				return err
			}
		} else {
			if err := p.problemExampleDataAccessor.DeleteProblemExampleOfProblem(ctx, uint64(problem.ID)); err != nil {
				return err
			}

			problemExampleList = lo.Map(*in.ExampleList, func(item rpc.ProblemExample, _ int) *db.ProblemExample {
				return &db.ProblemExample{
					OfProblemID: uint64(problem.ID),
					Input:       item.Input,
					Output:      item.Output,
				}
			})
			if err := p.problemExampleDataAccessor.WithDB(tx).CreateProblemExampleList(ctx, problemExampleList); err != nil {
				return err
			}
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
		token:                      p.token,
		role:                       p.role,
		accountDataAccessor:        p.accountDataAccessor.WithDB(db),
		problemDataAccessor:        p.problemDataAccessor.WithDB(db),
		problemExampleDataAccessor: p.problemExampleDataAccessor.WithDB(db),
		logger:                     p.logger,
		db:                         db,
		displayNameSanitizePolicy:  p.displayNameSanitizePolicy,
		descriptionSanitizePolicy:  p.descriptionSanitizePolicy,
	}
}
