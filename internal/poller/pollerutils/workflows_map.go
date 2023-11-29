package pollerutils

import "github.com/google/go-github/v56/github"

type WorkflowsMap map[int64]*github.Workflow

func (m WorkflowsMap) Save(wf *github.Workflow) {
	m[*wf.ID] = wf
}

func (m WorkflowsMap) Keys() []int64 {
	keys := make([]int64, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}

	return keys
}
