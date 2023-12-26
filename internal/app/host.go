package app

import (
	"context"

	"go.uber.org/zap"

	"github.com/robfig/cron/v3"

	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/http"
	"github.com/tranHieuDev23/cato/internal/handlers/jobs"
)

type Host struct {
	dbMigrator                                     db.Migrator
	createFirstAccountsJob                         jobs.CreateFirstAccounts
	scheduleSubmittedExecutingSubmissionToJudgeJob jobs.ScheduleSubmittedExecutingSubmissionToJudge
	revertExecutingSubmissionsJob                  jobs.RevertExecutingSubmissions
	httpServer                                     http.Server
	logger                                         *zap.Logger
	cron                                           *cron.Cron
	logicConfig                                    configs.Logic
}

func NewHost(
	dbMigrator db.Migrator,
	createFirstAccountsJob jobs.CreateFirstAccounts,
	scheduleSubmittedExecutingSubmissionToJudgeJob jobs.ScheduleSubmittedExecutingSubmissionToJudge,
	revertExecutingSubmissionsJob jobs.RevertExecutingSubmissions,
	httpServer http.Server,
	logger *zap.Logger,
	cron *cron.Cron,
	logicConfig configs.Logic,
) *Host {
	return &Host{
		dbMigrator:             dbMigrator,
		createFirstAccountsJob: createFirstAccountsJob,
		scheduleSubmittedExecutingSubmissionToJudgeJob: scheduleSubmittedExecutingSubmissionToJudgeJob,
		revertExecutingSubmissionsJob:                  revertExecutingSubmissionsJob,
		httpServer:                                     httpServer,
		logger:                                         logger,
		cron:                                           cron,
		logicConfig:                                    logicConfig,
	}
}

func (c Host) Start() error {
	if err := c.dbMigrator.Migrate(context.Background()); err != nil {
		return err
	}

	if err := c.createFirstAccountsJob.Run(); err != nil {
		return err
	}

	if err := c.scheduleSubmittedExecutingSubmissionToJudgeJob.Run(); err != nil {
		return err
	}

	if _, err := c.cron.AddFunc(c.logicConfig.RevertExecutingSubmissions.Schedule, func() {
		if err := c.revertExecutingSubmissionsJob.Run(); err != nil {
			c.logger.With(zap.Error(err)).Error("failed to run revert executing submission cronjob")
		}
	}); err != nil {
		return err
	}

	c.logger.Info("cron starting")
	c.cron.Start()

	return c.httpServer.Start()
}
