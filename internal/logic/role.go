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
	PermissionAccountsSelfRead     = gorbac.NewStdPermission("accounts.self.read")
	PermissionAccountsSelfWrite    = gorbac.NewStdPermission("accounts.self.write")
	PermissionAccountsAllRead      = gorbac.NewStdPermission("accounts.all.read")
	PermissionAccountsAllWrite     = gorbac.NewStdPermission("accounts.all.write")
	PermissionProblemsSelfRead     = gorbac.NewStdPermission("problems.self.read")
	PermissionProblemsSelfWrite    = gorbac.NewStdPermission("problems.self.write")
	PermissionProblemsAllRead      = gorbac.NewStdPermission("problems.all.read")
	PermissionProblemsAllWrite     = gorbac.NewStdPermission("problems.all.write")
	PermissionTestCasesSelfRead    = gorbac.NewStdPermission("test_cases.self.read")
	PermissionTestCasesSelfWrite   = gorbac.NewStdPermission("test_cases.self.write")
	PermissionTestCasesAllRead     = gorbac.NewStdPermission("test_cases.all.read")
	PermissionTestCasesAllWrite    = gorbac.NewStdPermission("test_cases.all.write")
	PermissionSubmissionsSelfRead  = gorbac.NewStdPermission("submissions.self.read")
	PermissionSubmissionsSelfWrite = gorbac.NewStdPermission("submissions.self.write")
	PermissionSubmissionsAllRead   = gorbac.NewStdPermission("submissions.all.read")
	PermissionSubmissionsAllWrite  = gorbac.NewStdPermission("submissions.all.write")
	PermissionSettingsWrite        = gorbac.NewStdPermission("settings.write")
)

type Role interface {
	AccountHasPermission(ctx context.Context, accountRole string, permissions ...gorbac.Permission) (bool, error)
}

type role struct {
	rbac   *gorbac.RBAC
	logger *zap.Logger
}

//nolint:errcheck // If the initialization fails, the process will exit right from the get-go
func initializeGoRBAC() *gorbac.RBAC {
	rbac := gorbac.New()

	goRBACRoleAdmin := gorbac.NewStdRole(string(db.AccountRoleAdmin))
	goRBACRoleAdmin.Assign(PermissionAccountsAllRead)
	goRBACRoleAdmin.Assign(PermissionAccountsAllWrite)
	goRBACRoleAdmin.Assign(PermissionProblemsAllRead)
	goRBACRoleAdmin.Assign(PermissionProblemsAllWrite)
	goRBACRoleAdmin.Assign(PermissionTestCasesAllRead)
	goRBACRoleAdmin.Assign(PermissionTestCasesAllWrite)
	goRBACRoleAdmin.Assign(PermissionSubmissionsAllRead)
	goRBACRoleAdmin.Assign(PermissionSubmissionsSelfWrite)
	goRBACRoleAdmin.Assign(PermissionSettingsWrite)

	goRBACRoleProblemSetter := gorbac.NewStdRole(string(db.AccountRoleProblemSetter))
	goRBACRoleProblemSetter.Assign(PermissionAccountsAllRead)
	goRBACRoleProblemSetter.Assign(PermissionAccountsSelfWrite)
	goRBACRoleProblemSetter.Assign(PermissionProblemsAllRead)
	goRBACRoleProblemSetter.Assign(PermissionProblemsSelfWrite)
	goRBACRoleProblemSetter.Assign(PermissionTestCasesAllRead)
	goRBACRoleProblemSetter.Assign(PermissionTestCasesSelfWrite)
	goRBACRoleProblemSetter.Assign(PermissionSubmissionsAllRead)

	goRBACRoleContestant := gorbac.NewStdRole(string(db.AccountRoleContestant))
	goRBACRoleContestant.Assign(PermissionAccountsSelfRead)
	goRBACRoleContestant.Assign(PermissionAccountsSelfWrite)
	goRBACRoleContestant.Assign(PermissionProblemsAllRead)
	goRBACRoleContestant.Assign(PermissionTestCasesAllRead)
	goRBACRoleContestant.Assign(PermissionSubmissionsAllRead)
	goRBACRoleContestant.Assign(PermissionSubmissionsSelfWrite)

	goRBACRoleWorker := gorbac.NewStdRole(string(db.AccountRoleWorker))
	goRBACRoleWorker.Assign(PermissionProblemsAllRead)
	goRBACRoleWorker.Assign(PermissionTestCasesAllRead)
	goRBACRoleWorker.Assign(PermissionSubmissionsAllRead)
	goRBACRoleWorker.Assign(PermissionSubmissionsAllWrite)

	rbac.Add(goRBACRoleAdmin)
	rbac.Add(goRBACRoleProblemSetter)
	rbac.Add(goRBACRoleContestant)
	rbac.Add(goRBACRoleWorker)

	rbac.SetParent(string(db.AccountRoleProblemSetter), string(db.AccountRoleAdmin))
	rbac.SetParent(string(db.AccountRoleContestant), string(db.AccountRoleAdmin))

	return rbac
}

func NewRole(logger *zap.Logger) Role {
	return &role{
		rbac:   initializeGoRBAC(),
		logger: logger,
	}
}

func (r role) AccountHasPermission(
	ctx context.Context,
	accountRole string,
	permissions ...gorbac.Permission,
) (bool, error) {
	logger := utils.LoggerWithContext(ctx, r.logger)

	accountRBACRole, _, err := r.rbac.Get(accountRole)
	if err != nil {
		if errors.Is(err, gorbac.ErrRoleNotExist) {
			logger.With(zap.String("account_role", accountRole)).Error("invalid account role")
		}

		return false, pjrpc.JRPCErrInternalError()
	}

	for i := range permissions {
		if accountRBACRole.Permit(permissions[i]) {
			return true, nil
		}
	}

	return false, nil
}
