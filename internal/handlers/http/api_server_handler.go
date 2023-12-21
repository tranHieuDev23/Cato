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

func (h apiServerHandler) CreateProblem(ctx context.Context, in *rpc.CreateProblemRequest) (*rpc.CreateProblemResponse, error) {
	panic("unimplemented")
}

func (h apiServerHandler) CreateSession(ctx context.Context, in *rpc.CreateSessionRequest) (*rpc.CreateSessionResponse, error) {
	panic("unimplemented")
}

func (h apiServerHandler) CreateSubmission(ctx context.Context, in *rpc.CreateSubmissionRequest) (*rpc.CreateSubmissionResponse, error) {
	panic("unimplemented")
}

func (h apiServerHandler) CreateUser(ctx context.Context, in *rpc.CreateUserRequest) (*rpc.CreateUserResponse, error) {
	panic("unimplemented")
}

func (h apiServerHandler) DeleteProblem(ctx context.Context, in *rpc.DeleteProblemRequest) (*rpc.DeleteProblemResponse, error) {
	panic("unimplemented")
}

func (h apiServerHandler) DeleteSession(ctx context.Context, in *rpc.DeleteSessionRequest) (*rpc.DeleteSessionResponse, error) {
	panic("unimplemented")
}

func (h apiServerHandler) DeleteSubmission(ctx context.Context, in *rpc.DeleteSubmissionRequest) (*rpc.DeleteSubmissionResponse, error) {
	panic("unimplemented")
}

func (h apiServerHandler) GetProblem(ctx context.Context, in *rpc.GetProblemRequest) (*rpc.GetProblemResponse, error) {
	panic("unimplemented")
}

func (h apiServerHandler) GetProblemSnippetList(ctx context.Context, in *rpc.GetProblemSnippetListRequest) (*rpc.GetProblemSnippetListResponse, error) {
	panic("unimplemented")
}

func (h apiServerHandler) GetProblemSubmissionSnippetList(ctx context.Context, in *rpc.GetProblemSubmissionSnippetListRequest) (*rpc.GetProblemSubmissionSnippetListResponse, error) {
	panic("unimplemented")
}

func (h apiServerHandler) GetSubmission(ctx context.Context, in *rpc.GetSubmissionRequest) (*rpc.GetSubmissionResponse, error) {
	panic("unimplemented")
}

func (h apiServerHandler) GetSubmissionSnippetList(ctx context.Context, in *rpc.GetSubmissionSnippetListRequest) (*rpc.GetSubmissionSnippetListResponse, error) {
	panic("unimplemented")
}

func (h apiServerHandler) GetUser(ctx context.Context, in *rpc.GetUserRequest) (*rpc.GetUserResponse, error) {
	panic("unimplemented")
}

func (h apiServerHandler) GetUserList(ctx context.Context, in *rpc.GetUserListRequest) (*rpc.GetUserListResponse, error) {
	panic("unimplemented")
}

func (h apiServerHandler) GetUserProblemSnippetList(ctx context.Context, in *rpc.GetUserProblemSnippetListRequest) (*rpc.GetUserProblemSnippetListResponse, error) {
	panic("unimplemented")
}

func (h apiServerHandler) GetUserSubmissionSnippetList(ctx context.Context, in *rpc.GetUserSubmissionSnippetListRequest) (*rpc.GetUserSubmissionSnippetListResponse, error) {
	panic("unimplemented")
}

func (h apiServerHandler) UpdateProblem(ctx context.Context, in *rpc.UpdateProblemRequest) (*rpc.UpdateProblemResponse, error) {
	panic("unimplemented")
}

func (h apiServerHandler) UpdateUser(ctx context.Context, in *rpc.UpdateUserRequest) (*rpc.UpdateUserResponse, error) {
	panic("unimplemented")
}
