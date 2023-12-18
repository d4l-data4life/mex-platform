package index

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/d4l-data4life/mex/mex/shared/constants"
	"github.com/d4l-data4life/mex/mex/shared/entities/erepo"
	E "github.com/d4l-data4life/mex/mex/shared/errstat"
	fieldUtils "github.com/d4l-data4life/mex/mex/shared/fields"
	"github.com/d4l-data4life/mex/mex/shared/hints"
	"github.com/d4l-data4life/mex/mex/shared/known/jobspb"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/solr"
	"github.com/d4l-data4life/mex/mex/shared/utils"

	"github.com/d4l-data4life/mex/mex/services/index/endpoints/index/pb"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/datamodel"
	kindHierarchy "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/hierarchy"
	kindLink "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/link"
)

// UpdateIndex load data from the DB into Solr.
func (svc *Service) UpdateIndex(ctx context.Context, request *pb.UpdateIndexRequest) (*pb.UpdateIndexResponse, error) {
	lock, err := svc.JobService.AcquireLock(ctx, SvcResourceName)
	if err != nil {
		return nil, E.MakeGRPCStatus(codes.AlreadyExists, "update: failed to acquire index lock; other job might be running", request).Err()
	}

	job, err := svc.JobService.CreateJob(ctx, &jobspb.CreateJobRequest{
		Title: "Repopulate Solr index",
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failure creating job:  %s", err.Error()))
	}

	// Create new independent context for the job copying the relevant values.
	// (The request's ctx will go out of scope before the job is done, so we cannot use it directly or as a parent context.)
	ctxJob := constants.NewContextWithValues(ctx, job.JobId)

	go func(ctx context.Context) {
		logJobError := func(message string) {
			svc.Log.Error(ctx, L.Message(message))
			_, err := svc.JobService.SetError(ctx, &jobspb.SetJobErrorRequest{Error: message, JobId: job.JobId})
			if err != nil {
				svc.Log.Warn(ctx, L.Messagef("could not set job error: %s", err.Error()))
			}
		}

		svc.Log.Info(ctx, L.Messagef("Solr data load: job started (%s)", job.JobId), L.Phase("job"))

		svc.JobService.SetStatusRunning(ctx, job.JobId)              //nolint:errcheck
		defer svc.JobService.SetStatusDone(ctx, job.JobId)           //nolint:errcheck
		defer svc.JobService.ReleaseLock(ctx, SvcResourceName, lock) //nolint:errcheck

		indexErr := svc.DoIndexUpdate(ctx, svc.TelemetryService)
		if indexErr != nil {
			logJobError(fmt.Sprintf("error during index population: %s", indexErr.Error()))
			return
		}

		svc.Log.Info(ctx, L.Messagef("Solr data load: job done (%s)", job.JobId), L.Phase("job"))
		svc.TelemetryService.Done()
	}(ctxJob)

	hints.HintHTTPStatusCode(ctx, http.StatusCreated)
	return &pb.UpdateIndexResponse{JobId: job.JobId}, nil
}

const progressReportSize = 1000

func (svc *Service) DoIndexUpdate(ctx context.Context, progressor utils.Progressor) error {
	if progressor == nil {
		progressor = &utils.NopProgressor{}
	}

	// Reset cache in Coding implementations
	for _, hook := range svc.SolrDataLoadHooks {
		hook.ResetCaches()
	}

	err := svc.iterateItems(ctx, progressor)
	return err
}

type iteratorState struct {
	currentItemID  string
	curItem        []datamodel.CurrentItemValue
	docBatch       []string
	docNoInBatch   int // No. of docs stored in current batch
	count          int // No of Solr docs XML objects created (regardless of whether they could also be uploaded)
	batchCount     int // No. of batches of documents successfully uploaded to Solr
	rowFailCount   int
	docFailCount   int
	batchFailCount int
}

func newIteratorState(batchSize int) iteratorState {
	return iteratorState{
		currentItemID:  "",
		docNoInBatch:   0,
		count:          0,
		batchCount:     0,
		rowFailCount:   0,
		docFailCount:   0,
		batchFailCount: 0,
		docBatch:       make([]string, batchSize),
		curItem:        []datamodel.CurrentItemValue{},
	}
}

/*
iterateItems pulls all information for the items to be indexed and sends them to Solr

The design with the complex state and the many small functions being called seem look a little overblown. However,
it was deliberately introduced to make the subtle logic "screechingly obvious" and testable. Future devs are
invited to search for a potentially better middle ground.
*/
func (svc *Service) iterateItems(ctx context.Context, progressor utils.Progressor) error {
	if progressor == nil {
		progressor = &utils.NopProgressor{}
	}

	// Construct the final query
	finalSQLQuery, err := svc.getSQLStatementForFieldValues(ctx, "")
	if err != nil {
		return err
	}

	startTime := time.Now()
	svc.Log.Info(ctx, L.Message("begin: indexable items query"))
	rows, err := svc.DB.Query(ctx, finalSQLQuery)
	if err != nil {
		svc.Log.Error(ctx, L.Messagef("failed to retrieve items to index from DB: %s", err.Error()))
		return err
	}
	svc.Log.Info(ctx, L.Messagef("end: indexable items query (%d ms)", time.Since(startTime).Milliseconds()))

	defer rows.Close()

	state := newIteratorState(solr.IndexBatchSize)

	i := 0
	for rows.Next() {
		var itemValue datamodel.CurrentItemValue
		if err := rows.Scan(
			&itemValue.ItemID,
			&itemValue.FieldName,
			&itemValue.FieldValue,
			&itemValue.Place,
			&itemValue.Revision,
			&itemValue.Language,
		); err != nil {
			// Skip rows that could not be read
			state = svc.handleSkippedRow(ctx, err, state)
			continue
		}

		state = svc.processItemValue(ctx, itemValue, state, solr.IndexBatchSize)

		// Update progress
		if i%progressReportSize == 0 {
			progressor.Progress("indexing", fmt.Sprintf("processed item values: %d", i))
		}
		i++
	}

	// Finish last, open item and force indexing of all remaining items
	state = svc.finishItem(ctx, state, solr.IndexBatchSize, true)

	svc.logReport(ctx, state)

	return nil
}

func isNewItemID(itemID string, state iteratorState) bool {
	return itemID != state.currentItemID
}

func curItemIsNonEmpty(state iteratorState) bool {
	return len(state.curItem) > 0
}

func batchIsFull(state iteratorState, batchSize int) bool {
	return state.docNoInBatch == batchSize
}

// processItemValue stores the value, updating index batch and triggering indexing if needed
func (svc *Service) processItemValue(ctx context.Context, itemValue datamodel.CurrentItemValue, state iteratorState, batchSize int) iteratorState {
	// Finish item when item ID changes in the list of values (also triggered by first row)
	if isNewItemID(itemValue.ItemID, state) {
		state = svc.finishItem(ctx, state, batchSize, false)
		state = startNewItem(itemValue.ItemID, state)
	}

	// Add the value to the current item
	state = addValueToItem(itemValue, state)
	return state
}

// finishItem add the item to the index batch, triggering indexing if needed
func (svc *Service) finishItem(ctx context.Context, state iteratorState, batchSize int, forceUpdate bool) iteratorState {
	if curItemIsNonEmpty(state) {
		state = svc.addItemToBatch(ctx, state)
	}
	if batchIsFull(state, batchSize) || forceUpdate {
		state = svc.indexBatch(ctx, state)
	}
	return state
}

func (svc *Service) logReport(ctx context.Context, state iteratorState) {
	svc.Log.Info(ctx, L.Messagef("loading done: %d document(s) in %d batch(es)", state.count, state.batchCount))
	if isErrorState(state) {
		svc.Log.Warn(ctx, L.Messagef("not all data was successfully uploaded: %d row(s), %d item(s), and %d batch(es) were skipped",
			state.rowFailCount, state.docFailCount, state.batchFailCount))
	}
}

func (svc *Service) handleSkippedRow(ctx context.Context, err error, state iteratorState) iteratorState {
	svc.Log.Warn(ctx, L.Messagef("failed to read row in item value return set: %s", err.Error()))
	state.rowFailCount++
	return state
}

func addValueToItem(value datamodel.CurrentItemValue, state iteratorState) iteratorState {
	state.curItem = append(state.curItem, value)
	return state
}

func startNewItem(itemID string, state iteratorState) iteratorState {
	state.currentItemID = itemID
	state.curItem = []datamodel.CurrentItemValue{}
	return state
}

// isErrorState checks if any actions failed
func isErrorState(state iteratorState) bool {
	return state.rowFailCount > 0 || state.docFailCount > 0 || state.batchFailCount > 0
}

// addItemToBatch builds the XML for a completed item and adds it to the current batch
func (svc *Service) addItemToBatch(ctx context.Context, state iteratorState) iteratorState {
	docXML, err := svc.buildDocumentXML(ctx, state.curItem)
	if err != nil {
		// Skip document if XML for Solr could not be built
		svc.Log.Warn(ctx, L.Messagef("failed to build document XML for Solr: %s", err.Error()))
		state.docFailCount++
	} else {
		state.count++
		state.docBatch[state.docNoInBatch] = docXML
		state.docNoInBatch++
	}
	return state
}

// indexBatch indexes a finished batch of items
func (svc *Service) indexBatch(ctx context.Context, state iteratorState) iteratorState {
	err := svc.Solr.AddDocuments(ctx, state.docBatch, state.docNoInBatch)
	if err != nil {
		// Skip batch that could not be uploaded to Solr
		svc.Log.Warn(ctx, L.Messagef("failed to upload document batch %d to Solr: %s", state.batchCount, err.Error()))
		state.batchFailCount++
	} else {
		state.batchCount++
	}
	// Clean buffer in preparation for next batch
	state.docBatch = make([]string, solr.IndexBatchSize)
	state.docNoInBatch = 0
	return state
}

// buildDocumentXML constructs the XML object to be sent to Solr to index the passed item
func (svc *Service) buildDocumentXML(ctx context.Context, item []datamodel.CurrentItemValue) (string, error) {
	if len(item) == 0 {
		return "", fmt.Errorf("cannot index empty item")
	}

	// Build XML <doc> tag for later Solr indexing
	var sb strings.Builder
	sb.WriteString("\t<doc>\n")

	idSeen := false
	for _, v := range item {
		fieldDef, err := svc.FieldRepo.GetFieldDefByName(ctx, v.FieldName)
		if err != nil {
			return "", err
		}

		hook := svc.SolrDataLoadHooks.GetHook(fieldDef.Kind())
		if hook == nil {
			return "", fmt.Errorf("hook not found for kind: %s", fieldDef.Kind())
		}

		var tags []string
		if v.FieldName == solr.DefaultUniqueKey {
			idSeen = true
			// Special handling for hardcoded ID field - do not generate any extra entries
			cleanedValue := utils.SanitizeXML(v.FieldValue)
			tags = []string{fmt.Sprintf("<field name=\"%s\">%s</field>", v.FieldName, cleanedValue)}
		} else {
			tags, err = hook.GenerateXMLFieldTags(ctx, fieldDef, v)
			if err != nil {
				return "", err
			}
		}
		for _, tag := range tags {
			sb.WriteString("\t\t" + tag + "\n")
		}
	}
	if !idSeen {
		return "", fmt.Errorf("generated XML documents for index has no entry for the Solr ID field '%s'", solr.DefaultUniqueKey)
	}

	sb.WriteString("\t</doc>\n")

	return sb.String(), nil
}

/*
getSQLStatementForFieldValues return a SQL query that will retrieve all field values to be indexed.
If the version ID is empty, the SQL returned will retrieve field values for all matching items.
Note that since the constructed queries always look up items in a view containing only the items representing
the most recent versions of the corresponding data, passing an item ID that corresponds to an older version of
the data will cause the result to be empty.
*/
func (svc *Service) getSQLStatementForFieldValues(ctx context.Context, newestVersionItemID string) (string, error) {
	focalEntityNames, err := svc.EntityRepo.GetEntityTypeNames(ctx, true)
	if err != nil {
		return "", err
	}
	focalEntityNames, err = erepo.QuoteAndSanitize(focalEntityNames)
	if err != nil {
		return "", E.MakeGRPCStatus(codes.InvalidArgument, "configured entity name malformed", E.Cause(err), E.DevMessagef("configured entity name malformed: '%s'", err.Error())).Err()
	}
	focalEntityNames = append(focalEntityNames, "':dummy:'")

	/*
		The final SQL statement must return
		1. the values for the normal configured fields for each (focal) item,
		2. the values for the pre-defined fields of an item (.e.g creation time or entity type)
		3. the values for the linked fields.

		We therefore construct the final query as a UNION of several different sub-queries, all of which returns
		results of the same row format: (item ID, field_name, field_value, place, revision, language)
		Each sub-query extracts a particular class of field values, restricting to fields on focal entities and to
		the items representing the most recent version of the corresponding data (and to a single item ID, if given).
		Finally, the combined (UNIONed) data is sorted so that field values from the same item (and within that, from
		the same field) are grouped together.
		This ensures that we can complete one item at the time as we scan through the output rows.
	*/

	// Separate SELECTs returning item values stored directly in the items table (predefined item fields)
	var idValues string
	var creationTimeValue string
	var entityNameValues string
	var businessIDValues string
	var configuredFieldValues string

	// Item ID given --> collect fields from the specified item
	idValues = sqlExpForItemFields(solr.DefaultUniqueKey, "item_id", focalEntityNames, newestVersionItemID)
	creationTimeValue = sqlExpForItemFields(solr.ItemCreatedAtField, "TO_CHAR(created_at, 'YYYY-MM-DD\"T\"HH24:MI:SSZ')", focalEntityNames, newestVersionItemID)
	entityNameValues = sqlExpForItemFields(solr.ItemEntityNameField, "entity_name", focalEntityNames, newestVersionItemID)
	businessIDValues = sqlExpForItemFields(solr.ItemBusinessIDField, "business_id", focalEntityNames, newestVersionItemID)

	fieldValueSelects := []string{
		idValues,
		creationTimeValue,
		entityNameValues,
		businessIDValues,
	}

	// Query returning normal (non-linked) field values
	configuredFieldNames, err := svc.FieldRepo.GetFieldDefNames(ctx)
	if err != nil {
		return "", err
	}
	configuredFieldNames, err = fieldUtils.QuoteAndSanitize(configuredFieldNames)
	if err != nil {
		return "", E.MakeGRPCStatus(codes.InvalidArgument, "configured field name malformed", E.Cause(err)).Err()
	}
	if len(configuredFieldNames) > 0 {
		configuredFieldValues = sqlExprForNormalFields(configuredFieldNames, focalEntityNames, newestVersionItemID)
		fieldValueSelects = append(fieldValueSelects, configuredFieldValues)
	}

	// Query returning linked field values
	linkingFieldDefs, err := svc.FieldRepo.GetFieldDefsByKind(ctx, kindLink.KindName)
	if err != nil {
		return "", err
	}
	hierarchyFieldDefs, err := svc.FieldRepo.GetFieldDefsByKind(ctx, kindHierarchy.KindName)
	if err != nil {
		return "", err
	}
	linkingFieldDefs = append(linkingFieldDefs, hierarchyFieldDefs...)
	if len(linkingFieldDefs) > 0 {
		linkFieldNames := make([]string, len(linkingFieldDefs))
		for i, fd := range linkingFieldDefs {
			linkFieldNames[i] = fd.Name()
		}
		linkFieldNames, err = fieldUtils.QuoteAndSanitize(linkFieldNames)
		if err != nil {
			return "", E.MakeGRPCStatus(codes.InvalidArgument, "configured linked field name malformed", E.Cause(err)).Err()
		}
		linkedFieldValues := sqlExpForLinkedFields(linkFieldNames, focalEntityNames, configuredFieldNames, newestVersionItemID)
		fieldValueSelects = append(fieldValueSelects, linkedFieldValues)
	}

	return createUnionSQL(fieldValueSelects), nil
}

// createUnionSQL joins the individual SELECT clauses into a single UNION and sorts the complete result
func createUnionSQL(fieldValueSelects []string) string {
	combinedFieldValues := strings.Join(fieldValueSelects, " UNION ")

	// Do a global sort on result to group (field, value) pairs by item and, within that, by field
	finalSQLQuery := fmt.Sprintf(`SELECT un.* FROM (%s) un ORDER BY un.item_id ASC, un.field_name ASC, un.place ASC`, combinedFieldValues)

	return finalSQLQuery
}

// sqlExpForItemFields builds the SELECT clause to retrieve an item property stored directly in the items table
func sqlExpForItemFields(fieldName string, fieldValueColumnExp string, focalEntityNames []string, newestVersionItemID string) string {
	focalEntityNamesClause := strings.Join(focalEntityNames, ",")
	itemIDConstraint := ""
	if newestVersionItemID != "" {
		itemIDConstraint = fmt.Sprintf(` AND item_id = '%s'`, newestVersionItemID)
	}

	return fmt.Sprintf(`(
	SELECT item_id AS item_id, '%s' AS field_name, %s AS field_value, 1 AS place, 1 AS revision, NULL AS language
	FROM latest_items_with_business_id
	WHERE entity_name IN (%s)%s
)`, fieldName, fieldValueColumnExp, focalEntityNamesClause, itemIDConstraint)
}

// sqlExprForNormalFields builds the SELECT clause for normal (user-configure) fields
func sqlExprForNormalFields(configuredFieldNames []string, focalEntityNames []string, newestVersionItemID string) string {
	configuredFieldNamesClause := strings.Join(configuredFieldNames, ",")
	focalEntityNamesClause := strings.Join(focalEntityNames, ",")
	itemIDConstraint := ""
	if newestVersionItemID != "" {
		itemIDConstraint = fmt.Sprintf(` AND civ.item_id = '%s'`, newestVersionItemID)
	}

	return fmt.Sprintf(`(
	SELECT civ.item_id, civ.field_name, civ.field_value, civ.place, civ.revision, civ.language
	FROM current_item_values civ
	JOIN latest_items_with_business_id liwbi
	ON liwbi.item_id = civ.item_id AND liwbi.entity_name IN (%s) AND civ.field_name IN (%s)%s
)`, focalEntityNamesClause, configuredFieldNamesClause, itemIDConstraint)
}

// sqlExpForLinkedFields builds the SELECT clause for linked fields
func sqlExpForLinkedFields(linkFieldNames []string, focalEntityNames []string, configuredFieldNames []string, newestVersionItemID string) string {
	linkFieldNamesClause := strings.Join(linkFieldNames, ",")
	itemIDConstraint := ""
	if newestVersionItemID != "" {
		itemIDConstraint = fmt.Sprintf(` AND civ_source.item_id = '%s'`, newestVersionItemID)
	}
	focalEntityNamesClause := strings.Join(focalEntityNames, ",")
	configuredFieldNamesClause := strings.Join(configuredFieldNames, ",")

	/*
		The constructed SQL does the following:
		1. Get the link fields in the newest versions of each focal item with a business ID (possibly restricted to a fixed item ID).
		2. For each such link field, find the most recent version of the item it links to (need not be focal).
		3. Get the values for all fields on this linked item.
		4. Build all possible linked fields with the attendant values.
		5. Restrict to the linked fields that are actually configured.
	*/
	return fmt.Sprintf(`(
	SELECT links.* FROM (
		SELECT civ_source.item_id as item_id, civ_source.field_name || '%s' || civ_target.field_name as field_name, civ_target.field_value as field_value,
			civ_target.place as place, civ_target.revision as revision, civ_target.language as language
		FROM current_item_values civ_source
		JOIN latest_items_with_business_id liwbi_source
		ON liwbi_source.item_id = civ_source.item_id AND liwbi_source.entity_name IN (%s) AND civ_source.field_name IN (%s)%s
		JOIN latest_items_with_business_id liwbi_target
		ON liwbi_target.business_id = civ_source.field_value
		JOIN current_item_values civ_target
		ON civ_target.item_id = liwbi_target.item_id
	) links
	WHERE links.field_name IN (%s)
)`, solr.LinkedFieldSeparator, focalEntityNamesClause, linkFieldNamesClause, itemIDConstraint, configuredFieldNamesClause)
}
