package log

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/d4l-data4life/mex/mex/shared/uuid"
)

type biEventEntry struct {
	entryBase

	EventID      string `json:"event-id"`
	ActivityType string `json:"activity-type"`
	Data         any    `json:"data,omitempty"`
}

func (e *biEventEntry) Type() EntryType {
	return LogEntryTypeBI
}

func (e *biEventEntry) JobID() string {
	return e.entryBase.IDs.JobID
}

func (e *biEventEntry) Bytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	err := enc.Encode(*e)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}

// type OpBIEventField = func(entry *opBIEventEntry)

func (l *MexLogger) baseBIEvent() *biEventEntry {
	return &biEventEntry{
		entryBase: entryBase{
			Timestamp:      time.Now(),
			Hostname:       l.hostName,
			ServiceName:    l.serviceName,
			ServiceVersion: l.serviceVersion,
			EventType:      "bi-event",
		},
	}
}

func (l *MexLogger) BIEvent(ctx context.Context, opts ...Opt) {
	entry := l.baseBIEvent()

	for _, opt := range opts {
		opt(entry)
	}

	entry.IDs = Locals(ctx)
	entry.EventID = uuid.MustNewV4()

	if l.redactPersonalFields {
		entry.IDs.UserID = RedactedUserID
	}

	l.Emit(entry)
}

func BIActivity(activityType string) Opt {
	return func(entry Entry) {
		if t, ok := entry.(*biEventEntry); ok {
			t.ActivityType = activityType
		}
	}
}

func BIData(data any) Opt {
	return func(entry Entry) {
		if t, ok := entry.(*biEventEntry); ok {
			t.Data = data
		}
	}
}
