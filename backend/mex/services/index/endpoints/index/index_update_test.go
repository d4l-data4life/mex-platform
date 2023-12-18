package index

import (
	"context"
	"reflect"
	"testing"

	sharedFields "github.com/d4l-data4life/mex/mex/shared/fields"
	"github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/solr"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/datamodel"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/frepo"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/hooks"
	kindstring "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/string"
	kindtext "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/text"
)

func TestIndexService_indexingLogic(t *testing.T) {

	fieldsRepo := frepo.NewMockedFieldRepo([]fields.BaseFieldDef{
		(&kindstring.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "category", Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
		(&kindtext.KindText{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "title", Kind: "text", IndexDef: &sharedFields.IndexDef{}}),
	})
	solrDataLoadHooks, _ := hooks.NewSolrDataLoadHooks(hooks.SolrDataLoadHooksConfig{})

	type counts struct {
		count          int
		batchCount     int
		rowFailCount   int
		docFailCount   int
		batchFailCount int
	}

	var tests = []struct {
		name                 string
		itemValues           []datamodel.CurrentItemValue
		useFailingSolrClient bool
		wantCounts           counts
		wantDocsUploaded     int
	}{
		{
			name: "All items are uploaded, even if there are not enough for a whole batch",
			itemValues: []datamodel.CurrentItemValue{
				{
					ID:         "1",
					ItemID:     "a",
					FieldName:  "category",
					FieldValue: "c1",
				},
				{
					ID:         "2",
					ItemID:     "a",
					FieldName:  "title",
					FieldValue: "First title",
				},
				{
					ID:         "3",
					ItemID:     "a",
					FieldName:  "title",
					FieldValue: "Second title",
				},
				{
					ID:         "4",
					ItemID:     "a",
					FieldName:  "id",
					FieldValue: "a",
				},
			},
			wantCounts: counts{
				count:          1,
				batchCount:     1,
				rowFailCount:   0,
				docFailCount:   0,
				batchFailCount: 0,
			},
			wantDocsUploaded: 1,
		},
		{
			name: "Multiple batches are created and all items uploaded if item count is above batch size",
			itemValues: []datamodel.CurrentItemValue{
				{
					ID:         "1",
					ItemID:     "a",
					FieldName:  "category",
					FieldValue: "c1",
				},
				{
					ID:         "2",
					ItemID:     "a",
					FieldName:  "title",
					FieldValue: "First title",
				},
				{
					ID:         "3",
					ItemID:     "a",
					FieldName:  "id",
					FieldValue: "a",
				},
				{
					ID:         "4",
					ItemID:     "b",
					FieldName:  "title",
					FieldValue: "Second title",
				},
				{
					ID:         "5",
					ItemID:     "b",
					FieldName:  "id",
					FieldValue: "b",
				},
				{
					ID:         "6",
					ItemID:     "c",
					FieldName:  "category",
					FieldValue: "c2",
				},
				{
					ID:         "7",
					ItemID:     "c",
					FieldName:  "id",
					FieldValue: "c",
				},
				{
					ID:         "8",
					ItemID:     "d",
					FieldName:  "category",
					FieldValue: "c1",
				},
				{
					ID:         "9",
					ItemID:     "d",
					FieldName:  "id",
					FieldValue: "d",
				},
				{
					ID:         "10",
					ItemID:     "e",
					FieldName:  "title",
					FieldValue: "Third title",
				},
				{
					ID:         "11",
					ItemID:     "e",
					FieldName:  "id",
					FieldValue: "e",
				},
			},
			wantCounts: counts{
				count:          5,
				batchCount:     3,
				rowFailCount:   0,
				docFailCount:   0,
				batchFailCount: 0,
			},
			wantDocsUploaded: 5,
		},
		{
			name: "If the Solr upload fails, the batch is skipped",
			itemValues: []datamodel.CurrentItemValue{
				{
					ID:         "1",
					ItemID:     "a",
					FieldName:  "category",
					FieldValue: "c1",
				},
				{
					ID:         "2",
					ItemID:     "a",
					FieldName:  "title",
					FieldValue: "First title",
				},
				{
					ID:         "3",
					ItemID:     "a",
					FieldName:  "title",
					FieldValue: "Second title",
				},
				{
					ID:         "4",
					ItemID:     "a",
					FieldName:  "id",
					FieldValue: "a",
				},
			},
			useFailingSolrClient: true,
			wantCounts: counts{
				count:          1,
				batchCount:     0,
				rowFailCount:   0,
				docFailCount:   0,
				batchFailCount: 1,
			},
			wantDocsUploaded: 0,
		},
		{
			name: "If an unknown fields causes XML creation to fail, the item is skipped",
			itemValues: []datamodel.CurrentItemValue{
				{
					ID:         "1",
					ItemID:     "a",
					FieldName:  "category",
					FieldValue: "c1",
				},
				{
					ID:         "2",
					ItemID:     "a",
					FieldName:  "title",
					FieldValue: "First title",
				},
				{
					ID:         "3",
					ItemID:     "a",
					FieldName:  "id",
					FieldValue: "a",
				},
				{
					ID:         "4",
					ItemID:     "b",
					FieldName:  "notAKnownField",
					FieldValue: "important",
				},
				{
					ID:         "5",
					ItemID:     "b",
					FieldName:  "id",
					FieldValue: "b",
				},
				{
					ID:         "6",
					ItemID:     "c",
					FieldName:  "title",
					FieldValue: "Second title",
				},
				{
					ID:         "7",
					ItemID:     "c",
					FieldName:  "id",
					FieldValue: "c",
				},
			},
			wantCounts: counts{
				count:          2,
				batchCount:     1,
				rowFailCount:   0,
				docFailCount:   1,
				batchFailCount: 0,
			},
			wantDocsUploaded: 2,
		},
	}
	for _, tt := range tests {
		var solrClient solr.MockClient
		if tt.useFailingSolrClient {
			solrClient = solr.NewMockClient(true, "fails", solr.ReturnVals{})
		} else {
			solrClient = solr.NewMockClient(false, "works", solr.ReturnVals{})
		}

		indexSvc := &Service{
			Log:               &log.NullLogger{},
			Solr:              &solrClient,
			FieldRepo:         fieldsRepo,
			SolrDataLoadHooks: solrDataLoadHooks,
		}
		t.Run(tt.name, func(t *testing.T) {
			testBatchSize := 2

			// To simulate iterating over DB rows, we iterate over items
			state := newIteratorState(testBatchSize)
			for _, val := range tt.itemValues {
				state = indexSvc.processItemValue(context.TODO(), val, state, testBatchSize)
			}
			// .. also triggering the final upload to match the full upload logic
			state = indexSvc.finishItem(context.TODO(), state, testBatchSize, true)

			finalCountState := counts{
				count:          state.count,
				batchCount:     state.batchCount,
				rowFailCount:   state.rowFailCount,
				docFailCount:   state.docFailCount,
				batchFailCount: state.batchFailCount,
			}
			if !reflect.DeepEqual(finalCountState, tt.wantCounts) {
				t.Errorf("processItemValue(): counts do not match, got = %v, want %v", finalCountState, tt.wantCounts)
			}
			if tt.wantDocsUploaded != solrClient.DocsUploaded {
				t.Errorf("processItemValue(): incorrect no. of docs uploaded, got = %v, want %v", solrClient.DocsUploaded, tt.wantDocsUploaded)
			}
		})
	}
}
