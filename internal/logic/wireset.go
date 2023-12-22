package logic

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewAccount,
	NewProblem,
	NewTestCase,
	NewSubmission,
)
