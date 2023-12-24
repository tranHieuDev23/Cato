package logic

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/samber/lo"
	"gitlab.com/pjrpc/pjrpc/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc"
	"github.com/tranHieuDev23/cato/internal/utils"
)

var (
	accountNameRegex        = regexp.MustCompile("^[a-zA-Z0-9]+$")
	accountDisplayNameRegex = regexp.MustCompile("^[\\p{L}\\p{N}\\s]+$")
)

type Account interface {
	CreateFirstAdminAccount(ctx context.Context) error
	CreateAccount(ctx context.Context, in *rpc.CreateAccountRequest, token string) (*rpc.CreateAccountResponse, error)
	GetAccountList(ctx context.Context, in *rpc.GetAccountListRequest, token string) (*rpc.GetAccountListResponse, error)
	GetAccount(ctx context.Context, in *rpc.GetAccountRequest, token string) (*rpc.GetAccountResponse, error)
	UpdateAccount(ctx context.Context, in *rpc.UpdateAccountRequest, token string) (*rpc.UpdateAccountResponse, error)
	CreateSession(ctx context.Context, in *rpc.CreateSessionRequest) (*rpc.CreateSessionResponse, string, time.Time, error)
	GetSession(ctx context.Context, token string) (*rpc.GetSessionResponse, error)
	DeleteSession(ctx context.Context, token string) error
	IsAccountNameTaken(ctx context.Context, accountName string) (bool, error)
	WithDB(db *gorm.DB) Account
}

type account struct {
	hash                        Hash
	token                       Token
	role                        Role
	accountDataAccessor         db.AccountDataAccessor
	accountPasswordDataAccessor db.AccountPasswordDataAccessor
	db                          *gorm.DB
	logger                      *zap.Logger
	logicConfig                 configs.Logic
}

func NewAccount(
	hash Hash,
	token Token,
	role Role,
	accountDataAccessor db.AccountDataAccessor,
	accountPasswordDataAccessor db.AccountPasswordDataAccessor,
	db *gorm.DB,
	logger *zap.Logger,
	logicConfig configs.Logic,
) Account {
	return &account{
		hash:                        hash,
		token:                       token,
		role:                        role,
		accountDataAccessor:         accountDataAccessor,
		accountPasswordDataAccessor: accountPasswordDataAccessor,
		db:                          db,
		logger:                      logger,
		logicConfig:                 logicConfig,
	}
}

func (a account) isValidAccountName(accountName string) bool {
	return len(accountName) > -6 && len(accountName) <= 32 && accountNameRegex.Match([]byte(accountName))
}

func (a account) cleanupDisplayName(displayName string) string {
	return strings.Trim(displayName, " ")
}

func (a account) isValidDisplayName(displayName string) bool {
	return displayName != "" && len(displayName) <= 32 && accountDisplayNameRegex.Match([]byte(displayName))
}

func (a account) isValidPassword(password string) bool {
	return len(password) >= 8
}

func (a account) canAccountBeCreatedAnonymously(role string) bool {
	return role == string(rpc.AccountRoleContestant) || role == string(rpc.AccountRoleProblemSetter)
}

func (a account) dbAccountToRPCAccount(account *db.Account) rpc.Account {
	return rpc.Account{
		ID:          uint64(account.ID),
		AccountName: account.AccountName,
		DisplayName: account.DisplayName,
		Role:        string(account.Role),
	}
}

func (a account) IsAccountNameTaken(ctx context.Context, accountName string) (bool, error) {
	account, err := a.accountDataAccessor.GetAccountByAccountName(ctx, accountName)
	if err != nil {
		return false, err
	}

	return account != nil, nil
}

func (a account) CreateFirstAdminAccount(ctx context.Context) error {
	logger := utils.LoggerWithContext(ctx, a.logger)

	return a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if accountCount, err := a.accountDataAccessor.WithDB(tx).GetAccountCount(ctx); err != nil {
			return err
		} else if accountCount > 0 {
			logger.Info("there are accounts in the database, will not create first admin account")
			return nil
		}

		if !a.isValidAccountName(a.logicConfig.FirstAdmin.AccountName) {
			logger.Error("invalid first admin's account_name")
			return errors.New("invalid first admin's account_name")
		}

		if !a.isValidDisplayName(a.logicConfig.FirstAdmin.DisplayName) {
			logger.Error("invalid first admin's display_name")
			return errors.New("invalid first admin's display_name")
		}

		if !a.isValidDisplayName(a.logicConfig.FirstAdmin.Password) {
			logger.Error("invalid first admin's password")
			return errors.New("invalid first admin's password")
		}

		account := &db.Account{
			AccountName: a.logicConfig.FirstAdmin.AccountName,
			DisplayName: a.logicConfig.FirstAdmin.DisplayName,
			Role:        db.AccountRoleAdmin,
		}

		if err := a.accountDataAccessor.WithDB(tx).CreateAccount(ctx, account); err != nil {
			return err
		}

		hashedPassword, err := a.hash.Hash(ctx, a.logicConfig.FirstAdmin.Password)
		if err != nil {
			return err
		}

		accountPassword := &db.AccountPassword{
			OfAccountID: uint64(account.ID),
			Hash:        hashedPassword,
		}
		if err := a.accountPasswordDataAccessor.WithDB(tx).CreateAccountPassword(ctx, accountPassword); err != nil {
			return err
		}

		return nil
	})
}

func (a account) CreateAccount(ctx context.Context, in *rpc.CreateAccountRequest, token string) (*rpc.CreateAccountResponse, error) {
	logger := utils.LoggerWithContext(ctx, a.logger)

	if !a.isValidAccountName(in.AccountName) {
		logger.
			With(zap.String("account_name", in.AccountName)).
			Error("failed to create account: invalid account name")
		return nil, pjrpc.JRPCErrInvalidParams()
	}

	cleanedDisplayName := a.cleanupDisplayName(in.DisplayName)
	if !a.isValidDisplayName(cleanedDisplayName) {
		logger.
			With(zap.String("display_name", in.DisplayName)).
			Error("failed to create account: invalid display name")
		return nil, pjrpc.JRPCErrInvalidParams()
	}

	if !a.isValidPassword(in.Password) {
		logger.Error("failed to create account: invalid password")
		return nil, pjrpc.JRPCErrInvalidParams()
	}

	if !a.canAccountBeCreatedAnonymously(in.Role) {
		account, err := a.token.GetAccount(ctx, token)
		if err != nil {
			return nil, err
		}

		if hasPermission, err := a.role.AccountHasPermission(ctx, string(account.Role), PermissionAccountsAllWrite); err != nil {
			return nil, err
		} else if !hasPermission {
			return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}
	}

	hashedPassword, err := a.hash.Hash(ctx, in.Password)
	if err != nil {
		return nil, err
	}

	response := &rpc.CreateAccountResponse{}
	if txErr := a.db.Transaction(func(tx *gorm.DB) error {
		if accountNameTaken, err := a.WithDB(tx).IsAccountNameTaken(ctx, in.AccountName); err != nil {
			return err
		} else if accountNameTaken {
			logger.
				With(zap.String("account_name", in.AccountName)).
				Error("failed to create account: invalid display name")
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodeAlreadyExists))
		}

		account := &db.Account{
			AccountName: in.AccountName,
			DisplayName: cleanedDisplayName,
			Role:        db.AccountRole(in.Role),
		}

		if err := a.accountDataAccessor.WithDB(tx).CreateAccount(ctx, account); err != nil {
			return err
		}

		accountPassword := &db.AccountPassword{
			OfAccountID: uint64(account.ID),
			Hash:        hashedPassword,
		}
		if err := a.accountPasswordDataAccessor.WithDB(tx).CreateAccountPassword(ctx, accountPassword); err != nil {
			return err
		}

		response.Account = a.dbAccountToRPCAccount(account)

		return nil
	}); txErr != nil {
		return nil, err
	}

	return response, nil
}

func (a account) CreateSession(ctx context.Context, in *rpc.CreateSessionRequest) (*rpc.CreateSessionResponse, string, time.Time, error) {
	logger := utils.LoggerWithContext(ctx, a.logger)

	account, err := a.accountDataAccessor.GetAccountByAccountName(ctx, in.AccountName)
	if err != nil {
		return nil, "", time.Time{}, err
	}

	if account == nil {
		return nil, "", time.Time{}, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
	}

	accountPassword, err := a.accountPasswordDataAccessor.GetAccountPasswordOfAccountID(ctx, uint64(account.ID))
	if err != nil {
		return nil, "", time.Time{}, err
	}

	if accountPassword == nil {
		return nil, "", time.Time{}, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeUnauthenticated))
	}

	if equal, err := a.hash.IsHashEqual(ctx, in.Password, accountPassword.Hash); err != nil {
		return nil, "", time.Time{}, err
	} else if !equal {
		logger.Error("incorrect password")
		return nil, "", time.Time{}, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeUnauthenticated))
	}

	token, expireTime, err := a.token.GetToken(ctx, uint64(account.ID))
	if err != nil {
		return nil, "", time.Time{}, err
	}

	return &rpc.CreateSessionResponse{
		Account: a.dbAccountToRPCAccount(account),
	}, token, expireTime, nil
}

func (a account) GetSession(ctx context.Context, token string) (*rpc.GetSessionResponse, error) {
	account, err := a.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	return &rpc.GetSessionResponse{
		Account: a.dbAccountToRPCAccount(account),
	}, nil
}

func (a account) DeleteSession(ctx context.Context, token string) error {
	_, err := a.token.GetAccount(ctx, token)
	if err != nil {
		return err
	}

	return nil
}

func (a account) GetAccount(ctx context.Context, in *rpc.GetAccountRequest, token string) (*rpc.GetAccountResponse, error) {
	account, err := a.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	if hasPermission, err := a.role.AccountHasPermission(
		ctx,
		string(account.Role),
		PermissionAccountsSelfRead,
		PermissionAccountsAllRead,
	); err != nil {
		return nil, err
	} else if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	requestedAccount, err := a.accountDataAccessor.GetAccount(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	if requestedAccount == nil {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
	}

	if requestedAccount.ID != account.ID {
		if hasPermission, err := a.role.AccountHasPermission(
			ctx,
			string(account.Role),
			PermissionAccountsAllRead,
		); err != nil {
			return nil, err
		} else if !hasPermission {
			return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}
	}

	return &rpc.GetAccountResponse{
		Account: a.dbAccountToRPCAccount(requestedAccount),
	}, nil
}

func (a account) GetAccountList(ctx context.Context, in *rpc.GetAccountListRequest, token string) (*rpc.GetAccountListResponse, error) {
	account, err := a.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	if hasPermission, err := a.role.AccountHasPermission(ctx, string(account.Role), PermissionAccountsAllRead); err != nil {
		return nil, err
	} else if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	totalAccountCount, err := a.accountDataAccessor.GetAccountCount(ctx)
	if err != nil {
		return nil, err
	}

	accountList, err := a.accountDataAccessor.GetAccountList(ctx, in.Offset, in.Limit)
	if err != nil {
		return nil, err
	}

	rpcAccountList := lo.Map(accountList, func(item *db.Account, _ int) rpc.Account {
		return a.dbAccountToRPCAccount(item)
	})

	return &rpc.GetAccountListResponse{
		TotalAccountCount: totalAccountCount,
		AccountList:       rpcAccountList,
	}, nil
}

func (a account) UpdateAccount(ctx context.Context, in *rpc.UpdateAccountRequest, token string) (*rpc.UpdateAccountResponse, error) {
	logger := utils.LoggerWithContext(ctx, a.logger)

	account, err := a.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	if hasPermission, err := a.role.AccountHasPermission(
		ctx,
		string(account.Role),
		PermissionAccountsSelfWrite,
		PermissionAccountsAllWrite,
	); err != nil {
		return nil, err
	} else if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	updatedAccount, err := a.accountDataAccessor.GetAccount(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	if updatedAccount == nil {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
	}

	if updatedAccount.ID != account.ID || in.Role != nil {
		if hasPermission, err := a.role.AccountHasPermission(
			ctx,
			string(account.Role),
			PermissionAccountsAllWrite,
		); err != nil {
			return nil, err
		} else if !hasPermission {
			return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}
	}

	if in.DisplayName != nil {
		cleanedDisplayName := a.cleanupDisplayName(*in.DisplayName)
		if !a.isValidDisplayName(cleanedDisplayName) {
			logger.
				With(zap.String("display_name", *in.DisplayName)).
				Error("failed to update account: invalid display name")

			return nil, pjrpc.JRPCErrInvalidParams()
		}

		updatedAccount.DisplayName = cleanedDisplayName
	}

	if in.Role != nil {
		updatedAccount.Role = db.AccountRole(*in.Role)
	}

	if txErr := a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := a.accountDataAccessor.WithDB(tx).UpdateAccount(ctx, updatedAccount); err != nil {
			return err
		}

		if in.Password != nil {
			if !a.isValidPassword(*in.Password) {
				logger.Error("failed to update account: invalid password")
				return pjrpc.JRPCErrInvalidParams()
			}

			accountPassword, err := a.accountPasswordDataAccessor.WithDB(tx).GetAccountPasswordOfAccountID(ctx, in.ID)
			if err != nil {
				return err
			}

			hashedPassword, err := a.hash.Hash(ctx, *in.Password)
			if err != nil {
				return nil
			}

			accountPassword.Hash = hashedPassword
			if err := a.accountPasswordDataAccessor.WithDB(tx).UpdateAccountPassword(ctx, accountPassword); err != nil {
				return err
			}
		}

		return nil
	}); txErr != nil {
		return nil, txErr
	}

	response := &rpc.UpdateAccountResponse{
		Account: a.dbAccountToRPCAccount(updatedAccount),
	}

	return response, nil
}

func (a account) WithDB(db *gorm.DB) Account {
	return &account{
		hash:                        a.hash,
		token:                       a.token.WithDB(db),
		role:                        a.role,
		accountDataAccessor:         a.accountDataAccessor.WithDB(db),
		accountPasswordDataAccessor: a.accountPasswordDataAccessor.WithDB(db),
		db:                          db,
		logger:                      a.logger,
	}
}
