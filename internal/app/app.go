package app

import (
	"context"
	"log/slog"

	"github.com/yandzee/wait-action/pkg/clock"
	"github.com/yandzee/wait-action/pkg/config"
	"github.com/yandzee/wait-action/pkg/poller"
	"github.com/yandzee/wait-action/pkg/tasks"

	"github.com/yandzee/wait-action/internal/github_client"
)

func Run(ctx context.Context, log *slog.Logger, cfg *config.Config) error {
	log.Info("running with config", slog.Group("config", cfg.LogAttrs()...))

	waitTasks, err := tasks.Parse(cfg.Workflows)
	if err != nil {
		return err
	}

	gh, err := github_client.New(log.WithGroup("github-client"), cfg.GithubToken)
	if err != nil {
		return err
	}

	return poller.New(log, cfg, &clock.StdClock{}, gh).Run(ctx, waitTasks)
}
