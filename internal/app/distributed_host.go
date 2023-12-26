package app

import (
	"context"

	"go.uber.org/zap"

	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/http"
	"github.com/tranHieuDev23/cato/internal/handlers/jobs"
)

type DistributedHostCato struct {
	dbMigrator             db.Migrator
	createFirstAccountsJob jobs.CreateFirstAccounts
	httpServer             http.Server
	logger                 *zap.Logger
}

func NewDistributedHostCato(
	dbMigrator db.Migrator,
	createFirstAccountsJob jobs.CreateFirstAccounts,
	httpServer http.Server,
	logger *zap.Logger,
) *DistributedHostCato {
	return &DistributedHostCato{
		dbMigrator:             dbMigrator,
		createFirstAccountsJob: createFirstAccountsJob,
		httpServer:             httpServer,
		logger:                 logger,
	}
}

func (c DistributedHostCato) Start() error {
	if err := c.dbMigrator.Migrate(context.Background()); err != nil {
		return err
	}

	if err := c.createFirstAccountsJob.Run(); err != nil {
		return err
	}

	return c.httpServer.Start()
}
