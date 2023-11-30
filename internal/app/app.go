package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/yandzee/wait-action/pkg/config"
	"github.com/yandzee/wait-action/pkg/poller"
	"github.com/yandzee/wait-action/pkg/tasks"

	"github.com/yandzee/wait-action/internal/github_client"
)

func Run(ctx context.Context, log *slog.Logger) error {
	workflows, err := tasks.Parse(os.Getenv("INPUT_WORKFLOWS"))
	if err != nil {
		return err
	}

	cfg, err := config.ParseEnv()
	if err != nil {
		return err
	}

	gh, err := github_client.New(log.WithGroup("github-client"), cfg.GithubToken)
	if err != nil {
		return err
	}

	return poller.New(log, cfg, gh).Run(ctx, workflows)
}
