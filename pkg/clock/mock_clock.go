package clock

import "time"

type MockClock struct {
	Current int64

	waiters waitersMap
}

type waitersMap map[int64][]chan time.Time

func (mc *MockClock) Now() int64 {
	return mc.Current
}

func (mc *MockClock) WaitChannel(d time.Duration) <-chan time.Time {
	waiter := make(chan time.Time, 1)
	targetTime := mc.Current + d.Nanoseconds()

	if mc.waiters == nil {
		mc.waiters = make(waitersMap)
	}

	mc.waiters[targetTime] = append(mc.waiters[targetTime], waiter)
	return waiter
}

func (mc *MockClock) Advance(d int64) int64 {
	mc.Current += d

	for targetTime, waiters := range mc.waiters {
		if targetTime > mc.Current {
			continue
		}

		for _, waiter := range waiters {
			waiter <- time.Unix(0, mc.Current)
		}

		delete(mc.waiters, targetTime)
	}

	return mc.Current
}
