package app

import (
	"context"
	"log/slog"

	"github.com/yandzee/wait-action/pkg/config"
	"github.com/yandzee/wait-action/pkg/poller"
	"github.com/yandzee/wait-action/pkg/tasks"

	"github.com/yandzee/wait-action/internal/github_client"
)

func Run(ctx context.Context, log *slog.Logger) error {
	cfg, err := config.ParseEnv()
	if err != nil {
		return err
	}

	log.Info("running with config", slog.Group("config", cfg.LogAttrs()...))

	workflows, err := tasks.Parse(cfg.Workflows)
	if err != nil {
		return err
	}

	gh, err := github_client.New(log.WithGroup("github-client"), cfg.GithubToken)
	if err != nil {
		return err
	}

	return poller.New(log, cfg, gh).Run(ctx, workflows)
}
