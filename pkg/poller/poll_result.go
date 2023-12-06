package poller

import (
	"log/slog"

	"github.com/yandzee/wait-action/pkg/github"
	"github.com/yandzee/wait-action/pkg/tasks"
)

type PollResult struct {
	log       *slog.Logger
	Workflows *JobStates[github.WorkflowMap]
}

func NewPollResult(log *slog.Logger) *PollResult {
	return &PollResult{
		log: log,
		Workflows: &JobStates[github.WorkflowMap]{
			Remaining: make(github.WorkflowMap),
			Done:      make(github.WorkflowMap),
			Failed:    make(github.WorkflowMap),
		},
	}
}

func (pd *PollResult) ApplyWorkflowRuns(
	matcher *tasks.WorkflowsMatcher,
	wfRuns github.WorkflowRuns,
) {
	for _, wfRun := range wfRuns {
		wf := wfRun.Workflow
		if !matcher.Matches(wf) {
			pd.log.Debug("workflow skipped by matcher", wfRun.LogAttrs()...)
			continue
		}

		switch {
		case !wfRun.Flags.IsFinished:
			pd.Workflows.Remaining[wfRun.WorkflowId] = wf

			pd.log.Info("remaining workflow run", wfRun.LogAttrs()...)
		case wfRun.Flags.IsSuccess:
			delete(pd.Workflows.Remaining, wfRun.WorkflowId)
			pd.Workflows.Done[wfRun.WorkflowId] = wf

			pd.log.Info("workflow run successfully completed", wfRun.LogAttrs()...)
		default:
			delete(pd.Workflows.Remaining, wfRun.WorkflowId)
			pd.Workflows.Failed[wfRun.WorkflowId] = wf

			pd.log.Info("workflow run failed", wfRun.LogAttrs()...)
		}
	}
}

func (pd *PollResult) HasRemaining() bool {
	return len(pd.Workflows.Remaining) > 0
}

func (pd *PollResult) HasFailures() bool {
	return len(pd.Workflows.Failed) > 0
}

func (pd *PollResult) LogAttrs() []any {
	return []any{
		slog.Group("workflows",
			slog.Int("remaining", len(pd.Workflows.Remaining)),
			slog.Int("done", len(pd.Workflows.Done)),
			slog.Int("failed", len(pd.Workflows.Failed)),
		),
	}
}
