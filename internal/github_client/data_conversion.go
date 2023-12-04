package github_client

import (
	"log/slog"

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
	rf := gh.runFlagsFromStatusOrConclusion(wf.GetStatus(), wf.GetConclusion())

	gh.log.Debug("constructing RunFlags from github workflow run",
		slog.Group("run-flags", rf.LogAttrs()...),
		slog.Group("gh.WorkflowRun", gh.workflowRunAttrs(wf)...),
	)

	return rf
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

func (gh *GithubClient) runFlagsFromStatusOrConclusion(st, con string) github.RunFlags {
	rf := github.RunFlags{
		IsFinished: false,
		IsSuccess:  false,
	}

	switch con {
	case "success":
		rf.IsFinished, rf.IsSuccess = true, true
	case "action_required":
		fallthrough
	case "cancelled":
		fallthrough
	case "failure":
		fallthrough
	case "skipped":
		fallthrough
	case "completed":
		fallthrough
	case "timed_out":
		rf.IsFinished, rf.IsSuccess = true, false
	}

	if rf.IsSuccess {
		return rf
	}

	switch st {
	case "neutral":
		fallthrough
	case "stale":
		fallthrough
	case "pending":
		fallthrough
	case "waiting":
		fallthrough
	case "requested":
		fallthrough
	case "in_progress":
		fallthrough
	case "queued":
		rf.IsFinished, rf.IsSuccess = false, false
	case "completed":
		rf.IsFinished = true
	case "action_required":
		fallthrough
	case "cancelled":
		fallthrough
	case "failure":
		fallthrough
	case "skipped":
		fallthrough
	case "timed_out":
		rf.IsFinished, rf.IsSuccess = true, false
	case "success":
		rf.IsFinished, rf.IsSuccess = true, true
	}

	return rf
}
