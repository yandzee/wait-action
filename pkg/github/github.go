package github

import "context"

type GithubClient interface {
	GetWorkflowRuns(context.Context, string, string, CommitSpec) ([]*WorkflowRun, error)
}
