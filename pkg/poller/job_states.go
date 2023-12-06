package poller

type JobStates[T any] struct {
	Remaining T
	Done      T
	Failed    T
}

func (st *JobStates[int]) Print() {}
