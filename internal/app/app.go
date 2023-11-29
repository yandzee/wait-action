package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/yandzee/wait-action/internal/config"
	"github.com/yandzee/wait-action/internal/github_client"
	"github.com/yandzee/wait-action/internal/poller"
	"github.com/yandzee/wait-action/internal/tasks"
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
