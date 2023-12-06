package app

import (
	"context"
	"log/slog"
	"os"

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

	result, err := poller.New(log, cfg, &clock.StdClock{}, gh).Run(ctx, waitTasks)
	if err != nil {
		return err
	}

	switch {
	case result.HasFailures():
		log.
			With(result.LogAttrs()...).
			Error("finished: some tasks failed")

		// NOTE: Terminating with error code to break Github actions run
		os.Exit(1)
	case result.HasRemaining():
		log.
			With(result.LogAttrs()...).
			Error("poller terminated with remaining tasks not being empty")

		panic("Poller has terminated prematurely")
	default:
		log.Info("finished: all tasks done")
	}

	return nil
}
