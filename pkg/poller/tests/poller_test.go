package tests

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/yandzee/wait-action/pkg/config"
	"github.com/yandzee/wait-action/pkg/github"
	"github.com/yandzee/wait-action/pkg/poller"
	"github.com/yandzee/wait-action/pkg/tasks"
)

func TestEmpty(t *testing.T) {
	ctx, p, desc := initPoller([][]TestWorkflowRun{})
	tasks := []tasks.WaitTask{}

	for i := 0; i < 100; i += 1 {
		isDone, hasFailures, err := p.Poll(ctx, desc, tasks)

		if err != nil {
			t.Fatalf("err is not nil: %s\n", err.Error())
		}

		if !isDone {
			t.Fatal("poll descriptor is not done")
		}

		if hasFailures {
			t.Fatal("has failures")
		}
	}
}

func TestOnFinishedSuccessWorkflows(t *testing.T) {
	wfPath := ".github/mocked-workflow.yaml"
	ctx, p, desc := initPoller([][]TestWorkflowRun{
		{
			{
				WorkflowId:    1,
				WorkflowRunId: 1,
				Path:          wfPath,
				IsFinished:    true,
				IsSuccess:     true,
			},
		},
	})

	tasks := []tasks.WaitTask{
		{
			Workflows: []string{wfPath},
		},
	}

	for i := 0; i < 100; i += 1 {
		isDone, hasFailures, err := p.Poll(ctx, desc, tasks)

		if err != nil {
			t.Fatalf("err is not nil: %s\n", err.Error())
		}

		if !isDone {
			t.Fatal("poll descriptor is not done")
		}

		if hasFailures {
			t.Fatal("has failures")
		}
	}
}

func TestOnFinishedFailedWorkflows(t *testing.T) {
	wfPath := ".github/mocked-workflow.yaml"
	ctx, p, desc := initPoller([][]TestWorkflowRun{
		{
			{
				WorkflowId:    1,
				WorkflowRunId: 1,
				Path:          wfPath,
				IsFinished:    true,
				IsSuccess:     false,
			},
		},
	})

	tasks := []tasks.WaitTask{
		{
			Workflows: []string{wfPath},
		},
	}

	for i := 0; i < 100; i += 1 {
		isDone, hasFailures, err := p.Poll(ctx, desc, tasks)

		if err != nil {
			t.Fatalf("err is not nil: %s\n", err.Error())
		}

		if !isDone {
			t.Fatal("poll descriptor is not done")
		}

		if !hasFailures {
			t.Fatal("no failures reported")
		}
	}
}

func initPoller(mockedRuns [][]TestWorkflowRun) (
	context.Context,
	*poller.Poller,
	*poller.PollDescriptor,
) {
	cfg := &config.Config{
		GithubToken: "",
		PollDelay:   0,
		RepoOwner:   "owner",
		Repo:        "repo",
		Head:        github.CommitSpec{},
		Workflows:   "",
	}

	ghClient := initMockedGithubClient(mockedRuns)
	p := poller.New(slog.New(slog.NewTextHandler(io.Discard, nil)), cfg, ghClient)

	return context.Background(), p, poller.NewPollDescriptor()
}
