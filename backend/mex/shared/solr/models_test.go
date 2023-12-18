package solr

import (
	"encoding/json"
	"testing"
)

func TestSolrFacetSet_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		facetSet SolrFacetSet
		want     string
		wantErr  bool
	}{
		{
			name: "causes error for unknown facet type",
			facetSet: SolrFacetSet{
				"facet1": {
					DetailedType: "notAType",
					Field:        "keyword",
					NumBuckets:   true,
					Limit:        10,
				},
			},
			wantErr: true,
		},
		{
			name: "correctly serializes relevant fields for terms facet",
			facetSet: SolrFacetSet{
				"facet1": {
					DetailedType: SolrTermsFacetType,
					Field:        "keyword",
					NumBuckets:   true,
					Limit:        10,
					Offset:       5,
					ExcludeTags:  []string{"tag1", "tag2"},
				},
			},
			want: `{"facet1":{"domain":{"excludeTags":["tag1","tag2"]},"field":"keyword","limit":10,"numBuckets":true,"offset":5,"type":"terms"}}`,
		},
		{
			name: "ignores irrelevant properties when serializing terms facet",
			facetSet: SolrFacetSet{
				"facet1": {
					DetailedType: SolrTermsFacetType,
					Field:        "keyword",
					NumBuckets:   true,
					Limit:        10,
					Offset:       5,
					StartString:  "1972-05-20T17:33:18.77Z",
					EndString:    "1974-05-20T17:33:18.77Z",
					GapString:    "+1YEARS",
					StatOp:       MinOperator,
				},
			},
			want: `{"facet1":{"field":"keyword","limit":10,"numBuckets":true,"offset":5,"type":"terms"}}`,
		},
		{
			name: "leaves out offset set to zero and numBuckets set to false when serializing terms facet",
			facetSet: SolrFacetSet{
				"facet1": {
					DetailedType: SolrTermsFacetType,
					Field:        "keyword",
					NumBuckets:   false,
					Limit:        10,
					StartString:  "1972-05-20T17:33:18.77Z",
					EndString:    "1974-05-20T17:33:18.77Z",
					GapString:    "+1YEARS",
				},
			},
			want: `{"facet1":{"field":"keyword","limit":10,"type":"terms"}}`,
		},
		{
			name: "correctly serializes relevant fields for range facet with dates",
			facetSet: SolrFacetSet{
				"facet1": {
					DetailedType: SolrStringRangeFacetType,
					Field:        "createdAt",
					StartString:  "1972-05-20T17:33:18.77Z",
					EndString:    "1974-05-20T17:33:18.77Z",
					GapString:    "+1YEARS",
					ExcludeTags:  []string{"tag1", "tag2"},
				},
			},
			want: `{"facet1":{"domain":{"excludeTags":["tag1","tag2"]},"end":"1974-05-20T17:33:18.77Z","field":"createdAt","gap":"+1YEARS","start":"1972-05-20T17:33:18.77Z","type":"range"}}`,
		},
		{
			name: "ignores irrelevant properties when serializing year-binned range facet with dates",
			facetSet: SolrFacetSet{
				"facet1": {
					DetailedType: SolrStringRangeFacetType,
					Field:        "createdAt",
					NumBuckets:   true,
					Limit:        10,
					Offset:       5,
					StartString:  "1972-05-20T17:33:18.77Z",
					EndString:    "1974-05-20T17:33:18.77Z",
					GapString:    "+1YEARS",
					StatOp:       MinOperator,
				},
			},
			want: `{"facet1":{"end":"1974-05-20T17:33:18.77Z","field":"createdAt","gap":"+1YEARS","start":"1972-05-20T17:33:18.77Z","type":"range"}}`,
		},
		{
			name: "correctly serializes relevant fields for stat facet",
			facetSet: SolrFacetSet{
				"facet1": {
					DetailedType: SolrStringStatFacetType,
					Field:        "createdAt",
					StatOp:       MaxOperator,
				},
			},
			want: `{"facet1":"max(createdAt)"}`,
		},
		{
			name: "does NOT exclude tag for stat facet",
			facetSet: SolrFacetSet{
				"facet1": {
					DetailedType: SolrStringStatFacetType,
					Field:        "createdAt",
					StatOp:       MaxOperator,
					ExcludeTags:  []string{"tag1", "tag2"},
				},
			},
			want: `{"facet1":"max(createdAt)"}`,
		},
		{
			name: "ignores irrelevant properties when serializing year-binned range facet with dates",
			facetSet: SolrFacetSet{
				"facet1": {
					DetailedType: SolrStringStatFacetType,
					Field:        "createdAt",
					StatOp:       MaxOperator,
					NumBuckets:   true,
					Limit:        10,
					Offset:       5,
					StartString:  "1972-05-20T17:33:18.77Z",
					EndString:    "1974-05-20T17:33:18.77Z",
					GapString:    "+1YEARS",
				},
			},
			want: `{"facet1":"max(createdAt)"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.facetSet)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				// No point in checking returned value if we got an expected error
				return
			}
			stringGot := string(got)
			if stringGot != tt.want {
				t.Errorf("MarshalJSON() error - got = %s, want %s", stringGot, tt.want)
			}
		})
	}
}
