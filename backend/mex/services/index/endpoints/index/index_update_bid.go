package index

import (
	"context"

	"github.com/d4l-data4life/mex/mex/shared/entities"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/utils"

	"github.com/d4l-data4life/mex/mex/services/index/endpoints/index/pb"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/datamodel"
)

func (svc *Service) IndexLatestItem(ctx context.Context, request *pb.IndexLatestItemRequest) (*pb.IndexLatestItemResponse, error) {
	queries := datamodel.New(svc.DB)

	// Get all the technical item IDs for the given business ID.
	items, err := queries.DbListItemsForBusinessId(ctx, request.BusinessId)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		svc.Log.Warn(ctx, L.Messagef("Indexing: could not find any stored items with business ID '%s'", request.BusinessId))
		return &pb.IndexLatestItemResponse{}, nil
	}

	// Do not index if the item is not of a focal type (as judged by looking at the newest version)
	newestVersionIndex := len(items) - 1
	newestItemVersion := items[newestVersionIndex]
	newestItemVersionID := newestItemVersion.ItemID
	isFocal, err := isOfFocalType(ctx, newestItemVersion, svc.EntityRepo)
	if err != nil {
		return nil, err
	}
	if !isFocal {
		return &pb.IndexLatestItemResponse{}, nil
	}

	svc.Log.Info(ctx, L.Messagef("Indexing: Will index the latest of %d version(s) of the entity with business ID %s", len(items), request.BusinessId))

	// Delete all versions of the item from Solr.
	// (Technically, there should be at most one, but trying to delete a non-existing document is fine in Solr.)
	err = svc.Solr.RemoveDocuments(ctx, utils.Map(items, func(x datamodel.ItemsWithBusinessID) string { return x.ItemID }))
	if err != nil {
		return nil, err
	}

	finalSQLQuery, err := svc.getSQLStatementForFieldValues(ctx, newestItemVersionID)
	if err != nil {
		return nil, err
	}
	rows, err := svc.DB.Query(ctx, finalSQLQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var latestItemValues []datamodel.CurrentItemValue
	rowFailCount := 0
	for rows.Next() {
		var i datamodel.CurrentItemValue
		if err := rows.Scan(
			&i.ItemID,
			&i.FieldName,
			&i.FieldValue,
			&i.Place,
			&i.Revision,
			&i.Language,
		); err != nil {
			// Skip rows that could not be read
			svc.Log.Warn(ctx, L.Messagef("failed to read row in item value return set: %s", err.Error()))
			rowFailCount++
			continue
		}
		latestItemValues = append(latestItemValues, i)
	}

	// Construct the Solr request payload.
	doc, err := svc.buildDocumentXML(ctx, latestItemValues)
	if err != nil {
		return nil, err
	}

	// And index the document.
	err = svc.Solr.AddDocuments(ctx, []string{doc}, 1)
	if err != nil {
		return nil, err
	}

	return &pb.IndexLatestItemResponse{}, nil
}

// isOfFocalType check if an item is of a focal item type
func isOfFocalType(ctx context.Context, item datamodel.ItemsWithBusinessID, entityRepo entities.EntityRepo) (bool, error) {
	// Check if this item is of a focal type and break if not (no indexing)
	focalEntityNames, err := entityRepo.GetEntityTypeNames(ctx, true)
	if err != nil {
		return false, err
	}
	return utils.Contains(focalEntityNames, item.EntityName), nil
}
