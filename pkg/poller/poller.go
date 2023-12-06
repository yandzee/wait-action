package poller

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/yandzee/wait-action/pkg/clock"
	"github.com/yandzee/wait-action/pkg/config"
	"github.com/yandzee/wait-action/pkg/github"
	"github.com/yandzee/wait-action/pkg/tasks"
)

type Poller[C clock.Clock] struct {
	log  *slog.Logger
	cfg  *config.Config
	clck C

	gh github.GithubClient
}

func New[C clock.Clock](
	log *slog.Logger,
	cfg *config.Config,
	clck C,
	gh github.GithubClient,
) *Poller[C] {
	return &Poller[C]{
		log:  log,
		cfg:  cfg,
		clck: clck,
		gh:   gh,
	}
}

func (p *Poller[C]) Run(ctx context.Context, t []tasks.WaitTask) error {
	// NOTE: Let's create so called "PollDescriptor" that is responsible for
	// tracking progress and saying if we are done
	desc := p.CreatePollDescriptor()

	for {
		// NOTE: Now we simply do poll iterations and on every such iteration
		// we are trying to input some new events/data into poll descriptor
		// regarding our progress
		isCompleted, hasFailures, err := p.Poll(ctx, desc, t)
		if err != nil {
			return fmt.Errorf("poll iteration failed: %s", err.Error())
		}

		if hasFailures {
			p.log.Info("Poller finished: some tasks are failed", desc.LogAttrs()...)
			return nil
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
		case <-p.clck.WaitChannel(p.cfg.PollDelay):
		}
	}
}

func (p *Poller[C]) Poll(
	ctx context.Context,
	desc *PollDescriptor,
	t []tasks.WaitTask,
) (bool, bool, error) {
	matcher := tasks.CreateWorkflowsMatcher(t)

	// NOTE: If matcher is trivial, we have no demand for waiting on workflows
	if matcher.IsTrivial() {
		return true, false, nil
	}

	workflowRuns, err := p.gh.GetWorkflowRuns(
		ctx,
		p.cfg.RepoOwner,
		p.cfg.Repo,
		p.cfg.Head,
	)

	if err != nil {
		return false, false, err
	}

	desc.ApplyWorkflowRuns(matcher, workflowRuns)
	return !desc.HasRemaining(), desc.HasFailures(), nil
}

func (p *Poller[C]) CreatePollDescriptor() *PollDescriptor {
	return NewPollDescriptor(p.log)
}
