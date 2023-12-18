package errhandler

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/d4l-data4life/mex/mex/shared/auth"
	E "github.com/d4l-data4life/mex/mex/shared/errstat"
	L "github.com/d4l-data4life/mex/mex/shared/log"
)

type errorFormatter func(ctx context.Context, err error) (int, *status.Status)

func grpcStatusErrorFormatter(ctx context.Context, err error) (int, *status.Status) {
	if err == nil {
		return 0, nil
	}

	st := status.Convert(err)
	if st == nil {
		// indicate unhandled
		return 0, nil
	}

	traceID := auth.GetTraceID(ctx)
	if traceID != "" {
		st, _ = st.WithDetails(&E.ErrorDetailTraceId{TraceId: traceID})
	}

	return runtime.HTTPStatusFromCode(st.Code()), st
}

func httpStatusErrorFormatter(ctx context.Context, err error) (int, *status.Status) {
	if err == nil {
		return 0, nil
	}

	if e, ok := err.(*runtime.HTTPStatusError); ok {
		st := status.New(codes.Internal, e.Error())
		st, _ = st.WithDetails(&E.ErrorDetailTraceId{TraceId: auth.GetTraceID(ctx)})
		return e.HTTPStatus, st
	}

	// indicate unhandled
	return 0, nil
}

type HTTPStatusCounter func(status int)

func CustomErrorHandler(log L.Logger, statusCounter HTTPStatusCounter) runtime.ErrorHandlerFunc {
	formatters := []errorFormatter{grpcStatusErrorFormatter, httpStatusErrorFormatter}

	if log == nil {
		panic("log is nil")
	}

	if statusCounter == nil {
		panic("HTTP status counter function is nil")
	}

	return func(ctx context.Context, mux *runtime.ServeMux, _ runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
		const fallbackNilErr = `{"message": "error handler is called but no error is given"}`
		const fallbackMarshalErr = `{"message": "failed to encode error message"}`

		w.Header().Set("content-type", "application/json")
		if err == nil {
			w.WriteHeader(http.StatusInternalServerError)
			statusCounter(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fallbackNilErr))
			return
		}

		// Iterate through all formatters and return the results of the first one
		// that returns a non-zero result.
		for _, formatter := range formatters {
			statusCode, msg := formatter(ctx, err)
			if statusCode != 0 {
				buf, err := protojson.Marshal(msg.Proto())
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					statusCounter(http.StatusInternalServerError)
					_, _ = io.WriteString(w, fallbackMarshalErr)
					return
				}

				// We redact pretty late here, but it is easier to implement than descending into a
				// Protobuf structure (that is, `msg`` above) and change respective values.
				buf, err = E.Redact(ctx, buf)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					statusCounter(http.StatusInternalServerError)
					_, _ = io.WriteString(w, fallbackMarshalErr)
					return
				}

				w.WriteHeader(statusCode)
				statusCounter(statusCode)
				_, _ = w.Write(buf)
				return
			}
		}

		w.WriteHeader(http.StatusInternalServerError)
		statusCounter(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fallbackMarshalErr))
	}
}
