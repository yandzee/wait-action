package tests

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/yandzee/wait-action/pkg/clock"
	"github.com/yandzee/wait-action/pkg/config"
	"github.com/yandzee/wait-action/pkg/github"
	"github.com/yandzee/wait-action/pkg/poller"
	"github.com/yandzee/wait-action/pkg/tasks"
)

const (
	wfPath1 = ".github/workflows/workflow-1.yaml"
	wfPath2 = ".github/workflows/workflow-2.yaml"
)

func TestEmpty(t *testing.T) {
	runTest(
		t,
		[]tasks.WaitTask{{Workflows: []string{}}},
		[][]TestWorkflowRun{},
		map[int][]bool{},
	)
}

func TestInitiallySuccessWorkflows(t *testing.T) {
	runTest(
		t,
		[]tasks.WaitTask{{Workflows: []string{wfPath1}}},
		[][]TestWorkflowRun{
			{
				{
					WorkflowId:    1,
					WorkflowRunId: 1,
					Path:          wfPath1,
					IsFinished:    true,
					IsSuccess:     true,
				},
			},
		},
		map[int][]bool{0: {true, false}},
	)
}

func TestInitiallyFailedWorkflows(t *testing.T) {
	runTest(
		t,
		[]tasks.WaitTask{{Workflows: []string{wfPath1}}},
		[][]TestWorkflowRun{
			{
				{
					WorkflowId:    1,
					WorkflowRunId: 1,
					Path:          wfPath1,
					IsFinished:    true,
					IsSuccess:     false,
				},
			},
		},
		map[int][]bool{0: {true, true}},
	)
}

func TestNonMatchingSuccessWorkflows(t *testing.T) {
	runTest(
		t,
		[]tasks.WaitTask{{Workflows: []string{wfPath1}}},
		[][]TestWorkflowRun{
			{
				{
					WorkflowId:    1,
					WorkflowRunId: 1,
					Path:          wfPath2,
					IsFinished:    false,
					IsSuccess:     false,
				},
			},
			{
				{
					WorkflowId:    1,
					WorkflowRunId: 2,
					Path:          wfPath2,
					IsFinished:    true,
					IsSuccess:     true,
				},
			},
		},
		map[int][]bool{},
	)
}

func TestNonMatchingFailedWorkflows(t *testing.T) {
	runTest(
		t,
		[]tasks.WaitTask{{Workflows: []string{wfPath1}}},
		[][]TestWorkflowRun{
			{
				{
					WorkflowId:    1,
					WorkflowRunId: 1,
					Path:          wfPath2,
					IsFinished:    false,
					IsSuccess:     false,
				},
			},
			{
				{
					WorkflowId:    1,
					WorkflowRunId: 2,
					Path:          wfPath2,
					IsFinished:    true,
					IsSuccess:     false,
				},
			},
		},
		map[int][]bool{0: {true, false}, 1: {true, false}},
	)
}

func TestMatchingSuccessWorkflows(t *testing.T) {
	runTest(
		t,
		[]tasks.WaitTask{{Workflows: []string{wfPath1}}},
		[][]TestWorkflowRun{
			{
				{
					WorkflowId:    1,
					WorkflowRunId: 1,
					Path:          wfPath1,
					IsFinished:    false,
					IsSuccess:     false,
				},
			},
			{
				{
					WorkflowId:    1,
					WorkflowRunId: 2,
					Path:          wfPath1,
					IsFinished:    true,
					IsSuccess:     true,
				},
			},
		},
		map[int][]bool{0: {false, false}, 1: {true, false}},
	)
}

func TestMatchingFailedWorkflows(t *testing.T) {
	runTest(
		t,
		[]tasks.WaitTask{{Workflows: []string{wfPath1}}},
		[][]TestWorkflowRun{
			{
				{
					WorkflowId:    1,
					WorkflowRunId: 1,
					Path:          wfPath1,
					IsFinished:    false,
					IsSuccess:     false,
				},
			},
			{
				{
					WorkflowId:    1,
					WorkflowRunId: 2,
					Path:          wfPath1,
					IsFinished:    true,
					IsSuccess:     false,
				},
			},
		},
		map[int][]bool{0: {false, false}, 1: {true, true}},
	)
}

func runTest(
	t *testing.T,
	wt []tasks.WaitTask,
	wfr [][]TestWorkflowRun,
	runExpectations map[int][]bool,
) {
	ctx, p, desc := initPoller(wfr)

	lastIsDone := !desc.HasRemaining()
	lastHasFailures := desc.HasFailures()
	var err error

	if !lastIsDone {
		t.Fatal("descriptor is not done initially")
	}

	if lastHasFailures {
		t.Fatal("descriptor has failures initially")
	}

	for idx := range wfr {
		lastIsDone, lastHasFailures, err = p.Poll(ctx, desc, wt)
		if err != nil {
			t.Fatalf("err is not nil: %s\n", err.Error())
		}

		expected, ok := runExpectations[idx]
		if !ok {
			continue
		}

		if lastIsDone != expected[0] {
			t.Fatalf(
				"run %d: done %v, expected: %v\n",
				idx,
				lastIsDone,
				expected[0],
			)
		}

		if lastHasFailures != expected[1] {
			t.Fatalf(
				"run %d: hasFailures: %v, expected: %v\n",
				idx,
				lastHasFailures,
				expected[1],
			)
		}
	}

	for i := 0; i < 100; i += 1 {
		isDone, hasFailures, err := p.Poll(ctx, desc, wt)
		if err != nil {
			t.Fatalf("After all runs: poll crashed: %s\n", err.Error())
		}

		if isDone != lastIsDone {
			t.Fatalf("After all runs: isDone %v, last: %v\n", isDone, lastIsDone)
		}

		if hasFailures != lastHasFailures {
			t.Fatalf("After all runs: hasFailures: %v, last: %v\n", hasFailures, lastHasFailures)
		}
	}
}

func initPoller(mockedRuns [][]TestWorkflowRun) (
	context.Context,
	*poller.Poller[*clock.MockClock],
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
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	p := poller.New(log, cfg, &clock.MockClock{}, ghClient)

	return context.Background(), p, p.CreatePollDescriptor()
}
