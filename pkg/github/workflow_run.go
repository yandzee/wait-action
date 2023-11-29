package github

type WorkflowRun struct {
	Id   int64
	Name string
	Path string

	Commit CommitSpec
	Flags  RunFlags
}

type WorkflowRuns []*WorkflowRun
