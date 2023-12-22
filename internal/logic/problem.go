package logic

import (
	"context"

	"go.uber.org/zap"

	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc"
)

type Problem interface {
	CreateProblem(ctx context.Context, in *rpc.CreateProblemRequest) (*rpc.CreateProblemResponse, error)
	GetProblemSnippetList(ctx context.Context, in *rpc.GetProblemSnippetListRequest) (*rpc.GetProblemSnippetListResponse, error)
	GetProblem(ctx context.Context, in *rpc.GetProblemRequest) (*rpc.GetProblemResponse, error)
	UpdateProblem(ctx context.Context, in *rpc.UpdateProblemRequest) (*rpc.UpdateProblemResponse, error)
	DeleteProblem(ctx context.Context, in *rpc.DeleteProblemRequest) (*rpc.DeleteProblemResponse, error)
	GetAccountProblemSnippetList(ctx context.Context, in *rpc.GetAccountProblemSnippetListRequest) (*rpc.GetAccountProblemSnippetListResponse, error)
}

type problem struct {
	problemDataAccessor db.ProblemDataAccessor
	logger              *zap.Logger
}

func NewProblem(
	problemDataAccessor db.ProblemDataAccessor,
	logger *zap.Logger,
) Problem {
	return &problem{
		problemDataAccessor: problemDataAccessor,
		logger:              logger,
	}
}

func (p problem) CreateProblem(ctx context.Context, in *rpc.CreateProblemRequest) (*rpc.CreateProblemResponse, error) {
	panic("unimplemented")
}

func (p problem) DeleteProblem(ctx context.Context, in *rpc.DeleteProblemRequest) (*rpc.DeleteProblemResponse, error) {
	panic("unimplemented")
}

func (p problem) GetAccountProblemSnippetList(ctx context.Context, in *rpc.GetAccountProblemSnippetListRequest) (*rpc.GetAccountProblemSnippetListResponse, error) {
	panic("unimplemented")
}

func (p problem) GetProblem(ctx context.Context, in *rpc.GetProblemRequest) (*rpc.GetProblemResponse, error) {
	panic("unimplemented")
}

func (p problem) GetProblemSnippetList(ctx context.Context, in *rpc.GetProblemSnippetListRequest) (*rpc.GetProblemSnippetListResponse, error) {
	panic("unimplemented")
}

func (p problem) UpdateProblem(ctx context.Context, in *rpc.UpdateProblemRequest) (*rpc.UpdateProblemResponse, error) {
	panic("unimplemented")
}
