package poller

import (
	"log/slog"

	"github.com/yandzee/wait-action/pkg/github"
	"github.com/yandzee/wait-action/pkg/tasks"
)

type PollDescriptor struct {
	Workflows *JobState[github.WorkflowMap]
}

type JobState[T any] struct {
	Remaining T
	Done      T
	Failed    T
}

func NewPollDescriptor() *PollDescriptor {
	return &PollDescriptor{
		Workflows: &JobState[github.WorkflowMap]{
			Remaining: make(github.WorkflowMap),
			Done:      make(github.WorkflowMap),
			Failed:    make(github.WorkflowMap),
		},
	}
}

func (pd *PollDescriptor) ApplyWorkflowRuns(
	matcher *tasks.WorkflowsMatcher,
	wfRuns github.WorkflowRuns,
) {
	for _, wfRun := range wfRuns {
		wf := wfRun.Workflow
		if !matcher.Matches(wf) {
			continue
		}

		if !wfRun.Flags.IsFinished {
			pd.Workflows.Remaining[wfRun.WorkflowId] = wf
		} else if wfRun.Flags.IsSuccess {
			delete(pd.Workflows.Remaining, wfRun.WorkflowId)
			pd.Workflows.Done[wfRun.WorkflowId] = wf
		} else {
			delete(pd.Workflows.Remaining, wfRun.WorkflowId)
			pd.Workflows.Failed[wfRun.WorkflowId] = wf
		}
	}
}

func (pd *PollDescriptor) HasRemaining() bool {
	return len(pd.Workflows.Remaining) > 0
}

func (pd *PollDescriptor) HasFailures() bool {
	return len(pd.Workflows.Failed) > 0
}

func (pd *PollDescriptor) LogAttrs() []any {
	return []any{
		slog.Group("workflows",
			slog.Int("remaining", len(pd.Workflows.Remaining)),
			slog.Int("done", len(pd.Workflows.Done)),
		),
	}
}
