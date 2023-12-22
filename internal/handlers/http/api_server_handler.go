package http

import (
	"context"

	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc/rpcserver"
)

type apiServerHandler struct {
}

func NewAPIServerHandler() rpcserver.APIServer {
	return &apiServerHandler{}
}

func (a apiServerHandler) CreateAccount(ctx context.Context, in *rpc.CreateAccountRequest) (*rpc.CreateAccountResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) CreateProblem(ctx context.Context, in *rpc.CreateProblemRequest) (*rpc.CreateProblemResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) CreateSession(ctx context.Context, in *rpc.CreateSessionRequest) (*rpc.CreateSessionResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) CreateSubmission(ctx context.Context, in *rpc.CreateSubmissionRequest) (*rpc.CreateSubmissionResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) CreateTestCase(ctx context.Context, in *rpc.CreateTestCaseRequest) (*rpc.CreateTestCaseResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) CreateTestCaseList(ctx context.Context, in *rpc.CreateTestCaseListRequest) (*rpc.CreateTestCaseListResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) DeleteProblem(ctx context.Context, in *rpc.DeleteProblemRequest) (*rpc.DeleteProblemResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) DeleteSession(ctx context.Context, in *rpc.DeleteSessionRequest) (*rpc.DeleteSessionResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) DeleteSubmission(ctx context.Context, in *rpc.DeleteSubmissionRequest) (*rpc.DeleteSubmissionResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) DeleteTestCase(ctx context.Context, in *rpc.DeleteTestCaseRequest) (*rpc.DeleteTestCaseResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) GetAccount(ctx context.Context, in *rpc.GetAccountRequest) (*rpc.GetAccountResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) GetAccountList(ctx context.Context, in *rpc.GetAccountListRequest) (*rpc.GetAccountListResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) GetAccountProblemSnippetList(ctx context.Context, in *rpc.GetAccountProblemSnippetListRequest) (*rpc.GetAccountProblemSnippetListResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) GetAccountSubmissionSnippetList(ctx context.Context, in *rpc.GetAccountSubmissionSnippetListRequest) (*rpc.GetAccountSubmissionSnippetListResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) GetProblem(ctx context.Context, in *rpc.GetProblemRequest) (*rpc.GetProblemResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) GetProblemSnippetList(ctx context.Context, in *rpc.GetProblemSnippetListRequest) (*rpc.GetProblemSnippetListResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) GetProblemSubmissionSnippetList(ctx context.Context, in *rpc.GetProblemSubmissionSnippetListRequest) (*rpc.GetProblemSubmissionSnippetListResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) GetProblemTestCaseSnippetList(ctx context.Context, in *rpc.GetProblemTestCaseSnippetListRequest) (*rpc.GetProblemTestCaseSnippetListResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) GetSubmission(ctx context.Context, in *rpc.GetSubmissionRequest) (*rpc.GetSubmissionResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) GetSubmissionSnippetList(ctx context.Context, in *rpc.GetSubmissionSnippetListRequest) (*rpc.GetSubmissionSnippetListResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) GetTestCase(ctx context.Context, in *rpc.GetTestCaseRequest) (*rpc.GetTestCaseResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) UpdateAccount(ctx context.Context, in *rpc.UpdateAccountRequest) (*rpc.UpdateAccountResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) UpdateProblem(ctx context.Context, in *rpc.UpdateProblemRequest) (*rpc.UpdateProblemResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) UpdateTestCase(ctx context.Context, in *rpc.UpdateTestCaseRequest) (*rpc.UpdateTestCaseResponse, error) {
	panic("unimplemented")
}
