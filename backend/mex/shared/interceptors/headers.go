package interceptors

import (
	"context"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"
)

func RemoveResponseHeaders(headersToRemove []string) func(context.Context, http.ResponseWriter, proto.Message) error {
	return func(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
		for _, header := range headersToRemove {
			w.Header().Del(header)
		}

		return nil
	}
}

func RewriteStatusCode() func(context.Context, http.ResponseWriter, proto.Message) error {
	return func(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
		md, ok := runtime.ServerMetadataFromContext(ctx)
		if !ok {
			return nil
		}

		if values := md.HeaderMD.Get("Mex-Hinted-Status-Code"); len(values) > 0 {
			code, err := strconv.Atoi(values[0])
			if err != nil {
				return err
			}

			delete(md.HeaderMD, "Mex-Hinted-Status-Code")
			delete(w.Header(), "Grpc-Metadata-Mex-Hinted-Status-Code")

			w.WriteHeader(code)
		}

		return nil
	}
}

func RewriteHintedHeaders() func(context.Context, http.ResponseWriter, proto.Message) error {
	return func(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
		md, ok := runtime.ServerMetadataFromContext(ctx)
		if !ok {
			return nil
		}

		if values := md.HeaderMD.Get("Mex-Hinted-Location"); len(values) > 0 {
			location := values[0]

			delete(md.HeaderMD, "Mex-Hinted-Location")
			delete(w.Header(), "Grpc-Metadata-Mex-Hinted-Location")

			w.Header().Set("Location", location)
		}

		if values := md.HeaderMD.Get("Mex-Hinted-Cache-Control"); len(values) > 0 {
			cacheControl := values[0]

			delete(md.HeaderMD, "Mex-Hinted-Cache-Control")
			delete(w.Header(), "Grpc-Metadata-Mex-Hinted-Cache-Control")

			w.Header().Set("Cache-Control", cacheControl)
		}

		if values := md.HeaderMD.Get("Mex-Hinted-Pragma"); len(values) > 0 {
			pragma := values[0]

			delete(md.HeaderMD, "Mex-Hinted-Pragma")
			delete(w.Header(), "Grpc-Metadata-Mex-Hinted-Pragma")

			w.Header().Set("Pragma", pragma)
		}

		return nil
	}
}

func SetResponseHeader(header string, value string) func(context.Context, http.ResponseWriter, proto.Message) error {
	return func(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
		w.Header().Set(header, value)
		return nil
	}
}
