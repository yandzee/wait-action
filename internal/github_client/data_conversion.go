package github_client

import (
	ghclient "github.com/google/go-github/v56/github"
	"github.com/yandzee/wait-action/pkg/github"
)

func (gh *GithubClient) convertWorkflowRuns(
	wfr []*ghclient.WorkflowRun,
	wfs github.Workflows,
) github.WorkflowRuns {
	converted := make(github.WorkflowRuns, 0, len(wfr))
	wfsMap := github.CreateWorkflowMap(wfs)

	for _, wf := range wfr {
		cwf := gh.convertWorkflowRun(wf)
		cwf.Workflow = wfsMap[cwf.WorkflowId]

		if cwf.Workflow == nil {
			gh.log.Error("workflow wasn't found for workflow run",
				"workflow-run-id", cwf.Id,
				"workflow-id", cwf.WorkflowId,
			)
		}

		converted = append(converted, cwf)
	}

	return converted
}

func (gh *GithubClient) convertWorkflowRun(wf *ghclient.WorkflowRun) *github.WorkflowRun {
	return &github.WorkflowRun{
		Id:         wf.GetID(),
		Name:       wf.GetName(),
		WorkflowId: wf.GetWorkflowID(),
		Workflow:   nil,
		Commit:     gh.commitSpecFromWorkflowRun(wf),
		Flags:      gh.runFlagsFromWorkflowRun(wf),
	}
}

func (gh *GithubClient) commitSpecFromWorkflowRun(wf *ghclient.WorkflowRun) github.CommitSpec {
	return github.CommitSpec{
		Sha:    wf.GetHeadSHA(),
		Branch: wf.GetHeadBranch(),
	}
}

func (gh *GithubClient) runFlagsFromWorkflowRun(wf *ghclient.WorkflowRun) github.RunFlags {
	isFinished := true
	isSuccess := false

	switch wf.GetConclusion() {
	case "pending":
		fallthrough
	case "waiting":
		fallthrough
	case "queued":
		isFinished = false
	case "success":
		isSuccess = true
	}

	return github.RunFlags{
		IsFinished: isFinished,
		IsSuccess:  isSuccess,
	}
}

func (gh *GithubClient) convertWorkflows(wfs []*ghclient.Workflow) []*github.Workflow {
	converted := make(github.Workflows, 0, len(wfs))

	for _, wf := range wfs {
		cwf := gh.convertWorkflow(wf)
		converted = append(converted, cwf)
	}

	return converted
}

func (gh *GithubClient) convertWorkflow(wf *ghclient.Workflow) *github.Workflow {
	return &github.Workflow{
		Id:   wf.GetID(),
		Name: wf.GetName(),
		Path: wf.GetPath(),
	}
}
