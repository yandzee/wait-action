package tests

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

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
	t.Parallel()

	ctx, cancel, p := initPoller(time.Second, wfr)
	defer cancel()

	expectedIdResults := map[int64]bool{}
	paths := map[string]struct{}{}

	for _, w := range wt {
		for _, wf := range w.Workflows {
			paths[wf] = struct{}{}
		}
	}

	for _, runs := range wfr {
		for _, run := range runs {
			if _, ok := paths[run.Path]; !ok {
				continue
			}

			expectedIdResults[run.WorkflowId] = run.IsSuccess
		}
	}

	expectedSuccess := true
	for _, success := range expectedIdResults {
		expectedSuccess = expectedSuccess && success
	}

	result, err := p.Run(ctx, wt)
	if err != nil {
		t.Fatalf("poller.Run() gives an error: %s %v\n", err.Error(), result.LogAttrs())
	}

	if result.HasFailures() && expectedSuccess {
		t.Fatalf("success is expected, results: %v\n", expectedIdResults)
	} else if !result.HasFailures() && !expectedSuccess {
		t.Fatalf("failure is expected, results: %v\n", expectedIdResults)
	}
}

func initPoller(timeout time.Duration, mockedRuns [][]TestWorkflowRun) (
	context.Context,
	context.CancelFunc,
	*poller.Poller[*NoWaitMockClock],
) {
	cfg := &config.Config{
		GithubToken:    "",
		PollDelay:      0,
		RepoOwner:      "owner",
		Repo:           "repo",
		Head:           github.CommitSpec{},
		Workflows:      "",
		IsDebugEnabled: true,
	}

	ghClient := initMockedGithubClient(mockedRuns)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	p := poller.New(log, cfg, &NoWaitMockClock{}, ghClient)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	return ctx, cancel, p
}
