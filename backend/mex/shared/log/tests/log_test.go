package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/d4l-data4life/mex/mex/shared/constants"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/log/emit"
)

type test struct {
	name  string
	call  func(l *L.MexLogger)
	check func(t *testing.T, logs []map[string]interface{})
}

var tests = []test{
	{
		name: "case 1",
		call: func(l *L.MexLogger) {
			ctx := context.WithValue(context.WithValue(
				context.Background(),
				constants.ContextKeyTraceID, "trace-0001"),
				constants.ContextKeyUserID, "user-0001")

			l.Info(ctx, L.Message("hello & goodbye"))
			l.Warn(ctx, L.Message("<world>"))
		},
		check: func(t *testing.T, logs []map[string]interface{}) {
			assert.Equal(t, 2, len(logs))

			assert.Equal(t, map[string]interface{}{
				"timestamp":       logs[0]["timestamp"], // ignore
				"event-type":      "generic",
				"message":         "hello & goodbye",
				"service-name":    "test-service",
				"service-version": "v1.0.0",
				"log-level":       "info",
				"tenant-id":       "",
				"trace-id":        "trace-0001",
				"hostname":        "",
				"user-id":         "00000000-0000-0000-0000-000000000000",
			}, logs[0])

			assert.Equal(t, map[string]interface{}{
				"timestamp":       logs[1]["timestamp"], // ignore
				"event-type":      "generic",
				"message":         "<world>",
				"service-name":    "test-service",
				"service-version": "v1.0.0",
				"log-level":       "warn",
				"tenant-id":       "",
				"trace-id":        "trace-0001",
				"hostname":        "",
				"user-id":         "00000000-0000-0000-0000-000000000000",
			}, logs[1])
		},
	},
	{
		name: "BI event",
		call: func(l *L.MexLogger) {
			ctx := context.WithValue(context.WithValue(
				context.Background(),
				constants.ContextKeyTraceID, "trace-0001"),
				constants.ContextKeyUserID, "user-0001")

			l.BIEvent(ctx, L.BIActivity("search"))
		},
		check: func(t *testing.T, logs []map[string]interface{}) {
			assert.Equal(t, 1, len(logs))

			assert.Equal(t, map[string]interface{}{
				"timestamp": logs[0]["timestamp"], // ignore
				"event-id":  logs[0]["event-id"],  // ignore

				"event-type":      "bi-event",
				"activity-type":   "search",
				"service-name":    "test-service",
				"service-version": "v1.0.0",
				"tenant-id":       "",
				"trace-id":        "trace-0001",
				"hostname":        "",
				"user-id":         "00000000-0000-0000-0000-000000000000",
			}, logs[0])
		},
	},
	{
		name: "HTTP request",
		call: func(l *L.MexLogger) {
			ctx := context.WithValue(context.WithValue(
				context.Background(),
				constants.ContextKeyTraceID, "trace-0001"),
				constants.ContextKeyUserID, "user-0001")

			r := httptest.NewRequest("POST", "/foo/bar?test=foo&state=1234", strings.NewReader("BODY BODY BODY"))
			r.Header.Set("Content-Type", "pdf")
			l.LogHTTPRequest(r.WithContext(ctx))
		},
		check: func(t *testing.T, logs []map[string]interface{}) {
			assert.Equal(t, 1, len(logs))

			assert.Equal(t, map[string]interface{}{
				"timestamp": logs[0]["timestamp"], // ignore

				"event-type":      "http-in-request",
				"service-name":    "test-service",
				"service-version": "v1.0.0",

				"real-ip":          "<redacted>",
				"req-ip":           "<redacted>",
				"req-method":       "POST",
				"content-type":     "pdf",
				"content-encoding": "",
				"req-url":          "/foo/bar?state=<redacted>&test=foo",
				"payload-length":   float64(-1),
				"req-body":         "",
				"req-form":         "",

				"hostname":  "",
				"tenant-id": "",

				"trace-id": "trace-0001",
				"user-id":  "00000000-0000-0000-0000-000000000000",
			}, logs[0])
		},
	},
	{
		name: "HTTP response, w/o request URL",
		call: func(l *L.MexLogger) {
			ctx := context.WithValue(context.WithValue(context.WithValue(
				context.Background(),
				constants.ContextKeyTraceID, "trace-0001"),
				constants.ContextKeyUserID, "user-0001"),
				constants.ContextKeyTenantID, "tenant-0001")

			r := httptest.NewRequest("POST", "/foo/bar?state=1234&test=foo", strings.NewReader("BODY BODY BODY"))
			r.Header.Set("Content-Type", "pdf")
			l.LogHTTPResponse(r.WithContext(ctx), 404, 12345, time.Now(), false)
		},
		check: func(t *testing.T, logs []map[string]interface{}) {
			assert.Equal(t, 1, len(logs))

			assert.Equal(t, map[string]interface{}{
				"timestamp": logs[0]["timestamp"], // ignore

				"event-type":      "http-in-response",
				"service-name":    "test-service",
				"service-version": "v1.0.0",

				"req-method":       "POST",
				"content-type":     "",
				"content-encoding": "",
				"req-url":          "",
				"payload-length":   float64(12345),

				"hostname":  "",
				"tenant-id": "tenant-0001",

				"trace-id":           "trace-0001",
				"user-id":            "00000000-0000-0000-0000-000000000000",
				"response-body":      "",
				"response-code":      float64(404),
				"roundtrip-duration": float64(0),
			}, logs[0])
		},
	},
	{
		name: "HTTP response, w/ request URL",
		call: func(l *L.MexLogger) {
			ctx := context.WithValue(context.WithValue(context.WithValue(
				context.Background(),
				constants.ContextKeyTraceID, "trace-0001"),
				constants.ContextKeyUserID, "user-0001"),
				constants.ContextKeyTenantID, "tenant-0001")

			r := httptest.NewRequest("POST", "/foo/bar?state=1234&test=foo", strings.NewReader("BODY BODY BODY"))
			r.Header.Set("Content-Type", "pdf")
			l.LogHTTPResponse(r.WithContext(ctx), 404, 12345, time.Now(), true)
		},
		check: func(t *testing.T, logs []map[string]interface{}) {
			assert.Equal(t, 1, len(logs))

			assert.Equal(t, map[string]interface{}{
				"timestamp": logs[0]["timestamp"], // ignore

				"event-type":      "http-in-response",
				"service-name":    "test-service",
				"service-version": "v1.0.0",

				"req-method":       "POST",
				"content-type":     "",
				"content-encoding": "",
				"req-url":          "/foo/bar?state=<redacted>&test=foo",
				"payload-length":   float64(12345),

				"hostname":  "",
				"tenant-id": "tenant-0001",

				"trace-id":           "trace-0001",
				"user-id":            "00000000-0000-0000-0000-000000000000",
				"response-body":      "",
				"response-code":      float64(404),
				"roundtrip-duration": float64(0),
			}, logs[0])
		},
	},
}

func Test_Redaction(t *testing.T) {

	for _, testCase := range tests {
		var buf bytes.Buffer
		l, err := L.New("test-service", "v1.0.0", &emit.WriterEmitter{Writer: &buf})
		if err != nil {
			t.Fail()
		}
		l.SetQueryRedactions([]string{"state"})

		testCase.call(l)

		var logs []map[string]interface{}
		lines := strings.Split(buf.String(), "\n")
		for i := 0; i < len(lines)-1; i++ {
			line := lines[i]
			var item map[string]interface{}
			err := json.Unmarshal([]byte(line), &item)
			if err != nil {
				t.Error(err)
			}
			logs = append(logs, item)
		}

		testCase.check(t, logs)
	}
}
