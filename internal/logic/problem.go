package logic

import (
	"context"
	"strings"

	"gitlab.com/pjrpc/pjrpc/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/mikespook/gorbac"
	"github.com/samber/lo"

	"github.com/microcosm-cc/bluemonday"

	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc/rpcclient"
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
	SyncProblemList(ctx context.Context) error
	WithDB(db *gorm.DB) Problem
}

type problem struct {
	token                           Token
	role                            Role
	testCase                        TestCase
	setting                         Setting
	accountDataAccessor             db.AccountDataAccessor
	problemDataAccessor             db.ProblemDataAccessor
	problemExampleDataAccessor      db.ProblemExampleDataAccessor
	problemTestCaseHashDataAccessor db.ProblemTestCaseHashDataAccessor
	testCaseDataAccessor            db.TestCaseDataAccessor
	submissionDataAccessor          db.SubmissionDataAccessor
	logger                          *zap.Logger
	db                              *gorm.DB
	apiClient                       rpcclient.APIClient
	logicConfig                     configs.Logic
	displayNameSanitizePolicy       *bluemonday.Policy
	descriptionSanitizePolicy       *bluemonday.Policy
}

func NewProblem(
	token Token,
	role Role,
	testCase TestCase,
	setting Setting,
	accountDataAccessor db.AccountDataAccessor,
	problemDataAccessor db.ProblemDataAccessor,
	problemExampleDataAccessor db.ProblemExampleDataAccessor,
	problemTestCaseHashDataAccessor db.ProblemTestCaseHashDataAccessor,
	testCaseDataAccessor db.TestCaseDataAccessor,
	submissionDataAccessor db.SubmissionDataAccessor,
	logger *zap.Logger,
	db *gorm.DB,
	apiClient rpcclient.APIClient,
	logicConfig configs.Logic,
) Problem {
	return &problem{
		token:                           token,
		role:                            role,
		testCase:                        testCase,
		setting:                         setting,
		accountDataAccessor:             accountDataAccessor,
		problemDataAccessor:             problemDataAccessor,
		problemExampleDataAccessor:      problemExampleDataAccessor,
		problemTestCaseHashDataAccessor: problemTestCaseHashDataAccessor,
		testCaseDataAccessor:            testCaseDataAccessor,
		submissionDataAccessor:          submissionDataAccessor,
		logger:                          logger,
		db:                              db,
		apiClient:                       apiClient,
		logicConfig:                     logicConfig,
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
	problemTestCaseHash *db.ProblemTestCaseHash,
) rpc.Problem {
	return rpc.Problem{
		UUID:        problem.UUID,
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
		CreatedTime:  uint64(problem.CreatedAt.UnixMilli()),
		UpdatedTime:  uint64(problem.UpdatedAt.UnixMilli()),
		TestCaseHash: problemTestCaseHash.Hash,
	}
}

func (p problem) dbProblemToRPCProblemSnippet(problem *db.Problem, author *db.Account) rpc.ProblemSnippet {
	return rpc.ProblemSnippet{
		UUID:        problem.UUID,
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

func (p problem) createProblem(
	ctx context.Context,
	uuid string,
	displayName string,
	authorAccountID uint64,
	description string,
	timeLimitInMillisecond uint64,
	memoryLimitInByte uint64,
	exampleList []rpc.ProblemExample,
) (*db.Problem, []*db.ProblemExample, *db.ProblemTestCaseHash, error) {
	displayName = p.cleanupDisplayName(displayName)
	if !p.isValidDisplayName(displayName) {
		return nil, make([]*db.ProblemExample, 0), nil, pjrpc.JRPCErrInvalidParams()
	}

	description = p.cleanupDescription(description)

	var (
		problem             *db.Problem
		problemExampleList  []*db.ProblemExample
		problemTestCaseHash *db.ProblemTestCaseHash
		err                 error
	)

	if txErr := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		problem = &db.Problem{
			UUID:                   uuid,
			DisplayName:            displayName,
			AuthorAccountID:        authorAccountID,
			Description:            description,
			TimeLimitInMillisecond: timeLimitInMillisecond,
			MemoryLimitInByte:      memoryLimitInByte,
		}
		err = p.problemDataAccessor.WithDB(tx).CreateProblem(ctx, problem)
		if err != nil {
			return err
		}

		problemExampleList = lo.Map(exampleList, func(item rpc.ProblemExample, _ int) *db.ProblemExample {
			return &db.ProblemExample{
				OfProblemID: uint64(problem.ID),
				Input:       utils.TrimSpaceRight(item.Input),
				Output:      utils.TrimSpaceRight(item.Output),
			}
		})
		err = p.problemExampleDataAccessor.WithDB(tx).CreateProblemExampleList(ctx, problemExampleList)
		if err != nil {
			return err
		}

		hash, hashErr := p.testCase.CalculateProblemTestCaseHash(ctx, uint64(problem.ID))
		if hashErr != nil {
			return hashErr
		}

		problemTestCaseHash = &db.ProblemTestCaseHash{
			OfProblemID: uint64(problem.ID),
			Hash:        hash,
		}
		err = p.problemTestCaseHashDataAccessor.WithDB(tx).CreateProblemTestCaseHash(ctx, problemTestCaseHash)
		if err != nil {
			return err
		}

		return nil
	}); txErr != nil {
		return nil, make([]*db.ProblemExample, 0), nil, txErr
	}

	return problem, problemExampleList, problemTestCaseHash, nil
}

func (p problem) CreateProblem(
	ctx context.Context,
	in *rpc.CreateProblemRequest,
	token string,
) (*rpc.CreateProblemResponse, error) {
	logger := utils.LoggerWithContext(ctx, p.logger)

	setting, err := p.setting.GetSetting(ctx)
	if err != nil {
		return nil, err
	}

	if setting.Problem.DisableProblemCreation {
		logger.Info("problem creation is disabled via setting")
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeUnavailable))
	}

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

	problem, problemExampleList, problemTestCaseHash, err := p.createProblem(
		ctx,
		uuid.NewString(),
		in.DisplayName,
		uint64(account.ID),
		in.Description,
		in.TimeLimitInMillisecond,
		in.MemoryLimitInByte,
		in.ExampleList,
	)
	if err != nil {
		return nil, err
	}

	return &rpc.CreateProblemResponse{
		Problem: p.dbProblemToRPCProblem(problem, account, problemExampleList, problemTestCaseHash),
	}, nil
}

func (p problem) DeleteProblem(ctx context.Context, in *rpc.DeleteProblemRequest, token string) error {
	logger := utils.LoggerWithContext(ctx, p.logger)

	account, err := p.token.GetAccount(ctx, token)
	if err != nil {
		return err
	}

	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		problem, problemErr := p.problemDataAccessor.WithDB(tx).GetProblemByUUID(ctx, in.UUID)
		if problemErr != nil {
			return problemErr
		}

		if problem == nil {
			logger.With(zap.String("uuid", in.UUID)).Error("cannot find problem")
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

	problem, err := p.problemDataAccessor.GetProblemByUUID(ctx, in.UUID)
	if err != nil {
		return nil, err
	}

	if problem == nil {
		logger.With(zap.String("uuid", in.UUID)).Error("cannot find problem")
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

	problemTestCaseHash, err := p.problemTestCaseHashDataAccessor.GetProblemTestCaseHashOfProblem(ctx, uint64(problem.ID))
	if err != nil {
		return nil, err
	}

	return &rpc.GetProblemResponse{
		Problem: p.dbProblemToRPCProblem(problem, author, problemExampleList, problemTestCaseHash),
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

func (p problem) updateProblemExampleList(
	ctx context.Context,
	problemID uint64,
	exampleList *[]rpc.ProblemExample,
	tx *gorm.DB,
) ([]*db.ProblemExample, error) {
	if exampleList == nil {
		return p.problemExampleDataAccessor.WithDB(tx).GetProblemExampleListOfProblem(ctx, problemID)
	}

	if err := p.problemExampleDataAccessor.WithDB(tx).DeleteProblemExampleOfProblem(ctx, problemID); err != nil {
		return make([]*db.ProblemExample, 0), err
	}

	problemExampleList := lo.Map(*exampleList, func(item rpc.ProblemExample, _ int) *db.ProblemExample {
		return &db.ProblemExample{OfProblemID: problemID, Input: item.Input, Output: item.Output}
	})
	if err := p.problemExampleDataAccessor.WithDB(tx).CreateProblemExampleList(ctx, problemExampleList); err != nil {
		return make([]*db.ProblemExample, 0), err
	}

	return problemExampleList, nil
}

func (p problem) updateProblem(
	ctx context.Context,
	problem *db.Problem,
	displayName *string,
	description *string,
	timeLimitInMillisecond *uint64,
	memoryLimitInByte *uint64,
	exampleList *[]rpc.ProblemExample,
	tx *gorm.DB,
) (*db.Problem, []*db.ProblemExample, error) {
	if displayName != nil {
		cleanedDisplayName := p.cleanupDisplayName(*displayName)
		if !p.isValidDisplayName(cleanedDisplayName) {
			return nil, make([]*db.ProblemExample, 0), pjrpc.JRPCErrInvalidParams()
		}

		problem.DisplayName = cleanedDisplayName
	}

	if description != nil {
		problem.Description = p.cleanupDescription(*description)
	}

	if timeLimitInMillisecond != nil {
		problem.TimeLimitInMillisecond = *timeLimitInMillisecond
	}

	if memoryLimitInByte != nil {
		problem.MemoryLimitInByte = *memoryLimitInByte
	}

	if err := p.problemDataAccessor.WithDB(tx).UpdateProblem(ctx, problem); err != nil {
		return nil, make([]*db.ProblemExample, 0), err
	}

	problemExampleList, err := p.updateProblemExampleList(ctx, uint64(problem.ID), exampleList, tx)
	if err != nil {
		return nil, make([]*db.ProblemExample, 0), err
	}

	return problem, problemExampleList, nil
}

func (p problem) UpdateProblem(
	ctx context.Context,
	in *rpc.UpdateProblemRequest,
	token string,
) (*rpc.UpdateProblemResponse, error) {
	logger := utils.LoggerWithContext(ctx, p.logger)

	setting, err := p.setting.GetSetting(ctx)
	if err != nil {
		return nil, err
	}

	if setting.Problem.DisableProblemUpdate {
		logger.Info("problem update is disabled via setting")
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeUnavailable))
	}

	account, err := p.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	response := &rpc.UpdateProblemResponse{}
	if txErr := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		problem, problemErr := p.problemDataAccessor.WithDB(tx).GetProblemByUUID(ctx, in.UUID)
		if problemErr != nil {
			return problemErr
		}

		if problem == nil {
			logger.With(zap.String("uuid", in.UUID)).Error("cannot find problem")
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

		problem, problemExampleList, problemErr := p.updateProblem(
			ctx,
			problem,
			in.DisplayName,
			in.Description,
			in.TimeLimitInMillisecond,
			in.MemoryLimitInByte,
			in.ExampleList,
			tx,
		)
		if problemErr != nil {
			return problemErr
		}

		author, problemErr := p.accountDataAccessor.GetAccount(ctx, problem.AuthorAccountID)
		if problemErr != nil {
			return problemErr
		}

		problemTestCaseHash, problemErr := p.problemTestCaseHashDataAccessor.
			GetProblemTestCaseHashOfProblem(ctx, uint64(problem.ID))
		if problemErr != nil {
			return problemErr
		}

		response.Problem = p.dbProblemToRPCProblem(problem, author, problemExampleList, problemTestCaseHash)

		return nil
	}); txErr != nil {
		return nil, err
	}

	return response, nil
}

func (p problem) syncProblem(ctx context.Context, problemUUID string) error {
	logger := utils.LoggerWithContext(ctx, p.logger).With(zap.String("problem_uuid", problemUUID))
	logger.Info("start syncing problem")
	defer func() { logger.Info("syncing problem completed") }()

	getProblemResponse, err := p.apiClient.GetProblem(ctx, &rpc.GetProblemRequest{UUID: problemUUID})
	if err != nil {
		return err
	}

	problem, err := p.problemDataAccessor.GetProblemByUUID(ctx, problemUUID)
	if err != nil {
		return err
	}

	if problem != nil {
		problemTestCaseHash, problemTestCaseHashErr := p.problemTestCaseHashDataAccessor.
			GetProblemTestCaseHashOfProblem(ctx, uint64(problem.ID))
		if problemTestCaseHashErr != nil {
			return problemTestCaseHashErr
		}

		if getProblemResponse.Problem.TestCaseHash == problemTestCaseHash.Hash {
			logger.Info("test case hash of problem is unchanged, will skip sync")
			return nil
		}
	}

	if problem == nil {
		logger.Info("problem not found locally, will create new")
		problem, _, _, err = p.createProblem(
			ctx,
			getProblemResponse.Problem.UUID,
			getProblemResponse.Problem.DisplayName,
			getProblemResponse.Problem.Author.ID,
			getProblemResponse.Problem.Description,
			getProblemResponse.Problem.TimeLimitInMillisecond,
			getProblemResponse.Problem.MemoryLimitInByte,
			getProblemResponse.Problem.ExampleList,
		)
		if err != nil {
			return err
		}
	} else {
		logger.Info("problem found locally, but hash changed, will update new")
		err = p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			problem, _, err = p.updateProblem(
				ctx,
				problem,
				&getProblemResponse.Problem.DisplayName,
				&getProblemResponse.Problem.Description,
				&getProblemResponse.Problem.TimeLimitInMillisecond,
				&getProblemResponse.Problem.MemoryLimitInByte,
				&getProblemResponse.Problem.ExampleList,
				tx,
			)

			return err
		})
		if err != nil {
			return err
		}
	}

	return utils.ExecuteUntilFirstError(
		func() error { return p.testCase.SyncProblemTestCaseList(ctx, problemUUID) },
		func() error { return p.testCase.UpsertProblemTestCaseHash(ctx, uint64(problem.ID)) },
	)
}

func (p problem) SyncProblemList(ctx context.Context) error {
	currentOffset := uint64(0)
	for {
		response, err := p.apiClient.GetProblemSnippetList(ctx, &rpc.GetProblemSnippetListRequest{
			Offset: currentOffset,
			Limit:  p.logicConfig.SyncProblem.GetProblemSnippetListBatchSize,
		})
		if err != nil {
			return err
		}

		if len(response.ProblemSnippetList) == 0 {
			return nil
		}

		for _, problemSnippet := range response.ProblemSnippetList {
			err = p.syncProblem(ctx, problemSnippet.UUID)
			if err != nil {
				return err
			}
		}

		currentOffset += uint64(len(response.ProblemSnippetList))
	}
}

func (p problem) WithDB(db *gorm.DB) Problem {
	return &problem{
		token:                           p.token.WithDB(db),
		role:                            p.role,
		setting:                         p.setting.WithDB(db),
		accountDataAccessor:             p.accountDataAccessor.WithDB(db),
		problemDataAccessor:             p.problemDataAccessor.WithDB(db),
		problemExampleDataAccessor:      p.problemExampleDataAccessor.WithDB(db),
		problemTestCaseHashDataAccessor: p.problemTestCaseHashDataAccessor.WithDB(db),
		testCase:                        p.testCase.WithDB(db),
		testCaseDataAccessor:            p.testCaseDataAccessor.WithDB(db),
		submissionDataAccessor:          p.submissionDataAccessor.WithDB(db),
		logger:                          p.logger,
		db:                              db,
		apiClient:                       p.apiClient,
		displayNameSanitizePolicy:       p.displayNameSanitizePolicy,
		descriptionSanitizePolicy:       p.descriptionSanitizePolicy,
		logicConfig:                     p.logicConfig,
	}
}
