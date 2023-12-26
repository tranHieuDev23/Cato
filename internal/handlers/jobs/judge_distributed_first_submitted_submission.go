package jobs

import (
	"context"

	"github.com/tranHieuDev23/cato/internal/logic"
)

type JudgeDistributedFirstSubmittedSubmission interface {
	Run() error
}

type judgeDistributedFirstSubmittedSubmission struct {
	judgeLogic logic.Judge
}

func NewJudgeDistributedFirstSubmittedSubmission(
	judgeLogic logic.Judge,
) JudgeDistributedFirstSubmittedSubmission {
	return &judgeDistributedFirstSubmittedSubmission{
		judgeLogic: judgeLogic,
	}
}

func (j judgeDistributedFirstSubmittedSubmission) Run() error {
	return j.judgeLogic.JudgeDistributedFirstSubmittedSubmission(context.Background())
}
