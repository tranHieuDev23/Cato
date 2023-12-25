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

	"github.com/mikespook/gorbac"
	"github.com/samber/lo"

	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc"
	"github.com/tranHieuDev23/cato/internal/utils"
)

const (
	testCaseDataSize       = 100
	testCaseInputFileName  = "input.txt"
	testCaseOutputFileName = "output.txt"
)

type TestCase interface {
	CalculateProblemTestCaseHash(ctx context.Context, problemID uint64) (string, error)
	UpsertProblemTestCaseHash(ctx context.Context, problemID uint64) error
	CreateTestCase(ctx context.Context, in *rpc.CreateTestCaseRequest, token string) (*rpc.CreateTestCaseResponse, error)
	CreateTestCaseList(ctx context.Context, in *rpc.CreateTestCaseListRequest, token string) error
	GetProblemTestCaseSnippetList(
		ctx context.Context,
		in *rpc.GetProblemTestCaseSnippetListRequest,
		token string,
	) (*rpc.GetProblemTestCaseSnippetListResponse, error)
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

func (t testCase) dbTestCaseToRPCTestCaseSnippet(
	testCase *db.TestCase,
	shouldHideInputOutput bool,
) rpc.TestCaseSnippet {
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
			hashList, hashListErr := t.testCaseDataAccessor.
				WithDB(tx).
				GetTestCaseHashListOfProblem(ctx, problemID, i, t.logicConfig.ProblemTestCaseHash.BatchSize)
			if hashListErr != nil {
				return hashListErr
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

		problemTestCaseHash, err := t.problemTestCaseHashDataAccessor.
			WithDB(tx).GetProblemTestCaseHashOfProblem(ctx, problemID)
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
		err = t.problemTestCaseHashDataAccessor.WithDB(tx).UpdateProblemTestCaseHash(ctx, problemTestCaseHash)
		if err != nil {
			return err
		}

		return nil
	})
}

func (t testCase) CreateTestCase(
	ctx context.Context,
	in *rpc.CreateTestCaseRequest,
	token string,
) (*rpc.CreateTestCaseResponse, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	account, err := t.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	response := &rpc.CreateTestCaseResponse{}
	if txErr := t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		problem, problemErr := t.problemDataAccessor.WithDB(tx).GetProblem(ctx, in.ProblemID)
		if problemErr != nil {
			return problemErr
		}

		if problem == nil {
			logger.With(zap.Uint64("problem_id", in.ProblemID)).Error("cannot find problem")
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
		}

		requiredPermissionList := []gorbac.Permission{PermissionTestCasesAllWrite}
		if problem.AuthorAccountID == uint64(account.ID) {
			requiredPermissionList = append(requiredPermissionList, PermissionTestCasesSelfWrite)
		}

		hasPermission, permissionErr := t.role.AccountHasPermission(ctx, string(account.Role), requiredPermissionList...)
		if permissionErr != nil {
			return permissionErr
		}
		if !hasPermission {
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}

		testCase := &db.TestCase{
			OfProblemID: in.ProblemID,
			Input:       utils.TrimSpaceRight(in.Input),
			Output:      utils.TrimSpaceRight(in.Output),
			IsHidden:    in.IsHidden,
			Hash:        t.calculateTestCaseHash(in.Input, in.Output),
		}

		problemErr = utils.ExecuteUntilFirstError(
			func() error { return t.testCaseDataAccessor.WithDB(tx).CreateTestCase(ctx, testCase) },
			func() error { return t.WithDB(tx).UpsertProblemTestCaseHash(ctx, in.ProblemID) },
		)
		if problemErr != nil {
			return problemErr
		}

		response.TestCaseSnippet = t.dbTestCaseToRPCTestCaseSnippet(testCase, false)

		return nil
	}); txErr != nil {
		return nil, txErr
	}

	return response, nil
}

type unzippedTestCase struct {
	Input  *string
	Output *string
}

func (t testCase) getFileDirectoryToUnzippedTestCaseMap(
	ctx context.Context,
	zipReader *zip.Reader,
) (map[string]*unzippedTestCase, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	fileDirectoryToUnzippedTestCaseMap := make(map[string]*unzippedTestCase)
	for i := range zipReader.File {
		if zipReader.File[i].FileInfo().IsDir() {
			continue
		}

		fileDirectory, fileName := filepath.Split(zipReader.File[i].Name)
		if fileName != testCaseInputFileName && fileName != testCaseOutputFileName {
			continue
		}

		if _, ok := fileDirectoryToUnzippedTestCaseMap[fileDirectory]; !ok {
			fileDirectoryToUnzippedTestCaseMap[fileDirectory] = &unzippedTestCase{}
		}

		fileReader, fileReaderErr := zipReader.File[i].Open()
		if fileReaderErr != nil {
			logger.With(zap.Error(fileReaderErr)).Error("failed to open file reader")
			return nil, pjrpc.JRPCErrInternalError()
		}

		fileContent, fileContentErr := io.ReadAll(fileReader)
		if fileContentErr != nil {
			logger.With(zap.Error(fileContentErr)).Error("failed to read file content")
			return nil, pjrpc.JRPCErrInternalError()
		}

		fileContentString := string(fileContent)
		if fileName == testCaseInputFileName {
			fileDirectoryToUnzippedTestCaseMap[fileDirectory].Input = &fileContentString
		} else {
			fileDirectoryToUnzippedTestCaseMap[fileDirectory].Output = &fileContentString
		}
	}

	return fileDirectoryToUnzippedTestCaseMap, nil
}

func (t testCase) getTestCaseListFromZippedData(
	ctx context.Context,
	problemID uint64,
	zippedTestData string,
) ([]*db.TestCase, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	decodedZippedTestData, err := base64.StdEncoding.DecodeString(zippedTestData)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to decode zipped test data")
		return nil, pjrpc.JRPCErrInternalError()
	}

	zippedTestDataReader, err := zip.NewReader(bytes.NewReader(decodedZippedTestData), int64(len(decodedZippedTestData)))
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to open zip reader")
		return nil, pjrpc.JRPCErrInternalError()
	}

	fileDirectoryToUnzippedTestCaseMap, err := t.getFileDirectoryToUnzippedTestCaseMap(ctx, zippedTestDataReader)
	if err != nil {
		return nil, err
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
			Input:       utils.TrimSpaceRight(*unzippedTestCase.Input),
			Output:      utils.TrimSpaceRight(*unzippedTestCase.Output),
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

	problem, err := t.problemDataAccessor.GetProblem(ctx, in.ProblemID)
	if err != nil {
		return err
	}

	if problem == nil {
		logger.With(zap.Uint64("problem_id", in.ProblemID)).Error("cannot find problem")
		return pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
	}

	requiredPermissionList := []gorbac.Permission{PermissionTestCasesAllWrite}
	if problem.AuthorAccountID == uint64(account.ID) {
		requiredPermissionList = append(requiredPermissionList, PermissionTestCasesSelfWrite)
	}

	hasPermission, problemErr := t.role.AccountHasPermission(ctx, string(account.Role), requiredPermissionList...)
	if problemErr != nil {
		return problemErr
	}
	if !hasPermission {
		return pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	testCaseList, err := t.getTestCaseListFromZippedData(ctx, in.ProblemID, in.ZippedTestData)
	if err != nil {
		return err
	}

	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return utils.ExecuteUntilFirstError(
			func() error { return t.testCaseDataAccessor.CreateTestCaseList(ctx, testCaseList) },
			func() error { return t.WithDB(tx).UpsertProblemTestCaseHash(ctx, in.ProblemID) },
		)
	})
}

func (t testCase) DeleteTestCase(ctx context.Context, in *rpc.DeleteTestCaseRequest, token string) error {
	logger := utils.LoggerWithContext(ctx, t.logger)

	account, err := t.token.GetAccount(ctx, token)
	if err != nil {
		return err
	}

	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		testCase, testCaseErr := t.testCaseDataAccessor.WithDB(tx).GetTestCase(ctx, in.ID)
		if testCaseErr != nil {
			return testCaseErr
		}

		if testCase == nil {
			logger.With(zap.Uint64("id", in.ID)).Error("cannot find test case")
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
		}

		problem, problemErr := t.problemDataAccessor.WithDB(tx).GetProblem(ctx, testCase.OfProblemID)
		if problemErr != nil {
			return problemErr
		}

		requiredPermissionList := []gorbac.Permission{PermissionTestCasesAllWrite}
		if problem.AuthorAccountID == uint64(account.ID) {
			requiredPermissionList = append(requiredPermissionList, PermissionTestCasesSelfWrite)
		}

		hasPermission, permissionErr := t.role.AccountHasPermission(ctx, string(account.Role), requiredPermissionList...)
		if permissionErr != nil {
			return permissionErr
		}
		if !hasPermission {
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}

		return utils.ExecuteUntilFirstError(
			func() error { return t.testCaseDataAccessor.WithDB(tx).DeleteTestCase(ctx, uint64(testCase.ID)) },
			func() error { return t.WithDB(tx).UpsertProblemTestCaseHash(ctx, uint64(problem.ID)) },
		)
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

	problem, err := t.problemDataAccessor.GetProblem(ctx, in.ProblemID)
	if err != nil {
		return nil, err
	}

	requiredPermissionList := []gorbac.Permission{PermissionTestCasesAllRead}
	if problem.AuthorAccountID == uint64(account.ID) {
		requiredPermissionList = append(requiredPermissionList, PermissionTestCasesSelfRead)
	}

	hasPermission, err := t.role.AccountHasPermission(ctx, string(account.Role), requiredPermissionList...)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
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
		TestCaseSnippetList: lo.Map[*db.TestCase, rpc.TestCaseSnippet](
			testCaseList,
			func(item *db.TestCase, _ int) rpc.TestCaseSnippet {
				return t.dbTestCaseToRPCTestCaseSnippet(item, account.Role == db.AccountRoleContestant)
			},
		),
	}, nil
}

func (t testCase) GetTestCase(
	ctx context.Context,
	in *rpc.GetTestCaseRequest,
	token string,
) (*rpc.GetTestCaseResponse, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	account, err := t.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	response := &rpc.GetTestCaseResponse{}
	if txErr := t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		testCase, testCaseErr := t.testCaseDataAccessor.WithDB(tx).GetTestCase(ctx, in.ID)
		if testCaseErr != nil {
			return testCaseErr
		}

		if testCase == nil {
			logger.With(zap.Uint64("id", in.ID)).Error("cannot find test case")
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
		}

		problem, problemErr := t.problemDataAccessor.WithDB(tx).GetProblem(ctx, testCase.OfProblemID)
		if err != nil {
			return problemErr
		}

		requiredPermissionList := []gorbac.Permission{PermissionTestCasesAllRead}
		if problem.AuthorAccountID == uint64(account.ID) {
			requiredPermissionList = append(requiredPermissionList, PermissionTestCasesSelfRead)
		}

		hasPermission, permissionErr := t.role.AccountHasPermission(ctx, string(account.Role), requiredPermissionList...)
		if permissionErr != nil {
			return permissionErr
		}
		if !hasPermission {
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}

		response.TestCase = t.dbTestCaseToRPCTestCase(testCase, account.Role == db.AccountRoleContestant)

		return nil
	}); txErr != nil {
		return nil, txErr
	}

	return response, nil
}

func (t testCase) applyUpdateTestCase(in *rpc.UpdateTestCaseRequest, testCase *db.TestCase) {
	if in.Input != nil {
		testCase.Input = utils.TrimSpaceRight(*in.Input)
	}

	if in.Output != nil {
		testCase.Output = utils.TrimSpaceRight(*in.Output)
	}

	if in.IsHidden != nil {
		testCase.IsHidden = *in.IsHidden
	}

	if in.Input != nil || in.Output != nil {
		testCase.Hash = t.calculateTestCaseHash(testCase.Input, testCase.Output)
	}
}

func (t testCase) UpdateTestCase(
	ctx context.Context,
	in *rpc.UpdateTestCaseRequest,
	token string,
) (*rpc.UpdateTestCaseResponse, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	account, err := t.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	response := &rpc.UpdateTestCaseResponse{}
	if txErr := t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		testCase, testCaseErr := t.testCaseDataAccessor.WithDB(tx).GetTestCase(ctx, in.ID)
		if testCaseErr != nil {
			return testCaseErr
		}

		if testCase == nil {
			logger.With(zap.Uint64("id", in.ID)).Error("cannot find test case")
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
		}

		problem, problemErr := t.problemDataAccessor.WithDB(tx).GetProblem(ctx, testCase.OfProblemID)
		if problemErr != nil {
			return problemErr
		}

		requiredPermissionList := []gorbac.Permission{PermissionTestCasesAllWrite}
		if problem.AuthorAccountID == uint64(account.ID) {
			requiredPermissionList = append(requiredPermissionList, PermissionTestCasesSelfWrite)
		}

		hasPermission, permissionErr := t.role.AccountHasPermission(ctx, string(account.Role), requiredPermissionList...)
		if permissionErr != nil {
			return permissionErr
		}
		if !hasPermission {
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}

		t.applyUpdateTestCase(in, testCase)

		err = utils.ExecuteUntilFirstError(
			func() error { return t.testCaseDataAccessor.WithDB(tx).UpdateTestCase(ctx, testCase) },
			func() error { return t.WithDB(tx).UpsertProblemTestCaseHash(ctx, uint64(problem.ID)) },
		)
		if err != nil {
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
