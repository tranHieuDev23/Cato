package app

import (
	"context"

	"go.uber.org/zap"

	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/http"
	"github.com/tranHieuDev23/cato/internal/handlers/jobs"
)

type Host struct {
	dbMigrator                                  db.Migrator
	createFirstAccountsJob                      jobs.CreateFirstAccounts
	scheduleSubmittedExecutingSubmissionToJudge jobs.ScheduleSubmittedExecutingSubmissionToJudge
	httpServer                                  http.Server
	logger                                      *zap.Logger
}

func NewHost(
	dbMigrator db.Migrator,
	createFirstAccountsJob jobs.CreateFirstAccounts,
	scheduleSubmittedExecutingSubmissionToJudge jobs.ScheduleSubmittedExecutingSubmissionToJudge,
	httpServer http.Server,
	logger *zap.Logger,
) *Host {
	return &Host{
		dbMigrator:             dbMigrator,
		createFirstAccountsJob: createFirstAccountsJob,
		scheduleSubmittedExecutingSubmissionToJudge: scheduleSubmittedExecutingSubmissionToJudge,
		httpServer: httpServer,
		logger:     logger,
	}
}

func (c Host) Start() error {
	if err := c.dbMigrator.Migrate(context.Background()); err != nil {
		return err
	}

	if err := c.createFirstAccountsJob.Run(); err != nil {
		return err
	}

	if err := c.scheduleSubmittedExecutingSubmissionToJudge.Run(); err != nil {
		return err
	}

	return c.httpServer.Start()
}
