package db

import "github.com/google/wire"

var WireSet = wire.NewSet(
	InitializeDB,
	NewTokenPublicKeyDataAccessor,
	NewMigrator,
	NewAccountDataAccessor,
	NewAccountPasswordDataAccessor,
	NewProblemDataAccessor,
	NewProblemExampleDataAccessor,
	NewTestCaseDataAccessor,
	NewProblemTestCaseHashDataAccessor,
	NewSubmissionDataAccessor,
)
