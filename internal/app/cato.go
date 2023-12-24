package app

import (
	"context"

	"go.uber.org/zap"

	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/http"
	"github.com/tranHieuDev23/cato/internal/handlers/jobs"
)

type Cato interface {
	Start() error
}

type cato struct {
	dbMigrator                 db.Migrator
	createFirstAdminAccountJob jobs.CreateFirstAdminAccount
	httpServer                 http.Server
	logger                     *zap.Logger
}

func NewCato(
	dbMigrator db.Migrator,
	createFirstAdminAccountJob jobs.CreateFirstAdminAccount,
	httpServer http.Server,
	logger *zap.Logger,
) Cato {
	return &cato{
		dbMigrator:                 dbMigrator,
		createFirstAdminAccountJob: createFirstAdminAccountJob,
		httpServer:                 httpServer,
		logger:                     logger,
	}
}

func (c cato) Start() error {
	if err := c.dbMigrator.Migrate(context.Background()); err != nil {
		return err
	}

	if err := c.createFirstAdminAccountJob.Run(); err != nil {
		return err
	}

	return c.httpServer.Start()
}
