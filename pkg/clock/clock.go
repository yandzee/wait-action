package clock

import "time"

type Clock interface {
	Now() int64
	WaitChannel(time.Duration) <-chan time.Time
}
