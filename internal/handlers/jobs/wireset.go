package jobs

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewCreateFirstAccounts,
	NewScheduleSubmittedExecutingSubmissionToJudge,
	NewSyncProblems,
	NewJudgeDistributedFirstSubmittedSubmission,
	NewRevertExecutingSubmissions,
)
