package notify

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"

	pbNotify "github.com/d4l-data4life/mex/mex/services/metadata/endpoints/notify/pb"
)

type customMarshaler struct {
	marshaler *runtime.JSONPb
}

func NewSendNotificationMarshaler() runtime.Marshaler {
	return &customMarshaler{
		marshaler: &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitUnpopulated: true,
			},
		},
	}
}

func (m *customMarshaler) Marshal(v interface{}) ([]byte, error) {
	return m.marshaler.Marshal(v)
}

func (m *customMarshaler) NewEncoder(w io.Writer) runtime.Encoder {
	return m.marshaler.NewEncoder(w)
}

func (m *customMarshaler) ContentType(v interface{}) string {
	return m.marshaler.ContentType(v)
}

// This function seems never to be called by the gRPC gateway.
// It rather uses the Decoder way of unmarshaling.
// We panic to see whether the method is called e.g. after some gRPC library update,
// in which case we would need to amend the logic.
func (*customMarshaler) Unmarshal(data []byte, target interface{}) error {
	panic("customMarshaler.Unmarshal")
}

func (m *customMarshaler) NewDecoder(r io.Reader) runtime.Decoder {
	return &customDecoder{reader: r}
}

type customDecoder struct {
	reader io.Reader
}

// Implement runtime.Decoder interface
func (d *customDecoder) Decode(target interface{}) error {
	buf, err := io.ReadAll(d.reader)
	if err != nil {
		return err
	}

	var mappedJSON map[string]any
	err = json.Unmarshal(buf, &mappedJSON)
	if err != nil {
		return err
	}

	templateInfo, err := parseTemplateInfo(mappedJSON["templateInfo"])
	if err != nil {
		return err
	}

	bufFormData, err := json.Marshal(mappedJSON["formData"])
	if err != nil {
		return err
	}

	notifyMsg := target.(*pbNotify.SendNotificationRequest)
	notifyMsg.TemplateInfo = templateInfo
	notifyMsg.FormData = string(bufFormData)

	return nil
}

func parseTemplateInfo(v any) (*pbNotify.TemplateInfo, error) {
	if v == nil {
		return nil, fmt.Errorf("cannot parse TemplateInfo: nil")
	}

	buf, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("error marshaling TemplateInfo: %s", err.Error())
	}

	var ti pbNotify.TemplateInfo
	err = protojson.Unmarshal(buf, &ti)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling template info: %s", err.Error())
	}

	return &ti, nil
}
