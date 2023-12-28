// Code generated by genpjrpc. DO NOT EDIT.
//  genpjrpc version: v0.4.0

package rpcclient

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.com/pjrpc/pjrpc/v2/client"

	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc"
)

// List of the client JSON-RPC methods.
const (
	JSONRPCMethodGetServerInfo_Client                                   = "get_server_info"
	JSONRPCMethodCreateAccount_Client                                   = "create_account"
	JSONRPCMethodGetAccountList_Client                                  = "get_account_list"
	JSONRPCMethodGetAccount_Client                                      = "get_account"
	JSONRPCMethodUpdateAccount_Client                                   = "update_account"
	JSONRPCMethodCreateSession_Client                                   = "create_session"
	JSONRPCMethodGetSession_Client                                      = "get_session"
	JSONRPCMethodDeleteSession_Client                                   = "delete_session"
	JSONRPCMethodCreateProblem_Client                                   = "create_problem"
	JSONRPCMethodGetProblemSnippetList_Client                           = "get_problem_snippet_list"
	JSONRPCMethodGetProblem_Client                                      = "get_problem"
	JSONRPCMethodUpdateProblem_Client                                   = "update_problem"
	JSONRPCMethodDeleteProblem_Client                                   = "delete_problem"
	JSONRPCMethodCreateTestCase_Client                                  = "create_test_case"
	JSONRPCMethodCreateTestCaseList_Client                              = "create_test_case_list"
	JSONRPCMethodGetProblemTestCaseSnippetList_Client                   = "get_problem_test_case_snippet_list"
	JSONRPCMethodGetTestCase_Client                                     = "get_test_case"
	JSONRPCMethodUpdateTestCase_Client                                  = "update_test_case"
	JSONRPCMethodDeleteTestCase_Client                                  = "delete_test_case"
	JSONRPCMethodGetAccountProblemSnippetList_Client                    = "get_account_problem_snippet_list"
	JSONRPCMethodCreateSubmission_Client                                = "create_submission"
	JSONRPCMethodGetSubmissionSnippetList_Client                        = "get_submission_snippet_list"
	JSONRPCMethodGetSubmission_Client                                   = "get_submission"
	JSONRPCMethodUpdateSubmission_Client                                = "update_submission"
	JSONRPCMethodDeleteSubmission_Client                                = "delete_submission"
	JSONRPCMethodGetAccountSubmissionSnippetList_Client                 = "get_account_submission_snippet_list"
	JSONRPCMethodGetProblemSubmissionSnippetList_Client                 = "get_problem_submission_snippet_list"
	JSONRPCMethodGetAccountProblemSubmissionSnippetList_Client          = "get_account_problem_submission_snippet_list"
	JSONRPCMethodGetAndUpdateFirstSubmittedSubmissionToExecuting_Client = "get_and_update_first_submitted_submission_to_executing"
	JSONRPCMethodUpdateSetting_Client                                   = "update_setting"
)

// APIClient is an API client for API service.
type APIClient interface {
	GetServerInfo(ctx context.Context, in *rpc.GetServerInfoRequest, mods ...client.Mod) (*rpc.GetServerInfoResponse, error)
	CreateAccount(ctx context.Context, in *rpc.CreateAccountRequest, mods ...client.Mod) (*rpc.CreateAccountResponse, error)
	GetAccountList(ctx context.Context, in *rpc.GetAccountListRequest, mods ...client.Mod) (*rpc.GetAccountListResponse, error)
	GetAccount(ctx context.Context, in *rpc.GetAccountRequest, mods ...client.Mod) (*rpc.GetAccountResponse, error)
	UpdateAccount(ctx context.Context, in *rpc.UpdateAccountRequest, mods ...client.Mod) (*rpc.UpdateAccountResponse, error)
	CreateSession(ctx context.Context, in *rpc.CreateSessionRequest, mods ...client.Mod) (*rpc.CreateSessionResponse, error)
	GetSession(ctx context.Context, in *rpc.GetSessionRequest, mods ...client.Mod) (*rpc.GetSessionResponse, error)
	DeleteSession(ctx context.Context, in *rpc.DeleteSessionRequest, mods ...client.Mod) (*rpc.DeleteSessionResponse, error)
	CreateProblem(ctx context.Context, in *rpc.CreateProblemRequest, mods ...client.Mod) (*rpc.CreateProblemResponse, error)
	GetProblemSnippetList(ctx context.Context, in *rpc.GetProblemSnippetListRequest, mods ...client.Mod) (*rpc.GetProblemSnippetListResponse, error)
	GetProblem(ctx context.Context, in *rpc.GetProblemRequest, mods ...client.Mod) (*rpc.GetProblemResponse, error)
	UpdateProblem(ctx context.Context, in *rpc.UpdateProblemRequest, mods ...client.Mod) (*rpc.UpdateProblemResponse, error)
	DeleteProblem(ctx context.Context, in *rpc.DeleteProblemRequest, mods ...client.Mod) (*rpc.DeleteProblemResponse, error)
	CreateTestCase(ctx context.Context, in *rpc.CreateTestCaseRequest, mods ...client.Mod) (*rpc.CreateTestCaseResponse, error)
	CreateTestCaseList(ctx context.Context, in *rpc.CreateTestCaseListRequest, mods ...client.Mod) (*rpc.CreateTestCaseListResponse, error)
	GetProblemTestCaseSnippetList(ctx context.Context, in *rpc.GetProblemTestCaseSnippetListRequest, mods ...client.Mod) (*rpc.GetProblemTestCaseSnippetListResponse, error)
	GetTestCase(ctx context.Context, in *rpc.GetTestCaseRequest, mods ...client.Mod) (*rpc.GetTestCaseResponse, error)
	UpdateTestCase(ctx context.Context, in *rpc.UpdateTestCaseRequest, mods ...client.Mod) (*rpc.UpdateTestCaseResponse, error)
	DeleteTestCase(ctx context.Context, in *rpc.DeleteTestCaseRequest, mods ...client.Mod) (*rpc.DeleteTestCaseResponse, error)
	GetAccountProblemSnippetList(ctx context.Context, in *rpc.GetAccountProblemSnippetListRequest, mods ...client.Mod) (*rpc.GetAccountProblemSnippetListResponse, error)
	CreateSubmission(ctx context.Context, in *rpc.CreateSubmissionRequest, mods ...client.Mod) (*rpc.CreateSubmissionResponse, error)
	GetSubmissionSnippetList(ctx context.Context, in *rpc.GetSubmissionSnippetListRequest, mods ...client.Mod) (*rpc.GetSubmissionSnippetListResponse, error)
	GetSubmission(ctx context.Context, in *rpc.GetSubmissionRequest, mods ...client.Mod) (*rpc.GetSubmissionResponse, error)
	UpdateSubmission(ctx context.Context, in *rpc.UpdateSubmissionRequest, mods ...client.Mod) (*rpc.UpdateSubmissionResponse, error)
	DeleteSubmission(ctx context.Context, in *rpc.DeleteSubmissionRequest, mods ...client.Mod) (*rpc.DeleteSubmissionResponse, error)
	GetAccountSubmissionSnippetList(ctx context.Context, in *rpc.GetAccountSubmissionSnippetListRequest, mods ...client.Mod) (*rpc.GetAccountSubmissionSnippetListResponse, error)
	GetProblemSubmissionSnippetList(ctx context.Context, in *rpc.GetProblemSubmissionSnippetListRequest, mods ...client.Mod) (*rpc.GetProblemSubmissionSnippetListResponse, error)
	GetAccountProblemSubmissionSnippetList(ctx context.Context, in *rpc.GetAccountProblemSubmissionSnippetListRequest, mods ...client.Mod) (*rpc.GetAccountProblemSubmissionSnippetListResponse, error)
	GetAndUpdateFirstSubmittedSubmissionToExecuting(ctx context.Context, in *rpc.GetAndUpdateFirstSubmittedSubmissionToExecutingRequest, mods ...client.Mod) (*rpc.GetAndUpdateFirstSubmittedSubmissionToExecutingResponse, error)
	UpdateSetting(ctx context.Context, in *rpc.UpdateSettingRequest, mods ...client.Mod) (*rpc.UpdateSettingResponse, error)
}

type implAPIClient struct {
	cl client.Invoker
}

// NewAPIClient returns new client implementation of the API service.
func NewAPIClient(cl client.Invoker) APIClient {
	return &implAPIClient{cl: cl}
}

func (c *implAPIClient) GetServerInfo(ctx context.Context, in *rpc.GetServerInfoRequest, mods ...client.Mod) (result *rpc.GetServerInfoResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.GetServerInfoResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodGetServerInfo_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodGetServerInfo_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) CreateAccount(ctx context.Context, in *rpc.CreateAccountRequest, mods ...client.Mod) (result *rpc.CreateAccountResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.CreateAccountResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodCreateAccount_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodCreateAccount_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) GetAccountList(ctx context.Context, in *rpc.GetAccountListRequest, mods ...client.Mod) (result *rpc.GetAccountListResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.GetAccountListResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodGetAccountList_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodGetAccountList_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) GetAccount(ctx context.Context, in *rpc.GetAccountRequest, mods ...client.Mod) (result *rpc.GetAccountResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.GetAccountResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodGetAccount_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodGetAccount_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) UpdateAccount(ctx context.Context, in *rpc.UpdateAccountRequest, mods ...client.Mod) (result *rpc.UpdateAccountResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.UpdateAccountResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodUpdateAccount_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodUpdateAccount_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) CreateSession(ctx context.Context, in *rpc.CreateSessionRequest, mods ...client.Mod) (result *rpc.CreateSessionResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.CreateSessionResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodCreateSession_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodCreateSession_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) GetSession(ctx context.Context, in *rpc.GetSessionRequest, mods ...client.Mod) (result *rpc.GetSessionResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.GetSessionResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodGetSession_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodGetSession_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) DeleteSession(ctx context.Context, in *rpc.DeleteSessionRequest, mods ...client.Mod) (result *rpc.DeleteSessionResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.DeleteSessionResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodDeleteSession_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodDeleteSession_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) CreateProblem(ctx context.Context, in *rpc.CreateProblemRequest, mods ...client.Mod) (result *rpc.CreateProblemResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.CreateProblemResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodCreateProblem_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodCreateProblem_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) GetProblemSnippetList(ctx context.Context, in *rpc.GetProblemSnippetListRequest, mods ...client.Mod) (result *rpc.GetProblemSnippetListResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.GetProblemSnippetListResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodGetProblemSnippetList_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodGetProblemSnippetList_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) GetProblem(ctx context.Context, in *rpc.GetProblemRequest, mods ...client.Mod) (result *rpc.GetProblemResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.GetProblemResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodGetProblem_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodGetProblem_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) UpdateProblem(ctx context.Context, in *rpc.UpdateProblemRequest, mods ...client.Mod) (result *rpc.UpdateProblemResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.UpdateProblemResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodUpdateProblem_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodUpdateProblem_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) DeleteProblem(ctx context.Context, in *rpc.DeleteProblemRequest, mods ...client.Mod) (result *rpc.DeleteProblemResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.DeleteProblemResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodDeleteProblem_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodDeleteProblem_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) CreateTestCase(ctx context.Context, in *rpc.CreateTestCaseRequest, mods ...client.Mod) (result *rpc.CreateTestCaseResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.CreateTestCaseResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodCreateTestCase_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodCreateTestCase_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) CreateTestCaseList(ctx context.Context, in *rpc.CreateTestCaseListRequest, mods ...client.Mod) (result *rpc.CreateTestCaseListResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.CreateTestCaseListResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodCreateTestCaseList_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodCreateTestCaseList_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) GetProblemTestCaseSnippetList(ctx context.Context, in *rpc.GetProblemTestCaseSnippetListRequest, mods ...client.Mod) (result *rpc.GetProblemTestCaseSnippetListResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.GetProblemTestCaseSnippetListResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodGetProblemTestCaseSnippetList_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodGetProblemTestCaseSnippetList_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) GetTestCase(ctx context.Context, in *rpc.GetTestCaseRequest, mods ...client.Mod) (result *rpc.GetTestCaseResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.GetTestCaseResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodGetTestCase_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodGetTestCase_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) UpdateTestCase(ctx context.Context, in *rpc.UpdateTestCaseRequest, mods ...client.Mod) (result *rpc.UpdateTestCaseResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.UpdateTestCaseResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodUpdateTestCase_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodUpdateTestCase_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) DeleteTestCase(ctx context.Context, in *rpc.DeleteTestCaseRequest, mods ...client.Mod) (result *rpc.DeleteTestCaseResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.DeleteTestCaseResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodDeleteTestCase_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodDeleteTestCase_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) GetAccountProblemSnippetList(ctx context.Context, in *rpc.GetAccountProblemSnippetListRequest, mods ...client.Mod) (result *rpc.GetAccountProblemSnippetListResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.GetAccountProblemSnippetListResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodGetAccountProblemSnippetList_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodGetAccountProblemSnippetList_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) CreateSubmission(ctx context.Context, in *rpc.CreateSubmissionRequest, mods ...client.Mod) (result *rpc.CreateSubmissionResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.CreateSubmissionResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodCreateSubmission_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodCreateSubmission_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) GetSubmissionSnippetList(ctx context.Context, in *rpc.GetSubmissionSnippetListRequest, mods ...client.Mod) (result *rpc.GetSubmissionSnippetListResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.GetSubmissionSnippetListResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodGetSubmissionSnippetList_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodGetSubmissionSnippetList_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) GetSubmission(ctx context.Context, in *rpc.GetSubmissionRequest, mods ...client.Mod) (result *rpc.GetSubmissionResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.GetSubmissionResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodGetSubmission_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodGetSubmission_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) UpdateSubmission(ctx context.Context, in *rpc.UpdateSubmissionRequest, mods ...client.Mod) (result *rpc.UpdateSubmissionResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.UpdateSubmissionResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodUpdateSubmission_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodUpdateSubmission_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) DeleteSubmission(ctx context.Context, in *rpc.DeleteSubmissionRequest, mods ...client.Mod) (result *rpc.DeleteSubmissionResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.DeleteSubmissionResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodDeleteSubmission_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodDeleteSubmission_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) GetAccountSubmissionSnippetList(ctx context.Context, in *rpc.GetAccountSubmissionSnippetListRequest, mods ...client.Mod) (result *rpc.GetAccountSubmissionSnippetListResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.GetAccountSubmissionSnippetListResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodGetAccountSubmissionSnippetList_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodGetAccountSubmissionSnippetList_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) GetProblemSubmissionSnippetList(ctx context.Context, in *rpc.GetProblemSubmissionSnippetListRequest, mods ...client.Mod) (result *rpc.GetProblemSubmissionSnippetListResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.GetProblemSubmissionSnippetListResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodGetProblemSubmissionSnippetList_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodGetProblemSubmissionSnippetList_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) GetAccountProblemSubmissionSnippetList(ctx context.Context, in *rpc.GetAccountProblemSubmissionSnippetListRequest, mods ...client.Mod) (result *rpc.GetAccountProblemSubmissionSnippetListResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.GetAccountProblemSubmissionSnippetListResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodGetAccountProblemSubmissionSnippetList_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodGetAccountProblemSubmissionSnippetList_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) GetAndUpdateFirstSubmittedSubmissionToExecuting(ctx context.Context, in *rpc.GetAndUpdateFirstSubmittedSubmissionToExecutingRequest, mods ...client.Mod) (result *rpc.GetAndUpdateFirstSubmittedSubmissionToExecutingResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.GetAndUpdateFirstSubmittedSubmissionToExecutingResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodGetAndUpdateFirstSubmittedSubmissionToExecuting_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodGetAndUpdateFirstSubmittedSubmissionToExecuting_Client, err)
	}

	return result, nil
}

func (c *implAPIClient) UpdateSetting(ctx context.Context, in *rpc.UpdateSettingRequest, mods ...client.Mod) (result *rpc.UpdateSettingResponse, err error) {
	gen, err := uuid.NewUUID()
	if err != nil {
		return result, fmt.Errorf("failed to create uuid generator: %w", err)
	}

	result = new(rpc.UpdateSettingResponse)

	err = c.cl.Invoke(ctx, gen.String(), JSONRPCMethodUpdateSetting_Client, in, result, mods...)
	if err != nil {
		return result, fmt.Errorf("failed to Invoke method %q: %w", JSONRPCMethodUpdateSetting_Client, err)
	}

	return result, nil
}
