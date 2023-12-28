package db

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/utils"
)

type gormZapLogger struct {
	logger *zap.Logger
}

func (la gormZapLogger) LogMode(_ logger.LogLevel) logger.Interface {
	return la
}

func (la gormZapLogger) Info(ctx context.Context, s string, args ...interface{}) {
	utils.LoggerWithContext(ctx, la.logger).Info(s, zap.Any("args", args))
}

func (la gormZapLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	utils.LoggerWithContext(ctx, la.logger).Warn(s, zap.Any("args", args))
}

func (la gormZapLogger) Error(ctx context.Context, s string, args ...interface{}) {
	utils.LoggerWithContext(ctx, la.logger).Error(s, zap.Any("args", args))
}

func (la gormZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	utils.LoggerWithContext(ctx, la.logger).Debug("SQL",
		zap.String("sql", sql),
		zap.Duration("elapsed", time.Since(begin)),
		zap.Int64("rows", rows),
		zap.Error(err),
	)
}

func InitializeDB(
	logger *zap.Logger,
	databaseConfig configs.Database,
) (*gorm.DB, error) {
	databaseFilePath, err := filepath.Abs(databaseConfig.FilePath)
	if err != nil {
		return nil, errors.New("failed to get abs database file path")
	}

	db, err := gorm.Open(
		sqlite.Open(fmt.Sprintf("file://%s", databaseFilePath)),
		&gorm.Config{
			Logger: gormZapLogger{logger: logger},
		},
	)
	if err != nil {
		return nil, err
	}

	migrator := NewMigrator(db, logger)
	err = migrator.Migrate(context.Background())
	if err != nil {
		return nil, err
	}

	return db, nil
}
