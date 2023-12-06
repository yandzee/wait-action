package poller

type State int

const (
	PendingState State = iota
	SuccessState
	FailedState
)

func (s State) String() string {
	switch s {
	case PendingState:
		return "pending"
	case FailedState:
		return "failed"
	case SuccessState:
		return "success"
	}

	return "unknown"
}
