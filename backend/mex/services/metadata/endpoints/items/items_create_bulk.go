package items

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/d4l-data4life/mex/mex/shared/auth"
	"github.com/d4l-data4life/mex/mex/shared/cfg"
	"github.com/d4l-data4life/mex/mex/shared/constants"
	"github.com/d4l-data4life/mex/mex/shared/db"
	"github.com/d4l-data4life/mex/mex/shared/errstat"
	"github.com/d4l-data4life/mex/mex/shared/hints"
	"github.com/d4l-data4life/mex/mex/shared/items"
	"github.com/d4l-data4life/mex/mex/shared/known/jobspb"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/uuid"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/canonical"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/datamodel"
	pbItems "github.com/d4l-data4life/mex/mex/services/metadata/endpoints/items/pb"
)

/*
CreateItemsBulk handles the incoming HTTP request for bulk creations and wraps the item creation in an asynchronous job.
It synchronously returns the job ID, allowing clients to monitor job progress. The item IDs of the created items are
fed into the job output before completing.
*/
func (svc *Service) CreateItemsBulk(ctx context.Context, request *pbItems.CreateItemsBulkRequest) (*pbItems.CreateItemsBulkResponse, error) {
	if len(request.Items) == 0 {
		return nil, errstat.MakeGRPCStatus(codes.InvalidArgument, "no items to be created passed", request).Err()
	}

	lock, err := svc.Jobber.AcquireLock(ctx, SvcResourceName)
	if err != nil {
		return nil, errstat.MakeGRPCStatus(codes.AlreadyExists, "failed to acquire items lock; other job might be running", request).Err()
	}

	job, err := svc.Jobber.CreateJob(ctx, &jobspb.CreateJobRequest{
		Title: "bulk item creation",
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failure creating job:  %s", err.Error()))
	}

	// Add the job ID to the context so that CreateItem can log created item IDs.
	// (The NewContextWithValues function puts more IDs into the new context than we need here, but that is okay.)
	ctxJob := constants.NewContextWithValues(ctx, job.JobId)

	go func(ctx context.Context) {
		svc.Log.Info(ctx, L.Messagef("bulk item creation: job started (%s)", job.JobId), L.Phase("job"))

		svc.Jobber.SetStatusRunning(ctx, job.JobId)              //nolint:errcheck
		defer svc.Jobber.SetStatusDone(ctx, job.JobId)           //nolint:errcheck
		defer svc.Jobber.ReleaseLock(ctx, SvcResourceName, lock) //nolint:errcheck

		// Override configured duplicate detection algorithm only if explicitly requested
		duplicateAlgorithm := svc.DuplicateDetectionAlgorithm
		if request.OverrideDuplicateAlgorithm {
			duplicateAlgorithm = request.DuplicateAlgorithm
		}

		// Create items
		itemResults, processErr := svc.doItemsCreate(ctx, createMultipleItemsInput{
			items:               request.Items,
			precomputedHashes:   []string{},
			duplicateAlgorithm:  duplicateAlgorithm,
			preventAnnouncement: true, // Do not announce the items created as part of a bulk load
		})
		if processErr != nil {
			svc.Log.Error(ctx, L.Message(processErr.Error()))
			_, err := svc.Jobber.SetError(ctx, &jobspb.SetJobErrorRequest{
				Error: processErr.Error(),
				JobId: job.JobId,
			})
			if err != nil {
				svc.Log.Warn(ctx, L.Messagef("could not set job error: %s", err.Error()))
			}
			return
		}

		svc.Log.Info(ctx, L.Messagef("bulk item creation job done (job ID: %s) - %d items created", job.JobId, len(itemResults)), L.Phase("job"))
	}(ctxJob)

	// Synchronous return
	hints.HintHTTPStatusCode(ctx, http.StatusCreated)
	return &pbItems.CreateItemsBulkResponse{JobId: job.JobId}, nil
}

type createSingleItemInput struct {
	Item                *items.Item
	preventAnnouncement bool
	duplicateAlgorithm  cfg.DuplicateDetectionAlgorithm
	hash                string
	dbTx                pgx.Tx
}

type createSingleItemResult struct {
	itemID               string
	aggregatedBusinessID string
}

type createMultipleItemsInput struct {
	items               []*items.Item
	precomputedHashes   []string
	duplicateAlgorithm  cfg.DuplicateDetectionAlgorithm
	preventAnnouncement bool
}

/*
doItemsCreate handles the DB transaction and duplication check associated with the creation of one or more new items.
Note that while the returned item IDs refers to the items created based on the information passed directly to this
method, the returned business IDs are those of any merged items created in response to the creation of a "primary" item.
*/
func (svc *Service) doItemsCreate(ctx context.Context, creationArgs createMultipleItemsInput) ([]createSingleItemResult, error) {
	tx, err := db.AcquireTx(ctx, svc.DB)
	if err != nil {
		return nil, fmt.Errorf("transaction creation failed: %s", err.Error())
	}

	/*
		Handle rollback and item announcement at end of function call.
		NOTE: The new item IDs must only be announced AFTER the transaction is committed, otherwise the items
		will not be visible to the aggregation logic (which runs asynchronously outside the transaction).
	*/
	txCommit := false
	var results []createSingleItemResult
	defer func() {
		if txCommit {
			if err := tx.Commit(ctx); err != nil {
				svc.Log.Error(ctx, L.Messagef("commit error: %s", err.Error()))
			} else {
				// Announce the IDs of the created items, if not prevented
				if !creationArgs.preventAnnouncement {
					for _, singleCreateResult := range results {
						svc.Announcer.AnnounceTechnicalItemID(singleCreateResult.itemID)
						if singleCreateResult.aggregatedBusinessID != "" {
							svc.Announcer.AnnounceBusinessItemID(singleCreateResult.aggregatedBusinessID)
						}
					}
				}
				// If there is a job service running, add the item IDs of the (primary) items created to it
				if svc.Jobber != nil {
					resultItemIDs := make([]string, len(results))
					for i, singleCreateResult := range results {
						resultItemIDs[i] = singleCreateResult.itemID
					}
					AddJobItemIDs(ctx, svc.Jobber, resultItemIDs)
				}
			}
		} else {
			if err := tx.Rollback(ctx); err != nil {
				svc.Log.Error(ctx, L.Messagef("rollback error: %s", err.Error()))
			}
		}
	}()

	// Ensure all item hashes are available, then detect duplicates
	itemHashesList, hashErr := getEffectiveHashList(creationArgs.items, creationArgs.precomputedHashes)
	if hashErr != nil {
		return nil, hashErr
	}
	duplicateCount, isDuplicateHash, duplicateErr := findDuplicates(ctx, duplicateArgs{
		hashes:             itemHashesList,
		duplicateAlgorithm: creationArgs.duplicateAlgorithm,
		dbTx:               tx,
	})
	if duplicateErr != nil {
		return nil, duplicateErr
	}
	svc.Log.Info(ctx, L.Messagef("hash check using algorithm '%s' completed - %d out of %d uploaded items classified as duplicates",
		svc.DuplicateDetectionAlgorithm, duplicateCount, len(creationArgs.items)))

	// Create all new items
	var createErr error
	var singleCreateResponse createSingleItemResult
	for i, item := range creationArgs.items {
		if isDuplicateHash[i] {
			continue // Skip duplicates
		}
		singleCreateResponse, createErr = svc.doSingleItemCreate(context.WithValue(ctx, db.ContextKeyTx, tx), createSingleItemInput{
			Item:                item,
			preventAnnouncement: creationArgs.preventAnnouncement,
			duplicateAlgorithm:  creationArgs.duplicateAlgorithm,
			hash:                itemHashesList[i],
			dbTx:                tx,
		})
		if createErr != nil {
			return nil, fmt.Errorf("error creating item %d/%d: %s", i+1, len(creationArgs.items), createErr.Error())
		}

		results = append(results, singleCreateResponse)
	}
	txCommit = true

	return results, nil
}

/*
doSingleItemCreate handles validation of the input for the creation of a new item, triggers aggregation where needed,
and structures the data passed for items creation. A non-empty hash value for the item must be passed, otherwise an
error is returned.
*/
func (svc *Service) doSingleItemCreate(ctx context.Context, input createSingleItemInput) (createSingleItemResult, error) {
	var aggregateItemBusinessID string

	user, err := auth.GetMexUser(ctx)
	if err != nil {
		return createSingleItemResult{}, err
	}
	if input.Item == nil {
		return createSingleItemResult{}, fmt.Errorf("input item is nil")
	}
	if input.Item.EntityType == "" {
		return createSingleItemResult{}, fmt.Errorf("entity type cannot be empty")
	}
	if input.hash == "" {
		return createSingleItemResult{}, fmt.Errorf("no hash value passed for new item")
	}

	entityType, err := svc.EntityRepo.GetEntityType(ctx, input.Item.EntityType)
	if err != nil {
		return createSingleItemResult{}, fmt.Errorf("unknown entity type: %s: %s", input.Item.EntityType, err.Error())
	}

	// Validate all the field values
	for _, v := range input.Item.Values {
		fieldDef, fieldsErr := svc.FieldRepo.GetFieldDefByName(ctx, v.FieldName)
		if fieldsErr != nil {
			return createSingleItemResult{}, fmt.Errorf("could not retrieve config for the field " + v.FieldName)
		}

		hook := svc.ItemCreationHooks.GetHook(fieldDef.Kind())
		if hook == nil {
			return createSingleItemResult{}, fmt.Errorf("no field value validation hook for kind: " + fieldDef.Kind())
		}

		if valErr := hook.ValidateFieldValue(ctx, fieldDef, v.FieldValue); valErr != nil {
			return createSingleItemResult{}, fmt.Errorf("invalid field value for: %s (%s)", v.FieldName, v.FieldValue)
		}
	}

	// Copy the input values and determine the place indexes of multi-values fields, updated business ID if needed
	businessID := input.Item.BusinessId // Could be "".
	var placerObj placer
	values := make([]*placedItemValue, len(input.Item.Values))
	for i, v := range input.Item.Values {
		values[i] = &placedItemValue{
			FieldName:  v.FieldName,
			FieldValue: v.FieldValue,
			Language:   v.Language,
			Place:      placerObj.nextPlaceOf(v.FieldName),
		}
		// Update business ID if (1) it is not already set and (2) we found the business ID field.
		// Hence, if the business ID field occurs multiple times, only the first instance will be used.
		if businessID == "" && v.FieldName == entityType.Config.BusinessIdFieldName {
			businessID = v.FieldValue
		}
	}

	// Create new item
	createItemArgs := singleCreateArgs{
		owner:               user.UserId,
		entityType:          entityType.Name,
		businessIDFieldName: entityType.Config.BusinessIdFieldName,
		businessID:          businessID,
		hash:                input.hash,
		values:              values,
		dbTx:                input.dbTx,
	}
	createdItem, err := createNewItem(ctx, createItemArgs)
	if err != nil {
		return createSingleItemResult{}, err
	}

	// Impute business IDs before aggregation so the views used in aggregation have access to imputed business IDs.
	_, imputeErr := imputeBusinessIDs(ctx, input.dbTx)
	if imputeErr != nil {
		return createSingleItemResult{}, imputeErr
	}

	// Run aggregation if required
	if entityType.Config.AggregationAlgorithm != "" {
		svc.Log.Info(ctx, L.Messagef("aggregating with algorithm: %s (%s)", entityType.Config.AggregationAlgorithm, input.Item.EntityType))
		var aggErr error
		aggregateItemBusinessID, aggErr = svc.runAggregationLogic(ctx, aggregationLogicArgs{
			sourceEntityType:    input.Item.EntityType,
			sourceBusinessID:    businessID,
			preventAnnouncement: input.preventAnnouncement,
			duplicateAlgorithm:  input.duplicateAlgorithm,
			dbTx:                input.dbTx,
		})
		if aggErr != nil {
			return createSingleItemResult{}, aggErr
		}
	}

	return createSingleItemResult{
		itemID:               createdItem.ID,
		aggregatedBusinessID: aggregateItemBusinessID,
	}, nil
}

type singleCreateArgs struct {
	dbTx                pgx.Tx
	owner               string
	entityType          string
	businessIDFieldName string
	businessID          string
	hash                string
	values              []*placedItemValue
}

// createNewItem handles the core creation of DB entities within the passed DB transaction.
func createNewItem(ctx context.Context, args singleCreateArgs) (datamodel.Item, error) {
	queries := datamodel.New(args.dbTx)

	now := time.Now()
	itemID := uuid.MustNewV4()
	createdItem, err := queries.DbCreateItem(ctx, datamodel.DbCreateItemParams{
		CreatedAt:           pgtype.Timestamptz{Time: now, Valid: true},
		ID:                  itemID,
		Owner:               args.owner,
		EntityName:          args.entityType,
		BusinessID:          db.TextFromString(args.businessID),
		BusinessIDFieldName: db.TextFromString(args.businessIDFieldName),
		Hash:                db.TextFromString(args.hash),
	})
	if err != nil {
		return datamodel.Item{}, status.Error(codes.Internal, fmt.Sprintf("failed to create new item in DB: %s", err.Error()))
	}

	for _, v := range args.values {
		_, valErr := queries.DbCreateItemValue(ctx, datamodel.DbCreateItemValueParams{
			CreatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
			ID:         uuid.MustNewV4(),
			Deleted:    false,
			FieldName:  v.FieldName,
			FieldValue: v.FieldValue,
			Language:   db.TextFromString(v.Language),
			Place:      v.Place,
			ItemID:     itemID,
		})
		if valErr != nil {
			return datamodel.Item{}, status.Error(codes.Internal, fmt.Sprintf("failed to create new item values in DB: %s", valErr.Error()))
		}
	}

	return createdItem, nil
}

type aggregationLogicArgs struct {
	sourceEntityType    string
	sourceBusinessID    string
	preventAnnouncement bool
	duplicateAlgorithm  cfg.DuplicateDetectionAlgorithm
	dbTx                pgx.Tx
}

// runAggregationLogic triggers aggregation based on a new item, returning the business ID of the created merged item.
func (svc *Service) runAggregationLogic(ctx context.Context, args aggregationLogicArgs) (string, error) {
	startTime := time.Now()

	aggregationResponse, err := svc.doItemsAggregation(context.WithValue(ctx, db.ContextKeyTx, args.dbTx), aggregationInfo{
		EntityType:          args.sourceEntityType,
		BusinessID:          args.sourceBusinessID,
		PreventAnnouncement: args.preventAnnouncement,
		DuplicateAlgorithm:  args.duplicateAlgorithm,
		DbTx:                args.dbTx,
	})
	if err != nil {
		return "", err
	}
	svc.Log.Info(ctx, L.Messagef("item create: aggregation finished in %d ms", time.Since(startTime).Milliseconds()))

	return aggregationResponse.NewBusinessId, nil
}

type duplicateArgs struct {
	hashes             []string
	duplicateAlgorithm cfg.DuplicateDetectionAlgorithm
	dbTx               pgx.Tx
}

/*
findDuplicates returns no. of duplicates found and a string slice of the same length as input the hash slice. In the
returned slice, the boolean value at a given position indicates of the corresponding hash is duplicate hash or not.
*/
func findDuplicates(ctx context.Context, args duplicateArgs) (int, []bool, error) {
	var err error
	queries := datamodel.New(args.dbTx)
	var foundHashesList []pgtype.Text
	switch args.duplicateAlgorithm {
	case cfg.DuplicateDetectionAlgorithm_SIMPLE:
		foundHashesList, err = queries.DbListHashesPresentSimple(ctx, args.hashes)
	case cfg.DuplicateDetectionAlgorithm_LATEST_ONLY:
		foundHashesList, err = queries.DbListHashesPresentLatestOnly(ctx, args.hashes)
	default:
		return 0, nil, fmt.Errorf("unknown duplicate detection algorithm specified")
	}
	if err != nil {
		return 0, nil, fmt.Errorf("item hash check failed: %s", err.Error())
	}

	foundHashesMap := make(map[string]struct{})
	for _, hash := range foundHashesList {
		if hash.Valid {
			foundHashesMap[hash.String] = struct{}{}
		}
	}
	duplicateCount := 0
	duplicateHashes := make([]bool, len(args.hashes))
	for j, hashVal := range args.hashes {
		if _, ok := foundHashesMap[hashVal]; ok {
			duplicateHashes[j] = true
			duplicateCount++
		}
	}

	return duplicateCount, duplicateHashes, nil
}

// imputeBusinessIDs trigger the imputation of item business IDs in the DB.
func imputeBusinessIDs(ctx context.Context, tx pgx.Tx) (int64, error) {
	queries := datamodel.New(tx)
	updateCount, err := queries.DbImputeBusinessId(ctx)
	if err != nil {
		return 0, fmt.Errorf("business ID imputation failed: %s", err.Error())
	}
	return updateCount, nil
}

// getEffectiveHashList returns the hash values for all items (in the same order as the items),
// using the pre-computed values if present and otherwise computing them.
func getEffectiveHashList(items []*items.Item, precomputedHashes []string) ([]string, error) {
	/*
		Cases:
		1. no hashes passed --> calculate hashes
		2. same nos. of items and hashes but all hashes are empty strings --> calculate hashes
		3. exactly one non-empty hash for each item given --> use the pre-computed hashes
		4. mixture of empty and non-empty pre-computed hashes passed --> error
		5. nos. of items and pre-computed hashes do not match  --> error
	*/
	itemHashesList := make([]string, len(items))
	switch len(precomputedHashes) {
	case 0:
		// Case (1)
		for i, item := range items {
			itemHashesList[i] = canonical.Fingerprint(item)
		}
	case len(items):
		emptyHashCount := 0
		for _, hash := range precomputedHashes {
			if hash == "" {
				emptyHashCount++
			}
		}
		switch emptyHashCount {
		case 0:
			// Case (3)
			itemHashesList = precomputedHashes
		case len(precomputedHashes):
			// Case (2)
			for i, item := range items {
				itemHashesList[i] = canonical.Fingerprint(item)
			}
		default:
			// Case (4)
			return nil, fmt.Errorf("pre-computed hashes is a mixture of empty and non-empty strings")
		}
	default:
		// Case (5)
		return nil, fmt.Errorf("no. of pre-computed hashes (%d) does not match no. of items (%d)", len(precomputedHashes), len(items))
	}

	return itemHashesList, nil
}

// Implement sort.Interface to sort item values by field name and place
type ByPlace []*placedItemValue

func (hList ByPlace) Len() int {
	return len(hList)
}

func (hList ByPlace) Swap(i int, j int) {
	hList[i], hList[j] = hList[j], hList[i]
}

func (hList ByPlace) Less(i int, j int) bool {
	if hList[i].FieldName == hList[j].FieldName {
		// For same field names, place determines the order.
		return hList[i].Place < hList[j].Place
	}
	// For non-equal field names, those names determine the order.
	return hList[i].FieldName < hList[j].FieldName
}

type placer struct {
	indexes map[string]int32
}

func (p *placer) nextPlaceOf(fieldName string) int32 {
	if p.indexes == nil {
		p.indexes = make(map[string]int32)
	}

	if idx, ok := p.indexes[fieldName]; ok {
		p.indexes[fieldName] = idx + 1
		return idx + 1
	}
	p.indexes[fieldName] = 1
	return 1
}
