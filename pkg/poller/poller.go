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

func (p *Poller[C]) Run(ctx context.Context, t []tasks.WaitTask) (*PollResult, error) {
	result := PollResult{
		log: p.log,
	}

	for {
		// NOTE: Now we simply do poll iterations and on every such iteration
		// we are trying to merge some new events/data into PollResult
		// regarding our progress
		presult, err := p.poll(ctx, t)
		if err != nil {
			return nil, fmt.Errorf("poller.Poll() failed: %s", err.Error())
		}

		// NOTE: Collecting all the result in single place
		result.MergeInPlace(presult)

		if result.HasFailures() || !result.HasRemaining() {
			return &result, nil
		}

		p.log.
			With(result.LogAttrs()...).
			With("delay", p.cfg.PollDelay.String()).
			Info("waiting before next poll")

		select {
		case <-ctx.Done():
			return &result, ctx.Err()
		case <-p.clck.WaitChannel(p.cfg.PollDelay):
		}
	}
}

func (p *Poller[C]) poll(ctx context.Context, t []tasks.WaitTask) (*PollResult, error) {
	matcher := tasks.CreateWorkflowsMatcher(t)
	result := &PollResult{
		log: p.log,
	}

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
