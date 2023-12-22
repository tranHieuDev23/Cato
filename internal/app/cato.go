package app

import (
	"context"

	"go.uber.org/zap"

	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/http"
)

type Cato interface {
	Start() error
}

type cato struct {
	dbMigrator db.Migrator
	httpServer http.Server
	logger     *zap.Logger
}

func NewCato(
	dbMigrator db.Migrator,
	httpServer http.Server,
	logger *zap.Logger,
) Cato {
	return &cato{
		dbMigrator: dbMigrator,
		httpServer: httpServer,
		logger:     logger,
	}
}

func (c cato) Start() error {
	if err := c.dbMigrator.Migrate(context.Background()); err != nil {
		return err
	}

	return c.httpServer.Start()
}
