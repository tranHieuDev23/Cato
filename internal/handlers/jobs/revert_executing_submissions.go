package jobs

import (
	"context"

	"github.com/tranHieuDev23/cato/internal/logic"
)

type RevertExecutingSubmissions interface {
	Run() error
}

type revertExecutingSubmissions struct {
	submissionLogic logic.Submission
}

func NewRevertExecutingSubmissions(
	submissionLogic logic.Submission,
) RevertExecutingSubmissions {
	return &revertExecutingSubmissions{
		submissionLogic: submissionLogic,
	}
}

func (j revertExecutingSubmissions) Run() error {
	return j.submissionLogic.UpdateExecutingSubmissionUpdatedBeforeThresholdToSubmitted(context.Background())
}
