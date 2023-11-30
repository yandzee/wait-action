package github

type WorkflowMap map[int64]*Workflow

func CreateWorkflowMap(wfs []*Workflow) WorkflowMap {
	m := make(WorkflowMap)

	for _, wf := range wfs {
		m[wf.Id] = wf
	}

	return m
}

