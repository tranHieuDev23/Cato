package http

import (
	"context"
	"net/http"
	"time"

	"gitlab.com/pjrpc/pjrpc/v2"

	"github.com/tranHieuDev23/cato/internal/handlers/http/middlewares"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc/rpcserver"
	"github.com/tranHieuDev23/cato/internal/logic"
)

type apiServerHandler struct {
	accountLogic logic.Account
}

func NewAPIServerHandler(
	accountLogic logic.Account,
) rpcserver.APIServer {
	return &apiServerHandler{
		accountLogic: accountLogic,
	}
}

func (a apiServerHandler) getAuthorizationBearerToken(ctx context.Context) string {
	contextData, ok := pjrpc.ContextGetData(ctx)
	if !ok {
		return ""
	}

	authorizationCookie, err := contextData.HTTPRequest.Cookie(middlewares.AuthorizationCookie)
	if err != nil {
		return ""
	}

	return authorizationCookie.Value
}

func (a apiServerHandler) setAuthorizationBearerToken(ctx context.Context, token string, expireTime time.Time) {
	contextData, ok := pjrpc.ContextGetData(ctx)
	if !ok {
		return
	}

	contextData.HTTPRequest.AddCookie(&http.Cookie{
		Name:     middlewares.AuthorizationCookie,
		Value:    token,
		HttpOnly: true,
		Expires:  expireTime,
		SameSite: http.SameSiteStrictMode,
	})
}

func (a apiServerHandler) CreateAccount(ctx context.Context, in *rpc.CreateAccountRequest) (*rpc.CreateAccountResponse, error) {
	token := a.getAuthorizationBearerToken(ctx)
	return a.accountLogic.CreateAccount(ctx, in, token)
}

func (a apiServerHandler) CreateProblem(ctx context.Context, in *rpc.CreateProblemRequest) (*rpc.CreateProblemResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) CreateSession(ctx context.Context, in *rpc.CreateSessionRequest) (*rpc.CreateSessionResponse, error) {
	response, token, expireTime, err := a.accountLogic.CreateSession(ctx, in)
	if err != nil {
		return nil, err
	}

	a.setAuthorizationBearerToken(ctx, token, expireTime)
	return response, err
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
	token := a.getAuthorizationBearerToken(ctx)
	if err := a.accountLogic.DeleteSession(ctx, token); err != nil {
		return nil, err
	}

	return &rpc.DeleteSessionResponse{}, nil
}

func (a apiServerHandler) DeleteSubmission(ctx context.Context, in *rpc.DeleteSubmissionRequest) (*rpc.DeleteSubmissionResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) DeleteTestCase(ctx context.Context, in *rpc.DeleteTestCaseRequest) (*rpc.DeleteTestCaseResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) GetAccount(ctx context.Context, in *rpc.GetAccountRequest) (*rpc.GetAccountResponse, error) {
	token := a.getAuthorizationBearerToken(ctx)
	return a.accountLogic.GetAccount(ctx, in, token)
}

func (a apiServerHandler) GetAccountList(ctx context.Context, in *rpc.GetAccountListRequest) (*rpc.GetAccountListResponse, error) {
	token := a.getAuthorizationBearerToken(ctx)
	return a.accountLogic.GetAccountList(ctx, in, token)
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
	token := a.getAuthorizationBearerToken(ctx)
	return a.accountLogic.UpdateAccount(ctx, in, token)
}

func (a apiServerHandler) UpdateProblem(ctx context.Context, in *rpc.UpdateProblemRequest) (*rpc.UpdateProblemResponse, error) {
	panic("unimplemented")
}

func (a apiServerHandler) UpdateTestCase(ctx context.Context, in *rpc.UpdateTestCaseRequest) (*rpc.UpdateTestCaseResponse, error) {
	panic("unimplemented")
}
