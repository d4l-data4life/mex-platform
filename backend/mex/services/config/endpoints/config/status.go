package config

import (
	"context"
	"time"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/d4l-data4life/mex/mex/shared/coll/accumap"
	"github.com/d4l-data4life/mex/mex/shared/known/statuspb"
	"github.com/d4l-data4life/mex/mex/shared/utils"

	pbConfig "github.com/d4l-data4life/mex/mex/services/config/endpoints/config/pb"
)

const defaultMaxAge = time.Minute

func (svc *Service) GetStatus(ctx context.Context, request *pbConfig.GetStatusRequest) (*pbConfig.GetStatusResponse, error) {
	am, err := svc.getAllServiceStatuses(ctx, defaultMaxAge)
	if err != nil {
		return nil, err
	}

	configHashes := map[string]struct{}{}

	ret := pbConfig.GetStatusResponse{
		Color:    statuspb.Color_GREEN,
		Statuses: []*statuspb.Status{},
	}

	for _, sts := range am {
		ret.Statuses = append(ret.Statuses, sts...)

		for _, st := range sts {
			ret.Color = minColor(ret.Color, st.Color)
			configHashes[st.ConfigHash] = struct{}{}
		}
	}

	ret.ConfigHashes = utils.KeysOfMap(configHashes)
	if len(ret.ConfigHashes) > 1 {
		ret.Color = statuspb.Color_RED
	}

	return &ret, nil
}

func (svc *Service) getAllServiceStatuses(ctx context.Context, maxAge time.Duration) (map[string][]*statuspb.Status, error) {
	result, err := svc.Redis.HGetAll(ctx, "status").Result()
	if err != nil {
		return nil, err
	}

	accu := func(bucket []*statuspb.Status, data *statuspb.Status) []*statuspb.Status {
		return append(bucket, data)
	}
	zero := func(data *statuspb.Status) []*statuspb.Status {
		return []*statuspb.Status{data}
	}
	keyer := func(data *statuspb.Status) string {
		return data.ServiceTag
	}

	am := accumap.NewAccumap(accu, zero, keyer)
	for _, v := range result {
		var msg statuspb.Status
		err := protojson.Unmarshal([]byte(v), &msg)
		if err != nil {
			return nil, err
		}

		// Only return statuses that are younger then maxAge.
		if time.Since(msg.LastReported.AsTime()) > maxAge {
			continue
		}

		am.Push(&msg)
	}

	return am.ToMap(), nil
}

func minColor(c1, c2 statuspb.Color) statuspb.Color {
	if c1 < c2 {
		return c1
	}
	return c2
}
