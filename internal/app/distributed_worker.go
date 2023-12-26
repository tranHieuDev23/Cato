package app

import (
	"context"

	"go.uber.org/zap"

	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/jobs"
)

type DistributedWorkerCato struct {
	dbMigrator                                  db.Migrator
	scheduleSubmittedExecutingSubmissionToJudge jobs.ScheduleSubmittedExecutingSubmissionToJudge
	logger                                      *zap.Logger
}

func NewDistributedWorkerCato(
	dbMigrator db.Migrator,
	scheduleSubmittedExecutingSubmissionToJudge jobs.ScheduleSubmittedExecutingSubmissionToJudge,
	logger *zap.Logger,
) *DistributedWorkerCato {
	return &DistributedWorkerCato{
		dbMigrator: dbMigrator,
		scheduleSubmittedExecutingSubmissionToJudge: scheduleSubmittedExecutingSubmissionToJudge,
		logger: logger,
	}
}

func (c DistributedWorkerCato) Start() error {
	if err := c.dbMigrator.Migrate(context.Background()); err != nil {
		return err
	}

	if err := c.scheduleSubmittedExecutingSubmissionToJudge.Run(); err != nil {
		return err
	}

	return nil
}
