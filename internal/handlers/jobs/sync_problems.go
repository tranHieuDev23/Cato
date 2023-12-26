package jobs

import (
	"context"

	"github.com/tranHieuDev23/cato/internal/logic"
)

type SyncProblems interface {
	Run() error
}

type syncProblems struct {
	problemLogic logic.Problem
}

func NewSyncProblems(
	problemLogic logic.Problem,
) SyncProblems {
	return &syncProblems{
		problemLogic: problemLogic,
	}
}

func (j syncProblems) Run() error {
	return j.problemLogic.SyncProblemList(context.Background())
}
