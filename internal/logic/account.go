package logic

import (
	"context"

	"go.uber.org/zap"

	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc"
)

type Account interface {
	CreateAccount(ctx context.Context, in *rpc.CreateAccountRequest) (*rpc.CreateAccountResponse, error)
	GetAccountList(ctx context.Context, in *rpc.GetAccountListRequest) (*rpc.GetAccountListResponse, error)
	GetAccount(ctx context.Context, in *rpc.GetAccountRequest) (*rpc.GetAccountResponse, error)
	UpdateAccount(ctx context.Context, in *rpc.UpdateAccountRequest) (*rpc.UpdateAccountResponse, error)
	CreateSession(ctx context.Context, in *rpc.CreateSessionRequest) (*rpc.CreateSessionResponse, error)
	DeleteSession(ctx context.Context, in *rpc.DeleteSessionRequest) (*rpc.DeleteSessionResponse, error)
}

type account struct {
	accountDataAccessor db.AccountDataAccessor
	logger              *zap.Logger
}

func NewAccount(
	accountDataAccessor db.AccountDataAccessor,
	logger *zap.Logger,
) Account {
	return &account{
		accountDataAccessor: accountDataAccessor,
		logger:              logger,
	}
}

func (a account) CreateAccount(ctx context.Context, in *rpc.CreateAccountRequest) (*rpc.CreateAccountResponse, error) {
	panic("unimplemented")
}

func (a account) CreateSession(ctx context.Context, in *rpc.CreateSessionRequest) (*rpc.CreateSessionResponse, error) {
	panic("unimplemented")
}

func (a account) DeleteSession(ctx context.Context, in *rpc.DeleteSessionRequest) (*rpc.DeleteSessionResponse, error) {
	panic("unimplemented")
}

func (a account) GetAccount(ctx context.Context, in *rpc.GetAccountRequest) (*rpc.GetAccountResponse, error) {
	panic("unimplemented")
}

func (a account) GetAccountList(ctx context.Context, in *rpc.GetAccountListRequest) (*rpc.GetAccountListResponse, error) {
	panic("unimplemented")
}

func (a account) UpdateAccount(ctx context.Context, in *rpc.UpdateAccountRequest) (*rpc.UpdateAccountResponse, error) {
	panic("unimplemented")
}
