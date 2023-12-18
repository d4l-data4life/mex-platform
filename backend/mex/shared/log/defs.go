package log

import (
	"context"
	"net/http"
	"time"
)

type EntryType int

type Entry interface {
	Type() EntryType
	JobID() string
	Bytes() []byte
}

type Opt = func(Entry)

type Emitter interface {
	Emit(Entry)
}

type LevelLogger interface {
	Error(context.Context, ...Opt)
	Warn(context.Context, ...Opt)
	Info(context.Context, ...Opt)
	Trace(context.Context, ...Opt)
}

type BILogger interface {
	BIEvent(context.Context, ...Opt)
}

type HTTPLogger interface {
	LogHTTPRequest(*http.Request)
	LogHTTPResponse(r *http.Request, status int, bodySize int, startTime time.Time, logRequestURL bool)
}

type Logger interface {
	LevelLogger
	BILogger
	HTTPLogger
}
