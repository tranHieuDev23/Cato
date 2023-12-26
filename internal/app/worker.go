package app

import (
	"context"
	"syscall"

	"go.uber.org/zap"

	"github.com/robfig/cron/v3"

	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/jobs"
	"github.com/tranHieuDev23/cato/internal/utils"
)

type Worker struct {
	dbMigrator                                  db.Migrator
	scheduleSubmittedExecutingSubmissionToJudge jobs.ScheduleSubmittedExecutingSubmissionToJudge
	syncProblems                                jobs.SyncProblems
	judgeDistributedFirstSubmittedSubmission    jobs.JudgeDistributedFirstSubmittedSubmission
	logger                                      *zap.Logger
	cron                                        *cron.Cron
	logicConfig                                 configs.Logic
}

func NewWorker(
	dbMigrator db.Migrator,
	scheduleSubmittedExecutingSubmissionToJudge jobs.ScheduleSubmittedExecutingSubmissionToJudge,
	syncProblems jobs.SyncProblems,
	judgeDistributedFirstSubmittedSubmission jobs.JudgeDistributedFirstSubmittedSubmission,
	logger *zap.Logger,
	cron *cron.Cron,
	logicConfig configs.Logic,
) *Worker {
	return &Worker{
		dbMigrator: dbMigrator,
		scheduleSubmittedExecutingSubmissionToJudge: scheduleSubmittedExecutingSubmissionToJudge,
		syncProblems:                             syncProblems,
		judgeDistributedFirstSubmittedSubmission: judgeDistributedFirstSubmittedSubmission,
		logger:                                   logger,
		cron:                                     cron,
		logicConfig:                              logicConfig,
	}
}

func (c Worker) Start() error {
	if err := c.dbMigrator.Migrate(context.Background()); err != nil {
		return err
	}

	if err := c.scheduleSubmittedExecutingSubmissionToJudge.Run(); err != nil {
		return err
	}

	if _, err := c.cron.AddFunc(c.logicConfig.SyncProblem.Schedule, func() {
		if err := c.syncProblems.Run(); err != nil {
			c.logger.With(zap.Error(err)).Error("failed to run sync problem cronjob")
		}
	}); err != nil {
		return err
	}

	if _, err := c.cron.AddFunc(c.logicConfig.Judge.Schedule, func() {
		if err := c.judgeDistributedFirstSubmittedSubmission.Run(); err != nil {
			c.logger.With(zap.Error(err)).Error("failed to run judge distributed first submitted submission cronjob")
		}
	}); err != nil {
		return err
	}

	c.logger.Info("cron starting")
	c.cron.Start()

	utils.BlockUntilSignal(syscall.SIGINT, syscall.SIGTERM)
	return nil
}
