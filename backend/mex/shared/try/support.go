package try

import (
	"context"

	L "github.com/d4l-data4life/mex/mex/shared/log"
)

type Supporter[T any] interface {
	Init(ctx context.Context, log L.Logger) (T, error)
}
