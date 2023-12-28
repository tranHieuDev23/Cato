package logic

import (
	"context"
	"encoding/json"

	"github.com/tranHieuDev23/cato/internal/dataaccess/cache"
	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc"
	"github.com/tranHieuDev23/cato/internal/utils"
	"gitlab.com/pjrpc/pjrpc/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Setting interface {
	GetSetting(ctx context.Context) (*rpc.Setting, error)
	UpdateSetting(ctx context.Context, in *rpc.UpdateSettingRequest, token string) (*rpc.UpdateSettingResponse, error)
	WithDB(db *gorm.DB) Setting
}

type setting struct {
	token               Token
	role                Role
	settingDataAccessor db.SettingDataAccessor
	settingCache        cache.Setting
	db                  *gorm.DB
	logger              *zap.Logger
}

func NewSetting(
	token Token,
	role Role,
	settingDataAccessor db.SettingDataAccessor,
	settingCache cache.Setting,
	db *gorm.DB,
	logger *zap.Logger,
) Setting {
	return &setting{
		token:               token,
		role:                role,
		settingDataAccessor: settingDataAccessor,
		settingCache:        settingCache,
		db:                  db,
		logger:              logger,
	}
}

func (s setting) GetSetting(ctx context.Context) (*rpc.Setting, error) {
	logger := utils.LoggerWithContext(ctx, s.logger)

	rpcSetting := new(rpc.Setting)

	cachedSettingBytes, err := s.settingCache.Get(ctx)
	if err == nil && cachedSettingBytes != nil {
		err = json.Unmarshal(cachedSettingBytes, rpcSetting)
		if err != nil {
			logger.
				With(zap.ByteString("data", cachedSettingBytes)).
				With(zap.Error(err)).
				Warn("failed to unmarshal setting bytes, will return the default")
			return rpcSetting, nil
		}
	}

	logger.With(zap.Error(err)).Warn("failed to get cached setting bytes from cache, will fail back to database")

	setting, err := s.settingDataAccessor.GetSetting(ctx)
	if err != nil {
		return nil, err
	}

	if setting == nil {
		logger.Info("no setting found in the database, will return the default")
		return rpcSetting, nil
	}

	err = s.settingCache.Set(ctx, setting.Data)
	if err != nil {
		logger.With(zap.Error(err)).Warn("failed to set cached setting bytes to cache")
	}

	err = json.Unmarshal(setting.Data, rpcSetting)
	if err != nil {
		logger.
			With(zap.ByteString("data", setting.Data)).
			With(zap.Error(err)).
			Warn("failed to unmarshal setting bytes, will return the default")
		return rpcSetting, nil
	}

	return rpcSetting, nil
}

func (s setting) UpdateSetting(
	ctx context.Context,
	in *rpc.UpdateSettingRequest,
	token string,
) (*rpc.UpdateSettingResponse, error) {
	logger := utils.LoggerWithContext(ctx, s.logger)

	account, err := s.token.GetAccount(ctx, token)
	if err != nil {
		return nil, err
	}

	hasPermission, err := s.role.AccountHasPermission(ctx, string(account.Role), PermissionSettingsWrite)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodePermissionDenied))
	}

	if txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		settingData, settingErr := json.Marshal(in.Setting)
		if settingErr != nil {
			logger.With(zap.Error(settingErr)).Error("failed to marshal setting")
			return settingErr
		}

		return utils.ExecuteUntilFirstError(
			func() error {
				return s.settingDataAccessor.WithDB(tx).UpsertSetting(ctx, &db.Setting{
					Data: settingData,
				})
			},
			func() error { return s.settingCache.Set(ctx, settingData) },
		)
	}); txErr != nil {
		return nil, txErr
	}

	return &rpc.UpdateSettingResponse{
		Setting: in.Setting,
	}, nil
}

func (s setting) WithDB(db *gorm.DB) Setting {
	return &setting{
		settingDataAccessor: s.settingDataAccessor.WithDB(db),
		settingCache:        s.settingCache,
		db:                  db,
		logger:              s.logger,
	}
}
