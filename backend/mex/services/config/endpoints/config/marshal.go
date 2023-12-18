package config

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"

	pbConfig "github.com/d4l-data4life/mex/mex/services/config/endpoints/config/pb"
)

type fileContentMarshaler struct {
	*runtime.HTTPBodyMarshaler
}

// This marshaler acts like the gRPC ServeMux's default '*' marshaler, except
// for responses of type GetFileResponse.
// In that case we return the raw content as the response body with the corresponding media/mime type.

func NewFileContentMarshaler(strictJSONParsing bool) runtime.Marshaler {
	discardUnknown := !strictJSONParsing
	return &fileContentMarshaler{
		HTTPBodyMarshaler: &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: discardUnknown,
				},
			},
		},
	}
}

func (m *fileContentMarshaler) Marshal(v interface{}) ([]byte, error) {
	if p, ok := v.(*pbConfig.GetFileResponse); ok {
		return p.Content, nil
	}
	return m.HTTPBodyMarshaler.Marshal(v)
}

func (m *fileContentMarshaler) ContentType(v interface{}) string {
	if p, ok := v.(*pbConfig.GetFileResponse); ok {
		return p.MimeType
	}
	return m.HTTPBodyMarshaler.ContentType(v)
}
