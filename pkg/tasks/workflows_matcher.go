package tasks

import "github.com/yandzee/wait-action/pkg/github"

type WorkflowsMatcher struct {
	wfFilter github.WorkflowFilter
}

func CreateWorkflowsMatcher(t []WaitTask) *WorkflowsMatcher {
	matcher := &WorkflowsMatcher{
		wfFilter: github.WorkflowFilter{
			Paths: []string{},
		},
	}

	for _, task := range t {
		matcher.wfFilter.AddPaths(task.Workflows)
	}

	return matcher
}

func (m *WorkflowsMatcher) Matches(wf *github.Workflow) bool {
	if m.IsTrivial() {
		return true
	}

	if wf == nil {
		return false
	}

	return m.wfFilter.PathMatches(wf.Path)
}

func (m *WorkflowsMatcher) IsTrivial() bool {
	return m.wfFilter.IsTrivial()
}
