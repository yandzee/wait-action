package github_client

import (
	"context"
	"log/slog"

	ghclient "github.com/google/go-github/v56/github"
	"github.com/yandzee/wait-action/pkg/github"
)

type GithubClient struct {
	log *slog.Logger
	cl  *ghclient.Client
}

func New(log *slog.Logger, token string) (*GithubClient, error) {
	return &GithubClient{
		log: log,
		cl:  ghclient.NewClient(nil).WithAuthToken(token),
	}, nil
}

func (gh *GithubClient) GetRepoWorkflows(
	ctx context.Context,
	repoOwner, repoName string,
	filters *github.WorkflowFilter,
) ([]*github.Workflow, error) {
	repoWorkflows, _, err := gh.cl.Actions.ListWorkflows(
		ctx,
		repoOwner,
		repoName,
		&ghclient.ListOptions{
			PerPage: 100,
		},
	)

	if err != nil {
		return nil, err
	}

	if filters == nil || len(filters.Paths) == 0 {
		return gh.convertWorkflows(repoWorkflows.Workflows), nil
	}

	filteredWorkflows := make([]*ghclient.Workflow, 0)
	for _, wf := range repoWorkflows.Workflows {
		if !filters.PathMatches(wf.GetPath()) {
			continue
		}

		filteredWorkflows = append(filteredWorkflows, wf)
	}

	return gh.convertWorkflows(filteredWorkflows), nil
}

func (gh *GithubClient) GetWorkflowRuns(
	ctx context.Context,
	repoOwner, repoName string,
	ref github.CommitSpec,
) ([]*github.WorkflowRun, error) {
	repoWorkflows, err := gh.GetRepoWorkflows(ctx, repoOwner, repoName, nil)
	if err != nil {
		return nil, err
	}

	opts := ghclient.ListWorkflowRunsOptions{
		ListOptions: ghclient.ListOptions{
			PerPage: 100,
		},
	}

	if ref.IsSha() {
		opts.HeadSHA = ref.Sha
	} else if ref.IsBranch() {
		opts.Branch = ref.Branch
	}

	workflowRuns, _, err := gh.cl.Actions.ListRepositoryWorkflowRuns(
		ctx,
		repoOwner,
		repoName,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	return gh.convertWorkflowRuns(workflowRuns.WorkflowRuns, repoWorkflows), nil
}
