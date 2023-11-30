package poller

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/yandzee/wait-action/pkg/config"
	"github.com/yandzee/wait-action/pkg/github"
	"github.com/yandzee/wait-action/pkg/tasks"
)

type Poller struct {
	log *slog.Logger
	cfg *config.Config

	gh github.GithubClient
}

func New(log *slog.Logger, cfg *config.Config, gh github.GithubClient) *Poller {
	return &Poller{
		log: log,
		cfg: cfg,
		gh:  gh,
	}
}

func (p *Poller) Run(ctx context.Context, t []tasks.WaitTask) error {
	// NOTE: Let's create so called "PollDescriptor" that is responsible for
	// tracking progress and saying if we are done
	desc := NewPollDescriptor()

	for {
		// NOTE: Now we simply do poll iterations and on every such iteration
		// we are trying to input some new events/data into poll descriptor
		// regarding our progress
		isCompleted, err := p.poll(ctx, desc, t)
		if err != nil {
			return fmt.Errorf("poll iteration failed: %s", err.Error())
		}

		if isCompleted {
			p.log.Info("Poller finished: all tasks done", desc.LogAttrs()...)
			return nil
		}

		p.log.
			With(desc.LogAttrs()...).
			With("delay", p.cfg.PollDelay.String()).
			Info("waiting")

		select {
		case <-ctx.Done():
			p.log.Warn("Poller context Done triggered", "err", ctx.Err().Error())

			return ctx.Err()
		case <-time.After(p.cfg.PollDelay):
		}
	}
}

func (p *Poller) poll(
	ctx context.Context,
	desc *PollDescriptor,
	t []tasks.WaitTask,
) (bool, error) {
	matcher := tasks.CreateWorkflowsMatcher(t)

	// NOTE: If matcher is trivial, we have no demand for waiting on workflows
	if matcher.IsTrivial() {
		return true, nil
	}

	workflowRuns, err := p.gh.GetWorkflowRuns(
		ctx,
		p.cfg.RepoOwner,
		p.cfg.Repo,
		p.cfg.Head,
	)

	if err != nil {
		return false, err
	}

	desc.ApplyWorkflowRuns(matcher, workflowRuns)
	return !desc.HasRemaining(), nil
}
