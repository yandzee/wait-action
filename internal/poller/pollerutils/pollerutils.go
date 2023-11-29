package pollerutils

import (
	"github.com/google/go-github/v56/github"
	"github.com/yandzee/wait-action/internal/tasks"
)

func FilterTaskWorkflows(wfs *github.Workflows, t []tasks.WaitTask) WorkflowsMap {
	matcher := tasks.CreateWorkflowsMatcher(t)
	m := make(WorkflowsMap)

	for _, wf := range wfs.Workflows {
		if matcher.Matches(wf) {
			m.Save(wf)
		}
	}

	return m
}
