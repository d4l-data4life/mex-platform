package items

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/d4l-data4life/mex/mex/shared/cfg"
	"github.com/d4l-data4life/mex/mex/shared/constants"
	"github.com/d4l-data4life/mex/mex/shared/db"
	"github.com/d4l-data4life/mex/mex/shared/entities"
	E "github.com/d4l-data4life/mex/mex/shared/errstat"
	"github.com/d4l-data4life/mex/mex/shared/items"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/solr"
	"github.com/d4l-data4life/mex/mex/shared/utils"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/datamodel"
	itemspb "github.com/d4l-data4life/mex/mex/services/metadata/endpoints/items/pb"
)

// AggregateItems handles the HTTP requests for creating a new aggregated (merged) item based on one or more extracted items tied together by a
// common ID. It organizes the DB transaction for the actual DB logic.
func (svc *Service) AggregateItems(ctx context.Context, request *itemspb.AggregateItemsRequest) (*itemspb.AggregateItemsResponse, error) {
	tx, err := db.AcquireTx(ctx, svc.DB)
	if err != nil {
		svc.Log.Error(ctx, L.Messagef("aggregation: transaction creation failed"))
		return nil, E.MakeGRPCStatus(codes.Internal, "transaction creation failed", E.Cause(err), request).Err()
	}

	txCommit := false
	defer func() {
		if txCommit {
			if commitErr := tx.Commit(ctx); commitErr != nil {
				svc.Log.Error(ctx, L.Messagef("commit error: %s", commitErr.Error()))
			}

			// Do not announce items as we may be inside a bulk creation.
			// Instead, we report back the new business ID and let the calling function take
			// responsibility of the announcement.
		} else {
			svc.Log.Warn(ctx, L.Message("rolling back result of aggregation"))
			if err := tx.Rollback(ctx); err != nil {
				svc.Log.Error(ctx, L.Messagef("rollback error: %s", err.Error()))
			}
		}
	}()

	result, aggErr := svc.doItemsAggregation(ctx, aggregationInfo{
		EntityType:          request.EntityType,
		BusinessID:          request.BusinessId,
		PreventAnnouncement: false, // Explicit aggregation triggers announcement of new items
		DuplicateAlgorithm:  svc.DuplicateDetectionAlgorithm,
		DbTx:                tx,
	})
	if aggErr != nil {
		svc.Log.Error(ctx, L.Messagef("aggregation: aggregation failed"))
		return nil, E.MakeGRPCStatus(codes.Internal, "aggregation: aggregation failed", E.Cause(aggErr), request).Err()
	}

	txCommit = true

	return result, nil
}

type aggregationInfo struct {
	EntityType          string
	BusinessID          string
	PreventAnnouncement bool
	DuplicateAlgorithm  cfg.DuplicateDetectionAlgorithm
	DbTx                pgx.Tx
}

// doItemsAggregation handles the aggregation logic, using the passed DB transaction for any DB interactions.
func (svc *Service) doItemsAggregation(ctx context.Context, args aggregationInfo) (*itemspb.AggregateItemsResponse, error) {
	if args.BusinessID == "" {
		svc.Log.Error(ctx, L.Messagef("aggregation: cannot aggregate: no business ID given"))
		return nil, E.MakeGRPCStatus(codes.InvalidArgument, "cannot aggregate: no business ID given").Err()
	}

	aggregationSourceEntityTypeInfo, err := svc.EntityRepo.GetEntityType(ctx, args.EntityType)
	if err != nil {
		svc.Log.Error(ctx, L.Messagef("aggregation: unknown entity type"))
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("unknown entity type: %s: %s", args.EntityType, err.Error()))
	}
	mergeConfig := aggregationSourceEntityTypeInfo.Config

	if mergeConfig.AggregationAlgorithm == "" {
		svc.Log.Error(ctx, L.Messagef("aggregation: cannot aggregate - no algorithm given"))
		return nil, E.MakeGRPCStatus(codes.InvalidArgument, "cannot aggregate: no algorithm given").Err()
	}

	if mergeConfig.AggregationAlgorithm != solr.SimpleAggregation && mergeConfig.AggregationAlgorithm != solr.SourcePartitionAggregation {
		svc.Log.Error(ctx, L.Messagef("aggregation: cannot aggregate - unsupported algorithm"))
		return nil, E.MakeGRPCStatus(codes.InvalidArgument, fmt.Sprintf("cannot aggregate: unsupported algorithm: %s", mergeConfig.AggregationAlgorithm)).Err()
	}
	aggregationTargetEntityTypeInfo, err := svc.EntityRepo.GetEntityType(ctx, mergeConfig.AggregationEntityType)
	if err != nil {
		svc.Log.Error(ctx, L.Messagef("aggregation: unknown merge target entity type"))
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("unknown merge target entity type: %s: %s", mergeConfig.AggregationEntityType, err.Error()))
	}
	targetBusinessIDField := aggregationTargetEntityTypeInfo.Config.BusinessIdFieldName

	linkFieldDefs, err := svc.FieldRepo.GetFieldDefsByKind(ctx, "link")
	if err != nil {
		return nil, err
	}

	linkFieldNames := make([]string, len(linkFieldDefs))
	for i, fd := range linkFieldDefs {
		linkFieldNames[i] = fd.Name()
	}

	queries := datamodel.New(args.DbTx)
	// Get field values of aggregated item
	concatenatedValues, foundSourceItemIDs, aggErr := svc.getMergedItemValues(ctx, args.BusinessID, queries, mergeConfig, targetBusinessIDField, linkFieldNames)
	if aggErr != nil {
		svc.Log.Error(ctx, L.Messagef("aggregation: failed to retrieve merge candidate field values"))
		return nil, E.MakeGRPCStatus(codes.Internal, "failed to retrieve merge candidate field values", E.Cause(aggErr)).Err()
	}

	// Build and carry out create args for merged item, manually setting the business ID
	newBusinessID := solr.CreateMergeID(args.BusinessID)
	createResponses, err := svc.doItemsCreate(context.WithValue(ctx, db.ContextKeyTx, args.DbTx), createMultipleItemsInput{
		items: []*items.Item{
			{
				EntityType: mergeConfig.AggregationEntityType,
				BusinessId: newBusinessID,
				Values:     concatenatedValues,
			},
		},
		duplicateAlgorithm:  args.DuplicateAlgorithm,
		preventAnnouncement: args.PreventAnnouncement,
	})
	if err != nil {
		svc.Log.Error(ctx, L.Messagef("aggregated item creation failed: %s", err.Error()))
		return nil, E.MakeGRPCStatus(codes.Internal, "aggregated item creation failed", E.Cause(err)).Err()
	}

	// The case len(createResponses) = 0 happens if the aggregation results in a duplicate of an existing item
	// --> no new item created --> empty response
	result := &itemspb.AggregateItemsResponse{}
	if len(createResponses) == 1 {
		// Link merged item to the source fragments via a relation
		sourceItemIDs := utils.KeysOfMap(foundSourceItemIDs)
		for _, sourceItemID := range sourceItemIDs {
			_, err = svc.DoRelationCreate(context.WithValue(ctx, db.ContextKeyTx, args.DbTx), &itemspb.CreateRelationRequest{
				Type:         constants.MergedItemToFragmentRelation,
				SourceItemId: createResponses[0].itemID,
				TargetItemId: sourceItemID,
			})
			if err != nil {
				return nil, err
			}
		}

		// NOTE! createResponses[0].aggregatedBusinessID will NOT have the same value as newBusinessID - rather, it will be
		// the business ID of any FURTHER merged item triggered by the creations of the (merged) item we just constructed.
		svc.Log.Info(ctx, L.Messagef("aggregated item created (item ID: %s, business ID: %s) by merging item(s) with item ID(s) %s",
			createResponses[0].itemID, newBusinessID, strings.Join(sourceItemIDs, ", ")))
		result = &itemspb.AggregateItemsResponse{
			AggregateItemId:   createResponses[0].itemID,
			AggregatedItemIds: sourceItemIDs,
			NewBusinessId:     newBusinessID,
		}
	} else if len(createResponses) > 1 {
		return nil, status.Error(codes.Internal, fmt.Sprintf("expected aggregation to lead to the creation of at most 1 new item, but %d were created", len(createResponses)))
	}

	return result, nil
}

// getMergedItemValues merges item values from difference source items.
func (svc *Service) getMergedItemValues(ctx context.Context, sourceBusinessID string, queries *datamodel.Queries, aggregationConfig *entities.EntityTypeConfig,
	targetBusinessIDField string, linkFieldNames []string,
) ([]*items.ItemValue, map[string]struct{}, error) {
	var concatenatedValues []*items.ItemValue
	var foundSourceItemIDs map[string]struct{}

	switch aggregationConfig.AggregationAlgorithm {
	case solr.SimpleAggregation:
		/*
			This is the first and simplest aggregation algorithm. It works as follows:
			1. Collect all items that declare the same merge target (sourceBusinessId)
			2. Partition these items based on the set of *all* fields of kind 'link'
			3. Pick the newest item from each partition
			4. Merge all picked items into a single merged item (no duplication elimination)
		*/
		svc.Log.Info(ctx, L.Message("aggregation: using simple algorithm"))
		aggregatedValues, aggregationErr := queries.DbSimpleAggregationByBusinessID(ctx, linkFieldNames, sourceBusinessID)
		if aggregationErr != nil {
			return nil, nil, aggregationErr
		}

		concatenatedValues, foundSourceItemIDs = aggregateValuesSimple(aggregatedValues, targetBusinessIDField)
	case solr.SourcePartitionAggregation:
		/*
			This is a more focussed and flexible aggregation algorithm that is also based on faster SQL queries.
			It allows specifying how to deal with duplicated field values. It works as follows:
			1. Collect all items that declare the same merge target (sourceBusinessId)
			2. Partition these items based on the value in a single, configured partition field (typically the data source ID)
			3. Pick the newest item from each partition
			4. Merge all picked items into a single merged item using a specific duplication handling strategy
		*/

		// Fetch merge candidate values
		svc.Log.Info(ctx, L.Message("aggregation: using source partition algorithm"))
		aggregationArgs := datamodel.DbAggregationCandidateValuesParams{
			PartitionField: aggregationConfig.PartitionFieldName,
			BusinessID: pgtype.Text{
				String: sourceBusinessID,
				Valid:  true,
			},
			SourceIDField: aggregationConfig.BusinessIdFieldName,
		}
		aggregatedValuesWithSource, aggCandidateErr := queries.DbAggregationCandidateValues(ctx, aggregationArgs)
		if aggCandidateErr != nil {
			return nil, nil, aggCandidateErr
		}

		var err error
		concatenatedValues, foundSourceItemIDs, err = doFlexibleMerge(aggregatedValuesWithSource, aggregationConfig.DuplicateStrategy, targetBusinessIDField)
		if err != nil {
			return nil, nil, err
		}
	default:
		return nil, nil, fmt.Errorf("unknown aggregation algorithm")
	}

	return concatenatedValues, foundSourceItemIDs, nil
}

func doFlexibleMerge(aggregatedValuesWithSource []datamodel.DbAggregationCandidateValuesRow,
	duplicateStrategy string, targetBusinessIDFieldName string,
) ([]*items.ItemValue, map[string]struct{}, error) {
	// Do the actual merge with the configured duplicate strategy
	var concatenatedValues []*items.ItemValue
	var foundSourceItemIDs map[string]struct{}
	if duplicateStrategy == "" {
		duplicateStrategy = solr.DefaultDuplicateStrategy
	}
	switch duplicateStrategy {
	case solr.KeepAllDuplicates:
		concatenatedValues, foundSourceItemIDs = keepAllDuplicatesMerge(aggregatedValuesWithSource, targetBusinessIDFieldName)
	case solr.RemoveAllDuplicates:
		concatenatedValues, foundSourceItemIDs = removeAllDuplicatesMerge(aggregatedValuesWithSource, targetBusinessIDFieldName)
	default:
		return nil, nil, fmt.Errorf("unknown duplication handling strategy")
	}
	return concatenatedValues, foundSourceItemIDs, nil
}

// aggregateValuesSimple concatenates all available values from each field
func aggregateValuesSimple(aggregationCandidateValues []datamodel.DbSimpleAggregationByBusinessIdRow, targetBusinessIDFieldName string) ([]*items.ItemValue, map[string]struct{}) {
	// Add the concatenation of all other fields
	sourceItemIDs := make(map[string]struct{})
	var concatenatedValues []*items.ItemValue

	// Add the concatenation of all other fields
	for _, value := range aggregationCandidateValues {
		/* Do not add values to the business ID field of the merged resource, even if the source had a
		(non-business-ID) field with the same name. The business ID field of the *source* should have been
		filtered out from the list of candidate values already. */
		if value.FieldName == targetBusinessIDFieldName {
			continue
		}
		concatenatedValues = append(concatenatedValues, &items.ItemValue{
			FieldName:  value.FieldName,
			FieldValue: value.FieldValue,
			Language:   value.Language.String,
		})
		sourceItemIDs[value.ItemID] = struct{}{}
	}
	return concatenatedValues, sourceItemIDs
}

// keepAllDuplicatesMerge keeps all duplicated field values in the output, regardless of where they came from
func keepAllDuplicatesMerge(aggregationCandidateEntries []datamodel.DbAggregationCandidateValuesRow, targetBusinessIDFieldName string) ([]*items.ItemValue, map[string]struct{}) {
	preMergeItemIDs := make(map[string]struct{})
	var curConcatenatedValues []*items.ItemValue
	curItemID := ""
	firstPartitionBusinessIDForItem := ""
	for _, entry := range aggregationCandidateEntries {
		/* Do not add values to the business ID field of the merged resource, even if the source had a
		(non-business-ID) field with the same name. The business ID field of the *source* should have been
		filtered out from the list of candidate values already. */
		if entry.FieldName == targetBusinessIDFieldName {
			continue
		}
		if entry.ItemID != curItemID {
			// When we get to a new item, record the first partition (source) seen for it
			curItemID = entry.ItemID
			firstPartitionBusinessIDForItem = entry.PartitionBusinessID
		} else if entry.PartitionBusinessID != firstPartitionBusinessIDForItem {
			/*
				Skip values if we are seeing the *same* item from a *different* partition (source).
				This covers the case where there a single pre-merge item is the newest item in multiple partitions
				because it listed multiple sources. NOTE: This logic depends on the required ordering of the input.
			*/
			continue
		}

		curConcatenatedValues = append(curConcatenatedValues, &items.ItemValue{
			FieldName:  entry.FieldName,
			FieldValue: entry.FieldValue,
			Language:   entry.Language.String,
		})
		preMergeItemIDs[entry.ItemID] = struct{}{}
	}
	return curConcatenatedValues, preMergeItemIDs
}

// removeAllDuplicatesMerge eliminates all duplicated field values in the output, regardless of where they came from
func removeAllDuplicatesMerge(aggregationCandidateEntries []datamodel.DbAggregationCandidateValuesRow, targetBusinessIDFieldName string,
) ([]*items.ItemValue, map[string]struct{}) {
	preMergeItemIDs := make(map[string]struct{})
	var curConcatenatedValues []*items.ItemValue
	valuesSeen := make(map[string]map[string]struct{})
	for _, entry := range aggregationCandidateEntries {
		/* Do not add values to the business ID field of the merged resource, even if the source had a
		(non-business-ID) field with the same name. The business ID field of the *source* should have been
		filtered out from the list of candidate values already. */
		if entry.FieldName != targetBusinessIDFieldName {
			if _, fieldSeen := valuesSeen[entry.FieldName]; !fieldSeen {
				// New field (and so also new field value)
				valuesSeen[entry.FieldName] = map[string]struct{}{entry.FieldValue: {}}
			} else {
				// Known field
				if _, valueSeen := valuesSeen[entry.FieldName][entry.FieldValue]; valueSeen {
					// Skip if we have ever seen this value for this field
					continue
				}
				// New field value - remember it
				valuesSeen[entry.FieldName][entry.FieldValue] = struct{}{}
			}

			curConcatenatedValues = append(curConcatenatedValues, &items.ItemValue{
				FieldName:  entry.FieldName,
				FieldValue: entry.FieldValue,
				Language:   entry.Language.String,
			})
		}
		preMergeItemIDs[entry.ItemID] = struct{}{}
	}
	return curConcatenatedValues, preMergeItemIDs
}
