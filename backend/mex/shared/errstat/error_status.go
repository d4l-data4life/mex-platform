package errstat

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoiface"
	"google.golang.org/protobuf/types/known/anypb"
)

func MakeGRPCStatus(code codes.Code, message string, details ...protoiface.MessageV1) *status.Status {
	st := status.New(code, message)

	if details == nil {
		return st
	}

	st, err := st.WithDetails(details...)
	if err != nil {
		panic("cannot add status detail")
	}

	return st
}

func MakeMexStatus(mexErrorCode MexErrorCode, message string, details ...protoiface.MessageV1) *status.Status {
	details = append(details, &ErrorDetailCode{Code: codeToStatus[mexErrorCode].MexErrorString})

	return MakeGRPCStatus(
		codeToStatus[mexErrorCode].GrpcCode,
		message,
		details...,
	)
}

func Cause(cause error) protoiface.MessageV1 {
	var anyVal anypb.Any

	if st, ok := status.FromError(cause); ok {
		err := anypb.MarshalFrom(&anyVal, st.Proto(), proto.MarshalOptions{})
		if err != nil {
			panic("cannot marshal status")
		}
	} else {
		err := anypb.MarshalFrom(&anyVal, &ErrorDetailReason{Reason: cause.Error()}, proto.MarshalOptions{})
		if err != nil {
			panic("cannot marshal error string")
		}
	}

	return &ErrorDetailCause{
		Cause: &anyVal,
	}
}

func DevMessage(message string) protoiface.MessageV1 {
	return &ErrorDetailDevMessage{DevMessage: message}
}

func DevMessagef(message string, a ...any) protoiface.MessageV1 {
	return &ErrorDetailDevMessage{DevMessage: fmt.Sprintf(message, a...)}
}

func CodeFrom(err error) codes.Code {
	if st, ok := status.FromError(err); ok {
		return st.Code()
	}
	return codes.Internal
}
