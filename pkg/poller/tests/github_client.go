package tests

import (
	"context"

	"github.com/yandzee/wait-action/pkg/github"
)

type MockedGithub struct {
	runs [][]TestWorkflowRun
}

func initMockedGithubClient(runs [][]TestWorkflowRun) github.GithubClient {
	return &MockedGithub{
		runs: runs,
	}
}

func (gh *MockedGithub) GetWorkflowRuns(
	ctx context.Context,
	repoOwner, repoName string,
	commit github.CommitSpec,
) ([]*github.WorkflowRun, error) {
	if len(gh.runs) == 0 {
		return []*github.WorkflowRun{}, nil
	}

	runs := gh.runs[0]
	gh.runs = gh.runs[1:]

	ghRuns := make([]*github.WorkflowRun, len(runs))
	for i, run := range runs {
		built := run.Build()
		ghRuns[i] = &built
	}

	return ghRuns, nil
}
