package tasks

import "github.com/google/go-github/v56/github"

type WorkflowsMatcher struct {
	pathsMap map[string]struct{}
}

func CreateWorkflowsMatcher(t []WaitTask) *WorkflowsMatcher {
	matcher := &WorkflowsMatcher{
		pathsMap: make(map[string]struct{}),
	}

	for _, task := range t {
		for _, workflowPath := range task.Workflows {
			matcher.pathsMap[workflowPath] = struct{}{}
		}
	}

	return matcher
}

func (m *WorkflowsMatcher) Matches(wf *github.Workflow) bool {
	if wf.Path == nil {
		return false
	}

	_, matches := m.pathsMap[*wf.Path]
	return matches
}
