package items

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	E "github.com/d4l-data4life/mex/mex/shared/errstat"
	"github.com/d4l-data4life/mex/mex/shared/hints"
	"github.com/d4l-data4life/mex/mex/shared/items"
	L "github.com/d4l-data4life/mex/mex/shared/log"

	itemspb "github.com/d4l-data4life/mex/mex/services/metadata/endpoints/items/pb"
)

type placedItemValue struct {
	FieldName  string
	FieldValue string
	Language   string
	Place      int32
}

/*
CreateItem handles the incoming HTTP request for item creation, forwarding to the same base functionality as bulk
creation. The returned item ID is the item ID of the item created directly on the basic of the item information in
the request. The returned business ID, however, is the business ID of the "secondary" item created
due to aggregation triggered by the uploaded item - if no aggregation was done, it will be set to the empty
string.
*/
func (svc *Service) CreateItem(ctx context.Context, request *itemspb.CreateItemRequest) (*itemspb.CreateItemResponse, error) {
	if request.Item == nil {
		return nil, E.MakeGRPCStatus(codes.InvalidArgument, "no item to be created passed").Err()
	}
	createResponses, err := svc.doItemsCreate(ctx, createMultipleItemsInput{
		items:               []*items.Item{request.Item},
		precomputedHashes:   []string{request.Hash},
		duplicateAlgorithm:  svc.DuplicateDetectionAlgorithm, // Use configured value - no override possible
		preventAnnouncement: request.PreventAnnouncement,     // Do not announce (and hence autoindex) new item
	})
	if err != nil {
		svc.Log.Error(ctx, L.Messagef("item creation failed: %s", err.Error()))
		return nil, E.MakeGRPCStatus(codes.Internal, "item creation failed", E.Cause(err), request).Err()
	}

	// The case len(createResponses) = 0 will occur if the item is a duplicate of an existing item
	// --> no item created --> empty response
	finalResponse := itemspb.CreateItemResponse{}
	if len(createResponses) == 1 {
		finalResponse = itemspb.CreateItemResponse{
			ItemId:     createResponses[0].itemID,
			BusinessId: createResponses[0].aggregatedBusinessID,
		}
		svc.Log.Info(ctx, L.Messagef("completed creation of new uploaded item - assigned item ID: %s", finalResponse.ItemId))
		if finalResponse.BusinessId != "" {
			svc.Log.Info(ctx, L.Messagef("new item also lead to the creation of a merged item (business ID %s)", finalResponse.BusinessId))
		}
	} else if len(createResponses) > 1 {
		return nil, status.Error(codes.Internal, fmt.Sprintf("expected at most 1 (primary) item to have been created, but got %d", len(createResponses)))
	}

	hints.HintHTTPStatusCode(ctx, http.StatusCreated)
	return &finalResponse, nil
}
