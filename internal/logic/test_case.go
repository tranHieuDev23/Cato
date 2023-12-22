package logic

import (
	"context"

	"go.uber.org/zap"

	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc"
)

type TestCase interface {
	CreateTestCase(ctx context.Context, in *rpc.CreateTestCaseRequest) (*rpc.CreateTestCaseResponse, error)
	CreateTestCaseList(ctx context.Context, in *rpc.CreateTestCaseListRequest) (*rpc.CreateTestCaseListResponse, error)
	GetProblemTestCaseSnippetList(ctx context.Context, in *rpc.GetProblemTestCaseSnippetListRequest) (*rpc.GetProblemTestCaseSnippetListResponse, error)
	GetTestCase(ctx context.Context, in *rpc.GetTestCaseRequest) (*rpc.GetTestCaseResponse, error)
	UpdateTestCase(ctx context.Context, in *rpc.UpdateTestCaseRequest) (*rpc.UpdateTestCaseResponse, error)
	DeleteTestCase(ctx context.Context, in *rpc.DeleteTestCaseRequest) (*rpc.DeleteTestCaseResponse, error)
}

type testCase struct {
	testCaseDataAccessor db.TestCaseDataAccessor
	logger               *zap.Logger
}

func NewTestCase(
	testCaseDataAccessor db.TestCaseDataAccessor,
	logger *zap.Logger,
) TestCase {
	return &testCase{
		testCaseDataAccessor: testCaseDataAccessor,
		logger:               logger,
	}
}

func (t testCase) CreateTestCase(ctx context.Context, in *rpc.CreateTestCaseRequest) (*rpc.CreateTestCaseResponse, error) {
	panic("unimplemented")
}

func (t testCase) CreateTestCaseList(ctx context.Context, in *rpc.CreateTestCaseListRequest) (*rpc.CreateTestCaseListResponse, error) {
	panic("unimplemented")
}

func (t testCase) DeleteTestCase(ctx context.Context, in *rpc.DeleteTestCaseRequest) (*rpc.DeleteTestCaseResponse, error) {
	panic("unimplemented")
}

func (t testCase) GetProblemTestCaseSnippetList(ctx context.Context, in *rpc.GetProblemTestCaseSnippetListRequest) (*rpc.GetProblemTestCaseSnippetListResponse, error) {
	panic("unimplemented")
}

func (t testCase) GetTestCase(ctx context.Context, in *rpc.GetTestCaseRequest) (*rpc.GetTestCaseResponse, error) {
	panic("unimplemented")
}

func (t testCase) UpdateTestCase(ctx context.Context, in *rpc.UpdateTestCaseRequest) (*rpc.UpdateTestCaseResponse, error) {
	panic("unimplemented")
}
