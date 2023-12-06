package poller

import (
	"github.com/yandzee/wait-action/pkg/github"
)

type WorkflowStates struct {
	Remaining github.WorkflowMap
	Failed    github.WorkflowMap
	Done      github.WorkflowMap
}

type State int

const (
	PendingState State = iota
	SuccessState
	FailedState
)

func (ws *WorkflowStates) ApplyRun(run *github.WorkflowRun) State {
	ws.ensureMaps()
	wf := run.Workflow
	st := PendingState

	switch {
	case !run.Flags.IsFinished:
		ws.Remaining[run.WorkflowId] = wf
	case run.Flags.IsSuccess:
		delete(ws.Remaining, run.WorkflowId)
		ws.Done[run.WorkflowId] = wf

		st = SuccessState
	default:
		delete(ws.Remaining, run.WorkflowId)
		ws.Failed[run.WorkflowId] = wf

		st = FailedState
	}

	return st
}

func (ws *WorkflowStates) HasRemaining() bool {
	return len(ws.Remaining) > 0
}

func (ws *WorkflowStates) HasFailures() bool {
	return len(ws.Failed) > 0
}

func (ws *WorkflowStates) ensureMaps() {
	if ws.Remaining == nil {
		ws.Remaining = make(github.WorkflowMap)
	}

	if ws.Failed == nil {
		ws.Failed = make(github.WorkflowMap)
	}

	if ws.Done == nil {
		ws.Done = make(github.WorkflowMap)
	}
}

func (s State) String() string {
	switch s {
	case PendingState:
		return "pending"
	case FailedState:
		return "failed"
	case SuccessState:
		return "success"
	}

	return "unknown"
}
