package async

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPromise(t *testing.T) {
	p := New(func(resolve Resolver[string], reject Rejecter) {
		time.Sleep(time.Second)

		resolve("foo")
	})

	s, err := p.Await()
	require.Nil(t, err)
	require.Equal(t, "foo", s)

	s, err = p.Await()
	require.Nil(t, err)
	require.Equal(t, "foo", s)
}

func TestPromiseMultiple(t *testing.T) {
	p := New(func(resolve Resolver[string], reject Rejecter) {
		time.Sleep(1 * time.Millisecond)

		resolve("foo")
	})

	var wg sync.WaitGroup

	for i := 0; i < 3000; i++ {
		wg.Add(1)

		go func() {
			s, err := p.Await()
			require.Nil(t, err)
			require.Equal(t, "foo", s)
			wg.Done()
		}()
	}

	fmt.Println("waiting...")
	wg.Wait()
	fmt.Println("done")
}
