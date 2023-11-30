package tests

import (
	"fmt"

	"github.com/yandzee/wait-action/pkg/github"
)

type TestWorkflowRun struct {
	WorkflowId    int64
	WorkflowRunId int64
	Path          string
	IsFinished    bool
	IsSuccess     bool
}

func (tw *TestWorkflowRun) Build() github.WorkflowRun {
	return github.WorkflowRun{
		Id:         tw.WorkflowRunId,
		Name:       fmt.Sprintf("wf-run-id-%d", tw.WorkflowRunId),
		WorkflowId: tw.WorkflowId,
		Workflow: &github.Workflow{
			Id:   tw.WorkflowId,
			Name: fmt.Sprintf("wf-id-%d", tw.WorkflowId),
			Path: tw.Path,
		},
		Commit: github.CommitSpec{},
		Flags: github.RunFlags{
			IsSuccess:  tw.IsSuccess,
			IsFinished: tw.IsFinished,
		},
	}
}
