package async

import (
	"sync"
)

// Simple promise implementation with the following behavior:
// - At any given time the promise is in one of the following states:
//   - pending
//   - resolved (with a value of type T)
//   - rejected (with an error)
// - A non-pending (that is, rejected or resolved) promise is also said to be settled.
// - A promise can be resolved with a value exactly once. It then cannot be rejected.
// - A promise can be rejected with an error exactly once. It then cannot be resolved.
// - Awaiting a promise blocks until the promise is settled and then returns its value or error.
// - Awaiting a promise multiple times or simultaneously is fine.

// Awaiting a promise returns always both the resolve and reject value.
// However, the implementation makes sure that only one is valid.
type Promise[T any] interface {
	Await() (T, error)
}

type prom[T any] struct {
	mu sync.Mutex
	wg sync.WaitGroup

	settled bool

	v T     // promise resolve value
	e error // promise reject error
}

type Resolver[T any] func(t T)

type Rejecter func(err error)

type PromiseFunc[T any] func(resolve Resolver[T], reject Rejecter)

func New[T any](f PromiseFunc[T]) Promise[T] {
	p := prom[T]{}

	// This indicates that a piece of works needs to be done: resolve or reject
	p.wg.Add(1)

	resolve := func(t T) {
		p.mu.Lock()
		defer p.mu.Unlock()

		if p.settled {
			panic("cannot resolve promise: already settled")
		}

		p.v = t
		p.settled = true
		p.wg.Done()
	}

	reject := func(err error) {
		p.mu.Lock()
		defer p.mu.Unlock()

		if p.settled {
			panic("cannot reject promise: already settled")
		}

		p.e = err
		p.settled = true
		p.wg.Done()
	}

	go f(resolve, reject)

	return &p
}

func (p *prom[T]) Await() (T, error) {
	p.wg.Wait()
	return p.v, p.e
}
