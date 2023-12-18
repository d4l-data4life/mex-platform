package items

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/d4l-data4life/mex/mex/shared/auth"
	"github.com/d4l-data4life/mex/mex/shared/db"
	"github.com/d4l-data4life/mex/mex/shared/hints"
	L "github.com/d4l-data4life/mex/mex/shared/log"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/datamodel"
	itemspb "github.com/d4l-data4life/mex/mex/services/metadata/endpoints/items/pb"
)

const listItemsLimit = 1000

func (svc *Service) ListItems(ctx context.Context, request *itemspb.ListItemsRequest) (*itemspb.ListItemsResponse, error) {
	queries := datamodel.New(svc.DB)

	var items []datamodel.ItemsNullableBusinessID
	var err error

	if request.EntityType == "" {
		items, err = queries.DbListItems(ctx, datamodel.DbListItemsParams{
			ItemID: request.Next,
			Limit:  listItemsLimit + 1,
		})
	} else {
		items, err = queries.DbListItemsOfType(ctx, datamodel.DbListItemsOfTypeParams{
			ItemID:     request.Next,
			Limit:      listItemsLimit + 1,
			EntityName: request.EntityType,
		})
	}
	if err != nil {
		return nil, err
	}

	size := 0
	next := ""
	if len(items) == listItemsLimit+1 {
		size = listItemsLimit // one less
		next = items[listItemsLimit].ItemID
	} else {
		size = len(items)
	}

	retItems := make([]*itemspb.ListItem, size)

	for i := 0; i < size; i++ {
		item := items[i]
		retItems[i] = &itemspb.ListItem{
			ItemId:     item.ItemID,
			Owner:      item.Owner,
			EntityType: item.EntityName,
			CreatedAt:  timestamppb.New(item.CreatedAt.Time),
			BusinessId: item.BusinessID.String,
		}
	}

	return &itemspb.ListItemsResponse{
		Items: retItems,
		Next:  next,
	}, nil
}

func (svc *Service) GetItem(ctx context.Context, request *itemspb.GetItemRequest) (*itemspb.GetItemResponse, error) {
	user, err := auth.GetMexUser(ctx)
	if err != nil {
		return nil, err
	}

	queries := datamodel.New(svc.DB)

	rawItems, err := queries.DbGetItem(ctx, request.ItemId)
	if err != nil {
		svc.Log.Error(ctx, L.Messagef("failed to fetch item from DB: % s", err.Error()))
		return nil, err
	}
	var item datamodel.Item
	switch len(rawItems) {
	case 0:
		hints.HintHTTPStatusCode(ctx, http.StatusNotFound)
		return &itemspb.GetItemResponse{}, nil
	case 1:
		item = rawItems[0]
	default:
		svc.Log.Error(ctx, L.Message("multiple items for the same item ID found"))
		return nil, fmt.Errorf("multiple items for the same item ID found")
	}

	values, err := queries.DbGetItemValues(ctx, request.ItemId)
	if err != nil {
		svc.Log.Error(ctx, L.Messagef("failed to fetch item values from DB: % s", err.Error()))
		return nil, err
	}

	err = queries.DbIncrementItemCounter(ctx, datamodel.DbIncrementItemCounterParams{
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ItemID:    request.ItemId,
		UserID:    user.UserId,
		Counts:    1,
	})
	if err != nil {
		svc.Log.Warn(ctx, L.Messagef("could not update item access counter: %s", err.Error()))
	}

	retValues := make([]*itemspb.GetItemResponse_FullItemValue, len(values))

	for i, v := range values {
		retValues[i] = &itemspb.GetItemResponse_FullItemValue{
			ItemValueId: v.ID,
			FieldName:   v.FieldName,
			FieldValue:  v.FieldValue,
			Place:       v.Place,
			Language:    v.Language.String,
			Revision:    v.Revision,
		}
	}

	return &itemspb.GetItemResponse{
		ItemId:     item.ID,
		EntityType: item.EntityName,
		Owner:      item.Owner,
		CreatedAt:  timestamppb.New(item.CreatedAt.Time),
		BusinessId: db.StringOrNil(item.BusinessID),
		Values:     retValues,
	}, nil
}
