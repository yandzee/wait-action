package github_client

import (
	"log/slog"

	ghclient "github.com/google/go-github/v56/github"
)

func (gh *GithubClient) workflowRunAttrs(wf *ghclient.WorkflowRun) []any {
	return []any{
		slog.Int64("id", wf.GetID()),
		slog.String("name", wf.GetName()),
		slog.String("status", wf.GetStatus()),
		slog.String("conclusion", wf.GetConclusion()),
		slog.String("status", wf.GetStatus()),
		slog.String("html-url", wf.GetHTMLURL()),
	}
}
