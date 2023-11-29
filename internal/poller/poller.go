package poller

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/go-github/v56/github"

	"github.com/yandzee/wait-action/internal/config"
	"github.com/yandzee/wait-action/internal/poller/pollerutils"
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
	desc := p.createPollDescriptor(t)

	for {
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
