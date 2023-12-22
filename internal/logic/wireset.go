package logic

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewHash,
	NewAccount,
	NewProblem,
	NewTestCase,
	NewSubmission,
)
