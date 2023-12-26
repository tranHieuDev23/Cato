package jobs

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewLocalCreateFirstAdminAccount,
	NewDistributedCreateFirstAdminAccount,
	NewLocalScheduleSubmittedExecutingSubmissionToJudge,
	NewDistributedScheduleSubmittedExecutingSubmissionToJudge,
)
