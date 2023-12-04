package github

import "log/slog"

type RunFlags struct {
	IsFinished bool
	IsSuccess  bool
}

type CommitSpec struct {
	Branch string
	Sha    string
}

func (cs *CommitSpec) LogAttrs() []any {
	return []any{
		slog.String("branch", cs.Branch),
		slog.String("sha", cs.Sha),
	}
}

func (cs *CommitSpec) IsSha() bool {
	return len(cs.Sha) > 0
}

func (cs *CommitSpec) IsBranch() bool {
	return len(cs.Branch) > 0
}

func (rf *RunFlags) LogAttrs() []any {
	return []any{
		slog.Bool("is-finished", rf.IsFinished),
		slog.Bool("is-success", rf.IsSuccess),
	}
}
