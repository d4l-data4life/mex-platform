package errstat

import (
	"context"
	"encoding/json"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/d4l-data4life/mex/mex/shared/constants"
)

func Redact(ctx context.Context, buf []byte) ([]byte, error) {
	// If the request is to be traced, do not redact anything.
	if constants.GetContextValueDefault(ctx, constants.ContextKeyTraceThis, "false") == "true" {
		return buf, nil
	}

	// Else: redact the dev message details recursively.
	var v map[string]any
	err := json.Unmarshal(buf, &v)
	if err != nil {
		return nil, err
	}

	var a anypb.Any
	_ = anypb.MarshalFrom(&a, &ErrorDetailDevMessage{}, proto.MarshalOptions{})
	redactHelper(v, a.TypeUrl)

	newBuf, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return newBuf, nil
}

func redactHelper(v any, typeURL string) {
	if vAsMap, ok := v.(map[string]any); ok {
		if vAsMap["@type"] == typeURL {
			for k := range vAsMap {
				if k != "@type" {
					vAsMap[k] = "<redacted>"
				}
			}
		}

		for _, v := range vAsMap {
			redactHelper(v, typeURL)
		}
	} else if vAsSlice, ok := v.([]any); ok {
		for _, v := range vAsSlice {
			redactHelper(v, typeURL)
		}
	}
}
