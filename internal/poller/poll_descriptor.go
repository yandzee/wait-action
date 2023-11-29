package poller

import (
	"log/slog"

	"github.com/yandzee/wait-action/pkg/github"
)

type PollDescriptor struct {
	Workflows *PollingTuple[github.WorkflowRuns]
}

func NewPollDescriptor() *PollDescriptor {
	return &PollDescriptor{
		Workflows: &PollingTuple[github.WorkflowRuns]{
			Remaining: github.WorkflowRuns{},
			Done:      github.WorkflowRuns{},
		},
	}
}

type PollingTuple[T any] struct {
	Remaining T
	Done      T
}

func (pd *PollDescriptor) LogAttrs() []any {
	return []any{
		slog.Group("workflows",
			slog.Int("remaining", len(pd.Workflows.Remaining)),
			slog.Int("done", len(pd.Workflows.Done)),
		),
	}
}
