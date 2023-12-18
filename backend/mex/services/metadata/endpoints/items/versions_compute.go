package items

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/datamodel"
	itemspb "github.com/d4l-data4life/mex/mex/services/metadata/endpoints/items/pb"
)

func (svc *Service) ComputeVersions(ctx context.Context, request *itemspb.ComputeVersionsRequest) (*itemspb.ComputeVersionsResponse, error) {
	queries := datamodel.New(svc.DB)

	versions, err := queries.DbComputeVersions(ctx, request.ItemId)
	if err != nil {
		return nil, err
	}

	if len(versions) == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("cannot find versions of item '%s' (it does either not exist or does not have a business ID", request.ItemId))
	}

	response := itemspb.ComputeVersionsResponse{
		Versions: []*itemspb.ComputeVersionsResponse_Version{},
	}

	for _, version := range versions {
		response.Versions = append(response.Versions, &itemspb.ComputeVersionsResponse_Version{
			ItemId:      version.ItemID,
			VersionDesc: fmt.Sprintf("v%d", version.Version),
			CreatedAt:   timestamppb.New(version.CreatedAt.Time),
		})
	}

	return &response, nil
}

func (svc *Service) ComputeVersionsByBusinessID(ctx context.Context, request *itemspb.ComputeVersionsByBusinessIdRequest) (*itemspb.ComputeVersionsByBusinessIdResponse, error) {
	queries := datamodel.New(svc.DB)

	versions, err := queries.DbListItemsForBusinessId(ctx, request.BusinessId)
	if err != nil {
		return nil, err
	}

	response := itemspb.ComputeVersionsByBusinessIdResponse{
		Versions: []*itemspb.ComputeVersionsByBusinessIdResponse_Version{},
	}

	for i, version := range versions {
		response.Versions = append(response.Versions, &itemspb.ComputeVersionsByBusinessIdResponse_Version{
			ItemId:      version.ItemID,
			VersionDesc: fmt.Sprintf("v%d", i+1),
			CreatedAt:   timestamppb.New(version.CreatedAt.Time),
		})
	}

	return &response, nil
}

func (svc *Service) ListAllVersions(ctx context.Context, request *itemspb.ListAllVersionsRequest) (*itemspb.ListAllVersionsResponse, error) {
	queries := datamodel.New(svc.DB)

	items, err := queries.DbListItemsWithBusinessId(ctx)
	if err != nil {
		return nil, err
	}

	response := itemspb.ListAllVersionsResponse{
		Versions: []*itemspb.ListAllVersionsResponse_Versions{},
	}

	prevBusinessID := ""
	currentVersionIndex := -1
	for _, item := range items {
		if prevBusinessID != item.BusinessID {
			response.Versions = append(response.Versions, &itemspb.ListAllVersionsResponse_Versions{
				BusinessId: item.BusinessID,
				ItemIds:    []string{item.ItemID},
			})
			prevBusinessID = item.BusinessID
			currentVersionIndex = len(response.Versions) - 1
		} else {
			response.Versions[currentVersionIndex].ItemIds = append(response.Versions[currentVersionIndex].ItemIds, item.ItemID)
		}
	}

	return &response, nil
}
