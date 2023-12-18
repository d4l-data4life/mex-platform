package items

import (
	"context"

	"github.com/d4l-data4life/mex/mex/shared/db"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/datamodel"
	itemspb "github.com/d4l-data4life/mex/mex/services/metadata/endpoints/items/pb"
)

func (svc *Service) ListRelations(ctx context.Context, request *itemspb.ListRelationsRequest) (*itemspb.ListRelationsResponse, error) {
	queries := datamodel.New(svc.DB)

	relations, err := queries.DbListRelations(ctx)
	if err != nil {
		return nil, err
	}

	retRelations := &itemspb.ListRelationsResponse{
		Relations: make(map[string]*itemspb.ListRelation),
	}

	for _, relation := range relations {
		retRelations.Relations[relation.ID] = &itemspb.ListRelation{
			RelationId:   relation.ID,
			SourceItemId: relation.SourceItemID,
			TargetItemId: relation.TargetItemID,
			InfoItemId:   db.StringOrNil(relation.InfoItemID),
			Type:         relation.Type,
		}
	}

	return retRelations, nil
}
