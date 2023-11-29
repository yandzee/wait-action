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

func (gh *GithubClient) GetWorkflowRuns(
	ctx context.Context,
	repoOwner, repoName string,
	ref github.CommitSpec,
) ([]*github.WorkflowRun, error) {
	return []*github.WorkflowRun{}, nil
}
