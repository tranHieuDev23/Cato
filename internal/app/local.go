package app

import (
	"context"

	"go.uber.org/zap"

	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/http"
	"github.com/tranHieuDev23/cato/internal/handlers/jobs"
)

type App interface {
	Start() error
}

type LocalCato struct {
	dbMigrator                                  db.Migrator
	createFirstAccountsJob                      jobs.CreateFirstAccounts
	scheduleSubmittedExecutingSubmissionToJudge jobs.ScheduleSubmittedExecutingSubmissionToJudge
	httpServer                                  http.Server
	logger                                      *zap.Logger
}

func NewLocalCato(
	dbMigrator db.Migrator,
	createFirstAccountsJob jobs.CreateFirstAccounts,
	scheduleSubmittedExecutingSubmissionToJudge jobs.ScheduleSubmittedExecutingSubmissionToJudge,
	httpServer http.Server,
	logger *zap.Logger,
) *LocalCato {
	return &LocalCato{
		dbMigrator:             dbMigrator,
		createFirstAccountsJob: createFirstAccountsJob,
		scheduleSubmittedExecutingSubmissionToJudge: scheduleSubmittedExecutingSubmissionToJudge,
		httpServer: httpServer,
		logger:     logger,
	}
}

func (c LocalCato) Start() error {
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
