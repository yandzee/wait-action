package poller

import (
	"context"
	"log/slog"

	"github.com/google/go-github/v56/github"

	"github.com/yandzee/wait-action/internal/config"
	"github.com/yandzee/wait-action/internal/tasks"
)

type Poller struct {
	log *slog.Logger
	cfg *config.Config

	gh *github.Client
}

func New(log *slog.Logger, cfg *config.Config) *Poller {
	client := github.NewClient(nil).WithAuthToken(cfg.GithubToken)

	return &Poller{
		log: log,
		cfg: cfg,
		gh:  client,
	}
}

func (p *Poller) Run(ctx context.Context, t []tasks.WaitTask) error {
	return nil
}
