package tests

import (
	"context"

	"github.com/yandzee/wait-action/pkg/github"
)

type MockedGithub struct {
	runs []TestWorkflowRun
}

func initMockedGithubClient(runs []TestWorkflowRun) github.GithubClient {
	return &MockedGithub{
		runs: runs,
	}
}

func (gh *MockedGithub) GetWorkflowRuns(
	ctx context.Context,
	repoOwner, repoName string,
	commit github.CommitSpec,
) ([]*github.WorkflowRun, error) {
	return nil, nil
}
