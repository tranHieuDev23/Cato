package jobs

import (
	"context"

	"github.com/tranHieuDev23/cato/internal/logic"
)

type CreateFirstAccounts interface {
	Run() error
}

type createFirstAccounts struct {
	accountLogic logic.Account
}

func NewCreateFirstAccounts(
	accountLogic logic.Account,
) CreateFirstAccounts {
	return &createFirstAccounts{
		accountLogic: accountLogic,
	}
}

func (j createFirstAccounts) Run() error {
	return j.accountLogic.CreateFirstAccounts(context.Background())
}
