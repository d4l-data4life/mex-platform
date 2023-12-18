package utils

import (
	"math"
	"time"
)

func ExponentialTicker(startDuration, maxDuration time.Duration, fn func()) chan<- struct{} {
	if startDuration > maxDuration {
		panic("start duration must not exceed max duration")
	}
	if fn == nil {
		panic("fn is nil")
	}

	currentDuration := startDuration

	ticker := time.NewTicker(currentDuration)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				ticker.Stop()
				fn()
				currentDuration = MinDuration(currentDuration*2, maxDuration)
				ticker = time.NewTicker(currentDuration)

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return quit
}

func MinDuration(a, b time.Duration) time.Duration {
	return time.Duration(math.Min(float64(a), float64(b)))
}
