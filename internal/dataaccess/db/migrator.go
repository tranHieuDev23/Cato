package db

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/tranHieuDev23/cato/internal/utils"
)

type Migrator interface {
	Migrate(ctx context.Context) error
}

type migrator struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewMigrator(
	db *gorm.DB,
	logger *zap.Logger,
) Migrator {
	return &migrator{
		db:     db,
		logger: logger,
	}
}

func (m migrator) Migrate(ctx context.Context) error {
	if err := m.db.AutoMigrate(
		new(Account),
		new(AccountPassword),
		new(TokenPublicKey),
		new(Problem),
		new(ProblemExample),
		new(TestCase),
		new(ProblemTestCaseHash),
		new(Submission),
	); err != nil {
		utils.LoggerWithContext(ctx, m.logger).
			With(zap.Error(err)).
			Error("failed to migrate database")
		return err
	}

	return nil
}
