package github

type WorkflowFilter struct {
	Paths []string

	pathsMap map[string]struct{}
}

func (wff *WorkflowFilter) PathMatches(p string) bool {
	if wff == nil || len(wff.Paths) == 0 {
		return true
	}

	if len(wff.pathsMap) == 0 {
		wff.buildPathsMap()
	}

	_, ok := wff.pathsMap[p]
	return ok
}

func (wff *WorkflowFilter) buildPathsMap() {
	wff.pathsMap = make(map[string]struct{})

	for _, p := range wff.Paths {
		wff.pathsMap[p] = struct{}{}
	}
}
