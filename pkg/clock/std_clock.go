package clock

import "time"

type StdClock struct{}

func (sc *StdClock) WaitChannel(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func (sc *StdClock) Now() int64 {
	return time.Now().UnixNano()
}
