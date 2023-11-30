package github

type RunFlags struct {
	IsFinished bool
	IsSuccess  bool
}

type CommitSpec struct {
	Branch string
	Sha    string
}

func (cs *CommitSpec) IsSha() bool {
	return len(cs.Sha) > 0
}

func (cs *CommitSpec) IsBranch() bool {
	return len(cs.Branch) > 0
}
