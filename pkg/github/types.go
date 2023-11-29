package github

type RunFlags struct {
	IsFinished bool
	IsSuccess  bool
}

type CommitSpec struct {
	Branch string
	Sha    string
}
