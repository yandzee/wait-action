package poller

import (
	"log/slog"

	"github.com/yandzee/wait-action/pkg/github"
	"github.com/yandzee/wait-action/pkg/tasks"
)

type PollDescriptor struct {
	Workflows *PollingTuple[github.WorkflowRuns]
}

type PollingTuple[T any] struct {
	Remaining T
	Done      T
}

func NewPollDescriptor() *PollDescriptor {
	return &PollDescriptor{
		Workflows: &PollingTuple[github.WorkflowRuns]{
			Remaining: github.WorkflowRuns{},
			Done:      github.WorkflowRuns{},
		},
	}
}

func (pd *PollDescriptor) ApplyWorkflowRuns(
	matcher *tasks.WorkflowsMatcher,
	wfRuns github.WorkflowRuns,
) {
	for _, wfRun := range wfRuns {
		if !matcher.Matches(wfRun.Workflow) {
			continue
		}
	}
}

func (pd *PollDescriptor) LogAttrs() []any {
	return []any{
		slog.Group("workflows",
			slog.Int("remaining", len(pd.Workflows.Remaining)),
			slog.Int("done", len(pd.Workflows.Done)),
		),
	}
}
