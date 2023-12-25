package app

import (
	"go.uber.org/zap"

	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
)

type DistributedWorkerCato struct {
	dbMigrator db.Migrator
	logger     *zap.Logger
}

func NewDistributedWorkerCato(
	dbMigrator db.Migrator,
	logger *zap.Logger,
) *DistributedWorkerCato {
	return &DistributedWorkerCato{
		dbMigrator: dbMigrator,
		logger:     logger,
	}
}

func (c DistributedWorkerCato) Start() error {
	return nil
}
