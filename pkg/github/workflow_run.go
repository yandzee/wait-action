package github

import "log/slog"

type WorkflowRun struct {
	Id      int64
	Name    string
	HtmlURL string

	WorkflowId int64
	Workflow   *Workflow

	Commit CommitSpec
	Flags  RunFlags
}

type WorkflowRuns []*WorkflowRun

func (wfr *WorkflowRun) LogAttrs() []any {
	attrs := []any{
		slog.Int64("id", wfr.Id),
		slog.String("name", wfr.Name),
		slog.String("html-url", wfr.HtmlURL),
		slog.Int64("wf-id", wfr.WorkflowId),
		slog.Group("flag", wfr.Flags.LogAttrs()...),
		slog.Group("commit", wfr.Commit.LogAttrs()...),
	}

	if wfr.Workflow != nil {
		attrs = append(attrs, slog.Group("wf", wfr.Workflow.LogAttrs()...))
	}

	return attrs
}
