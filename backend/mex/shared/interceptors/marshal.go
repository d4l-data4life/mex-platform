package interceptors

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

func NewMarshaler(strictJSONParsing bool) runtime.Marshaler {
	return &runtime.HTTPBodyMarshaler{
		Marshaler: &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitUnpopulated: true, // emit default zero values
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				AllowPartial:   false,
				DiscardUnknown: !strictJSONParsing,
			},
		},
	}
}
