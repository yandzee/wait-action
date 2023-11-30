package github

type WorkflowFilter struct {
	Paths []string

	pathsMap map[string]struct{}
}

func (wff *WorkflowFilter) PathMatches(p string) bool {
	if wff == nil || len(wff.Paths) == 0 {
		return true
	}

	wff.ensurePathsMap()

	_, ok := wff.pathsMap[p]
	return ok
}

func (wff *WorkflowFilter) AddPaths(ps []string) {
	wff.ensurePathsMap()

	wff.Paths = append(wff.Paths, ps...)

	for _, p := range ps {
		wff.pathsMap[p] = struct{}{}
	}
}

func (wff *WorkflowFilter) ensurePathsMap() {
	if wff.pathsMap != nil {
		return
	}

	wff.pathsMap = make(map[string]struct{})

	for _, p := range wff.Paths {
		wff.pathsMap[p] = struct{}{}
	}
}

func (wff *WorkflowFilter) IsTrivial() bool {
	return len(wff.Paths) == 0
}
