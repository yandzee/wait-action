package poller

import (
	"github.com/yandzee/wait-action/pkg/github"
)

type WorkflowStates struct {
	Remaining github.WorkflowMap
	Failed    github.WorkflowMap
	Done      github.WorkflowMap
}

func (ws *WorkflowStates) HasRemaining() bool {
	return len(ws.Remaining) > 0
}

func (ws *WorkflowStates) HasFailures() bool {
	return len(ws.Failed) > 0
}

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

func (ws *WorkflowStates) Merge(rhs *WorkflowStates) {
	ws.ensureMaps()

	ws.mergeMap(ws.Remaining, rhs.Remaining)
	ws.dropEntriesIfInMap(rhs.Remaining, ws.Done, ws.Failed)

	ws.mergeMap(ws.Done, rhs.Done)
	ws.dropEntriesIfInMap(rhs.Done, ws.Remaining, ws.Failed)

	ws.mergeMap(ws.Failed, rhs.Failed)
	ws.dropEntriesIfInMap(rhs.Failed, ws.Remaining, ws.Done)
}

func (ws *WorkflowStates) dropEntriesIfInMap(
	m github.WorkflowMap, targets ...github.WorkflowMap,
) {
	for _, target := range targets {
		for id := range target {
			if _, ok := m[id]; ok {
				delete(target, id)
			}
		}
	}
}

func (ws *WorkflowStates) mergeMap(dst, src github.WorkflowMap) {
	for id, wf := range src {
		dst[id] = wf
	}
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
