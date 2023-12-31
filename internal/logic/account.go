package logic

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/mikespook/gorbac"
	"github.com/samber/lo"
	"gitlab.com/pjrpc/pjrpc/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc"
	"github.com/tranHieuDev23/cato/internal/utils"
)

const (
	accountDisplayNameMaxLength = 32
	accountNameMinLength        = 8
)

var (
	accountNameRegex        = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	accountDisplayNameRegex = regexp.MustCompile(`^[\p{L}\p{N}\s]+$`)
)

type Account interface {
	CreateFirstAccounts(ctx context.Context) error
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
	setting                     Setting
	accountDataAccessor         db.AccountDataAccessor
	accountPasswordDataAccessor db.AccountPasswordDataAccessor
	db                          *gorm.DB
	logger                      *zap.Logger
	logicConfig                 configs.Logic
	appArguments                utils.Arguments
	displayNameSanitizePolicy   *bluemonday.Policy
}

func NewAccount(
	hash Hash,
	token Token,
	role Role,
	setting Setting,
	accountDataAccessor db.AccountDataAccessor,
	accountPasswordDataAccessor db.AccountPasswordDataAccessor,
	db *gorm.DB,
	logger *zap.Logger,
	logicConfig configs.Logic,
	appArguments utils.Arguments,
) Account {
	return &account{
		hash:                        hash,
		token:                       token,
		role:                        role,
		setting:                     setting,
		accountDataAccessor:         accountDataAccessor,
		accountPasswordDataAccessor: accountPasswordDataAccessor,
		db:                          db,
		logger:                      logger,
		logicConfig:                 logicConfig,
		appArguments:                appArguments,
		displayNameSanitizePolicy:   bluemonday.StrictPolicy(),
	}
}

func (a account) isValidAccountName(accountName string) bool {
	return len(accountName) >= 6 && len(accountName) <= 32 && accountNameRegex.MatchString(accountName)
}

func (a account) cleanupDisplayName(displayName string) string {
	displayName = strings.Trim(displayName, " ")
	displayName = a.displayNameSanitizePolicy.Sanitize(displayName)
	return displayName
}

func (a account) isValidDisplayName(displayName string) bool {
	return displayName != "" &&
		len(displayName) <= accountDisplayNameMaxLength &&
		accountDisplayNameRegex.MatchString(displayName)
}

func (a account) isValidPassword(password string) bool {
	return len(password) >= accountNameMinLength
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

func (a account) createFirstAccount(
	ctx context.Context,
	tx *gorm.DB,
	firstAccountConfig configs.FirstAccount,
	role db.AccountRole,
) error {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.String("role", string(role)))
	logger.Info("creating first account")

	if !a.isValidAccountName(firstAccountConfig.AccountName) {
		logger.Error("invalid first account's account_name")
		return errors.New("invalid first account's account_name")
	}

	if !a.isValidDisplayName(firstAccountConfig.DisplayName) {
		logger.Error("invalid first account's display_name")
		return errors.New("invalid first account's display_name")
	}

	if !a.isValidPassword(firstAccountConfig.Password) {
		logger.Error("invalid first account's password")
		return errors.New("invalid first account's password")
	}

	account := &db.Account{
		AccountName: firstAccountConfig.AccountName,
		DisplayName: firstAccountConfig.DisplayName,
		Role:        role,
	}

	if err := a.accountDataAccessor.WithDB(tx).CreateAccount(ctx, account); err != nil {
		return err
	}

	hashedPassword, err := a.hash.Hash(ctx, firstAccountConfig.Password)
	if err != nil {
		return err
	}

	accountPassword := &db.AccountPassword{
		OfAccountID: uint64(account.ID),
		Hash:        hashedPassword,
	}
	err = a.accountPasswordDataAccessor.WithDB(tx).CreateAccountPassword(ctx, accountPassword)
	if err != nil {
		return err
	}

	return nil
}

func (a account) CreateFirstAccounts(ctx context.Context) error {
	logger := utils.LoggerWithContext(ctx, a.logger)

	return a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		accountCount, err := a.accountDataAccessor.WithDB(tx).GetAccountCount(ctx)
		if err != nil {
			return err
		}

		if accountCount > 0 {
			logger.Info("there are accounts in the database, will not create first admin account")
			return nil
		}

		err = a.createFirstAccount(ctx, tx, a.logicConfig.FirstAccounts.Admin, db.AccountRoleAdmin)
		if err != nil {
			logger.With(zap.Error(err)).Error("failed to create first admin")
			return err
		}

		if a.appArguments.Distributed {
			err = a.createFirstAccount(ctx, tx, a.logicConfig.FirstAccounts.Worker, db.AccountRoleWorker)
			if err != nil {
				logger.With(zap.Error(err)).Error("failed to create first worker")
				return err
			}
		}

		return nil
	})
}

func (a account) getCreatedDBAccount(ctx context.Context, in *rpc.CreateAccountRequest) (*db.Account, error) {
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

	if !a.appArguments.Distributed && in.Role == string(rpc.AccountRoleWorker) {
		logger.Error("failed to create account: trying to create worker account on local server")
		return nil, pjrpc.JRPCErrInvalidParams()
	}

	return &db.Account{
		AccountName: in.AccountName,
		DisplayName: cleanedDisplayName,
		Role:        db.AccountRole(in.Role),
	}, nil
}

func (a account) CreateAccount(
	ctx context.Context,
	in *rpc.CreateAccountRequest,
	token string,
) (*rpc.CreateAccountResponse, error) {
	logger := utils.LoggerWithContext(ctx, a.logger)

	setting, err := a.setting.GetSetting(ctx)
	if err != nil {
		return nil, err
	}

	if setting.Account.DisableAccountCreation {
		logger.Info("account creation is disabled via setting")
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeUnavailable))
	}

	if !a.canAccountBeCreatedAnonymously(in.Role) {
		account, accountErr := a.token.GetAccount(ctx, token)
		if accountErr != nil {
			return nil, accountErr
		}

		hasPermission, accountErr := a.role.AccountHasPermission(ctx, string(account.Role), PermissionAccountsAllWrite)
		if accountErr != nil {
			return nil, accountErr
		}
		if !hasPermission {
			return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
		}
	}

	account, err := a.getCreatedDBAccount(ctx, in)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := a.hash.Hash(ctx, in.Password)
	if err != nil {
		return nil, err
	}

	response := &rpc.CreateAccountResponse{}
	if txErr := a.db.Transaction(func(tx *gorm.DB) error {
		accountNameTaken, isAccountNameTakenErr := a.WithDB(tx).IsAccountNameTaken(ctx, in.AccountName)
		if isAccountNameTakenErr != nil {
			return isAccountNameTakenErr
		}

		if accountNameTaken {
			logger.
				With(zap.String("account_name", in.AccountName)).
				Error("failed to create account: account name is already taken")
			return pjrpc.JRPCErrServerError(int(rpc.ErrorCodeAlreadyExists))
		}

		err = utils.ExecuteUntilFirstError(
			func() error {
				return a.accountDataAccessor.WithDB(tx).CreateAccount(ctx, account)
			},
			func() error {
				accountPassword := &db.AccountPassword{OfAccountID: uint64(account.ID), Hash: hashedPassword}
				return a.accountPasswordDataAccessor.WithDB(tx).CreateAccountPassword(ctx, accountPassword)
			},
		)
		if err != nil {
			return err
		}

		response.Account = a.dbAccountToRPCAccount(account)
		return nil
	}); txErr != nil {
		return nil, err
	}

	return response, nil
}

func (a account) CreateSession(
	ctx context.Context,
	in *rpc.CreateSessionRequest,
) (*rpc.CreateSessionResponse, string, time.Time, error) {
	logger := utils.LoggerWithContext(ctx, a.logger)

	setting, err := a.setting.GetSetting(ctx)
	if err != nil {
		return nil, "", time.Time{}, err
	}

	account, err := a.accountDataAccessor.GetAccountByAccountName(ctx, in.AccountName)
	if err != nil {
		return nil, "", time.Time{}, err
	}

	if account == nil {
		return nil, "", time.Time{}, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
	}

	if account.Role == db.AccountRoleContestant && setting.Account.DisableSessionCreationForContestant {
		logger.Info("session creation for contestant is disabled via setting")
		return nil, "", time.Time{}, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeUnavailable))
	}

	if account.Role == db.AccountRoleProblemSetter && setting.Account.DisableSessionCreationForProblemSetter {
		logger.Info("session creation for contestant is disabled via setting")
		return nil, "", time.Time{}, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeUnavailable))
	}

	accountPassword, err := a.accountPasswordDataAccessor.GetAccountPasswordOfAccountID(ctx, uint64(account.ID))
	if err != nil {
		return nil, "", time.Time{}, err
	}

	if accountPassword == nil {
		return nil, "", time.Time{}, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeUnauthenticated))
	}

	equal, err := a.hash.IsHashEqual(ctx, in.Password, accountPassword.Hash)
	if err != nil {
		return nil, "", time.Time{}, err
	}
	if !equal {
		logger.Error("incorrect password")
		return nil, "", time.Time{}, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeUnauthenticated))
	}

	token, expireTime, err := a.token.GetToken(ctx, uint64(account.ID))
	if err != nil {
		return nil, "", time.Time{}, err
	}

	return &rpc.CreateSessionResponse{
		Account: a.dbAccountToRPCAccount(account),
		Token:   token,
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

func (a account) GetAccount(
	ctx context.Context,
	in *rpc.GetAccountRequest,
	token string,
) (*rpc.GetAccountResponse, error) {
	account, err := a.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	requestedAccount, err := a.accountDataAccessor.GetAccount(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	if requestedAccount == nil {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
	}

	requiredPermissionList := []gorbac.Permission{PermissionAccountsAllRead}
	if requestedAccount.ID == account.ID {
		requiredPermissionList = append(requiredPermissionList, PermissionAccountsSelfRead)
	}

	hasPermission, err := a.role.AccountHasPermission(ctx, string(account.Role), requiredPermissionList...)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	return &rpc.GetAccountResponse{
		Account: a.dbAccountToRPCAccount(requestedAccount),
	}, nil
}

func (a account) GetAccountList(
	ctx context.Context,
	in *rpc.GetAccountListRequest,
	token string,
) (*rpc.GetAccountListResponse, error) {
	account, err := a.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	hasPermission, err := a.role.AccountHasPermission(ctx, string(account.Role), PermissionAccountsAllRead)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
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

func (a account) applyUpdateAccount(in *rpc.UpdateAccountRequest, account *db.Account) error {
	if in.DisplayName != nil {
		cleanedDisplayName := a.cleanupDisplayName(*in.DisplayName)
		if !a.isValidDisplayName(cleanedDisplayName) {
			return pjrpc.JRPCErrInvalidParams()
		}

		account.DisplayName = cleanedDisplayName
	}

	if in.Role != nil {
		account.Role = db.AccountRole(*in.Role)
	}

	return nil
}

func (a account) updateAccountPassword(
	ctx context.Context,
	in *rpc.UpdateAccountRequest,
	tx *gorm.DB,
) error {
	if in.Password == nil {
		return nil
	}

	logger := utils.LoggerWithContext(ctx, a.logger)
	if !a.isValidPassword(*in.Password) {
		logger.Error("failed to update account: invalid password")
		return pjrpc.JRPCErrInvalidParams()
	}

	accountPassword, err := a.accountPasswordDataAccessor.
		WithDB(tx).
		GetAccountPasswordOfAccountID(ctx, in.ID)
	if err != nil {
		return err
	}

	hashedPassword, err := a.hash.Hash(ctx, *in.Password)
	if err != nil {
		return err
	}

	accountPassword.Hash = hashedPassword

	return a.accountPasswordDataAccessor.WithDB(tx).UpdateAccountPassword(ctx, accountPassword)
}

func (a account) UpdateAccount(
	ctx context.Context,
	in *rpc.UpdateAccountRequest,
	token string,
) (*rpc.UpdateAccountResponse, error) {
	logger := utils.LoggerWithContext(ctx, a.logger)

	setting, err := a.setting.GetSetting(ctx)
	if err != nil {
		return nil, err
	}

	if setting.Account.DisableAccountUpdate {
		logger.Info("account update is disabled via setting")
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeUnavailable))
	}

	if !a.appArguments.Distributed && in.Role != nil && *in.Role == string(rpc.AccountRoleWorker) {
		logger.Error("failed to update account: trying to update account to worker on local server")
		return nil, pjrpc.JRPCErrInvalidParams()
	}

	account, err := a.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	updatedAccount, err := a.accountDataAccessor.GetAccount(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	if updatedAccount == nil {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeNotFound))
	}

	requiredPermissionList := []gorbac.Permission{PermissionAccountsAllWrite}
	if updatedAccount.ID == account.ID && in.Role == nil {
		requiredPermissionList = append(requiredPermissionList, PermissionAccountsSelfWrite)
	}

	hasPermission, err := a.role.AccountHasPermission(ctx, string(account.Role), requiredPermissionList...)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	err = a.applyUpdateAccount(in, updatedAccount)
	if err != nil {
		return nil, err
	}

	if txErr := a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return utils.ExecuteUntilFirstError(
			func() error { return a.accountDataAccessor.WithDB(tx).UpdateAccount(ctx, updatedAccount) },
			func() error { return a.updateAccountPassword(ctx, in, tx) },
		)
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
		setting:                     a.setting.WithDB(db),
		accountDataAccessor:         a.accountDataAccessor.WithDB(db),
		accountPasswordDataAccessor: a.accountPasswordDataAccessor.WithDB(db),
		db:                          db,
		logger:                      a.logger,
		appArguments:                a.appArguments,
		logicConfig:                 a.logicConfig,
		displayNameSanitizePolicy:   a.displayNameSanitizePolicy,
	}
}
