package logic

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"hash/fnv"
	"io"
	"path/filepath"
	"sort"

	"gitlab.com/pjrpc/pjrpc/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/samber/lo"

	"github.com/tranHieuDev23/cato/internal/configs"
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
	CalculateProblemTestCaseHash(ctx context.Context, problemID uint64) (string, error)
	UpsertProblemTestCaseHash(ctx context.Context, problemID uint64) error
	CreateTestCase(ctx context.Context, in *rpc.CreateTestCaseRequest, token string) (*rpc.CreateTestCaseResponse, error)
	CreateTestCaseList(ctx context.Context, in *rpc.CreateTestCaseListRequest, token string) error
	GetProblemTestCaseSnippetList(ctx context.Context, in *rpc.GetProblemTestCaseSnippetListRequest, token string) (*rpc.GetProblemTestCaseSnippetListResponse, error)
	GetTestCase(ctx context.Context, in *rpc.GetTestCaseRequest, token string) (*rpc.GetTestCaseResponse, error)
	UpdateTestCase(ctx context.Context, in *rpc.UpdateTestCaseRequest, token string) (*rpc.UpdateTestCaseResponse, error)
	DeleteTestCase(ctx context.Context, in *rpc.DeleteTestCaseRequest, token string) error
	WithDB(db *gorm.DB) TestCase
}

type testCase struct {
	token                           Token
	role                            Role
	problemDataAccessor             db.ProblemDataAccessor
	testCaseDataAccessor            db.TestCaseDataAccessor
	problemTestCaseHashDataAccessor db.ProblemTestCaseHashDataAccessor
	db                              *gorm.DB
	logger                          *zap.Logger
	logicConfig                     configs.Logic
}

func NewTestCase(
	token Token,
	role Role,
	problemDataAccessor db.ProblemDataAccessor,
	testCaseDataAccessor db.TestCaseDataAccessor,
	problemTestCaseHashDataAccessor db.ProblemTestCaseHashDataAccessor,
	db *gorm.DB,
	logger *zap.Logger,
	logicConfig configs.Logic,
) TestCase {
	return &testCase{
		token:                           token,
		role:                            role,
		problemDataAccessor:             problemDataAccessor,
		testCaseDataAccessor:            testCaseDataAccessor,
		problemTestCaseHashDataAccessor: problemTestCaseHashDataAccessor,
		db:                              db,
		logger:                          logger,
		logicConfig:                     logicConfig,
	}
}

func (t testCase) getTextSnippet(text string) string {
	if len(text) < testCaseDataSize {
		return text
	}

	return text[:testCaseDataSize] + "..."
}

func (t testCase) dbTestCaseToRPCTestCase(testCase *db.TestCase, shouldHideInputOutput bool) rpc.TestCase {
	if shouldHideInputOutput && testCase.IsHidden {
		return rpc.TestCase{
			ID:       uint64(testCase.ID),
			IsHidden: testCase.IsHidden,
		}
	}

	return rpc.TestCase{
		ID:       uint64(testCase.ID),
		Input:    testCase.Input,
		Output:   testCase.Output,
		IsHidden: testCase.IsHidden,
	}
}

func (t testCase) dbTestCaseToRPCTestCaseSnippet(testCase *db.TestCase, shouldHideInputOutput bool) rpc.TestCaseSnippet {
	if shouldHideInputOutput && testCase.IsHidden {
		return rpc.TestCaseSnippet{
			ID:       uint64(testCase.ID),
			IsHidden: testCase.IsHidden,
		}
	}

	return rpc.TestCaseSnippet{
		ID:       uint64(testCase.ID),
		Input:    t.getTextSnippet(testCase.Input),
		Output:   t.getTextSnippet(testCase.Output),
		IsHidden: testCase.IsHidden,
	}
}

func (t testCase) calculateTestCaseHash(input, output string) string {
	fnvHash := fnv.New64a()
	fnvHash.Write([]byte(input))
	fnvHash.Write([]byte{0})
	fnvHash.Write([]byte(output))
	fnvHash.Write([]byte{0})
	return base64.StdEncoding.EncodeToString(fnvHash.Sum(nil))
}

func (t testCase) CalculateProblemTestCaseHash(ctx context.Context, problemID uint64) (string, error) {
	fnvHash := fnv.New64a()
	if txErr := t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		totalTestCaseCount, err := t.testCaseDataAccessor.WithDB(tx).GetTestCaseCountOfProblem(ctx, problemID)
		if err != nil {
			return err
		}

		for i := uint64(0); i < totalTestCaseCount; i += t.logicConfig.ProblemTestCaseHash.BatchSize {
			hashList, err := t.testCaseDataAccessor.
				WithDB(tx).
				GetTestCaseHashListOfProblem(ctx, problemID, i, t.logicConfig.ProblemTestCaseHash.BatchSize)
			if err != nil {
				return err
			}

			for _, hash := range hashList {
				fnvHash.Write([]byte(hash))
				fnvHash.Write([]byte{0})
			}
		}

		return nil
	}); txErr != nil {
		return "", txErr
	}

	return base64.StdEncoding.EncodeToString(fnvHash.Sum(nil)), nil
}

func (t testCase) UpsertProblemTestCaseHash(ctx context.Context, problemID uint64) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		hash, err := t.WithDB(tx).CalculateProblemTestCaseHash(ctx, problemID)
		if err != nil {
			return err
		}

		problemTestCaseHash, err := t.problemTestCaseHashDataAccessor.WithDB(tx).GetProblemTestCaseHashOfProblem(ctx, problemID)
		if err != nil {
			return err
		}

		if problemTestCaseHash == nil {
			return t.problemTestCaseHashDataAccessor.WithDB(tx).CreateProblemTestCaseHash(ctx, &db.ProblemTestCaseHash{
				OfProblemID: problemID,
				Hash:        hash,
			})
		}

		problemTestCaseHash.Hash = hash
		if err := t.problemTestCaseHashDataAccessor.WithDB(tx).UpdateProblemTestCaseHash(ctx, problemTestCaseHash); err != nil {
			return err
		}

		return nil
	})
}

func (t testCase) CreateTestCase(ctx context.Context, in *rpc.CreateTestCaseRequest, token string) (*rpc.CreateTestCaseResponse, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	account, err := t.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	if hasPermission, err := t.role.AccountHasPermission(
		ctx,
		string(account.Role),
		PermissionTestCasesSelfWrite,
		PermissionTestCasesAllWrite,
	); err != nil {
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

		if problem.AuthorAccountID != uint64(account.ID) {
			if hasPermission, err := t.role.AccountHasPermission(
				ctx,
				string(account.Role),
				PermissionTestCasesAllWrite,
			); err != nil {
				return err
			} else if !hasPermission {
				return pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
			}
		}

		testCase := &db.TestCase{
			OfProblemID: in.ProblemID,
			Input:       in.Input,
			Output:      in.Output,
			IsHidden:    in.IsHidden,
			Hash:        t.calculateTestCaseHash(in.Input, in.Output),
		}
		if err := t.testCaseDataAccessor.WithDB(tx).CreateTestCase(ctx, testCase); err != nil {
			return err
		}

		if err := t.WithDB(tx).UpsertProblemTestCaseHash(ctx, in.ProblemID); err != nil {
			return err
		}

		response.TestCaseSnippet = t.dbTestCaseToRPCTestCaseSnippet(testCase, false)

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
			Hash:        t.calculateTestCaseHash(*unzippedTestCase.Input, *unzippedTestCase.Output),
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

	if hasPermission, err := t.role.AccountHasPermission(
		ctx,
		string(account.Role),
		PermissionTestCasesSelfWrite,
		PermissionTestCasesAllWrite,
	); err != nil {
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

	if problem.AuthorAccountID != uint64(account.ID) {
		if hasPermission, err := t.role.AccountHasPermission(
			ctx,
			string(account.Role),
			PermissionTestCasesAllWrite,
		); err != nil {
			return err
		} else if !hasPermission {
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}
	}

	testCaseList, err := t.getTestCaseListFromZippedData(ctx, in.ProblemID, []byte(in.ZippedTestData))
	if err != nil {
		return err
	}

	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := t.testCaseDataAccessor.CreateTestCaseList(ctx, testCaseList); err != nil {
			return err
		}

		if err := t.WithDB(tx).UpsertProblemTestCaseHash(ctx, in.ProblemID); err != nil {
			return err
		}

		return nil
	})
}

func (t testCase) DeleteTestCase(ctx context.Context, in *rpc.DeleteTestCaseRequest, token string) error {
	logger := utils.LoggerWithContext(ctx, t.logger)

	account, err := t.token.GetAccount(ctx, token)
	if err != nil {
		return err
	}

	if hasPermission, err := t.role.AccountHasPermission(
		ctx,
		string(account.Role),
		PermissionTestCasesSelfWrite,
		PermissionTestCasesAllWrite,
	); err != nil {
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

		if problem.AuthorAccountID != uint64(account.ID) {
			if hasPermission, err := t.role.AccountHasPermission(
				ctx,
				string(account.Role),
				PermissionTestCasesAllWrite,
			); err != nil {
				return err
			} else if !hasPermission {
				return pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
			}
		}

		if err := t.WithDB(tx).UpsertProblemTestCaseHash(ctx, uint64(problem.ID)); err != nil {
			return err
		}

		return t.testCaseDataAccessor.WithDB(tx).DeleteTestCase(ctx, uint64(testCase.ID))
	})
}

func (t testCase) GetProblemTestCaseSnippetList(
	ctx context.Context,
	in *rpc.GetProblemTestCaseSnippetListRequest,
	token string,
) (*rpc.GetProblemTestCaseSnippetListResponse, error) {
	account, err := t.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	if hasPermission, err := t.role.AccountHasPermission(
		ctx,
		string(account.Role),
		PermissionTestCasesSelfRead,
		PermissionTestCasesAllRead,
	); err != nil {
		return nil, err
	} else if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	problem, err := t.problemDataAccessor.GetProblem(ctx, in.ProblemID)
	if err != nil {
		return nil, err
	}

	if problem.AuthorAccountID != uint64(account.ID) {
		if hasPermission, err := t.role.AccountHasPermission(
			ctx,
			string(account.Role),
			PermissionTestCasesAllRead,
		); err != nil {
			return nil, err
		} else if !hasPermission {
			return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}
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
			return t.dbTestCaseToRPCTestCaseSnippet(item, account.Role == db.AccountRoleContestant)
		}),
	}, nil
}

func (t testCase) GetTestCase(ctx context.Context, in *rpc.GetTestCaseRequest, token string) (*rpc.GetTestCaseResponse, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	account, err := t.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	if hasPermission, err := t.role.AccountHasPermission(
		ctx,
		string(account.Role),
		PermissionTestCasesSelfRead,
		PermissionTestCasesAllRead,
	); err != nil {
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

		if problem.AuthorAccountID != uint64(account.ID) {
			if hasPermission, err := t.role.AccountHasPermission(
				ctx,
				string(account.Role),
				PermissionTestCasesAllRead,
			); err != nil {
				return err
			} else if !hasPermission {
				return pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
			}
		}

		response.TestCase = t.dbTestCaseToRPCTestCase(testCase, account.Role == db.AccountRoleContestant)

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

	if hasPermission, err := t.role.AccountHasPermission(
		ctx,
		string(account.Role),
		PermissionTestCasesSelfWrite,
		PermissionTestCasesAllWrite,
	); err != nil {
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

		if problem.AuthorAccountID != uint64(account.ID) {
			if hasPermission, err := t.role.AccountHasPermission(
				ctx,
				string(account.Role),
				PermissionTestCasesAllWrite,
			); err != nil {
				return err
			} else if !hasPermission {
				return pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
			}
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

		if in.Input != nil || in.Output != nil {
			testCase.Hash = t.calculateTestCaseHash(testCase.Input, testCase.Output)
		}

		if err := t.testCaseDataAccessor.WithDB(tx).UpdateTestCase(ctx, testCase); err != nil {
			return err
		}

		if err := t.WithDB(tx).UpsertProblemTestCaseHash(ctx, uint64(problem.ID)); err != nil {
			return err
		}

		response.TestCaseSnippet = t.dbTestCaseToRPCTestCaseSnippet(testCase, false)

		return nil
	}); txErr != nil {
		return nil, txErr
	}

	return response, nil
}

func (t testCase) WithDB(db *gorm.DB) TestCase {
	return &testCase{
		token:                           t.token.WithDB(db),
		role:                            t.role,
		problemDataAccessor:             t.problemDataAccessor.WithDB(db),
		testCaseDataAccessor:            t.testCaseDataAccessor.WithDB(db),
		problemTestCaseHashDataAccessor: t.problemTestCaseHashDataAccessor.WithDB(db),
		db:                              db,
		logger:                          t.logger,
		logicConfig:                     t.logicConfig,
	}
}
