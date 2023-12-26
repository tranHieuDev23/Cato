// Code generated by genpjrpc. DO NOT EDIT.
//  genpjrpc version: v0.4.0

package rpcserver

import (
	"context"
	"encoding/json"
	"fmt"

	pjrpc "gitlab.com/pjrpc/pjrpc/v2"
	"gitlab.com/pjrpc/pjrpc/v2/pjson"

	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc"
)

// List of the server JSON-RPC methods.
const (
	JSONRPCMethodGetServerInfo                          = "get_server_info"
	JSONRPCMethodCreateAccount                          = "create_account"
	JSONRPCMethodGetAccountList                         = "get_account_list"
	JSONRPCMethodGetAccount                             = "get_account"
	JSONRPCMethodUpdateAccount                          = "update_account"
	JSONRPCMethodCreateSession                          = "create_session"
	JSONRPCMethodGetSession                             = "get_session"
	JSONRPCMethodDeleteSession                          = "delete_session"
	JSONRPCMethodCreateProblem                          = "create_problem"
	JSONRPCMethodGetProblemSnippetList                  = "get_problem_snippet_list"
	JSONRPCMethodGetProblem                             = "get_problem"
	JSONRPCMethodUpdateProblem                          = "update_problem"
	JSONRPCMethodDeleteProblem                          = "delete_problem"
	JSONRPCMethodCreateTestCase                         = "create_test_case"
	JSONRPCMethodCreateTestCaseList                     = "create_test_case_list"
	JSONRPCMethodGetProblemTestCaseSnippetList          = "get_problem_test_case_snippet_list"
	JSONRPCMethodGetTestCase                            = "get_test_case"
	JSONRPCMethodUpdateTestCase                         = "update_test_case"
	JSONRPCMethodDeleteTestCase                         = "delete_test_case"
	JSONRPCMethodGetAccountProblemSnippetList           = "get_account_problem_snippet_list"
	JSONRPCMethodCreateSubmission                       = "create_submission"
	JSONRPCMethodGetSubmissionSnippetList               = "get_submission_snippet_list"
	JSONRPCMethodGetSubmission                          = "get_submission"
	JSONRPCMethodDeleteSubmission                       = "delete_submission"
	JSONRPCMethodGetAccountSubmissionSnippetList        = "get_account_submission_snippet_list"
	JSONRPCMethodGetProblemSubmissionSnippetList        = "get_problem_submission_snippet_list"
	JSONRPCMethodGetAccountProblemSubmissionSnippetList = "get_account_problem_submission_snippet_list"
)

// APIServer is an API server for API service.
type APIServer interface {
	GetServerInfo(ctx context.Context, in *rpc.GetServerInfoRequest) (*rpc.GetServerInfoResponse, error)
	CreateAccount(ctx context.Context, in *rpc.CreateAccountRequest) (*rpc.CreateAccountResponse, error)
	GetAccountList(ctx context.Context, in *rpc.GetAccountListRequest) (*rpc.GetAccountListResponse, error)
	GetAccount(ctx context.Context, in *rpc.GetAccountRequest) (*rpc.GetAccountResponse, error)
	UpdateAccount(ctx context.Context, in *rpc.UpdateAccountRequest) (*rpc.UpdateAccountResponse, error)
	CreateSession(ctx context.Context, in *rpc.CreateSessionRequest) (*rpc.CreateSessionResponse, error)
	GetSession(ctx context.Context, in *rpc.GetSessionRequest) (*rpc.GetSessionResponse, error)
	DeleteSession(ctx context.Context, in *rpc.DeleteSessionRequest) (*rpc.DeleteSessionResponse, error)
	CreateProblem(ctx context.Context, in *rpc.CreateProblemRequest) (*rpc.CreateProblemResponse, error)
	GetProblemSnippetList(ctx context.Context, in *rpc.GetProblemSnippetListRequest) (*rpc.GetProblemSnippetListResponse, error)
	GetProblem(ctx context.Context, in *rpc.GetProblemRequest) (*rpc.GetProblemResponse, error)
	UpdateProblem(ctx context.Context, in *rpc.UpdateProblemRequest) (*rpc.UpdateProblemResponse, error)
	DeleteProblem(ctx context.Context, in *rpc.DeleteProblemRequest) (*rpc.DeleteProblemResponse, error)
	CreateTestCase(ctx context.Context, in *rpc.CreateTestCaseRequest) (*rpc.CreateTestCaseResponse, error)
	CreateTestCaseList(ctx context.Context, in *rpc.CreateTestCaseListRequest) (*rpc.CreateTestCaseListResponse, error)
	GetProblemTestCaseSnippetList(ctx context.Context, in *rpc.GetProblemTestCaseSnippetListRequest) (*rpc.GetProblemTestCaseSnippetListResponse, error)
	GetTestCase(ctx context.Context, in *rpc.GetTestCaseRequest) (*rpc.GetTestCaseResponse, error)
	UpdateTestCase(ctx context.Context, in *rpc.UpdateTestCaseRequest) (*rpc.UpdateTestCaseResponse, error)
	DeleteTestCase(ctx context.Context, in *rpc.DeleteTestCaseRequest) (*rpc.DeleteTestCaseResponse, error)
	GetAccountProblemSnippetList(ctx context.Context, in *rpc.GetAccountProblemSnippetListRequest) (*rpc.GetAccountProblemSnippetListResponse, error)
	CreateSubmission(ctx context.Context, in *rpc.CreateSubmissionRequest) (*rpc.CreateSubmissionResponse, error)
	GetSubmissionSnippetList(ctx context.Context, in *rpc.GetSubmissionSnippetListRequest) (*rpc.GetSubmissionSnippetListResponse, error)
	GetSubmission(ctx context.Context, in *rpc.GetSubmissionRequest) (*rpc.GetSubmissionResponse, error)
	DeleteSubmission(ctx context.Context, in *rpc.DeleteSubmissionRequest) (*rpc.DeleteSubmissionResponse, error)
	GetAccountSubmissionSnippetList(ctx context.Context, in *rpc.GetAccountSubmissionSnippetListRequest) (*rpc.GetAccountSubmissionSnippetListResponse, error)
	GetProblemSubmissionSnippetList(ctx context.Context, in *rpc.GetProblemSubmissionSnippetListRequest) (*rpc.GetProblemSubmissionSnippetListResponse, error)
	GetAccountProblemSubmissionSnippetList(ctx context.Context, in *rpc.GetAccountProblemSubmissionSnippetListRequest) (*rpc.GetAccountProblemSubmissionSnippetListResponse, error)
}

type regAPI struct {
	svc APIServer
}

// RegisterAPIServer registers rpc handlers with middlewares in the server router.
func RegisterAPIServer(srv pjrpc.Registrator, svc APIServer, middlewares ...pjrpc.Middleware) {
	r := &regAPI{svc: svc}

	srv.RegisterMethod(JSONRPCMethodGetServerInfo, r.regGetServerInfo)
	srv.RegisterMethod(JSONRPCMethodCreateAccount, r.regCreateAccount)
	srv.RegisterMethod(JSONRPCMethodGetAccountList, r.regGetAccountList)
	srv.RegisterMethod(JSONRPCMethodGetAccount, r.regGetAccount)
	srv.RegisterMethod(JSONRPCMethodUpdateAccount, r.regUpdateAccount)
	srv.RegisterMethod(JSONRPCMethodCreateSession, r.regCreateSession)
	srv.RegisterMethod(JSONRPCMethodGetSession, r.regGetSession)
	srv.RegisterMethod(JSONRPCMethodDeleteSession, r.regDeleteSession)
	srv.RegisterMethod(JSONRPCMethodCreateProblem, r.regCreateProblem)
	srv.RegisterMethod(JSONRPCMethodGetProblemSnippetList, r.regGetProblemSnippetList)
	srv.RegisterMethod(JSONRPCMethodGetProblem, r.regGetProblem)
	srv.RegisterMethod(JSONRPCMethodUpdateProblem, r.regUpdateProblem)
	srv.RegisterMethod(JSONRPCMethodDeleteProblem, r.regDeleteProblem)
	srv.RegisterMethod(JSONRPCMethodCreateTestCase, r.regCreateTestCase)
	srv.RegisterMethod(JSONRPCMethodCreateTestCaseList, r.regCreateTestCaseList)
	srv.RegisterMethod(JSONRPCMethodGetProblemTestCaseSnippetList, r.regGetProblemTestCaseSnippetList)
	srv.RegisterMethod(JSONRPCMethodGetTestCase, r.regGetTestCase)
	srv.RegisterMethod(JSONRPCMethodUpdateTestCase, r.regUpdateTestCase)
	srv.RegisterMethod(JSONRPCMethodDeleteTestCase, r.regDeleteTestCase)
	srv.RegisterMethod(JSONRPCMethodGetAccountProblemSnippetList, r.regGetAccountProblemSnippetList)
	srv.RegisterMethod(JSONRPCMethodCreateSubmission, r.regCreateSubmission)
	srv.RegisterMethod(JSONRPCMethodGetSubmissionSnippetList, r.regGetSubmissionSnippetList)
	srv.RegisterMethod(JSONRPCMethodGetSubmission, r.regGetSubmission)
	srv.RegisterMethod(JSONRPCMethodDeleteSubmission, r.regDeleteSubmission)
	srv.RegisterMethod(JSONRPCMethodGetAccountSubmissionSnippetList, r.regGetAccountSubmissionSnippetList)
	srv.RegisterMethod(JSONRPCMethodGetProblemSubmissionSnippetList, r.regGetProblemSubmissionSnippetList)
	srv.RegisterMethod(JSONRPCMethodGetAccountProblemSubmissionSnippetList, r.regGetAccountProblemSubmissionSnippetList)

	srv.With(middlewares...)
}

func (r *regAPI) regGetServerInfo(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.GetServerInfoRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.GetServerInfo(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed GetServerInfo: %w", err)
	}

	return res, nil
}

func (r *regAPI) regCreateAccount(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.CreateAccountRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.CreateAccount(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed CreateAccount: %w", err)
	}

	return res, nil
}

func (r *regAPI) regGetAccountList(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.GetAccountListRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.GetAccountList(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed GetAccountList: %w", err)
	}

	return res, nil
}

func (r *regAPI) regGetAccount(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.GetAccountRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.GetAccount(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed GetAccount: %w", err)
	}

	return res, nil
}

func (r *regAPI) regUpdateAccount(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.UpdateAccountRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.UpdateAccount(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed UpdateAccount: %w", err)
	}

	return res, nil
}

func (r *regAPI) regCreateSession(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.CreateSessionRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.CreateSession(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed CreateSession: %w", err)
	}

	return res, nil
}

func (r *regAPI) regGetSession(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.GetSessionRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.GetSession(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed GetSession: %w", err)
	}

	return res, nil
}

func (r *regAPI) regDeleteSession(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.DeleteSessionRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.DeleteSession(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed DeleteSession: %w", err)
	}

	return res, nil
}

func (r *regAPI) regCreateProblem(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.CreateProblemRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.CreateProblem(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed CreateProblem: %w", err)
	}

	return res, nil
}

func (r *regAPI) regGetProblemSnippetList(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.GetProblemSnippetListRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.GetProblemSnippetList(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed GetProblemSnippetList: %w", err)
	}

	return res, nil
}

func (r *regAPI) regGetProblem(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.GetProblemRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.GetProblem(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed GetProblem: %w", err)
	}

	return res, nil
}

func (r *regAPI) regUpdateProblem(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.UpdateProblemRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.UpdateProblem(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed UpdateProblem: %w", err)
	}

	return res, nil
}

func (r *regAPI) regDeleteProblem(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.DeleteProblemRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.DeleteProblem(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed DeleteProblem: %w", err)
	}

	return res, nil
}

func (r *regAPI) regCreateTestCase(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.CreateTestCaseRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.CreateTestCase(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed CreateTestCase: %w", err)
	}

	return res, nil
}

func (r *regAPI) regCreateTestCaseList(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.CreateTestCaseListRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.CreateTestCaseList(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed CreateTestCaseList: %w", err)
	}

	return res, nil
}

func (r *regAPI) regGetProblemTestCaseSnippetList(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.GetProblemTestCaseSnippetListRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.GetProblemTestCaseSnippetList(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed GetProblemTestCaseSnippetList: %w", err)
	}

	return res, nil
}

func (r *regAPI) regGetTestCase(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.GetTestCaseRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.GetTestCase(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed GetTestCase: %w", err)
	}

	return res, nil
}

func (r *regAPI) regUpdateTestCase(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.UpdateTestCaseRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.UpdateTestCase(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed UpdateTestCase: %w", err)
	}

	return res, nil
}

func (r *regAPI) regDeleteTestCase(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.DeleteTestCaseRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.DeleteTestCase(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed DeleteTestCase: %w", err)
	}

	return res, nil
}

func (r *regAPI) regGetAccountProblemSnippetList(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.GetAccountProblemSnippetListRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.GetAccountProblemSnippetList(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed GetAccountProblemSnippetList: %w", err)
	}

	return res, nil
}

func (r *regAPI) regCreateSubmission(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.CreateSubmissionRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.CreateSubmission(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed CreateSubmission: %w", err)
	}

	return res, nil
}

func (r *regAPI) regGetSubmissionSnippetList(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.GetSubmissionSnippetListRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.GetSubmissionSnippetList(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed GetSubmissionSnippetList: %w", err)
	}

	return res, nil
}

func (r *regAPI) regGetSubmission(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.GetSubmissionRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.GetSubmission(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed GetSubmission: %w", err)
	}

	return res, nil
}

func (r *regAPI) regDeleteSubmission(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.DeleteSubmissionRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.DeleteSubmission(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed DeleteSubmission: %w", err)
	}

	return res, nil
}

func (r *regAPI) regGetAccountSubmissionSnippetList(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.GetAccountSubmissionSnippetListRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.GetAccountSubmissionSnippetList(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed GetAccountSubmissionSnippetList: %w", err)
	}

	return res, nil
}

func (r *regAPI) regGetProblemSubmissionSnippetList(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.GetProblemSubmissionSnippetListRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.GetProblemSubmissionSnippetList(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed GetProblemSubmissionSnippetList: %w", err)
	}

	return res, nil
}

func (r *regAPI) regGetAccountProblemSubmissionSnippetList(ctx context.Context, params json.RawMessage) (any, error) {
	in := new(rpc.GetAccountProblemSubmissionSnippetListRequest)
	if len(params) != 0 {
		if err := pjson.Unmarshal(params, in); err != nil {
			return nil, pjrpc.JRPCErrParseError("failed to parse params")
		}
	}

	res, err := r.svc.GetAccountProblemSubmissionSnippetList(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed GetAccountProblemSubmissionSnippetList: %w", err)
	}

	return res, nil
}
