package db

import (
	"context"

	"gorm.io/gorm"
)

type Migrator interface {
	Migrate(ctx context.Context) error
}

type migrator struct {
	db *gorm.DB
}

func NewMigrator(db *gorm.DB) Migrator {
	return &migrator{
		db: db,
	}
}

func (m migrator) Migrate(ctx context.Context) error {
	if err := m.db.AutoMigrate(
		new(Account),
		new(AccountPassword),
		new(Problem),
		new(ProblemExample),
		new(TestCase),
		new(Submission),
	); err != nil {
		return err
	}

	return nil
}
