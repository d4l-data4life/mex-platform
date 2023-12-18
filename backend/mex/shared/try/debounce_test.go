package try

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// This test creates a number of actions of which all odd-numbered ones always fail.
// We record the number of fails and check them after queue flush.
func Test_Debounce(t *testing.T) {
	const count = 10
	const maxAttempts = 5
	fails := make(map[int]int)

	makeAction := func(i int, willFail bool) func() error {
		return func() error {
			if willFail {
				time.Sleep(20 * time.Millisecond)
				fails[i]++
				return fmt.Errorf("error in action %d", i)
			}
			time.Sleep(30 * time.Millisecond)
			return nil
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	q, done := NewRequeue(ctx, nil, maxAttempts, 10*time.Millisecond)

	for i := 0; i < count; i++ {
		// Every odd action will always fail
		assert.True(t, q.Push(makeAction(i, i%2 == 1)))
	}

	cancel()

	// The below sleep allows the context cancellation to reach the queue.
	// Without the sleep we will get the 99 action queued and handled.
	// This is not a problem (it is basically a case like the action in the above loop),
	// but at this point we want to test the closing of the queue after cancellation.
	time.Sleep(2 * time.Millisecond)
	assert.False(t, q.Push(makeAction(99, true)))

	// Wait for the queue to flush.
	<-done

	for i := 0; i < count; i++ {
		if i%2 == 1 {
			assert.Equal(t, maxAttempts, fails[i])
		}
	}
}
