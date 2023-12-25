package http

import (
	"context"
	"net/http"
	"time"

	"gitlab.com/pjrpc/pjrpc/v2"
	"go.uber.org/zap"

	validator "github.com/go-playground/validator/v10"

	"github.com/tranHieuDev23/cato/internal/handlers/http/middlewares"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc/rpcserver"
	"github.com/tranHieuDev23/cato/internal/logic"
	"github.com/tranHieuDev23/cato/internal/utils"
)

type apiServerHandler struct {
	accountLogic    logic.Account
	problemLogic    logic.Problem
	testCaseLogic   logic.TestCase
	submissionLogic logic.Submission
	logger          *zap.Logger
	validate        *validator.Validate
}

func NewAPIServerHandler(
	accountLogic logic.Account,
	problemLogic logic.Problem,
	testCaseLogic logic.TestCase,
	submissionLogic logic.Submission,
	logger *zap.Logger,
) rpcserver.APIServer {
	validate := validator.New()
	return &apiServerHandler{
		accountLogic:    accountLogic,
		problemLogic:    problemLogic,
		testCaseLogic:   testCaseLogic,
		submissionLogic: submissionLogic,
		logger:          logger,
		validate:        validate,
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

	contextData.HTTPRequest.Header.Del("Cookie")
	contextData.HTTPRequest.AddCookie(&http.Cookie{
		Name:     middlewares.AuthorizationCookie,
		Value:    token,
		HttpOnly: true,
		Expires:  expireTime,
		SameSite: http.SameSiteStrictMode,
	})
}

func (a apiServerHandler) validateRequest(
	ctx context.Context,
	in any,
) error {
	logger := utils.LoggerWithContext(ctx, a.logger)
	if err := a.validate.StructCtx(ctx, in); err != nil {
		logger.With(zap.Error(err)).Error("invalid request params")
		return pjrpc.JRPCErrInvalidParams()
	}

	return nil
}

func (a apiServerHandler) CreateAccount(
	ctx context.Context,
	in *rpc.CreateAccountRequest,
) (*rpc.CreateAccountResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	token := a.getAuthorizationBearerToken(ctx)
	return a.accountLogic.CreateAccount(ctx, in, token)
}

func (a apiServerHandler) CreateProblem(
	ctx context.Context,
	in *rpc.CreateProblemRequest,
) (*rpc.CreateProblemResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	token := a.getAuthorizationBearerToken(ctx)
	return a.problemLogic.CreateProblem(ctx, in, token)
}

func (a apiServerHandler) CreateSession(
	ctx context.Context,
	in *rpc.CreateSessionRequest,
) (*rpc.CreateSessionResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	response, token, expireTime, err := a.accountLogic.CreateSession(ctx, in)
	if err != nil {
		return nil, err
	}

	a.setAuthorizationBearerToken(ctx, token, expireTime)
	return response, err
}

func (a apiServerHandler) CreateSubmission(
	ctx context.Context,
	in *rpc.CreateSubmissionRequest,
) (*rpc.CreateSubmissionResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	token := a.getAuthorizationBearerToken(ctx)
	return a.submissionLogic.CreateSubmission(ctx, in, token)
}

func (a apiServerHandler) CreateTestCase(
	ctx context.Context,
	in *rpc.CreateTestCaseRequest,
) (*rpc.CreateTestCaseResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	token := a.getAuthorizationBearerToken(ctx)
	return a.testCaseLogic.CreateTestCase(ctx, in, token)
}

func (a apiServerHandler) CreateTestCaseList(
	ctx context.Context,
	in *rpc.CreateTestCaseListRequest,
) (*rpc.CreateTestCaseListResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	token := a.getAuthorizationBearerToken(ctx)
	if err := a.testCaseLogic.CreateTestCaseList(ctx, in, token); err != nil {
		return nil, err
	}

	return &rpc.CreateTestCaseListResponse{}, nil
}

func (a apiServerHandler) DeleteProblem(
	ctx context.Context,
	in *rpc.DeleteProblemRequest,
) (*rpc.DeleteProblemResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	token := a.getAuthorizationBearerToken(ctx)
	if err := a.problemLogic.DeleteProblem(ctx, in, token); err != nil {
		return nil, err
	}

	return &rpc.DeleteProblemResponse{}, nil
}

func (a apiServerHandler) DeleteSession(
	ctx context.Context,
	in *rpc.DeleteSessionRequest,
) (*rpc.DeleteSessionResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	token := a.getAuthorizationBearerToken(ctx)
	if err := a.accountLogic.DeleteSession(ctx, token); err != nil {
		return nil, err
	}

	a.setAuthorizationBearerToken(ctx, "", time.Unix(0, 0))
	return &rpc.DeleteSessionResponse{}, nil
}

func (a apiServerHandler) DeleteSubmission(
	ctx context.Context,
	in *rpc.DeleteSubmissionRequest,
) (*rpc.DeleteSubmissionResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	token := a.getAuthorizationBearerToken(ctx)
	if err := a.submissionLogic.DeleteSubmission(ctx, in, token); err != nil {
		return nil, err
	}

	return &rpc.DeleteSubmissionResponse{}, nil
}

func (a apiServerHandler) DeleteTestCase(
	ctx context.Context,
	in *rpc.DeleteTestCaseRequest,
) (*rpc.DeleteTestCaseResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	token := a.getAuthorizationBearerToken(ctx)
	if err := a.testCaseLogic.DeleteTestCase(ctx, in, token); err != nil {
		return nil, err
	}

	return &rpc.DeleteTestCaseResponse{}, nil
}

func (a apiServerHandler) GetAccount(
	ctx context.Context,
	in *rpc.GetAccountRequest,
) (*rpc.GetAccountResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	token := a.getAuthorizationBearerToken(ctx)
	return a.accountLogic.GetAccount(ctx, in, token)
}

func (a apiServerHandler) GetAccountList(
	ctx context.Context,
	in *rpc.GetAccountListRequest,
) (*rpc.GetAccountListResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	token := a.getAuthorizationBearerToken(ctx)
	return a.accountLogic.GetAccountList(ctx, in, token)
}

func (a apiServerHandler) GetAccountProblemSnippetList(
	ctx context.Context,
	in *rpc.GetAccountProblemSnippetListRequest,
) (*rpc.GetAccountProblemSnippetListResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	token := a.getAuthorizationBearerToken(ctx)
	return a.problemLogic.GetAccountProblemSnippetList(ctx, in, token)
}

func (a apiServerHandler) GetAccountProblemSubmissionSnippetList(
	ctx context.Context,
	in *rpc.GetAccountProblemSubmissionSnippetListRequest,
) (*rpc.GetAccountProblemSubmissionSnippetListResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	token := a.getAuthorizationBearerToken(ctx)
	return a.submissionLogic.GetAccountProblemSubmissionSnippetList(ctx, in, token)
}

func (a apiServerHandler) GetAccountSubmissionSnippetList(
	ctx context.Context,
	in *rpc.GetAccountSubmissionSnippetListRequest,
) (*rpc.GetAccountSubmissionSnippetListResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	token := a.getAuthorizationBearerToken(ctx)
	return a.submissionLogic.GetAccountSubmissionSnippetList(ctx, in, token)
}

func (a apiServerHandler) GetSession(
	ctx context.Context,
	in *rpc.GetSessionRequest,
) (*rpc.GetSessionResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	token := a.getAuthorizationBearerToken(ctx)
	return a.accountLogic.GetSession(ctx, token)
}

func (a apiServerHandler) GetProblem(
	ctx context.Context,
	in *rpc.GetProblemRequest,
) (*rpc.GetProblemResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	token := a.getAuthorizationBearerToken(ctx)
	return a.problemLogic.GetProblem(ctx, in, token)
}

func (a apiServerHandler) GetProblemSnippetList(
	ctx context.Context,
	in *rpc.GetProblemSnippetListRequest,
) (*rpc.GetProblemSnippetListResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	token := a.getAuthorizationBearerToken(ctx)
	return a.problemLogic.GetProblemSnippetList(ctx, in, token)
}

func (a apiServerHandler) GetProblemSubmissionSnippetList(
	ctx context.Context,
	in *rpc.GetProblemSubmissionSnippetListRequest,
) (*rpc.GetProblemSubmissionSnippetListResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	token := a.getAuthorizationBearerToken(ctx)
	return a.submissionLogic.GetProblemSubmissionSnippetList(ctx, in, token)
}

func (a apiServerHandler) GetProblemTestCaseSnippetList(
	ctx context.Context,
	in *rpc.GetProblemTestCaseSnippetListRequest,
) (*rpc.GetProblemTestCaseSnippetListResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	token := a.getAuthorizationBearerToken(ctx)
	return a.testCaseLogic.GetProblemTestCaseSnippetList(ctx, in, token)
}

func (a apiServerHandler) GetSubmission(
	ctx context.Context,
	in *rpc.GetSubmissionRequest,
) (*rpc.GetSubmissionResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	token := a.getAuthorizationBearerToken(ctx)
	return a.submissionLogic.GetSubmission(ctx, in, token)
}

func (a apiServerHandler) GetSubmissionSnippetList(
	ctx context.Context,
	in *rpc.GetSubmissionSnippetListRequest,
) (*rpc.GetSubmissionSnippetListResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	token := a.getAuthorizationBearerToken(ctx)
	return a.submissionLogic.GetSubmissionSnippetList(ctx, in, token)
}

func (a apiServerHandler) GetTestCase(
	ctx context.Context,
	in *rpc.GetTestCaseRequest,
) (*rpc.GetTestCaseResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	token := a.getAuthorizationBearerToken(ctx)
	return a.testCaseLogic.GetTestCase(ctx, in, token)
}

func (a apiServerHandler) UpdateAccount(
	ctx context.Context,
	in *rpc.UpdateAccountRequest,
) (*rpc.UpdateAccountResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	token := a.getAuthorizationBearerToken(ctx)
	return a.accountLogic.UpdateAccount(ctx, in, token)
}

func (a apiServerHandler) UpdateProblem(
	ctx context.Context,
	in *rpc.UpdateProblemRequest,
) (*rpc.UpdateProblemResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	token := a.getAuthorizationBearerToken(ctx)
	return a.problemLogic.UpdateProblem(ctx, in, token)
}

func (a apiServerHandler) UpdateTestCase(
	ctx context.Context,
	in *rpc.UpdateTestCaseRequest,
) (*rpc.UpdateTestCaseResponse, error) {
	if err := a.validateRequest(ctx, in); err != nil {
		return nil, err
	}

	token := a.getAuthorizationBearerToken(ctx)
	return a.testCaseLogic.UpdateTestCase(ctx, in, token)
}

type LocalAPIServerHandler rpcserver.APIServer

func NewLocalAPIServerHandler(
	accountLogic logic.Account,
	problemLogic logic.Problem,
	testCaseLogic logic.TestCase,
	submissionLogic logic.LocalSubmission,
	logger *zap.Logger,
) LocalAPIServerHandler {
	return NewAPIServerHandler(accountLogic, problemLogic, testCaseLogic, submissionLogic, logger)
}

type DistributedAPIServerHandler rpcserver.APIServer

func NewDistributedAPIServerHandler(
	accountLogic logic.Account,
	problemLogic logic.Problem,
	testCaseLogic logic.TestCase,
	submissionLogic logic.DistributedSubmission,
	logger *zap.Logger,
) DistributedAPIServerHandler {
	return NewAPIServerHandler(accountLogic, problemLogic, testCaseLogic, submissionLogic, logger)
}
