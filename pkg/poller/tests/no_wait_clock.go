package tests

import (
	"time"

	"github.com/yandzee/wait-action/pkg/clock"
)

type NoWaitMockClock struct {
	clock.MockClock
}

func (nwmc *NoWaitMockClock) WaitChannel(d time.Duration) <-chan time.Time {
	ch := nwmc.MockClock.WaitChannel(d)
	nwmc.MockClock.Advance(d.Nanoseconds())

	return ch
}
