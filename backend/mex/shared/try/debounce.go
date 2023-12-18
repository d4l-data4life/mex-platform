package try

import (
	"context"
	"sync"
	"time"

	L "github.com/d4l-data4life/mex/mex/shared/log"
)

// A requeue is a queue implementing the following behavior:
// - You push in an action representing a unit of self-contained work.
//   This function must return an error if unsuccessful, or nil if fine.
// - The action is sent to a queue and processed in arrival order.
// - If the action is unsuccessful, it is re-queued.
// - There is a maximum number of re-queuing.
// - There is a delay between execution of queue actions.
// - The context `ctx` of the maker function should be cancelable.
//   Upon context cancellation, the queue does not accept any new actions and
//   flushes its remaining contents.
// - The `done` channel is used to signal the end of flushing.

type Action = func() error

type Attempt struct {
	action   Action
	attempts uint
}

type Requeue struct {
	mu sync.Mutex

	q      chan Attempt
	wg     *sync.WaitGroup
	log    L.LevelLogger
	closed bool
}

// The sending queue logic is borrowed and adapted from here: https://github.com/gesundheitscloud/go-svc/blob/master/pkg/client/bievents.go
func NewRequeue(ctx context.Context, log L.LevelLogger, maxAttempts uint, pause time.Duration) (*Requeue, <-chan struct{}) {
	if log == nil {
		log = &L.NullLogger{}
	}

	// This signal channel is returned so others can wait for the queue flushing on service shutdown.
	done := make(chan struct{})

	q := make(chan Attempt)
	r := Requeue{q: q, wg: &sync.WaitGroup{}, log: log}

	go func() {
		for {
			currentAttempt := <-q
			err := currentAttempt.Exec()
			if err != nil {
				if currentAttempt.attempts < maxAttempts {
					log.Info(ctx, L.Messagef("current / max: %d / %d, requeue action", currentAttempt.attempts, maxAttempts))
					go func(a Attempt) {
						q <- a // may block
					}(currentAttempt)
				} else {
					log.Warn(ctx, L.Messagef("requeue: too many failed attempts for action; last error: %s", err.Error()))
					r.wg.Done()
				}
				time.Sleep(pause)
			} else {
				r.wg.Done()
			}
		}
	}()

	go func() {
		<-ctx.Done()

		// After setting closed to true we will not accept any new actions.
		// The below lock is not really necessary. Even if the assignment is interleaved with
		// the addition of a new action, this would be reflected in the wait group.
		// But this way we make explicit that we are closing the queue.
		r.mu.Lock()
		r.closed = true
		r.mu.Unlock()

		log.Info(ctx, L.Message("executing: remaining actions"))
		r.wg.Wait() // wait for remaining actions to be executed, if any
		log.Info(ctx, L.Message("executed: remaining actions"))

		close(done) // signal that we are done executing remaining actions
	}()

	return &r, done
}

func (r *Requeue) Push(a Action) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return false
	}

	r.wg.Add(1)
	go func() {
		r.q <- Attempt{
			action:   a,
			attempts: 0,
		}
	}()
	return true
}

func (a *Attempt) Exec() error {
	a.attempts++
	return a.action()
}
