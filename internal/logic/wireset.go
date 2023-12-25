package logic

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewHash,
	NewToken,
	NewRole,
	NewAccount,
	NewProblem,
	NewTestCase,
	NewSubmission,
	NewCompile,
	NewTestCaseRun,
	NewJudge,
)
