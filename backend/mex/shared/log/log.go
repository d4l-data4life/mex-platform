package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/d4l-data4life/mex/mex/shared/constants"
)

type level string

const (
	LevelFatal   level = "fatal"
	LevelError   level = "error"
	LevelWarning level = "warn"
	LevelInfo    level = "info"
	LevelTrace   level = "trace"
	LevelAudit   level = "audit"

	RedactedUserID = "00000000-0000-0000-0000-000000000000"
	RedactedNote   = "<redacted>"

	LogEntryTypeLeveled = 1
	LogEntryTypeBI      = 2
	LogEntryTypeHTTP    = 3
)

type IDs struct {
	TraceID   string `json:"trace-id"`
	UserID    string `json:"user-id,omitempty"`
	JobID     string `json:"job-id,omitempty"`
	TenantID  string `json:"tenant-id"`
	RequestID string `json:"request-id,omitempty"`
}

type entryBase struct {
	IDs
	Timestamp      time.Time `json:"timestamp"`
	ServiceName    string    `json:"service-name"`
	ServiceVersion string    `json:"service-version"`
	Hostname       string    `json:"hostname"`
	EventType      string    `json:"event-type"`
	Color          string    `json:"color,omitempty"`
}

type levelEntry struct {
	entryBase

	LogLevel level  `json:"log-level"`
	Message  string `json:"message"`
	Error    string `json:"error,omitempty"`
	Phase    string `json:"phase,omitempty"`
}

func (e *levelEntry) Type() EntryType {
	return LogEntryTypeLeveled
}

func (e *levelEntry) JobID() string {
	return e.entryBase.IDs.JobID
}

func (e *levelEntry) Bytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(buf)
	err := enc.Encode(*e)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}

type MexLogger struct {
	hostName       string
	serviceName    string
	serviceVersion string

	redactPersonalFields bool
	redactQueryParams    []string

	emitters []Emitter
}

func New(serviceName string, version string, initialEmitter Emitter) (*MexLogger, error) {
	l := MexLogger{
		hostName:       os.Getenv("HOSTNAME"),
		serviceName:    serviceName,
		serviceVersion: version,

		redactPersonalFields: true,       // redact by default
		redactQueryParams:    []string{}, // no redaction of HTTP query parameters by default

		emitters: []Emitter{initialEmitter},
	}

	return &l, nil
}

func (l *MexLogger) SetRedact(state bool) {
	l.redactPersonalFields = state
}

func (l *MexLogger) SetQueryRedactions(paramsToRedact []string) {
	l.redactQueryParams = paramsToRedact
}

func (l *MexLogger) AddEmitter(emitter Emitter) {
	l.emitters = append(l.emitters, emitter)
}

func (l *MexLogger) Emit(e Entry) {
	for _, emitter := range l.emitters {
		emitter.Emit(e)
	}
}

func (l *MexLogger) baseOpEntry() *levelEntry {
	return &levelEntry{
		entryBase: entryBase{
			Timestamp:      time.Now(),
			Hostname:       l.hostName,
			ServiceName:    l.serviceName,
			ServiceVersion: l.serviceVersion,
			EventType:      "generic",
		},
		LogLevel: LevelInfo, // default level
	}
}

func (l *MexLogger) Trace(ctx context.Context, opts ...Opt) {
	if constants.GetContextValueDefault(ctx, constants.ContextKeyTraceThis, "false") == "false" {
		return
	}

	entry := l.baseOpEntry()
	entry.LogLevel = LevelTrace

	for _, opt := range opts {
		opt(entry)
	}

	if ctx != nil {
		entry.IDs = Locals(ctx)
	}

	if l.redactPersonalFields {
		entry.IDs.UserID = RedactedUserID
	}

	l.Emit(entry)
}

func (l *MexLogger) Info(ctx context.Context, opts ...Opt) {
	entry := l.baseOpEntry()
	entry.LogLevel = LevelInfo

	for _, opt := range opts {
		opt(entry)
	}

	if ctx != nil {
		entry.IDs = Locals(ctx)
	}

	if l.redactPersonalFields {
		entry.IDs.UserID = RedactedUserID
	}

	l.Emit(entry)
}

func (l *MexLogger) Warn(ctx context.Context, opts ...Opt) {
	entry := l.baseOpEntry()
	entry.LogLevel = LevelWarning

	for _, opt := range opts {
		opt(entry)
	}

	if ctx != nil {
		entry.IDs = Locals(ctx)
	}

	if l.redactPersonalFields {
		entry.IDs.UserID = RedactedUserID
	}

	l.Emit(entry)
}

func (l *MexLogger) Error(ctx context.Context, opts ...Opt) {
	entry := l.baseOpEntry()
	entry.LogLevel = LevelError

	for _, opt := range opts {
		opt(entry)
	}

	if ctx != nil {
		entry.IDs = Locals(ctx)
	}

	if l.redactPersonalFields {
		entry.IDs.UserID = RedactedUserID
	}

	l.Emit(entry)
}

func Messagef(message string, a ...any) Opt {
	return func(e Entry) {
		if t, ok := e.(*levelEntry); ok {
			t.Message = fmt.Sprintf(message, a...)
		}
	}
}

func Message(message string) Opt {
	return func(e Entry) {
		if t, ok := e.(*levelEntry); ok {
			t.Message = message
		}
	}
}

func Phase(phase string) Opt {
	return func(e Entry) {
		if t, ok := e.(*levelEntry); ok {
			t.Phase = phase
		}
	}
}

func Reason(reason interface{}) Opt {
	return func(e Entry) {
		if t, ok := e.(*levelEntry); ok {
			t.Error = emptyIfNil(reason)
		}
	}
}

func EventType(eventType string) Opt {
	return func(e Entry) {
		if t, ok := e.(*levelEntry); ok {
			t.EventType = eventType
		}
	}
}

func Color(color string) Opt {
	return func(e Entry) {
		if t, ok := e.(*levelEntry); ok {
			t.Color = color
		}
	}
}

func Locals(ctx context.Context) IDs {
	return IDs{
		RequestID: constants.GetContextValueDefault(ctx, constants.ContextKeyRequestID, ""),
		JobID:     constants.GetContextValueDefault(ctx, constants.ContextKeyJobID, ""),
		TraceID:   constants.GetContextValueDefault(ctx, constants.ContextKeyTraceID, ""),
		UserID:    constants.GetContextValueDefault(ctx, constants.ContextKeyUserID, ""),
		TenantID:  constants.GetContextValueDefault(ctx, constants.ContextKeyTenantID, ""),
	}
}

var (
	PhaseStartup  = Phase("startup")
	PhaseShutdown = Phase("shutdown")
)

func emptyIfNil(obj any) string {
	if obj == nil {
		return ""
	}

	switch v := obj.(type) {
	case error:
		return v.Error()
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", obj)
	}
}
