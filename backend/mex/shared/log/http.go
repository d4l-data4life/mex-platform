package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/d4l-data4life/mex/mex/shared/constants"
	"github.com/d4l-data4life/mex/mex/shared/utils"
)

type httpInRequestEntry struct {
	entryBase

	ReqIP  string `json:"req-ip"`
	RealIP string `json:"real-ip"`

	ReqMethod string `json:"req-method"`
	ReqURL    string `json:"req-url"`

	ContentType     string `json:"content-type"`
	ContentEncoding string `json:"content-encoding"`
	PayloadLength   int64  `json:"payload-length"`

	ReqBody string `json:"req-body"`
	ReqForm string `json:"req-form"`
}

type httpInResponseEntry struct {
	entryBase

	ReqMethod string `json:"req-method"`
	ReqURL    string `json:"req-url"`

	ResponseCode    int    `json:"response-code"`
	ContentType     string `json:"content-type"`
	ContentEncoding string `json:"content-encoding"`
	PayloadLength   int64  `json:"payload-length"`

	ResponseBody string `json:"response-body"`

	Duration int64 `json:"roundtrip-duration"`
}

func (e *httpInRequestEntry) Type() EntryType {
	return LogEntryTypeHTTP
}

func (e *httpInRequestEntry) JobID() string {
	return e.entryBase.IDs.JobID
}

func (e *httpInRequestEntry) Bytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(buf)
	err := enc.Encode(*e)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}

func (e *httpInResponseEntry) Type() EntryType {
	return LogEntryTypeHTTP
}

func (e *httpInResponseEntry) JobID() string {
	return e.entryBase.IDs.JobID
}

func (e *httpInResponseEntry) Bytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(buf)
	err := enc.Encode(*e)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}

func (l *MexLogger) LogHTTPRequest(r *http.Request) {
	entry := &httpInRequestEntry{
		entryBase: entryBase{
			Timestamp:      time.Now(),
			Hostname:       l.hostName,
			ServiceName:    l.serviceName,
			ServiceVersion: l.serviceVersion,
			EventType:      "http-in-request",

			IDs: IDs{
				TraceID:   constants.GetContextValueDefault(r.Context(), constants.ContextKeyTraceID, ""),
				JobID:     constants.GetContextValueDefault(r.Context(), constants.ContextKeyJobID, ""),
				RequestID: constants.GetContextValueDefault(r.Context(), constants.ContextKeyRequestID, ""),

				// UserID and TenantID not available here as they would come from a JWT header which is not
				// evaluated/validated at this point in time.
			},
		},

		ReqIP:  r.RemoteAddr,
		RealIP: r.Header.Get(constants.HTTPHeaderRealIP),

		ReqMethod: r.Method,
		ReqURL:    redactURI(r.RequestURI, l.redactQueryParams),

		ContentType:     r.Header.Get("Content-Type"),
		ContentEncoding: r.Header.Get("Content-Encoding"),
		PayloadLength:   parseInt(r.Header.Get("Content-Length")),
	}

	if l.redactPersonalFields {
		entry.IDs.UserID = RedactedUserID
		entry.RealIP = RedactedNote
		entry.ReqIP = RedactedNote
	}

	l.Emit(entry)
}

func (l *MexLogger) LogHTTPResponse(r *http.Request, status int, bodySize int, startTime time.Time, logRequestURL bool) {
	reqURL := ""
	if logRequestURL {
		reqURL = redactURI(r.RequestURI, l.redactQueryParams)
	}

	entry := &httpInResponseEntry{
		entryBase: entryBase{
			Timestamp:      time.Now(),
			Hostname:       l.hostName,
			ServiceName:    l.serviceName,
			ServiceVersion: l.serviceVersion,
			EventType:      "http-in-response",

			IDs: IDs{
				TraceID:   constants.GetContextValueDefault(r.Context(), constants.ContextKeyTraceID, ""),
				JobID:     constants.GetContextValueDefault(r.Context(), constants.ContextKeyJobID, ""),
				RequestID: constants.GetContextValueDefault(r.Context(), constants.ContextKeyRequestID, ""),

				UserID:   constants.GetContextValueDefault(r.Context(), constants.ContextKeyUserID, ""),
				TenantID: constants.GetContextValueDefault(r.Context(), constants.ContextKeyTenantID, ""),
			},
		},

		ReqMethod: r.Method,

		ReqURL: reqURL,

		ResponseCode:  status,
		PayloadLength: int64(bodySize),

		Duration: time.Since(startTime).Milliseconds(),
	}

	if l.redactPersonalFields {
		entry.IDs.UserID = RedactedUserID
	}

	l.Emit(entry)
}

func parseInt(s string) int64 {
	n, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}
	return int64(n)
}

func redactURI(uri string, redactQueryParams []string) string {
	u, err := url.Parse(uri)
	if err != nil {
		return "<error parsing URI>"
	}

	query := []string{}

	for k, v := range u.Query() {
		if utils.Contains(redactQueryParams, k) {
			query = append(query, fmt.Sprintf("%s=<redacted>", k))
		} else {
			query = append(query, fmt.Sprintf("%s=%s", k, strings.Join(v, ",")))
		}
	}

	if len(query) == 0 {
		return fmt.Sprintf("%s%s", u.Host, u.Path)
	}

	sort.Strings(query)
	return fmt.Sprintf("%s%s?%s", u.Host, u.Path, strings.Join(query, "&"))
}

func NewRequestResponseLoggingMiddleware(log HTTPLogger, excludePatterns []string, logRequestURL bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Bypass logs for excluded paths
			for _, pattern := range excludePatterns {
				if pattern == r.URL.Path {
					next.ServeHTTP(w, r)
					return
				}
			}

			start := time.Now()
			log.LogHTTPRequest(r)

			nw := NewPeekingResponseWriter(w)
			next.ServeHTTP(nw, r)

			log.LogHTTPResponse(r, nw.StatusCode, nw.BodyCount, start, logRequestURL)
		})
	}
}

// This struct and its below methods just relay the data to the
// embedded http.ResponseWriter. Doing so it counts the amount
// of data and remembers the status code.
type PeekingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	BodyCount  int
}

func NewPeekingResponseWriter(w http.ResponseWriter) *PeekingResponseWriter {
	return &PeekingResponseWriter{w, http.StatusOK, 0}
}

func (lrw *PeekingResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *PeekingResponseWriter) Write(data []byte) (int, error) {
	lrw.BodyCount += len(data)
	return lrw.ResponseWriter.Write(data)
}
