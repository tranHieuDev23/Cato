package logic

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewHash,
	NewToken,
	NewRole,
	NewAccount,
	NewProblem,
	NewTestCase,
	NewLocalSubmission,
	NewDistributedSubmission,
	NewCompile,
	NewTestCaseRun,
	NewLocalJudge,
	NewDistributedJudge,
)
