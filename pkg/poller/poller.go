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
	result := PollResult{
		log: p.log,
	}

	for {
		// NOTE: Now we simply do poll iterations and on every such iteration
		// we are trying to merge some new events/data into PollResult
		// regarding our progress
		presult, err := p.Poll(ctx, t)
		if err != nil {
			return fmt.Errorf("Poll() failed: %s", err.Error())
		}

		result.MergeInPlace(presult)

		// if hasFailures {
		// 	p.log.Info("Poller finished: some tasks are failed", result.LogAttrs()...)
		// 	return nil
		// }
		//
		// if isCompleted {
		// 	p.log.Info("Poller finished: all tasks done", desc.LogAttrs()...)
		// 	return nil
		// }
		//
		// p.log.
		// 	With(desc.LogAttrs()...).
		// 	With("delay", p.cfg.PollDelay.String()).
		// 	Info("waiting")

		select {
		case <-ctx.Done():
			p.log.Warn("Poller context Done triggered", "err", ctx.Err().Error())

			return ctx.Err()
		case <-p.clck.WaitChannel(p.cfg.PollDelay):
		}
	}
}

func (p *Poller[C]) Poll(ctx context.Context, t []tasks.WaitTask) (*PollResult, error) {
	matcher := tasks.CreateWorkflowsMatcher(t)
	result := &PollResult{}

	// NOTE: If matcher is trivial, we have no demand for waiting on workflows
	if matcher.IsTrivial() {
		return result, nil
	}

	workflowRuns, err := p.gh.GetWorkflowRuns(
		ctx,
		p.cfg.RepoOwner,
		p.cfg.Repo,
		p.cfg.Head,
	)

	if err != nil {
		return nil, err
	}

	result.ApplyWorkflowRuns(matcher, workflowRuns)
	return result, nil
}
