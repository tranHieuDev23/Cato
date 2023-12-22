package logic

import (
	"context"

	"go.uber.org/zap"

	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc"
)

type Submission interface {
	CreateSubmission(ctx context.Context, in *rpc.CreateSubmissionRequest, token string) (*rpc.CreateSubmissionResponse, error)
	GetSubmissionSnippetList(ctx context.Context, in *rpc.GetSubmissionSnippetListRequest, token string) (*rpc.GetSubmissionSnippetListResponse, error)
	GetSubmission(ctx context.Context, in *rpc.GetSubmissionRequest, token string) (*rpc.GetSubmissionResponse, error)
	DeleteSubmission(ctx context.Context, in *rpc.DeleteSubmissionRequest, token string) error
	GetAccountSubmissionSnippetList(ctx context.Context, in *rpc.GetAccountSubmissionSnippetListRequest, token string) (*rpc.GetAccountSubmissionSnippetListResponse, error)
	GetProblemSubmissionSnippetList(ctx context.Context, in *rpc.GetProblemSubmissionSnippetListRequest, token string) (*rpc.GetProblemSubmissionSnippetListResponse, error)
}

type submission struct {
	token                  Token
	role                   Role
	problemDataAccessor    db.ProblemDataAccessor
	submissionDataAccessor db.SubmissionDataAccessor
	logger                 *zap.Logger
}

func NewSubmission(
	token Token,
	role Role,
	problemDataAccessor db.ProblemDataAccessor,
	submissionDataAccessor db.SubmissionDataAccessor,
	logger *zap.Logger,
) Submission {
	return &submission{
		token:                  token,
		role:                   role,
		problemDataAccessor:    problemDataAccessor,
		submissionDataAccessor: submissionDataAccessor,
		logger:                 logger,
	}
}

func (s submission) CreateSubmission(ctx context.Context, in *rpc.CreateSubmissionRequest, token string) (*rpc.CreateSubmissionResponse, error) {
	panic("unimplemented")
}

func (s submission) DeleteSubmission(ctx context.Context, in *rpc.DeleteSubmissionRequest, token string) error {
	panic("unimplemented")
}

func (s submission) GetAccountSubmissionSnippetList(ctx context.Context, in *rpc.GetAccountSubmissionSnippetListRequest, token string) (*rpc.GetAccountSubmissionSnippetListResponse, error) {
	panic("unimplemented")
}

func (s submission) GetProblemSubmissionSnippetList(ctx context.Context, in *rpc.GetProblemSubmissionSnippetListRequest, token string) (*rpc.GetProblemSubmissionSnippetListResponse, error) {
	panic("unimplemented")
}

func (s submission) GetSubmission(ctx context.Context, in *rpc.GetSubmissionRequest, token string) (*rpc.GetSubmissionResponse, error) {
	panic("unimplemented")
}

func (s submission) GetSubmissionSnippetList(ctx context.Context, in *rpc.GetSubmissionSnippetListRequest, token string) (*rpc.GetSubmissionSnippetListResponse, error) {
	panic("unimplemented")
}
