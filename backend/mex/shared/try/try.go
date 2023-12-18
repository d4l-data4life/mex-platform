package try

import (
	"context"
	"fmt"

	L "github.com/d4l-data4life/mex/mex/shared/log"
)

type Task[V any] struct {
	Log           L.Logger
	Desc          string
	Phase         string
	PauseStrategy PauseStrategy
	Func          func() (V, error)
}

func Try[V any](ctx context.Context, task Task[V]) (V, error) {
	if task.Log == nil {
		panic("task must have a logger")
	}

	task.Log.Info(ctx, L.Messagef("trying:  %s", task.Desc), L.Phase(task.Phase))

	var i uint
	for i = 0; ; i++ {
		task.Log.Info(ctx, L.Messagef("trying:  %s (attempt %d)", task.Desc, i+1), L.Phase(task.Phase))

		v, err := task.Func()

		if err == nil {
			task.Log.Info(ctx, L.Messagef("success: %s (attempt %d)", task.Desc, i+1), L.Phase(task.Phase))
			return v, nil
		}

		task.Log.Warn(ctx, L.Messagef("failed:  %s (attempt %d), waiting before retry (error: %s)", task.Desc, i+1, err.Error()), L.Phase(task.Phase))
		if task.PauseStrategy.Pause() != nil {
			break
		}
	}

	var zero V
	return zero, fmt.Errorf("exceeded attempts when trying: %s", task.Desc)
}
