package app

import (
	"context"
	"syscall"

	"go.uber.org/zap"

	"github.com/go-co-op/gocron/v2"

	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/jobs"
	"github.com/tranHieuDev23/cato/internal/utils"
)

type DistributedWorker struct {
	dbMigrator                                  db.Migrator
	scheduleSubmittedExecutingSubmissionToJudge jobs.ScheduleSubmittedExecutingSubmissionToJudge
	syncProblems                                jobs.SyncProblems
	logger                                      *zap.Logger
	scheduler                                   gocron.Scheduler
	logicConfig                                 configs.Logic
}

func NewDistributedWorker(
	dbMigrator db.Migrator,
	scheduleSubmittedExecutingSubmissionToJudge jobs.ScheduleSubmittedExecutingSubmissionToJudge,
	syncProblems jobs.SyncProblems,
	logger *zap.Logger,
	scheduler gocron.Scheduler,
	logicConfig configs.Logic,
) *DistributedWorker {
	return &DistributedWorker{
		dbMigrator: dbMigrator,
		scheduleSubmittedExecutingSubmissionToJudge: scheduleSubmittedExecutingSubmissionToJudge,
		syncProblems: syncProblems,
		logger:       logger,
		scheduler:    scheduler,
		logicConfig:  logicConfig,
	}
}

func (c DistributedWorker) Start() error {
	if err := c.dbMigrator.Migrate(context.Background()); err != nil {
		return err
	}

	if err := c.scheduleSubmittedExecutingSubmissionToJudge.Run(); err != nil {
		return err
	}

	if _, err := c.scheduler.NewJob(
		gocron.CronJob(c.logicConfig.SyncProblem.Schedule, false),
		gocron.NewTask(func() {
			if err := c.syncProblems.Run(); err != nil {
				c.logger.With(zap.Error(err)).Error("failed to run sync problem cronjob")
			}
		}),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	); err != nil {
		return err
	}

	c.logger.Info("scheduler starting")
	c.scheduler.Start()

	utils.BlockUntilSignal(syscall.SIGINT, syscall.SIGTERM)

	return nil
}
