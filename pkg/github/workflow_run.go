package github

type WorkflowRun struct {
	Id   int64
	Name string

	WorkflowId int64
	Workflow   *Workflow

	Commit CommitSpec
	Flags  RunFlags
}

type WorkflowRuns []*WorkflowRun
