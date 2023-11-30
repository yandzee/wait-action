package github

type Workflow struct {
	Id   int64
	Name string
	Path string
}

type Workflows []*Workflow
