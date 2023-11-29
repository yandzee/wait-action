package poller

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/yandzee/wait-action/pkg/github"

	"github.com/yandzee/wait-action/internal/config"
	"github.com/yandzee/wait-action/internal/poller/pollerutils"
	"github.com/yandzee/wait-action/internal/tasks"
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
	desc := p.createPollDescriptor(t)

	for {
		// NOTE: Now we simply do poll iterations and on every such iteration
		// we are trying to input some new events/data into poll descriptor
		// regarding our progress
		isCompleted, err := p.poll(desc)
		if err != nil {
			return fmt.Errorf("failed to do poll iteration: %s", err.Error())
		}

		if isCompleted {
			p.log.Info("Poller finished: all tasks done", desc.LogAttrs())
			return nil
		}

		p.log.Info("Poller still waits for remaining tasks",
			"delay", p.cfg.PollDelay.String(), desc.RemainingLogAttrs(),
		)

		select {
		case <-ctx.Done():
			p.log.Warn("Poller context Done triggered", "err", ctx.Err().Error())

			return ctx.Err()
		case <-time.After(p.cfg.PollDelay):
		}
	}
	// watchedWorkflows, _, err := p.getWorkflowsFromTasks(ctx, t)
	// if err != nil {
	// 	return err
	// }
	//
	// if len(watchedWorkflows) == 0 {
	// 	return fmt.Errorf("workflows not found")
	// }
	//
	// for {
	// 	workflowsRuns, err := p.getWorkflowRunsByIds(watchedWorkflows.Keys())
	// 	if err != nil {
	// 		return fmt.Errorf("failed to get workflow runs: %s", err.Error())
	// 	}
	//
	// }
	//
	// return nil
}

func (p *Poller) createPollDescriptor(t []tasks.WaitTask) *PollDescriptor {
	desc := new(PollDescriptor)

	return desc
}

func (p *Poller) getWorkflowsFromTasks(
	ctx context.Context, t []tasks.WaitTask,
) (pollerutils.WorkflowsMap, *github.Workflows, error) {
	repoWorkflows, _, err := p.gh.Actions.ListWorkflows(
		ctx,
		p.cfg.RepoOwner,
		p.cfg.Repo,
		&github.ListOptions{
			PerPage: 100,
		},
	)

	if err != nil {
		return nil, nil, err
	}

	workflowsMap := pollerutils.FilterTaskWorkflows(repoWorkflows, t)
	return workflowsMap, repoWorkflows, nil
}
