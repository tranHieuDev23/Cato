package jobs

import (
	"context"

	"github.com/tranHieuDev23/cato/internal/logic"
)

type CreateFirstAdminAccount interface {
	Run() error
}

type createFirstAdminAccount struct {
	accountLogic logic.Account
}

func NewCreateFirstAdminAccount(
	accountLogic logic.Account,
) CreateFirstAdminAccount {
	return &createFirstAdminAccount{
		accountLogic: accountLogic,
	}
}

func (j createFirstAdminAccount) Run() error {
	return j.accountLogic.CreateFirstAdminAccount(context.Background())
}

type LocalCreateFirstAdminAccount CreateFirstAdminAccount

func NewLocalCreateFirstAdminAccount(
	accountLogic logic.LocalAccount,
) LocalCreateFirstAdminAccount {
	return NewCreateFirstAdminAccount(accountLogic)
}

type DistributedCreateFirstAdminAccount CreateFirstAdminAccount

func NewDistributedCreateFirstAdminAccount(
	accountLogic logic.DistributedAccount,
) DistributedCreateFirstAdminAccount {
	return NewCreateFirstAdminAccount(accountLogic)
}
