package items

import (
	"context"
	"math"
	"net/http"

	"github.com/d4l-data4life/mex/mex/shared/constants"
	"github.com/d4l-data4life/mex/mex/shared/hints"
	"github.com/d4l-data4life/mex/mex/shared/log"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/datamodel"
	itemspb "github.com/d4l-data4life/mex/mex/services/metadata/endpoints/items/pb"
)

func (svc *Service) DeleteItem(ctx context.Context, request *itemspb.DeleteItemRequest) (*itemspb.DeleteItemResponse, error) {
	queries := datamodel.New(svc.DB)

	err := queries.DbDeleteItem(ctx, request.ItemId)
	if err != nil {
		return nil, err
	}

	hints.HintHTTPStatusCode(ctx, http.StatusNoContent)
	return &itemspb.DeleteItemResponse{}, nil
}

func (svc *Service) DeleteItems(ctx context.Context, request *itemspb.DeleteItemsRequest) (*itemspb.DeleteItemsResponse, error) {
	itemsToDeleteMap := make(map[string]bool)
	for _, id := range request.ItemIds {
		itemsToDeleteMap[id] = true
	}

	queries := datamodel.New(svc.DB)

	for _, bID := range request.BusinessIds {
		versions, err := queries.DbListItemsForBusinessId(ctx, bID)
		if err != nil {
			return nil, err
		}
		for _, vItem := range versions {
			itemsToDeleteMap[vItem.ItemID] = true
			// If cascading is on, find & delete related fragments using the corresponding relation
			if request.Cascade {
				fragmentIds, fragmentErr := queries.DbGetRelationTargetsForSourceAndType(ctx,
					datamodel.DbGetRelationTargetsForSourceAndTypeParams{
						SourceItemID: vItem.ItemID,
						Type:         constants.MergedItemToFragmentRelation,
					},
				)
				if fragmentErr != nil {
					return nil, fragmentErr
				}
				for _, fID := range fragmentIds {
					itemsToDeleteMap[fID] = true
				}
			}
		}
	}

	if len(itemsToDeleteMap) == 0 {
		svc.Log.Info(ctx, log.Messagef("nothing to delete"))
		hints.HintHTTPStatusCode(ctx, http.StatusNoContent)
		return &itemspb.DeleteItemsResponse{}, nil
	}

	itemsToDelete := make([]string, len(itemsToDeleteMap))
	i := 0
	for iID := range itemsToDeleteMap {
		itemsToDelete[i] = iID
		i++
	}
	rawAffectedRowNo, err := queries.DbDeleteItems(ctx, itemsToDelete)
	if err != nil {
		return nil, err
	}
	if int(rawAffectedRowNo) != len(itemsToDeleteMap) {
		svc.Log.Warn(ctx, log.Messagef("attempted to delete %d items but only updated %d DB rows", len(itemsToDeleteMap), rawAffectedRowNo))
	}

	/*
		To ensure JSON serialization as a number, we must cast to int32.
		Overflow is indicated by a negative value.
	*/
	var affectedRowNo int32
	if rawAffectedRowNo > math.MaxInt32 || rawAffectedRowNo < math.MinInt32 {
		affectedRowNo = -1
	} else {
		affectedRowNo = int32(rawAffectedRowNo)
	}
	// We return a 200 if we also return a body, cf. RFC 7231 sec. 4.3.5
	hints.HintHTTPStatusCode(ctx, http.StatusOK)
	return &itemspb.DeleteItemsResponse{
		DeleteItemIds: itemsToDelete,
		RowsModified:  affectedRowNo,
	}, nil
}

func (svc *Service) DeleteAllItems(ctx context.Context, request *itemspb.DeleteAllItemsRequest) (*itemspb.DeleteAllItemsResponse, error) {
	queries := datamodel.New(svc.DB)

	err := queries.DbDeleteAllItems(ctx)
	if err != nil {
		return nil, err
	}

	hints.HintHTTPStatusCode(ctx, http.StatusNoContent)
	return &itemspb.DeleteAllItemsResponse{}, nil
}
