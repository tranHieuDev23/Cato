package db

import (
	"context"
	"errors"

	"github.com/tranHieuDev23/cato/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Setting struct {
	gorm.Model
	Data []byte
}

type SettingDataAccessor interface {
	UpsertSetting(ctx context.Context, setting *Setting) error
	GetSetting(ctx context.Context) (*Setting, error)
	WithDB(db *gorm.DB) SettingDataAccessor
}

type settingDataAccessor struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewSettingDataAccessor(
	db *gorm.DB,
	logger *zap.Logger,
) SettingDataAccessor {
	return &settingDataAccessor{
		db:     db,
		logger: logger,
	}
}

func (a settingDataAccessor) GetSetting(ctx context.Context) (*Setting, error) {
	logger := utils.LoggerWithContext(ctx, a.logger)
	db := a.db.WithContext(ctx)
	setting := new(Setting)
	if err := db.First(setting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		logger.With(zap.Error(err)).Error("failed to get setting")
		return nil, err
	}

	return setting, nil
}

func (a settingDataAccessor) UpsertSetting(ctx context.Context, setting *Setting) error {
	logger := utils.LoggerWithContext(ctx, a.logger)

	return a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		existingSetting, err := a.WithDB(tx).GetSetting(ctx)
		if err != nil {
			return err
		}

		if existingSetting == nil {
			err = tx.Create(setting).Error
			if err != nil {
				logger.With(zap.Error(err)).Error("failed to create setting record")
				return err
			}

			return nil
		}

		existingSetting.Data = setting.Data
		err = tx.Save(existingSetting).Error
		if err != nil {
			logger.With(zap.Error(err)).Error("failed to create setting record")
			return err
		}

		return nil
	})
}

func (a settingDataAccessor) WithDB(db *gorm.DB) SettingDataAccessor {
	return &settingDataAccessor{
		db:     db,
		logger: a.logger,
	}
}
