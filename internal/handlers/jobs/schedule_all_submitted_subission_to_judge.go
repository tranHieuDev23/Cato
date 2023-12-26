package jobs

import (
	"context"

	"github.com/tranHieuDev23/cato/internal/logic"
)

type ScheduleSubmittedExecutingSubmissionToJudge interface {
	Run() error
}

type scheduleSubmittedExecutingSubmissionToJudge struct {
	submissionLogic logic.Submission
}

func NewScheduleSubmittedExecutingSubmissionToJudge(
	submissionLogic logic.Submission,
) ScheduleSubmittedExecutingSubmissionToJudge {
	return &scheduleSubmittedExecutingSubmissionToJudge{
		submissionLogic: submissionLogic,
	}
}

func (j scheduleSubmittedExecutingSubmissionToJudge) Run() error {
	return j.submissionLogic.ScheduleSubmittedExecutingSubmissionToJudge(context.Background())
}
