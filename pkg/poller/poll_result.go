package poller

import (
	"log/slog"

	"github.com/yandzee/wait-action/pkg/github"
	"github.com/yandzee/wait-action/pkg/tasks"
)

type PollResult struct {
	log       *slog.Logger
	Workflows WorkflowStates
}

func NewPollResult(log *slog.Logger) *PollResult {
	return &PollResult{
		log:       log,
		Workflows: WorkflowStates{},
	}
}

func (pr *PollResult) ApplyWorkflowRuns(
	matcher *tasks.WorkflowsMatcher,
	runs github.WorkflowRuns,
) {
	for _, run := range runs {
		wf := run.Workflow
		if !matcher.Matches(wf) {
			pr.log.Debug("workflow skipped by matcher", run.LogAttrs()...)
			continue
		}

		state := pr.Workflows.ApplyRun(run)
		pr.log.
			With(run.LogAttrs()...).
			Info("workflow run state updated", slog.String("state", state.String()))
	}
}

func (pr *PollResult) MergeInPlace(rhs *PollResult) {
	pr.Workflows.Merge(&rhs.Workflows)
}

func (pr *PollResult) HasRemaining() bool {
	return pr.Workflows.HasRemaining()
}

func (pr *PollResult) HasFailures() bool {
	return pr.Workflows.HasFailures()
}

func (pr *PollResult) LogAttrs() []any {
	return []any{
		slog.Group("workflows",
			slog.Int("remaining", len(pr.Workflows.Remaining)),
			slog.Int("done", len(pr.Workflows.Done)),
			slog.Int("failed", len(pr.Workflows.Failed)),
		),
	}
}
