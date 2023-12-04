package github

import "log/slog"

type Workflow struct {
	Id   int64
	Name string
	Path string
}

type Workflows []*Workflow

func (wf *Workflow) LogAttrs() []any {
	return []any{
		slog.Int64("id", wf.Id),
		slog.String("name", wf.Name),
		slog.String("path", wf.Path),
	}
}
