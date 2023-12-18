package items

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/d4l-data4life/mex/mex/shared/auth"
	"github.com/d4l-data4life/mex/mex/shared/constants"
	"github.com/d4l-data4life/mex/mex/shared/db"
	"github.com/d4l-data4life/mex/mex/shared/items"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/uuid"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/datamodel"
	itemspb "github.com/d4l-data4life/mex/mex/services/metadata/endpoints/items/pb"
)

func (svc *Service) CreateRelation(ctx context.Context, request *itemspb.CreateRelationRequest) (*itemspb.CreateRelationResponse, error) {
	if validationErr := checkRelationType(request.Type); validationErr != nil {
		return nil, validationErr
	}

	return svc.DoRelationCreate(ctx, request)
}

func (svc *Service) DoRelationCreate(ctx context.Context, request *itemspb.CreateRelationRequest) (*itemspb.CreateRelationResponse, error) {
	user, err := auth.GetMexUser(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := db.AcquireTx(ctx, svc.DB)
	if err != nil {
		return nil, status.Error(codes.Internal, "transaction creation failed")
	}

	var infoItemID *string // may remain nil
	if len(request.Values) > 0 {
		svc.Log.Info(ctx, L.Messagef("Storing information item associated with new relation"))
		itemResponse, err := svc.CreateItem(context.WithValue(ctx, db.ContextKeyTx, tx), &itemspb.CreateItemRequest{
			Item: &items.Item{
				EntityType: "relation:data", // TODO: externalize me
				Values:     request.Values,
			},
		})
		if err != nil {
			svc.Log.Error(ctx, L.Message("Item creation failed - initiating rollback!"))
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				return nil, fmt.Errorf("item creation failed (%s) and rollback failed (%s)", err.Error(), rollbackErr.Error())
			}
			return nil, err
		}
		infoItemID = &itemResponse.ItemId
	}

	queries := datamodel.New(tx)
	now := time.Now()

	relationID := uuid.MustNewV4()
	svc.Log.Info(ctx, L.Messagef("Creating relation of type %s between source item with ID %s and target item with"+
		" ID %s", request.Type, request.SourceItemId, request.TargetItemId))
	_, err = queries.DbCreateRelation(ctx, datamodel.DbCreateRelationParams{
		CreatedAt:    pgtype.Timestamptz{Time: now, Valid: true},
		ID:           relationID,
		Owner:        user.UserId,
		Type:         request.Type,
		SourceItemID: request.SourceItemId,
		TargetItemID: request.TargetItemId,
		InfoItemID:   db.TextFromStringPtr(infoItemID),
	})
	if err != nil {
		svc.Log.Error(ctx, L.Message("Link creation failed - initiating rollback!"))
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			return nil, fmt.Errorf("relation creation failed (%s) and rollback failed (%s)", err.Error(), rollbackErr.Error())
		}
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &itemspb.CreateRelationResponse{
		RelationId: relationID,
		InfoItemId: infoItemID,
	}, nil
}

//nolint:lll
func (svc Service) CreateRelationsFromBusinessIds(ctx context.Context, request *itemspb.CreateRelationsFromBusinessIdsRequest) (*itemspb.CreateRelationsFromBusinessIdsResponse, error) {
	if validationErr := checkRelationType(request.RelationType); validationErr != nil {
		return nil, validationErr
	}

	user, err := auth.GetMexUser(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := db.AcquireTx(ctx, svc.DB)
	if err != nil {
		return nil, err
	}

	queries := datamodel.New(tx)

	svc.Log.Info(ctx, L.Messagef("Creating relations of type %s based on the link field '%s' in the item with ID %s",
		request.RelationType, request.SourceItemFieldName, request.SourceItemId))

	insertedRows, err := queries.DbCreateRelationsFromBusinessIDs(ctx, datamodel.DbCreateRelationsFromBusinessIDsParams{
		Owner:     user.UserId,
		Type:      request.RelationType,
		ItemID:    request.SourceItemId,
		FieldName: request.SourceItemFieldName,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			_ = tx.Commit(ctx)
			return &itemspb.CreateRelationsFromBusinessIdsResponse{
				Inserted: 0,
			}, nil
		}
		svc.Log.Error(ctx, L.Message("Link creation failed - initiating rollback!"))
		_ = tx.Rollback(ctx)
		return nil, err
	}

	_ = tx.Commit(ctx)

	return &itemspb.CreateRelationsFromBusinessIdsResponse{
		Inserted: int32(insertedRows),
	}, nil
}

//nolint:lll
func (svc Service) CreateRelationsFromOriginalItems(ctx context.Context, request *itemspb.CreateRelationsFromOriginalItemsRequest) (*itemspb.CreateRelationsFromOriginalItemsResponse, error) {
	if validationErr := checkRelationType(request.RelationType); validationErr != nil {
		return nil, validationErr
	}

	user, err := auth.GetMexUser(ctx)
	if err != nil {
		return nil, err
	}

	queries := datamodel.New(svc.DB)

	linkFieldDefs, err := svc.FieldRepo.GetFieldDefsByKind(ctx, "link")
	if err != nil {
		return nil, err
	}

	linkFieldNames := make([]string, len(linkFieldDefs))
	for i, fd := range linkFieldDefs {
		linkFieldNames[i] = fd.Name()
	}

	provenanceItems, err := queries.DbItemsByLinks(ctx, linkFieldNames, request.BusinessId)
	if err != nil {
		if err == sql.ErrNoRows {
			return &itemspb.CreateRelationsFromOriginalItemsResponse{
				Inserted: 0,
			}, nil
		}
		return nil, err
	}

	svc.Log.Info(ctx, L.Messagef("%#v", provenanceItems))

	now := time.Now()

	for _, item := range provenanceItems {
		relationID := uuid.MustNewV4()
		svc.Log.Info(ctx, L.Messagef("Creating relations of type %s between the source item with ID %s and the"+
			" target item with ID %s", request.RelationType, request.SourceItemId, item.SourceItemID))
		_, err = queries.DbCreateRelation(ctx, datamodel.DbCreateRelationParams{
			CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
			ID:        relationID,
			Owner:     user.UserId,

			Type:         request.RelationType,
			SourceItemID: request.SourceItemId,
			TargetItemID: item.SourceItemID,
		})
		if err != nil {
			svc.Log.Error(ctx, L.Message("Link creation failed"))
			return nil, err
		}
	}

	return &itemspb.CreateRelationsFromOriginalItemsResponse{
		Inserted: int32(len(provenanceItems)),
	}, nil
}

func checkRelationType(relType string) error {
	if relType == constants.MergedItemToFragmentRelation {
		return fmt.Errorf("relations of type '%s' are reserved for internal use and cannot be manually created", constants.MergedItemToFragmentRelation)
	}
	return nil
}
