package items

import (
	"reflect"
	"testing"

	"github.com/d4l-data4life/mex/mex/shared/items"
	"github.com/d4l-data4life/mex/mex/shared/solr"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/datamodel"
)

func Test_aggregateValuesSimple(t *testing.T) {
	tests := []struct {
		name                       string
		aggregationCandidateValues []datamodel.DbSimpleAggregationByBusinessIdRow
		targetBusinessIdFieldName  string
		wantItemValues             []*items.ItemValue
		wantSourceIDs              map[string]struct{}
	}{
		{
			name: "Concatenates values in all fields from the different items",
			aggregationCandidateValues: []datamodel.DbSimpleAggregationByBusinessIdRow{
				{
					ItemID:     "item1",
					FieldName:  "field1",
					FieldValue: "a",
				},
				{
					ItemID:     "item1",
					FieldName:  "field2",
					FieldValue: "r",
				},
				{
					ItemID:     "item2",
					FieldName:  "field1",
					FieldValue: "b",
				},
				{
					ItemID:     "item2",
					FieldName:  "field2",
					FieldValue: "s",
				},
			},
			wantItemValues: []*items.ItemValue{
				{
					FieldName:  "field1",
					FieldValue: "a",
				},
				{
					FieldName:  "field2",
					FieldValue: "r",
				},
				{
					FieldName:  "field1",
					FieldValue: "b",
				},
				{
					FieldName:  "field2",
					FieldValue: "s",
				},
			},
			wantSourceIDs: map[string]struct{}{"item1": {}, "item2": {}},
		},
		{
			name: "Does not eliminate duplicated values within a single item",
			aggregationCandidateValues: []datamodel.DbSimpleAggregationByBusinessIdRow{
				{
					ItemID:     "item1",
					FieldName:  "field1",
					FieldValue: "a",
				},
				{
					ItemID:     "item1",
					FieldName:  "field1",
					FieldValue: "a",
				},
				{
					ItemID:     "item2",
					FieldName:  "field1",
					FieldValue: "b",
				},
			},
			wantItemValues: []*items.ItemValue{
				{
					FieldName:  "field1",
					FieldValue: "a",
				},
				{
					FieldName:  "field1",
					FieldValue: "a",
				},
				{
					FieldName:  "field1",
					FieldValue: "b",
				},
			},
			wantSourceIDs: map[string]struct{}{"item1": {}, "item2": {}},
		},
		{
			name: "Does not eliminate duplicated values across different items",
			aggregationCandidateValues: []datamodel.DbSimpleAggregationByBusinessIdRow{
				{
					ItemID:     "item1",
					FieldName:  "field1",
					FieldValue: "a",
				},
				{
					ItemID:     "item2",
					FieldName:  "field1",
					FieldValue: "a",
				},
			},
			wantItemValues: []*items.ItemValue{
				{
					FieldName:  "field1",
					FieldValue: "a",
				},
				{
					FieldName:  "field1",
					FieldValue: "a",
				},
			},
			wantSourceIDs: map[string]struct{}{"item1": {}, "item2": {}},
		},
		{
			name: "Skips the business ID field of the merge target",
			aggregationCandidateValues: []datamodel.DbSimpleAggregationByBusinessIdRow{
				{
					ItemID:     "item1",
					FieldName:  "targetIdField",
					FieldValue: "x",
				},
				{
					ItemID:     "item1",
					FieldName:  "field1",
					FieldValue: "a",
				},
				{
					ItemID:     "item2",
					FieldName:  "field1",
					FieldValue: "b",
				},
				{
					ItemID:     "item2",
					FieldName:  "targetIdField",
					FieldValue: "y",
				},
			},
			wantItemValues: []*items.ItemValue{
				{
					FieldName:  "field1",
					FieldValue: "a",
				},
				{
					FieldName:  "field1",
					FieldValue: "b",
				},
			},
			targetBusinessIdFieldName: "targetIdField",
			wantSourceIDs:             map[string]struct{}{"item1": {}, "item2": {}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotItemValues, gotSourceIDs := aggregateValuesSimple(tt.aggregationCandidateValues, tt.targetBusinessIdFieldName)
			if !reflect.DeepEqual(gotItemValues, tt.wantItemValues) {
				t.Errorf("aggregateValuesSimple() got = %v, want %v", gotItemValues, tt.wantItemValues)
			}
			if !reflect.DeepEqual(gotSourceIDs, tt.wantSourceIDs) {
				t.Errorf("aggregateValuesSimple() got2 = %v, want %v", gotSourceIDs, tt.wantSourceIDs)
			}
		})
	}
}

func Test_aggregateValuesFlexible(t *testing.T) {
	tests := []struct {
		name                       string
		aggregationCandidateValues []datamodel.DbAggregationCandidateValuesRow
		targetBusinessIdFieldName  string
		duplicateStrategy          string
		wantItemValues             []*items.ItemValue
		wantSourceIDs              map[string]struct{}
		wantErr                    bool
	}{
		{
			name: "keepall strategy: Concatenates values in all fields from the different items/sources",
			aggregationCandidateValues: []datamodel.DbAggregationCandidateValuesRow{
				{
					FieldName:  "field1",
					ItemID:     "item1",
					FieldValue: "a",
				},
				{
					FieldName:  "field1",
					ItemID:     "item1",
					FieldValue: "b",
				},
				{
					FieldName:           "field1",
					ItemID:              "item3",
					PartitionBusinessID: "source2",
					FieldValue:          "c",
				},
				{
					FieldName:           "field2",
					ItemID:              "item1",
					PartitionBusinessID: "source1",
					FieldValue:          "r",
				},
				{
					FieldName:           "field2",
					ItemID:              "item2",
					PartitionBusinessID: "source2",
					FieldValue:          "s",
				},
				{
					FieldName:           "field2",
					ItemID:              "item3",
					PartitionBusinessID: "source3",
					FieldValue:          "t",
				},
			},
			duplicateStrategy: solr.KeepAllDuplicates,
			wantItemValues: []*items.ItemValue{
				{
					FieldName:  "field1",
					FieldValue: "a",
				},
				{
					FieldName:  "field1",
					FieldValue: "b",
				},
				{
					FieldName:  "field1",
					FieldValue: "c",
				},
				{
					FieldName:  "field2",
					FieldValue: "r",
				},
				{
					FieldName:  "field2",
					FieldValue: "s",
				},
				{
					FieldName:  "field2",
					FieldValue: "t",
				},
			},
			wantSourceIDs: map[string]struct{}{"item1": {}, "item2": {}, "item3": {}},
		},
		{
			name: "keepall strategy: Does not eliminate duplicated values within a single item",
			aggregationCandidateValues: []datamodel.DbAggregationCandidateValuesRow{
				{
					FieldName:           "field1",
					ItemID:              "item1",
					PartitionBusinessID: "source1",
					FieldValue:          "a",
				},
				{
					FieldName:           "field1",
					ItemID:              "item1",
					PartitionBusinessID: "source1",
					FieldValue:          "a",
				},
				{
					FieldName:           "field2",
					ItemID:              "item2",
					PartitionBusinessID: "source2",
					FieldValue:          "r",
				},
				{
					FieldName:           "field2",
					ItemID:              "item3",
					PartitionBusinessID: "source3",
					FieldValue:          "s",
				},
			},
			duplicateStrategy: solr.KeepAllDuplicates,
			wantItemValues: []*items.ItemValue{
				{
					FieldName:  "field1",
					FieldValue: "a",
				},
				{
					FieldName:  "field1",
					FieldValue: "a",
				},
				{
					FieldName:  "field2",
					FieldValue: "r",
				},
				{
					FieldName:  "field2",
					FieldValue: "s",
				},
			},
			wantSourceIDs: map[string]struct{}{"item1": {}, "item2": {}, "item3": {}},
		},
		{
			name: "keepall strategy: Does not eliminate duplicated values across different items from the same source",
			aggregationCandidateValues: []datamodel.DbAggregationCandidateValuesRow{
				{
					FieldName:           "field1",
					ItemID:              "item1",
					PartitionBusinessID: "source1",
					FieldValue:          "a",
				},
				{
					FieldName:           "field1",
					ItemID:              "item2",
					PartitionBusinessID: "source1",
					FieldValue:          "a",
				},
				{
					FieldName:           "field2",
					ItemID:              "item3",
					PartitionBusinessID: "source3",
					FieldValue:          "r",
				},
			},
			duplicateStrategy: solr.KeepAllDuplicates,
			wantItemValues: []*items.ItemValue{
				{
					FieldName:  "field1",
					FieldValue: "a",
				},
				{
					FieldName:  "field1",
					FieldValue: "a",
				},
				{
					FieldName:  "field2",
					FieldValue: "r",
				},
			},
			wantSourceIDs: map[string]struct{}{"item1": {}, "item2": {}, "item3": {}},
		},
		{
			name: "keepall strategy: Does not eliminate duplicated values across different items from different sources",
			aggregationCandidateValues: []datamodel.DbAggregationCandidateValuesRow{
				{
					FieldName:           "field1",
					ItemID:              "item1",
					PartitionBusinessID: "source1",
					FieldValue:          "a",
				},
				{
					FieldName:           "field1",
					ItemID:              "item2",
					PartitionBusinessID: "source2",
					FieldValue:          "a",
				},
				{
					FieldName:           "field2",
					ItemID:              "item2",
					PartitionBusinessID: "source2",
					FieldValue:          "r",
				},
				{
					FieldName:           "field2",
					ItemID:              "item3",
					PartitionBusinessID: "source3",
					FieldValue:          "s",
				},
			},
			duplicateStrategy: solr.KeepAllDuplicates,
			wantItemValues: []*items.ItemValue{
				{
					FieldName:  "field1",
					FieldValue: "a",
				},
				{
					FieldName:  "field1",
					FieldValue: "a",
				},
				{
					FieldName:  "field2",
					FieldValue: "r",
				},
				{
					FieldName:  "field2",
					FieldValue: "s",
				},
			},
			wantSourceIDs: map[string]struct{}{"item1": {}, "item2": {}, "item3": {}},
		},
		{
			name: "keepall strategy: Skips the business ID field of the merge target",
			aggregationCandidateValues: []datamodel.DbAggregationCandidateValuesRow{
				{
					FieldName:           "field1",
					ItemID:              "item1",
					PartitionBusinessID: "source1",
					FieldValue:          "a",
				},
				{
					FieldName:           "targetIdField",
					ItemID:              "item1",
					PartitionBusinessID: "source1",
					FieldValue:          "x",
				},
				{
					FieldName:           "field1",
					ItemID:              "item2",
					PartitionBusinessID: "source2",
					FieldValue:          "b",
				},
				{
					FieldName:           "targetIdField",
					ItemID:              "item2",
					PartitionBusinessID: "source2",
					FieldValue:          "x",
				},
				{
					FieldName:           "field2",
					ItemID:              "item2",
					PartitionBusinessID: "source2",
					FieldValue:          "r",
				},
			},
			duplicateStrategy: solr.KeepAllDuplicates,
			wantItemValues: []*items.ItemValue{
				{
					FieldName:  "field1",
					FieldValue: "a",
				},
				{
					FieldName:  "field1",
					FieldValue: "b",
				},
				{
					FieldName:  "field2",
					FieldValue: "r",
				},
			},
			targetBusinessIdFieldName: "targetIdField",
			wantSourceIDs:             map[string]struct{}{"item1": {}, "item2": {}},
		},
		{
			name: "keepall strategy: If the same item occurs in multiple partitions (e.g. raw item coming from multiple source systems), the values are only added once",
			aggregationCandidateValues: []datamodel.DbAggregationCandidateValuesRow{
				{
					FieldName:           "field1",
					ItemID:              "item1",
					PartitionBusinessID: "source1",
					FieldValue:          "a",
				},
				{
					FieldName:           "field1",
					ItemID:              "item1",
					PartitionBusinessID: "source2",
					FieldValue:          "a",
				},
				{
					FieldName:           "field2",
					ItemID:              "item1",
					PartitionBusinessID: "source1",
					FieldValue:          "b",
				},
				{
					FieldName:           "field2",
					ItemID:              "item1",
					PartitionBusinessID: "source2",
					FieldValue:          "b",
				},
			},
			duplicateStrategy: solr.KeepAllDuplicates,
			wantItemValues: []*items.ItemValue{
				{
					FieldName:  "field1",
					FieldValue: "a",
				},
				{
					FieldName:  "field2",
					FieldValue: "b",
				},
			},
			wantSourceIDs: map[string]struct{}{"item1": {}},
		},
		{
			name: "removeall strategy: Concatenates values in all fields from the different items/sources",
			aggregationCandidateValues: []datamodel.DbAggregationCandidateValuesRow{
				{
					FieldName:           "field1",
					ItemID:              "item1",
					PartitionBusinessID: "source1",
					FieldValue:          "a",
				},
				{
					FieldName:           "field1",
					ItemID:              "item2",
					PartitionBusinessID: "source1",
					FieldValue:          "b",
				},
				{
					FieldName:           "field1",
					ItemID:              "item3",
					PartitionBusinessID: "source2",
					FieldValue:          "c",
				},
				{
					FieldName:           "field2",
					ItemID:              "item1",
					PartitionBusinessID: "source1",
					FieldValue:          "r",
				},
				{
					FieldName:           "field2",
					ItemID:              "item2",
					PartitionBusinessID: "source2",
					FieldValue:          "s",
				},
				{
					FieldName:           "field2",
					ItemID:              "item3",
					PartitionBusinessID: "source3",
					FieldValue:          "t",
				},
			},
			duplicateStrategy: solr.RemoveAllDuplicates,
			wantItemValues: []*items.ItemValue{
				{
					FieldName:  "field1",
					FieldValue: "a",
				},
				{
					FieldName:  "field1",
					FieldValue: "b",
				},
				{
					FieldName:  "field1",
					FieldValue: "c",
				},
				{
					FieldName:  "field2",
					FieldValue: "r",
				},
				{
					FieldName:  "field2",
					FieldValue: "s",
				},
				{
					FieldName:  "field2",
					FieldValue: "t",
				},
			},
			wantSourceIDs: map[string]struct{}{"item1": {}, "item2": {}, "item3": {}},
		},
		{
			name: "removeall strategy: eliminates duplicated values within a single item",
			aggregationCandidateValues: []datamodel.DbAggregationCandidateValuesRow{
				{
					FieldName:           "field1",
					PartitionBusinessID: "source1",
					FieldValue:          "a",
					ItemID:              "item1",
				},
				{
					FieldName:           "field1",
					PartitionBusinessID: "source1",
					FieldValue:          "a",
					ItemID:              "item1",
				},
				{
					FieldName:           "field2",
					PartitionBusinessID: "source2",
					FieldValue:          "r",
					ItemID:              "item2",
				},
				{
					FieldName:           "field2",
					PartitionBusinessID: "source2",
					FieldValue:          "s",
					ItemID:              "item3",
				},
			},
			duplicateStrategy: solr.RemoveAllDuplicates,
			wantItemValues: []*items.ItemValue{
				{
					FieldName:  "field1",
					FieldValue: "a",
				},
				{
					FieldName:  "field2",
					FieldValue: "r",
				},
				{
					FieldName:  "field2",
					FieldValue: "s",
				},
			},
			wantSourceIDs: map[string]struct{}{"item1": {}, "item2": {}, "item3": {}},
		},
		{
			name: "removeall strategy: eliminates duplicated values across different items from the same source",
			aggregationCandidateValues: []datamodel.DbAggregationCandidateValuesRow{
				{
					FieldName:           "field1",
					PartitionBusinessID: "source1",
					FieldValue:          "a",
					ItemID:              "item1",
				},
				{
					FieldName:           "field1",
					PartitionBusinessID: "source1",
					FieldValue:          "a",
					ItemID:              "item2",
				},
				{
					FieldName:           "field2",
					PartitionBusinessID: "source2",
					FieldValue:          "r",
					ItemID:              "item2",
				},
				{
					FieldName:           "field2",
					PartitionBusinessID: "source2",
					FieldValue:          "s",
					ItemID:              "item3",
				},
			},
			duplicateStrategy: solr.RemoveAllDuplicates,
			wantItemValues: []*items.ItemValue{
				{
					FieldName:  "field1",
					FieldValue: "a",
				},
				{
					FieldName:  "field2",
					FieldValue: "r",
				},
				{
					FieldName:  "field2",
					FieldValue: "s",
				},
			},
			wantSourceIDs: map[string]struct{}{"item1": {}, "item2": {}, "item3": {}},
		},
		{
			name: "removeall strategy: eliminates duplicated values across different items from different sources",
			aggregationCandidateValues: []datamodel.DbAggregationCandidateValuesRow{
				{
					FieldName:           "field1",
					PartitionBusinessID: "source1",
					FieldValue:          "a",
					ItemID:              "item1",
				},
				{
					FieldName:           "field1",
					PartitionBusinessID: "source2",
					FieldValue:          "a",
					ItemID:              "item2",
				},
				{
					FieldName:           "field2",
					PartitionBusinessID: "source2",
					FieldValue:          "r",
					ItemID:              "item2",
				},
				{
					FieldName:           "field2",
					PartitionBusinessID: "source2",
					FieldValue:          "s",
					ItemID:              "item3",
				},
			},
			duplicateStrategy: solr.RemoveAllDuplicates,
			wantItemValues: []*items.ItemValue{
				{
					FieldName:  "field1",
					FieldValue: "a",
				},
				{
					FieldName:  "field2",
					FieldValue: "r",
				},
				{
					FieldName:  "field2",
					FieldValue: "s",
				},
			},
			wantSourceIDs: map[string]struct{}{"item1": {}, "item2": {}, "item3": {}},
		},
		{
			name: "removeall strategy: Skips the business ID field of the merge target",
			aggregationCandidateValues: []datamodel.DbAggregationCandidateValuesRow{
				{
					FieldName:           "field1",
					PartitionBusinessID: "source1",
					FieldValue:          "a",
					ItemID:              "item1",
				},
				{
					FieldName:           "targetIdField",
					PartitionBusinessID: "source1",
					FieldValue:          "x",
					ItemID:              "item1",
				},
				{
					FieldName:           "field1",
					PartitionBusinessID: "source2",
					FieldValue:          "b",
					ItemID:              "item2",
				},
				{
					FieldName:           "targetIdField",
					PartitionBusinessID: "source2",
					FieldValue:          "x",
					ItemID:              "item2",
				},
				{
					FieldName:           "field2",
					PartitionBusinessID: "source2",
					FieldValue:          "r",
					ItemID:              "item2",
				},
			},
			duplicateStrategy: solr.RemoveAllDuplicates,
			wantItemValues: []*items.ItemValue{
				{
					FieldName:  "field1",
					FieldValue: "a",
				},
				{
					FieldName:  "field1",
					FieldValue: "b",
				},
				{
					FieldName:  "field2",
					FieldValue: "r",
				},
			},
			targetBusinessIdFieldName: "targetIdField",
			wantSourceIDs:             map[string]struct{}{"item1": {}, "item2": {}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotItemValues, gotSourceIDs, gotErr := doFlexibleMerge(tt.aggregationCandidateValues, tt.duplicateStrategy, tt.targetBusinessIdFieldName)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("mergeValuesFlexible() got error %v, wanted error = %v, ", gotErr, tt.wantErr)
			}
			if gotErr == nil {
				if !reflect.DeepEqual(gotItemValues, tt.wantItemValues) {
					t.Errorf("mergeValuesFlexible() gotItemValues = %v, want %v", gotItemValues, tt.wantItemValues)
				}
				if !reflect.DeepEqual(gotSourceIDs, tt.wantSourceIDs) {
					t.Errorf("mergeValuesFlexible() gotSourceIDs = %v, want %v", gotSourceIDs, tt.wantSourceIDs)
				}
			}
		})
	}
}
