package logic

import (
	"context"
	"errors"

	"github.com/mikespook/gorbac"
	"gitlab.com/pjrpc/pjrpc/v2"
	"go.uber.org/zap"

	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/utils"
)

var (
	PermissionAccountsRead     = gorbac.NewStdPermission("accounts.read")
	PermissionAccountsWrite    = gorbac.NewStdPermission("accounts.write")
	PermissionProblemsRead     = gorbac.NewStdPermission("problems.read")
	PermissionProblemsWrite    = gorbac.NewStdPermission("problems.write")
	PermissionTestCasesRead    = gorbac.NewStdPermission("test_cases.read")
	PermissionTestCasesWrite   = gorbac.NewStdPermission("test_cases.write")
	PermissionSubmissionsRead  = gorbac.NewStdPermission("submissions.read")
	PermissionSubmissionsWrite = gorbac.NewStdPermission("submissions.write")
)

type Role interface {
	AccountCanAccessResource(ctx context.Context, accountID uint64, accountRole string, resourceOwnerID uint64) (bool, error)
	AccountHasPermission(ctx context.Context, accountRole string, permission gorbac.Permission) (bool, error)
}

type role struct {
	rbac   *gorbac.RBAC
	logger *zap.Logger
}

func initializeGoRBAC() *gorbac.RBAC {
	rbac := gorbac.New()

	goRBACRoleAdmin := gorbac.NewStdRole(string(db.AccountRoleAdmin))
	goRBACRoleProblemSetter := gorbac.NewStdRole(string(db.AccountRoleProblemSetter))
	goRBACRoleContestant := gorbac.NewStdRole(string(db.AccountRoleContestant))
	goRBACRoleWorker := gorbac.NewStdRole(string(db.AccountRoleWorker))

	rbac.Add(goRBACRoleAdmin)
	rbac.Add(goRBACRoleProblemSetter)
	rbac.Add(goRBACRoleContestant)
	rbac.Add(goRBACRoleWorker)

	rbac.SetParent(string(db.AccountRoleProblemSetter), string(db.AccountRoleAdmin))
	rbac.SetParent(string(db.AccountRoleContestant), string(db.AccountRoleAdmin))

	goRBACRoleAdmin.Assign(PermissionAccountsRead)
	goRBACRoleProblemSetter.Assign(PermissionAccountsRead)
	goRBACRoleContestant.Assign(PermissionAccountsRead)

	goRBACRoleAdmin.Assign(PermissionAccountsWrite)
	goRBACRoleProblemSetter.Assign(PermissionAccountsWrite)
	goRBACRoleContestant.Assign(PermissionAccountsWrite)

	goRBACRoleProblemSetter.Assign(PermissionProblemsRead)
	goRBACRoleProblemSetter.Assign(PermissionProblemsWrite)

	goRBACRoleProblemSetter.Assign(PermissionTestCasesRead)
	goRBACRoleProblemSetter.Assign(PermissionTestCasesWrite)

	goRBACRoleContestant.Assign(PermissionTestCasesRead)

	goRBACRoleContestant.Assign(PermissionTestCasesRead)
	goRBACRoleContestant.Assign(PermissionTestCasesWrite)

	return rbac
}

func NewRole(logger *zap.Logger) Role {
	return &role{
		rbac:   initializeGoRBAC(),
		logger: logger,
	}
}

func (r role) AccountCanAccessResource(
	ctx context.Context,
	accountID uint64,
	accountRole string,
	resourceOwnerID uint64,
) (bool, error) {
	logger := utils.LoggerWithContext(ctx, r.logger)

	if _, _, err := r.rbac.Get(accountRole); err != nil {
		if errors.Is(err, gorbac.ErrRoleNotExist) {
			logger.With(zap.String("account_role", accountRole)).Error("invalid account role")
		}

		return false, pjrpc.JRPCErrInternalError()
	}

	if accountID == resourceOwnerID {
		return true, nil
	}

	return accountRole == string(db.AccountRoleAdmin), nil
}

func (r role) AccountHasPermission(ctx context.Context, accountRole string, permission gorbac.Permission) (bool, error) {
	logger := utils.LoggerWithContext(ctx, r.logger)

	accountRBACRole, _, err := r.rbac.Get(accountRole)
	if err != nil {
		if errors.Is(err, gorbac.ErrRoleNotExist) {
			logger.With(zap.String("account_role", accountRole)).Error("invalid account role")
		}

		return false, pjrpc.JRPCErrInternalError()
	}

	return accountRBACRole.Permit(permission), nil
}
