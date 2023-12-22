package logic

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"path/filepath"
	"sort"

	"gitlab.com/pjrpc/pjrpc/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/samber/lo"

	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc"
	"github.com/tranHieuDev23/cato/internal/utils"
)

const (
	testCaseDataSize       = 250
	testCaseInputFileName  = "input.txt"
	testCaseOutputFileName = "output.txt"
)

type TestCase interface {
	CreateTestCase(ctx context.Context, in *rpc.CreateTestCaseRequest, token string) (*rpc.CreateTestCaseResponse, error)
	CreateTestCaseList(ctx context.Context, in *rpc.CreateTestCaseListRequest, token string) error
	GetProblemTestCaseSnippetList(ctx context.Context, in *rpc.GetProblemTestCaseSnippetListRequest, token string) (*rpc.GetProblemTestCaseSnippetListResponse, error)
	GetTestCase(ctx context.Context, in *rpc.GetTestCaseRequest, token string) (*rpc.GetTestCaseResponse, error)
	UpdateTestCase(ctx context.Context, in *rpc.UpdateTestCaseRequest, token string) (*rpc.UpdateTestCaseResponse, error)
	DeleteTestCase(ctx context.Context, in *rpc.DeleteTestCaseRequest, token string) error
}

type testCase struct {
	token                Token
	role                 Role
	problemDataAccessor  db.ProblemDataAccessor
	testCaseDataAccessor db.TestCaseDataAccessor
	db                   *gorm.DB
	logger               *zap.Logger
}

func NewTestCase(
	token Token,
	role Role,
	problemDataAccessor db.ProblemDataAccessor,
	testCaseDataAccessor db.TestCaseDataAccessor,
	db *gorm.DB,
	logger *zap.Logger,
) TestCase {
	return &testCase{
		token:                token,
		role:                 role,
		problemDataAccessor:  problemDataAccessor,
		testCaseDataAccessor: testCaseDataAccessor,
		db:                   db,
		logger:               logger,
	}
}

func (t testCase) getTextSnippet(text string) string {
	if len(text) < testCaseDataSize {
		return text
	}

	return text[:testCaseDataSize] + "..."
}

func (t testCase) dbTestCaseToRPCTestCaseSnippet(testCase *db.TestCase) rpc.TestCaseSnippet {
	return rpc.TestCaseSnippet{
		ID:       uint64(testCase.ID),
		Input:    t.getTextSnippet(testCase.Input),
		Output:   t.getTextSnippet(testCase.Output),
		IsHidden: testCase.IsHidden,
	}
}

func (t testCase) CreateTestCase(ctx context.Context, in *rpc.CreateTestCaseRequest, token string) (*rpc.CreateTestCaseResponse, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	account, err := t.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	if hasPermission, err := t.role.AccountHasPermission(ctx, string(account.Role), PermissionTestCasesWrite); err != nil {
		return nil, err
	} else if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	response := &rpc.CreateTestCaseResponse{}
	if txErr := t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		problem, err := t.problemDataAccessor.WithDB(tx).GetProblem(ctx, in.ProblemID)
		if err != nil {
			return err
		}

		if problem == nil {
			logger.With(zap.Uint64("problem_id", in.ProblemID)).Error("cannot find problem")
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
		}

		if hasAccess, err := t.role.AccountCanAccessResource(
			ctx,
			uint64(account.ID),
			string(account.Role),
			problem.AuthorAccountID,
		); err != nil {
			return err
		} else if !hasAccess {
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}

		testCase := &db.TestCase{
			OfProblemID: in.ProblemID,
			Input:       in.Input,
			Output:      in.Output,
			IsHidden:    in.IsHidden,
		}
		if err := t.testCaseDataAccessor.WithDB(tx).CreateTestCase(ctx, testCase); err != nil {
			return err
		}

		response.TestCaseSnippet = t.dbTestCaseToRPCTestCaseSnippet(testCase)

		return nil
	}); txErr != nil {
		return nil, txErr
	}

	return response, nil
}

func (t testCase) getTestCaseListFromZippedData(ctx context.Context, problemID uint64, zippedData []byte) ([]*db.TestCase, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	zippedDataReader, err := zip.NewReader(bytes.NewReader(zippedData), int64(len(zippedData)))
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to open zip reader")
		return nil, pjrpc.JRPCErrInternalError()
	}

	type unzippedTestCase struct {
		Input  *string
		Output *string
	}

	fileDirectoryToUnzippedTestCaseMap := make(map[string]*unzippedTestCase)

	for i := range zippedDataReader.File {
		fileInfo := zippedDataReader.File[i].FileInfo()
		if fileInfo.IsDir() {
			continue
		}

		filePath := fileInfo.Name()
		fileDirectory, fileName := filepath.Split(filePath)
		if fileName != testCaseInputFileName && fileName != testCaseOutputFileName {
			continue
		}

		if _, ok := fileDirectoryToUnzippedTestCaseMap[fileDirectory]; !ok {
			fileDirectoryToUnzippedTestCaseMap[fileDirectory] = &unzippedTestCase{}
		}

		fileReader, err := zippedDataReader.File[i].Open()
		if err != nil {
			logger.With(zap.Error(err)).Error("failed to open file reader")
			return nil, pjrpc.JRPCErrInternalError()
		}

		fileContent, err := io.ReadAll(fileReader)
		if err != nil {
			logger.With(zap.Error(err)).Error("failed to read file content")
			return nil, pjrpc.JRPCErrInternalError()
		}

		fileContentString := string(fileContent)
		if fileName == testCaseInputFileName {
			fileDirectoryToUnzippedTestCaseMap[fileDirectory].Input = &fileContentString
		} else {
			fileDirectoryToUnzippedTestCaseMap[fileDirectory].Output = &fileContentString
		}
	}

	fileDirectoryToUnzippedTestCaseEntryList := lo.Entries[string, *unzippedTestCase](fileDirectoryToUnzippedTestCaseMap)
	sort.Slice(fileDirectoryToUnzippedTestCaseEntryList, func(i, j int) bool {
		return fileDirectoryToUnzippedTestCaseEntryList[i].Key < fileDirectoryToUnzippedTestCaseEntryList[j].Key
	})

	testCaseList := make([]*db.TestCase, 0)
	for i := range fileDirectoryToUnzippedTestCaseEntryList {
		unzippedTestCase := fileDirectoryToUnzippedTestCaseEntryList[i].Value
		if unzippedTestCase.Input == nil || unzippedTestCase.Output == nil {
			continue
		}

		testCaseList = append(testCaseList, &db.TestCase{
			OfProblemID: problemID,
			Input:       *unzippedTestCase.Input,
			Output:      *unzippedTestCase.Output,
			IsHidden:    true,
		})
	}

	return testCaseList, nil
}

func (t testCase) CreateTestCaseList(ctx context.Context, in *rpc.CreateTestCaseListRequest, token string) error {
	logger := utils.LoggerWithContext(ctx, t.logger)

	account, err := t.token.GetAccount(ctx, token)
	if err != nil {
		return err
	}

	if hasPermission, err := t.role.AccountHasPermission(ctx, string(account.Role), PermissionTestCasesWrite); err != nil {
		return err
	} else if !hasPermission {
		return pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	problem, err := t.problemDataAccessor.GetProblem(ctx, in.ProblemID)
	if err != nil {
		return err
	}

	if problem == nil {
		logger.With(zap.Uint64("problem_id", in.ProblemID)).Error("cannot find problem")
		return pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
	}

	if hasAccess, err := t.role.AccountCanAccessResource(
		ctx,
		uint64(account.ID),
		string(account.Role),
		problem.AuthorAccountID,
	); err != nil {
		return err
	} else if !hasAccess {
		return pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	testCaseList, err := t.getTestCaseListFromZippedData(ctx, in.ProblemID, []byte(in.ZippedTestData))
	if err != nil {
		return err
	}

	return t.testCaseDataAccessor.CreateTestCaseList(ctx, testCaseList)
}

func (t testCase) DeleteTestCase(ctx context.Context, in *rpc.DeleteTestCaseRequest, token string) error {
	logger := utils.LoggerWithContext(ctx, t.logger)

	account, err := t.token.GetAccount(ctx, token)
	if err != nil {
		return err
	}

	if hasPermission, err := t.role.AccountHasPermission(ctx, string(account.Role), PermissionTestCasesWrite); err != nil {
		return err
	} else if !hasPermission {
		return pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		testCase, err := t.testCaseDataAccessor.WithDB(tx).GetTestCase(ctx, in.ID)
		if err != nil {
			return err
		}

		if testCase == nil {
			logger.With(zap.Uint64("id", in.ID)).Error("cannot find test case")
		}

		problem, err := t.problemDataAccessor.WithDB(tx).GetProblem(ctx, testCase.OfProblemID)
		if err != nil {
			return err
		}

		if hasAccess, err := t.role.AccountCanAccessResource(
			ctx,
			uint64(account.ID),
			string(account.Role),
			problem.AuthorAccountID,
		); err != nil {
			return err
		} else if !hasAccess {
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}

		return t.testCaseDataAccessor.WithDB(tx).DeleteTestCase(ctx, uint64(testCase.ID))
	})
}

func (t testCase) GetProblemTestCaseSnippetList(ctx context.Context, in *rpc.GetProblemTestCaseSnippetListRequest, token string) (*rpc.GetProblemTestCaseSnippetListResponse, error) {
	account, err := t.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	if hasPermission, err := t.role.AccountHasPermission(ctx, string(account.Role), PermissionTestCasesRead); err != nil {
		return nil, err
	} else if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	problem, err := t.problemDataAccessor.GetProblem(ctx, in.ProblemID)
	if err != nil {
		return nil, err
	}

	if hasAccess, err := t.role.AccountCanAccessResource(
		ctx,
		uint64(account.ID),
		string(account.Role),
		problem.AuthorAccountID,
	); err != nil {
		return nil, err
	} else if !hasAccess {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	testCaseCount, err := t.testCaseDataAccessor.GetTestCaseCountOfProblem(ctx, in.ProblemID)
	if err != nil {
		return nil, err
	}

	testCaseList, err := t.testCaseDataAccessor.GetTestCaseListOfProblem(ctx, in.ProblemID, in.Offset, in.Limit)
	if err != nil {
		return nil, err
	}

	return &rpc.GetProblemTestCaseSnippetListResponse{
		TotalTestCaseCount: testCaseCount,
		TestCaseSnippetList: lo.Map[*db.TestCase, rpc.TestCaseSnippet](testCaseList, func(item *db.TestCase, _ int) rpc.TestCaseSnippet {
			return t.dbTestCaseToRPCTestCaseSnippet(item)
		}),
	}, nil
}

func (t testCase) GetTestCase(ctx context.Context, in *rpc.GetTestCaseRequest, token string) (*rpc.GetTestCaseResponse, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	account, err := t.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	if hasPermission, err := t.role.AccountHasPermission(ctx, string(account.Role), PermissionTestCasesRead); err != nil {
		return nil, err
	} else if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	response := &rpc.GetTestCaseResponse{}
	if txErr := t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		testCase, err := t.testCaseDataAccessor.WithDB(tx).GetTestCase(ctx, in.ID)
		if err != nil {
			return err
		}

		if testCase == nil {
			logger.With(zap.Uint64("id", in.ID)).Error("cannot find test case")
		}

		problem, err := t.problemDataAccessor.WithDB(tx).GetProblem(ctx, testCase.OfProblemID)
		if err != nil {
			return err
		}

		if hasAccess, err := t.role.AccountCanAccessResource(
			ctx,
			uint64(account.ID),
			string(account.Role),
			problem.AuthorAccountID,
		); err != nil {
			return err
		} else if !hasAccess {
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}

		response.TestCase = rpc.TestCase(t.dbTestCaseToRPCTestCaseSnippet(testCase))

		return nil
	}); txErr != nil {
		return nil, txErr
	}

	return response, nil
}

func (t testCase) UpdateTestCase(ctx context.Context, in *rpc.UpdateTestCaseRequest, token string) (*rpc.UpdateTestCaseResponse, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	account, err := t.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	if hasPermission, err := t.role.AccountHasPermission(ctx, string(account.Role), PermissionTestCasesWrite); err != nil {
		return nil, err
	} else if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	response := &rpc.UpdateTestCaseResponse{}
	if txErr := t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		testCase, err := t.testCaseDataAccessor.WithDB(tx).GetTestCase(ctx, in.ID)
		if err != nil {
			return err
		}

		if testCase == nil {
			logger.With(zap.Uint64("id", in.ID)).Error("cannot find test case")
		}

		problem, err := t.problemDataAccessor.WithDB(tx).GetProblem(ctx, testCase.OfProblemID)
		if err != nil {
			return err
		}

		if hasAccess, err := t.role.AccountCanAccessResource(
			ctx,
			uint64(account.ID),
			string(account.Role),
			problem.AuthorAccountID,
		); err != nil {
			return err
		} else if !hasAccess {
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}

		if in.Input != nil {
			testCase.Input = *in.Input
		}

		if in.Output != nil {
			testCase.Output = *in.Output
		}

		if in.IsHidden != nil {
			testCase.IsHidden = *in.IsHidden
		}

		if err := t.testCaseDataAccessor.WithDB(tx).UpdateTestCase(ctx, testCase); err != nil {
			return err
		}

		response.TestCaseSnippet = t.dbTestCaseToRPCTestCaseSnippet(testCase)

		return nil
	}); txErr != nil {
		return nil, txErr
	}

	return response, nil
}
