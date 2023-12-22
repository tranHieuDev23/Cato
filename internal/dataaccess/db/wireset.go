package db

import "github.com/google/wire"

var WireSet = wire.NewSet(
	InitializeDB,
	NewMigrator,
	NewAccountDataAccessor,
	NewAccountPasswordDataAccessor,
	NewProblemDataAccessor,
	NewProblemExampleDataAccessor,
	NewTestCaseDataAccessor,
	NewSubmissionDataAccessor,
)
