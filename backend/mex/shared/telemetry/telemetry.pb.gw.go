// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: shared/telemetry/telemetry.proto

/*
Package telemetry is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package telemetry

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Suppress "imported and not used" errors
var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray
var _ = metadata.Join

func request_Telemetry_LivenessProbe_0(ctx context.Context, marshaler runtime.Marshaler, client TelemetryClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq LivenessRequest
	var metadata runtime.ServerMetadata

	msg, err := client.LivenessProbe(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_Telemetry_LivenessProbe_0(ctx context.Context, marshaler runtime.Marshaler, server TelemetryServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq LivenessRequest
	var metadata runtime.ServerMetadata

	msg, err := server.LivenessProbe(ctx, &protoReq)
	return msg, metadata, err

}

func request_Telemetry_ReadinessProbe_0(ctx context.Context, marshaler runtime.Marshaler, client TelemetryClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ReadinessRequest
	var metadata runtime.ServerMetadata

	msg, err := client.ReadinessProbe(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_Telemetry_ReadinessProbe_0(ctx context.Context, marshaler runtime.Marshaler, server TelemetryServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ReadinessRequest
	var metadata runtime.ServerMetadata

	msg, err := server.ReadinessProbe(ctx, &protoReq)
	return msg, metadata, err

}

// RegisterTelemetryHandlerServer registers the http handlers for service Telemetry to "mux".
// UnaryRPC     :call TelemetryServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features to stop working. Consider using RegisterTelemetryHandlerFromEndpoint instead.
func RegisterTelemetryHandlerServer(ctx context.Context, mux *runtime.ServeMux, server TelemetryServer) error {

	mux.Handle("GET", pattern_Telemetry_LivenessProbe_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/d4l.mex.telemetry.Telemetry/LivenessProbe", runtime.WithHTTPPathPattern("/probes/liveness"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_Telemetry_LivenessProbe_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Telemetry_LivenessProbe_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Telemetry_ReadinessProbe_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/d4l.mex.telemetry.Telemetry/ReadinessProbe", runtime.WithHTTPPathPattern("/probes/readiness"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_Telemetry_ReadinessProbe_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Telemetry_ReadinessProbe_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

// RegisterTelemetryHandlerFromEndpoint is same as RegisterTelemetryHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterTelemetryHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.DialContext(ctx, endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterTelemetryHandler(ctx, mux, conn)
}

// RegisterTelemetryHandler registers the http handlers for service Telemetry to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterTelemetryHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterTelemetryHandlerClient(ctx, mux, NewTelemetryClient(conn))
}

// RegisterTelemetryHandlerClient registers the http handlers for service Telemetry
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "TelemetryClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "TelemetryClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "TelemetryClient" to call the correct interceptors.
func RegisterTelemetryHandlerClient(ctx context.Context, mux *runtime.ServeMux, client TelemetryClient) error {

	mux.Handle("GET", pattern_Telemetry_LivenessProbe_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/d4l.mex.telemetry.Telemetry/LivenessProbe", runtime.WithHTTPPathPattern("/probes/liveness"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_Telemetry_LivenessProbe_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Telemetry_LivenessProbe_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Telemetry_ReadinessProbe_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/d4l.mex.telemetry.Telemetry/ReadinessProbe", runtime.WithHTTPPathPattern("/probes/readiness"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_Telemetry_ReadinessProbe_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Telemetry_ReadinessProbe_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_Telemetry_LivenessProbe_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"probes", "liveness"}, ""))

	pattern_Telemetry_ReadinessProbe_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"probes", "readiness"}, ""))
)

var (
	forward_Telemetry_LivenessProbe_0 = runtime.ForwardResponseMessage

	forward_Telemetry_ReadinessProbe_0 = runtime.ForwardResponseMessage
)
