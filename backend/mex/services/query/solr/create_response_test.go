package solr

import (
	"context"
	"encoding/json"
	"reflect"
	"sort"
	"testing"

	"github.com/d4l-data4life/mex/mex/shared/testutils"

	sharedFields "github.com/d4l-data4life/mex/mex/shared/fields"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/searchconfig"
	"github.com/d4l-data4life/mex/mex/shared/searchconfig/screpo"
	"github.com/d4l-data4life/mex/mex/shared/solr"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/frepo"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/hooks"
	kind_number "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/number"
	kind_string "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/string"
	kind_text "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/text"
	kind_timestamp "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/timestamp"

	"github.com/d4l-data4life/mex/mex/services/query/endpoints/search/pb"
)

func Test_CreateResponse_basics(t *testing.T) {
	tests := []struct {
		name         string
		facets       []*solr.Facet
		diagnostics  *solr.Diagnostics
		solrResponse *solr.QueryResponse
		want         *pb.SearchResponse
		wantErr      bool
	}{
		{
			name:        "nil Solr response causes error",
			diagnostics: &solr.Diagnostics{},
			wantErr:     true,
		},
		{
			name:        "Empty Solr response is correctly mapped",
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{},
			},
			want: pb.NewSearchResponse(),
		},
		{
			name: "Diagnostic information is passed on",
			diagnostics: &solr.Diagnostics{
				ParsingSucceeded: false,
				ParsingErrors:    []string{"blip", "blup"},
			},
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{},
			},
			want: &pb.SearchResponse{
				Highlights: make([]*solr.Highlight, 0),
				Facets:     make([]*solr.FacetResult, 0),
				Items:      make([]*solr.DocItem, 0),
				Diagnostics: &solr.Diagnostics{
					ParsingSucceeded: false,
					ParsingErrors:    []string{"blip", "blup"},
				},
			},
		},
		{
			name:        "No. of matches is correctly mapped",
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFound: 6,
				},
			},
			want: &pb.SearchResponse{
				NumFound:   6,
				Facets:     make([]*solr.FacetResult, 0),
				Highlights: make([]*solr.Highlight, 0),
				Items:      make([]*solr.DocItem, 0),
				Diagnostics: &solr.Diagnostics{
					ParsingSucceeded: true,
				},
			},
		},
		{
			name:        "Start index is correctly mapped",
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					Start: 6,
				},
			},
			want: &pb.SearchResponse{
				Start:      6,
				Facets:     make([]*solr.FacetResult, 0),
				Highlights: make([]*solr.Highlight, 0),
				Items:      make([]*solr.DocItem, 0),
				Diagnostics: &solr.Diagnostics{
					ParsingSucceeded: true,
				},
			},
		},
		{
			name:        "NumFoundExact index is correctly mapped",
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
				},
			},
			want: &pb.SearchResponse{
				NumFoundExact: true,
				Facets:        make([]*solr.FacetResult, 0),
				Highlights:    make([]*solr.Highlight, 0),
				Items:         make([]*solr.DocItem, 0),
				Diagnostics: &solr.Diagnostics{
					ParsingSucceeded: true,
				},
			},
		},
		{
			name:        "MaxScore index is correctly mapped",
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					MaxScore: 3.14,
				},
			},
			want: &pb.SearchResponse{
				MaxScore:   3.14,
				Facets:     make([]*solr.FacetResult, 0),
				Highlights: make([]*solr.Highlight, 0),
				Items:      make([]*solr.DocItem, 0),
				Diagnostics: &solr.Diagnostics{
					ParsingSucceeded: true,
				},
			},
		},
	}

	log := &L.NullLogger{}
	postQueryHooks, _ := hooks.NewPostQueryHooks(hooks.PostQueryHooksConfig{})

	createResponseConstraintFieldsRepo := frepo.NewMockedFieldRepo([]fields.BaseFieldDef{
		(&kind_string.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "category", Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
		(&kind_string.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "type", Kind: "string",
			IndexDef: &sharedFields.IndexDef{}}),
		(&kind_string.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "system", Kind: "string",
			IndexDef: &sharedFields.IndexDef{}}),
		(&kind_text.KindText{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "title", Kind: "text", IndexDef: &sharedFields.IndexDef{}}),
		(&kind_timestamp.KindTimestamp{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "publishingDate",
			Kind: "timestamp", IndexDef: &sharedFields.IndexDef{}}),
	})
	createResponseConstraintSearchConfigRepo := screpo.NewMockSearchConfigRepo([]*searchconfig.SearchConfigObject{})

	for _, tt := range tests {
		opts := QueryEngineOptions{
			Log:              log,
			FieldRepo:        createResponseConstraintFieldsRepo,
			SearchConfigRepo: createResponseConstraintSearchConfigRepo,
			PostQueryHooks:   postQueryHooks,
		}
		qe, _ := newQueryEngine(&constantConverterNonPhrase, opts)
		t.Run(tt.name, func(t *testing.T) {
			got, err := qe.CreateResponse(context.TODO(), tt.solrResponse, tt.facets, tt.diagnostics)
			if (err != nil) != tt.wantErr {
				t.Errorf("createResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createResponse() got: %v but wanted: %v", got, tt.want)
			}
		})
	}
}

func Test_CreateResponse_exact_facets(t *testing.T) {
	tests := []struct {
		name         string
		facets       []*solr.Facet
		diagnostics  *solr.Diagnostics
		solrResponse *solr.QueryResponse
		want         *pb.SearchResponse
		wantErr      bool
	}{
		{
			name: "Exact (terms) facets: if nothing is returned for a facet, " +
				"the facet is returned with empty bucket array",
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			facets: []*solr.Facet{
				{
					Type: solr.MexExactFacetType,
					Axis: "category",
				},
			},
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      0,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Facets: map[string]interface{}{},
			},
			want: &pb.SearchResponse{
				NumFound:      0,
				NumFoundExact: true,
				Start:         0,
				MaxScore:      100,
				Items:         make([]*solr.DocItem, 0),
				Facets: []*solr.FacetResult{
					{
						Type:     solr.MexExactFacetType,
						Axis:     "category",
						BucketNo: 0,
						Buckets:  []*solr.FacetBucket{},
					},
				},
				Highlights: make([]*solr.Highlight, 0),
				Diagnostics: &solr.Diagnostics{
					ParsingSucceeded: true,
				},
			},
		},
		{
			name: "Exact (terms) facets: for a facet field with no exact-equivalent field, " +
				"the facet type, field, and total bucket count is correctly mapped",
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			facets: []*solr.Facet{
				{
					Type: solr.MexExactFacetType,
					Axis: "category",
				},
			},
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      14,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Facets: map[string]interface{}{
					"count": float64(14),
					createFacetName("category", ""): map[string]interface{}{
						"numBuckets": float64(2),
						"buckets": []interface{}{
							map[string]interface{}{
								"val":   "cheap",
								"count": float64(7),
							},
							map[string]interface{}{
								"val":   "expensive",
								"count": float64(7),
							},
						},
					},
				},
			},
			want: &pb.SearchResponse{
				NumFound:      14,
				NumFoundExact: true,
				Start:         0,
				MaxScore:      100,
				Items:         make([]*solr.DocItem, 0),
				Facets: []*solr.FacetResult{
					{
						Type:     solr.MexExactFacetType,
						Axis:     "category",
						BucketNo: 2,
						Buckets: []*solr.FacetBucket{
							{
								Value: "cheap",
								Count: 7,
							},
							{
								Value: "expensive",
								Count: 7,
							},
						},
					},
				},
				Highlights: make([]*solr.Highlight, 0),
				Diagnostics: &solr.Diagnostics{
					ParsingSucceeded: true,
				},
			},
		},
		{
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			name: "Exact (terms) facets: for a facet field WITH an exact-equivalent field, " +
				"the field and total bucket count is correctly mapped",
			facets: []*solr.Facet{
				{
					Type: solr.MexExactFacetType,
					Axis: "title",
				},
			},
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      14,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Facets: map[string]interface{}{
					"count": float64(14),
					createFacetName("title", ""): map[string]interface{}{
						"numBuckets": float64(2),
						"buckets": []interface{}{
							map[string]interface{}{
								"val":   "Important stuff!",
								"count": float64(7),
							},
							map[string]interface{}{
								"val":   "NewFieldMapper discovery!",
								"count": float64(7),
							},
						},
					},
				},
			},
			want: &pb.SearchResponse{
				NumFound:      14,
				NumFoundExact: true,
				Start:         0,
				MaxScore:      100,
				Items:         make([]*solr.DocItem, 0),
				Facets: []*solr.FacetResult{
					{
						Type:     solr.MexExactFacetType,
						Axis:     "title",
						BucketNo: 2,
						Buckets: []*solr.FacetBucket{
							{
								Value: "Important stuff!",
								Count: 7,
							},
							{
								Value: "NewFieldMapper discovery!",
								Count: 7,
							},
						},
					},
				},
				Highlights: make([]*solr.Highlight, 0),
				Diagnostics: &solr.Diagnostics{
					ParsingSucceeded: true,
				},
			},
		},
		{
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			name:        "Exact (terms) facets: error occurs if total bucket count is missing for an exact (terms) facet",
			facets: []*solr.Facet{
				{
					Type: solr.MexExactFacetType,
					Axis: "type",
				},
			},
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      14,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Facets: map[string]interface{}{
					"count": float64(14),
					createFacetName("type", ""): map[string]interface{}{
						"buckets": []interface{}{
							map[string]interface{}{
								"val":   "cheap",
								"count": float64(7),
							},
							map[string]interface{}{
								"val":   "expensive",
								"count": float64(7),
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			name:        "Exact (terms) facets: facets are returned in order requested",
			facets: []*solr.Facet{
				{
					Type: solr.MexExactFacetType,
					Axis: "type",
				},
				{
					Type: solr.MexExactFacetType,
					Axis: "category",
				},
				{
					Type: solr.MexExactFacetType,
					Axis: "system",
				},
			},
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      14,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Facets: map[string]interface{}{
					"count": float64(14),
					createFacetName("category", ""): map[string]interface{}{
						"numBuckets": float64(2),
						"buckets": []interface{}{
							map[string]interface{}{
								"val":   "tools",
								"count": float64(10),
							},
							map[string]interface{}{
								"val":   "plants",
								"count": float64(4),
							},
						},
					},
					createFacetName("system", ""): map[string]interface{}{
						"numBuckets": float64(2),
						"buckets": []interface{}{
							map[string]interface{}{
								"val":   "MacOS",
								"count": float64(7),
							},
							map[string]interface{}{
								"val":   "Linux",
								"count": float64(7),
							},
						},
					},
					createFacetName("type", ""): map[string]interface{}{
						"numBuckets": float64(2),
						"buckets": []interface{}{
							map[string]interface{}{
								"val":   "cheap",
								"count": float64(7),
							},
							map[string]interface{}{
								"val":   "expensive",
								"count": float64(7),
							},
						},
					},
				},
			},
			want: &pb.SearchResponse{
				NumFound:      14,
				NumFoundExact: true,
				Start:         0,
				MaxScore:      100,
				Items:         make([]*solr.DocItem, 0),
				Facets: []*solr.FacetResult{
					{
						Type:     solr.MexExactFacetType,
						Axis:     "type",
						BucketNo: 2,
						Buckets: []*solr.FacetBucket{
							{
								Value: "cheap",
								Count: 7,
							},
							{
								Value: "expensive",
								Count: 7,
							},
						},
					},
					{
						Type:     solr.MexExactFacetType,
						Axis:     "category",
						BucketNo: 2,
						Buckets: []*solr.FacetBucket{
							{
								Value: "tools",
								Count: 10,
							},
							{
								Value: "plants",
								Count: 4,
							},
						},
					},
					{
						Type:     solr.MexExactFacetType,
						Axis:     "system",
						BucketNo: 2,
						Buckets: []*solr.FacetBucket{
							{
								Value: "MacOS",
								Count: 7,
							},
							{
								Value: "Linux",
								Count: 7,
							},
						},
					},
				},
				Highlights: make([]*solr.Highlight, 0),
				Diagnostics: &solr.Diagnostics{
					ParsingSucceeded: true,
				},
			},
		},
	}

	log := &L.NullLogger{}
	postQueryHooks, _ := hooks.NewPostQueryHooks(hooks.PostQueryHooksConfig{})

	createResponseConstraintFieldsRepo := frepo.NewMockedFieldRepo([]fields.BaseFieldDef{
		(&kind_string.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "category", Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
		(&kind_string.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "type", Kind: "string",
			IndexDef: &sharedFields.IndexDef{}}),
		(&kind_string.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "system", Kind: "string",
			IndexDef: &sharedFields.IndexDef{}}),
		(&kind_text.KindText{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "title", Kind: "text", IndexDef: &sharedFields.IndexDef{}}),
		(&kind_timestamp.KindTimestamp{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "publishingDate",
			Kind: "timestamp", IndexDef: &sharedFields.IndexDef{}}),
	})

	createResponseConstraintSearchConfigRepo := screpo.NewMockSearchConfigRepo([]*searchconfig.SearchConfigObject{
		{
			Type:   solr.MexSearchFocusType,
			Name:   "default",
			Fields: []string{"category"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "category",
			Fields: []string{"category"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "title",
			Fields: []string{"title"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "type",
			Fields: []string{"type"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "system",
			Fields: []string{"system"},
		},
	})

	for _, tt := range tests {
		opts := QueryEngineOptions{
			Log:              log,
			FieldRepo:        createResponseConstraintFieldsRepo,
			SearchConfigRepo: createResponseConstraintSearchConfigRepo,
			PostQueryHooks:   postQueryHooks,
		}
		qe, _ := newQueryEngine(&constantConverterNonPhrase, opts)
		t.Run(tt.name, func(t *testing.T) {
			got, err := qe.CreateResponse(context.TODO(), tt.solrResponse, tt.facets, tt.diagnostics)
			if (err != nil) != tt.wantErr {
				t.Errorf("createResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createResponse() got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func Test_CreateResponse_year_range_facets(t *testing.T) {
	tests := []struct {
		name         string
		facets       []*solr.Facet
		diagnostics  *solr.Diagnostics
		solrResponse *solr.QueryResponse
		want         *pb.SearchResponse
		wantErr      bool
	}{
		{
			name:        "Year-range facets: if Solr return nothing, a facet with empty buckets is returned",
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			facets: []*solr.Facet{
				{
					Type: solr.MexYearRangeFacetType,
					Axis: "publishingDate",
				},
			},
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      0,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Facets: map[string]interface{}{
					"count": float64(0),
					createFacetName("publishingDate", ""): map[string]interface{}{
						"numBuckets": float64(0),
						"buckets":    []interface{}{},
					},
				},
			},
			want: &pb.SearchResponse{
				NumFound:      0,
				NumFoundExact: true,
				Start:         0,
				MaxScore:      100,
				Items:         make([]*solr.DocItem, 0),
				Facets: []*solr.FacetResult{
					{
						Type:     solr.MexYearRangeFacetType,
						Axis:     "publishingDate",
						BucketNo: 0,
						Buckets:  []*solr.FacetBucket{},
					},
				},
				Highlights: make([]*solr.Highlight, 0),
				Diagnostics: &solr.Diagnostics{
					ParsingSucceeded: true,
				},
			},
		},
		{
			name:        "Year-range facets: single facet is correctly handled",
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			facets: []*solr.Facet{
				{
					Type: solr.MexYearRangeFacetType,
					Axis: "publishingDate",
				},
			},
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      14,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Facets: map[string]interface{}{
					"count": float64(14),
					createFacetName("publishingDate", ""): map[string]interface{}{
						"numBuckets": float64(2),
						"buckets": []interface{}{
							map[string]interface{}{
								"val":   "2020-01-01T00:00:00Z",
								"count": float64(7),
							},
							map[string]interface{}{
								"val":   "2021-01-01T00:00:00Z",
								"count": float64(3),
							},
						},
					},
				},
			},
			want: &pb.SearchResponse{
				NumFound:      14,
				NumFoundExact: true,
				Start:         0,
				MaxScore:      100,
				Items:         make([]*solr.DocItem, 0),
				Facets: []*solr.FacetResult{
					{
						Type:     solr.MexYearRangeFacetType,
						Axis:     "publishingDate",
						BucketNo: 2,
						Buckets: []*solr.FacetBucket{
							{
								Value: "2020-01-01T00:00:00Z",
								Count: 7,
							},
							{
								Value: "2021-01-01T00:00:00Z",
								Count: 3,
							},
						},
					},
				},
				Highlights: make([]*solr.Highlight, 0),
				Diagnostics: &solr.Diagnostics{
					ParsingSucceeded: true,
				},
			},
		},
		{
			name:        "Year-range facets: multiple facets are correctly handled and returned in the order requested",
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			facets: []*solr.Facet{
				{
					Type: solr.MexYearRangeFacetType,
					Axis: "publishingDate",
				},
				{
					Type: solr.MexYearRangeFacetType,
					Axis: "createdAt",
				},
			},
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      14,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Facets: map[string]interface{}{
					"count": float64(14),
					createFacetName("createdAt", ""): map[string]interface{}{
						"numBuckets": float64(2),
						"buckets": []interface{}{
							map[string]interface{}{
								"val":   "2010-01-01T00:00:00Z",
								"count": float64(7),
							},
							map[string]interface{}{
								"val":   "2011-01-01T00:00:00Z",
								"count": float64(3),
							},
						},
					},
					createFacetName("publishingDate", ""): map[string]interface{}{
						"numBuckets": float64(2),
						"buckets": []interface{}{
							map[string]interface{}{
								"val":   "2020-01-01T00:00:00Z",
								"count": float64(5),
							},
							map[string]interface{}{
								"val":   "2021-01-01T00:00:00Z",
								"count": float64(6),
							},
						},
					},
				},
			},
			want: &pb.SearchResponse{
				NumFound:      14,
				NumFoundExact: true,
				Start:         0,
				MaxScore:      100,
				Items:         make([]*solr.DocItem, 0),
				Facets: []*solr.FacetResult{
					{
						Type:     solr.MexYearRangeFacetType,
						Axis:     "publishingDate",
						BucketNo: 2,
						Buckets: []*solr.FacetBucket{
							{
								Value: "2020-01-01T00:00:00Z",
								Count: 5,
							},
							{
								Value: "2021-01-01T00:00:00Z",
								Count: 6,
							},
						},
					},
					{
						Type:     solr.MexYearRangeFacetType,
						Axis:     "createdAt",
						BucketNo: 2,
						Buckets: []*solr.FacetBucket{
							{
								Value: "2010-01-01T00:00:00Z",
								Count: 7,
							},
							{
								Value: "2011-01-01T00:00:00Z",
								Count: 3,
							},
						},
					},
				},
				Highlights: make([]*solr.Highlight, 0),
				Diagnostics: &solr.Diagnostics{
					ParsingSucceeded: true,
				},
			},
		},
		{
			name:        "Year-range facets: no error occurs if total bucket count is not returned",
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			facets: []*solr.Facet{
				{
					Type: solr.MexYearRangeFacetType,
					Axis: "publishingDate",
				},
			},
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      14,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Facets: map[string]interface{}{
					createFacetName("publishingDate", ""): map[string]interface{}{
						"numBuckets": float64(2),
						"buckets": []interface{}{
							map[string]interface{}{
								"val":   "2020-01-01T00:00:00Z",
								"count": float64(7),
							},
							map[string]interface{}{
								"val":   "2021-01-01T00:00:00Z",
								"count": float64(3),
							},
						},
					},
				},
			},
			want: &pb.SearchResponse{
				NumFound:      14,
				NumFoundExact: true,
				Start:         0,
				MaxScore:      100,
				Items:         make([]*solr.DocItem, 0),
				Facets: []*solr.FacetResult{
					{
						Type:     solr.MexYearRangeFacetType,
						Axis:     "publishingDate",
						BucketNo: 2,
						Buckets: []*solr.FacetBucket{
							{
								Value: "2020-01-01T00:00:00Z",
								Count: 7,
							},
							{
								Value: "2021-01-01T00:00:00Z",
								Count: 3,
							},
						},
					},
				},
				Highlights: make([]*solr.Highlight, 0),
				Diagnostics: &solr.Diagnostics{
					ParsingSucceeded: true,
				},
			},
		},
	}

	log := &L.NullLogger{}
	postQueryHooks, _ := hooks.NewPostQueryHooks(hooks.PostQueryHooksConfig{})

	createResponseConstraintFieldsRepo := frepo.NewMockedFieldRepo([]fields.BaseFieldDef{
		(&kind_string.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "category", Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
		(&kind_string.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "type", Kind: "string",
			IndexDef: &sharedFields.IndexDef{}}),
		(&kind_string.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "system", Kind: "string",
			IndexDef: &sharedFields.IndexDef{}}),
		(&kind_text.KindText{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "title", Kind: "text", IndexDef: &sharedFields.IndexDef{}}),
		(&kind_timestamp.KindTimestamp{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "publishingDate",
			Kind: "timestamp", IndexDef: &sharedFields.IndexDef{}}),
	})

	createResponseConstraintSearchConfigRepo := screpo.NewMockSearchConfigRepo([]*searchconfig.SearchConfigObject{
		{
			Type:   solr.MexSearchFocusType,
			Name:   "default",
			Fields: []string{"category"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "publishingDate",
			Fields: []string{"publishingDate"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "createdAt",
			Fields: []string{"createdAt"},
		},
	})

	for _, tt := range tests {
		opts := QueryEngineOptions{
			Log:              log,
			FieldRepo:        createResponseConstraintFieldsRepo,
			SearchConfigRepo: createResponseConstraintSearchConfigRepo,
			PostQueryHooks:   postQueryHooks,
		}
		qe, _ := newQueryEngine(&constantConverterNonPhrase, opts)
		t.Run(tt.name, func(t *testing.T) {
			got, err := qe.CreateResponse(context.TODO(), tt.solrResponse, tt.facets, tt.diagnostics)
			if (err != nil) != tt.wantErr {
				t.Errorf("createResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createResponse() got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func Test_CreateResponse_string_stat_facets(t *testing.T) {
	tests := []struct {
		name         string
		facets       []*solr.Facet
		diagnostics  *solr.Diagnostics
		solrResponse *solr.QueryResponse
		want         *pb.SearchResponse
		wantErr      bool
	}{
		{
			name:        "String stat facets: single facet correctly returned",
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			facets: []*solr.Facet{
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDate",
					StatName: "min_facet",
					StatOp:   solr.MinOperator,
				},
			},
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      14,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Facets: map[string]interface{}{
					createFacetName("publishingDate", "min_facet"): "2022-05-20T17:33:18Z",
				},
			},
			want: &pb.SearchResponse{
				NumFound:      14,
				NumFoundExact: true,
				Start:         0,
				MaxScore:      100,
				Items:         make([]*solr.DocItem, 0),
				Facets: []*solr.FacetResult{
					{
						Type:             solr.MexStringStatFacetType,
						Axis:             "publishingDate",
						StatName:         "min_facet",
						StringStatResult: "2022-05-20T17:33:18Z",
					},
				},
				Highlights: make([]*solr.Highlight, 0),
				Diagnostics: &solr.Diagnostics{
					ParsingSucceeded: true,
				},
			},
		},
		{
			name:        "String stat facets: is Solr returns nothing for a facet, an empty string facet is returned",
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			facets: []*solr.Facet{
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDate",
					StatName: "min_facet",
					StatOp:   solr.MinOperator,
				},
			},
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      0,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Facets: map[string]interface{}{},
			},
			want: &pb.SearchResponse{
				NumFound:      0,
				NumFoundExact: true,
				Start:         0,
				MaxScore:      100,
				Items:         make([]*solr.DocItem, 0),
				Facets: []*solr.FacetResult{
					{
						Type:             solr.MexStringStatFacetType,
						Axis:             "publishingDate",
						StatName:         "min_facet",
						StringStatResult: "",
					},
				},
				Highlights: make([]*solr.Highlight, 0),
				Diagnostics: &solr.Diagnostics{
					ParsingSucceeded: true,
				},
			},
		},
		{
			name: "String stat facets: multiple facets on the same field correctly returned, " +
				"in order requested",
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			facets: []*solr.Facet{
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDate",
					StatName: "min_PA_facet",
					StatOp:   solr.MinOperator,
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDate",
					StatName: "max_PA_facet",
					StatOp:   solr.MaxOperator,
				},
			},
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      14,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Facets: map[string]interface{}{
					createFacetName("publishingDate", "max_PA_facet"): "2022-03-20T17:33:18Z",
					createFacetName("publishingDate", "min_PA_facet"): "2020-05-20T17:33:18Z",
				},
			},
			want: &pb.SearchResponse{
				NumFound:      14,
				NumFoundExact: true,
				Start:         0,
				MaxScore:      100,
				Items:         make([]*solr.DocItem, 0),
				Facets: []*solr.FacetResult{
					{
						Type:             solr.MexStringStatFacetType,
						Axis:             "publishingDate",
						StatName:         "min_PA_facet",
						StringStatResult: "2020-05-20T17:33:18Z",
					},
					{
						Type:             solr.MexStringStatFacetType,
						Axis:             "publishingDate",
						StatName:         "max_PA_facet",
						StringStatResult: "2022-03-20T17:33:18Z",
					},
				},
				Highlights: make([]*solr.Highlight, 0),
				Diagnostics: &solr.Diagnostics{
					ParsingSucceeded: true,
				},
			},
		},
		{
			name: "String stat facets: multiple facets on different fields correctly returned, " +
				"in order requested",
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			facets: []*solr.Facet{
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDate",
					StatName: "min_PA_facet",
					StatOp:   solr.MinOperator,
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "createdAt",
					StatName: "max_CA_facet",
					StatOp:   solr.MaxOperator,
				},
			},
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      14,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Facets: map[string]interface{}{
					createFacetName("createdAt", "max_CA_facet"):      "2022-03-20T17:33:18Z",
					createFacetName("publishingDate", "min_PA_facet"): "2020-05-20T17:33:18Z",
				},
			},
			want: &pb.SearchResponse{
				NumFound:      14,
				NumFoundExact: true,
				Start:         0,
				MaxScore:      100,
				Items:         make([]*solr.DocItem, 0),
				Facets: []*solr.FacetResult{
					{
						Type:             solr.MexStringStatFacetType,
						Axis:             "publishingDate",
						StatName:         "min_PA_facet",
						StringStatResult: "2020-05-20T17:33:18Z",
					},
					{
						Type:             solr.MexStringStatFacetType,
						Axis:             "createdAt",
						StatName:         "max_CA_facet",
						StringStatResult: "2022-03-20T17:33:18Z",
					},
				},
				Highlights: make([]*solr.Highlight, 0),
				Diagnostics: &solr.Diagnostics{
					ParsingSucceeded: true,
				},
			},
		},
	}

	log := &L.NullLogger{}
	postQueryHooks, _ := hooks.NewPostQueryHooks(hooks.PostQueryHooksConfig{})

	createResponseConstraintFieldsRepo := frepo.NewMockedFieldRepo([]fields.BaseFieldDef{
		(&kind_string.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "category", Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
		(&kind_string.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "type", Kind: "string",
			IndexDef: &sharedFields.IndexDef{}}),
		(&kind_string.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "system", Kind: "string",
			IndexDef: &sharedFields.IndexDef{}}),
		(&kind_text.KindText{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "title", Kind: "text", IndexDef: &sharedFields.IndexDef{}}),
		(&kind_timestamp.KindTimestamp{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "publishingDate",
			Kind: "timestamp", IndexDef: &sharedFields.IndexDef{}}),
	})

	createResponseConstraintSearchConfigRepo := screpo.NewMockSearchConfigRepo([]*searchconfig.SearchConfigObject{})

	for _, tt := range tests {
		opts := QueryEngineOptions{
			Log:              log,
			FieldRepo:        createResponseConstraintFieldsRepo,
			SearchConfigRepo: createResponseConstraintSearchConfigRepo,
			PostQueryHooks:   postQueryHooks,
		}
		qe, _ := newQueryEngine(&constantConverterNonPhrase, opts)
		t.Run(tt.name, func(t *testing.T) {
			got, err := qe.CreateResponse(context.TODO(), tt.solrResponse, tt.facets, tt.diagnostics)
			if (err != nil) != tt.wantErr {
				t.Errorf("createResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createResponse() got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func Test_CreateResponse_mixed_facets(t *testing.T) {

	tests := []struct {
		name         string
		facets       []*solr.Facet
		diagnostics  *solr.Diagnostics
		solrResponse *solr.QueryResponse
		want         *pb.SearchResponse
		wantErr      bool
	}{
		{
			name:        "Mixed facet types: exact, year range, and stat facets can be used together",
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			facets: []*solr.Facet{
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDate",
					StatName: "min_PA_facet",
					StatOp:   solr.MinOperator,
				},
				{
					Type: solr.MexExactFacetType,
					Axis: "category",
				},
				{
					Type: solr.MexYearRangeFacetType,
					Axis: "publishingDate",
				},
			},
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      14,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Facets: map[string]interface{}{
					"count": float64(14),
					createFacetName("category", ""): map[string]interface{}{
						"numBuckets": float64(2),
						"buckets": []interface{}{
							map[string]interface{}{
								"val":   "tools",
								"count": float64(10),
							},
							map[string]interface{}{
								"val":   "plants",
								"count": float64(4),
							},
						},
					},
					createFacetName("publishingDate", "min_PA_facet"): "2020-05-20T17:33:18Z",
					createFacetName("publishingDate", ""): map[string]interface{}{
						"numBuckets": float64(2),
						"buckets": []interface{}{
							map[string]interface{}{
								"val":   "2020-01-01T00:00:00Z",
								"count": float64(7),
							},
							map[string]interface{}{
								"val":   "2021-01-01T00:00:00Z",
								"count": float64(3),
							},
						},
					},
				},
			},
			want: &pb.SearchResponse{
				NumFound:      14,
				NumFoundExact: true,
				Start:         0,
				MaxScore:      100,
				Items:         make([]*solr.DocItem, 0),
				Facets: []*solr.FacetResult{
					{
						Type:             solr.MexStringStatFacetType,
						Axis:             "publishingDate",
						StatName:         "min_PA_facet",
						StringStatResult: "2020-05-20T17:33:18Z",
					},
					{
						Type:     solr.MexExactFacetType,
						Axis:     "category",
						BucketNo: 2,
						Buckets: []*solr.FacetBucket{
							{
								Value: "tools",
								Count: 10,
							},
							{
								Value: "plants",
								Count: 4,
							},
						},
					},
					{
						Type:     solr.MexYearRangeFacetType,
						Axis:     "publishingDate",
						BucketNo: 2,
						Buckets: []*solr.FacetBucket{
							{
								Value: "2020-01-01T00:00:00Z",
								Count: 7,
							},
							{
								Value: "2021-01-01T00:00:00Z",
								Count: 3,
							},
						},
					},
				},
				Highlights: make([]*solr.Highlight, 0),
				Diagnostics: &solr.Diagnostics{
					ParsingSucceeded: true,
				},
			},
		},
	}

	log := &L.NullLogger{}
	postQueryHooks, _ := hooks.NewPostQueryHooks(hooks.PostQueryHooksConfig{})

	createResponseConstraintFieldsRepo := frepo.NewMockedFieldRepo([]fields.BaseFieldDef{
		(&kind_string.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "category", Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
		(&kind_string.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "type", Kind: "string",
			IndexDef: &sharedFields.IndexDef{}}),
		(&kind_string.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "system", Kind: "string",
			IndexDef: &sharedFields.IndexDef{}}),
		(&kind_text.KindText{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "title", Kind: "text", IndexDef: &sharedFields.IndexDef{}}),
		(&kind_timestamp.KindTimestamp{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "publishingDate",
			Kind: "timestamp", IndexDef: &sharedFields.IndexDef{}}),
	})

	createResponseConstraintSearchConfigRepo := screpo.NewMockSearchConfigRepo([]*searchconfig.SearchConfigObject{
		{
			Type:   solr.MexSearchFocusType,
			Name:   "default",
			Fields: []string{"category"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "publishingDate",
			Fields: []string{"publishingDate"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "category",
			Fields: []string{"category"},
		},
	})

	for _, tt := range tests {
		opts := QueryEngineOptions{
			Log:              log,
			FieldRepo:        createResponseConstraintFieldsRepo,
			SearchConfigRepo: createResponseConstraintSearchConfigRepo,
			PostQueryHooks:   postQueryHooks,
		}
		qe, _ := newQueryEngine(&constantConverterNonPhrase, opts)
		t.Run(tt.name, func(t *testing.T) {
			got, err := qe.CreateResponse(context.TODO(), tt.solrResponse, tt.facets, tt.diagnostics)
			v, _ := json.Marshal(got.Facets[0])
			if (err != nil) != tt.wantErr && v == nil {
				t.Errorf("createResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createResponse() got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func Test_CreateResponse_highlighting(t *testing.T) {

	genericTitleField, _ := solr.GetLangSpecificFieldName("title", solr.GenericLangAbbrev)
	deTitleField, _ := solr.GetLangSpecificFieldName("title", solr.GermanLangAbbrev)
	enTitleField, _ := solr.GetLangSpecificFieldName("title", solr.EnglishLangAbbrev)
	unanalyzedTitleField := solr.GetRawBackingFieldName("title")
	unanalyzedDescriptionField := solr.GetRawBackingFieldName("description")

	tests := []struct {
		name         string
		facets       []*solr.Facet
		diagnostics  *solr.Diagnostics
		solrResponse *solr.QueryResponse
		checks       []testutils.ResponseCheck
		wantErr      bool
	}{
		{
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			name:        "Highlights: empty highlight result is correctly handled",
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      14,
					Start:         0,
					MaxScore:      100,
				},
				Highlighting: map[string]interface{}{},
			},
			checks: []testutils.ResponseCheck{testutils.CheckHighlightResponse([]*solr.Highlight{})},
		},
		{
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			name:        "Highlights: entries for items are ordered by item ID",
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      14,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Highlighting: map[string]interface{}{
					"xyz789": map[string]interface{}{
						"category": []interface{}{
							"Again <em>smt</em> different",
						},
					},
					"def456": map[string]interface{}{
						"category": []interface{}{
							"<em>smt</em> look different this time",
						},
					},
					"abc123": map[string]interface{}{
						"category": []interface{}{
							"<em>smt</em> altogether different",
							"<em>smt</em> never seen before",
						},
					},
				},
			},
			checks: []testutils.ResponseCheck{testutils.CheckHighlightResponse([]*solr.Highlight{
				{
					ItemId: "abc123",
					Matches: []*solr.FieldHighlight{
						{
							FieldName: "category",
							Snippets: []string{
								"<em>smt</em> altogether different",
								"<em>smt</em> never seen before",
							},
						},
					},
				},
				{
					ItemId: "def456",
					Matches: []*solr.FieldHighlight{
						{
							FieldName: "category",
							Snippets: []string{
								"<em>smt</em> look different this time",
							},
						},
					},
				},
				{
					ItemId: "xyz789",
					Matches: []*solr.FieldHighlight{
						{
							FieldName: "category",
							Snippets: []string{
								"Again <em>smt</em> different",
							},
						},
					},
				},
			})},
		},
		{
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			name: "Highlights: highlight in language-specific field underlying text fields are mapped back to" +
				" the base field name with appropriate language tag",
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      14,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Highlighting: map[string]interface{}{
					"genericLangItemId": map[string]interface{}{
						genericTitleField: []interface{}{
							"Again <em>smt</em> different",
						},
					},
					"deItemId": map[string]interface{}{
						deTitleField: []interface{}{
							"<em>etwas</em> sieht anders aus",
						},
					},
					"enItemId": map[string]interface{}{
						enTitleField: []interface{}{
							"<em>smt</em> altogether different",
							"<em>smt</em> never seen before",
						},
					},
				},
			},
			checks: []testutils.ResponseCheck{testutils.CheckHighlightResponse([]*solr.Highlight{
				{
					ItemId: "deItemId",
					Matches: []*solr.FieldHighlight{
						{
							FieldName: "title",
							Language:  solr.GermanLangAbbrev,
							Snippets: []string{
								"<em>etwas</em> sieht anders aus",
							},
						},
					},
				},
				{
					ItemId: "enItemId",
					Matches: []*solr.FieldHighlight{
						{
							FieldName: "title",
							Language:  solr.EnglishLangAbbrev,
							Snippets: []string{
								"<em>smt</em> altogether different",
								"<em>smt</em> never seen before",
							},
						},
					},
				},
				{
					ItemId: "genericLangItemId",
					Matches: []*solr.FieldHighlight{
						{
							FieldName: "title",
							Snippets: []string{
								"Again <em>smt</em> different",
							},
						},
					},
				},
			})},
		},
		{
			name: "Highlights: highlights in the unanalyzed backing field are mapped back to" +
				" the base field name with no language tag",
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      14,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Highlighting: map[string]interface{}{
					"unanalyzedItemId": map[string]interface{}{
						unanalyzedTitleField: []interface{}{
							"Again <em>smt</em> different",
						},
					},
				},
			},
			checks: []testutils.ResponseCheck{testutils.CheckHighlightResponse([]*solr.Highlight{
				{
					ItemId: "unanalyzedItemId",
					Matches: []*solr.FieldHighlight{
						{
							FieldName: "title",
							Snippets: []string{
								"Again <em>smt</em> different",
							},
						},
					},
				},
			})},
		},
		{
			name:        "Highlights: if exactly the same snippet is returned by multiple fields mapping to the same MEx field, only one copy is returned (with language, if possible)",
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      14,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Highlighting: map[string]interface{}{
					"itemId": map[string]interface{}{
						genericTitleField: []interface{}{
							"I've <em>seen this</em> before",
						},
						deTitleField: []interface{}{
							"I've <em>seen this</em> before",
						},
						unanalyzedTitleField: []interface{}{
							"I've <em>seen this</em> before",
						},
					},
				},
			},
			checks: []testutils.ResponseCheck{testutils.CheckHighlightResponse([]*solr.Highlight{
				{
					ItemId: "itemId",
					Matches: []*solr.FieldHighlight{
						{
							FieldName: "title",
							Language:  solr.GermanLangAbbrev,
							Snippets: []string{
								"I've <em>seen this</em> before",
							},
						},
					},
				},
			})},
		},
		{
			name:        "Highlights: identical snippets in different fields are not removed",
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      14,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Highlighting: map[string]interface{}{
					"itemId": map[string]interface{}{
						genericTitleField: []interface{}{
							"I've <em>seen this</em> before",
						},
						enTitleField: []interface{}{
							"I've <em>seen this</em> before",
						},
						unanalyzedTitleField: []interface{}{
							"I've <em>seen this</em> before",
						},
						// Maps to different MEx field
						unanalyzedDescriptionField: []interface{}{
							"I've <em>seen this</em> before",
						},
					},
				},
			},
			checks: []testutils.ResponseCheck{testutils.CheckHighlightResponse([]*solr.Highlight{
				{
					ItemId: "itemId",
					Matches: []*solr.FieldHighlight{
						{
							FieldName: "title",
							Language:  solr.EnglishLangAbbrev,
							Snippets: []string{
								"I've <em>seen this</em> before",
							},
						},
						{
							FieldName: "description",
							Snippets: []string{
								"I've <em>seen this</em> before",
							},
						},
					},
				},
			})},
		},
		{
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			name:        "Highlights: empty list of snippets for items or fields are removed",
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      14,
					Start:         0,
					MaxScore:      100,
				},
				Highlighting: map[string]interface{}{
					"abc123": map[string]interface{}{
						"category": []interface{}{
							"contains a <em>highlight</em>",
						},
						"type": []interface{}{},
					},
					"def456": map[string]interface{}{},
				},
			},
			checks: []testutils.ResponseCheck{testutils.CheckHighlightResponse([]*solr.Highlight{
				{
					ItemId: "abc123",
					Matches: []*solr.FieldHighlight{
						{
							FieldName: "category",
							Snippets: []string{
								"contains a <em>highlight</em>",
							},
						},
					},
				},
			})},
		},
	}

	log := &L.NullLogger{}
	postQueryHooks, _ := hooks.NewPostQueryHooks(hooks.PostQueryHooksConfig{})

	createResponseConstraintFieldsRepo := frepo.NewMockedFieldRepo([]fields.BaseFieldDef{
		(&kind_text.KindText{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "title", Kind: "text",
			IndexDef: &sharedFields.IndexDef{MultiValued: true}}),
		(&kind_text.KindText{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "description", Kind: "text",
			IndexDef: &sharedFields.IndexDef{MultiValued: true}}),
		(&kind_string.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "category", Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
		(&kind_string.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "type", Kind: "string",
			IndexDef: &sharedFields.IndexDef{}}),
	})

	createResponseConstraintSearchConfigRepo := screpo.NewMockSearchConfigRepo([]*searchconfig.SearchConfigObject{})

	for _, tt := range tests {
		opts := QueryEngineOptions{
			Log:              log,
			FieldRepo:        createResponseConstraintFieldsRepo,
			SearchConfigRepo: createResponseConstraintSearchConfigRepo,
			PostQueryHooks:   postQueryHooks,
		}
		qe, _ := newQueryEngine(&constantConverterNonPhrase, opts)
		t.Run(tt.name, func(t *testing.T) {
			got, err := qe.CreateResponse(context.TODO(), tt.solrResponse, tt.facets, tt.diagnostics)
			if (err != nil) != tt.wantErr {
				t.Errorf("createResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, check := range tt.checks {
				check(got, t)
			}
		})
	}
}

func Test_CreateResponse_matches(t *testing.T) {

	genericTitleField, _ := solr.GetLangSpecificFieldName("title", "")
	deTitleField, _ := solr.GetLangSpecificFieldName("title", solr.GermanLangAbbrev)
	enTitleField, _ := solr.GetLangSpecificFieldName("title", solr.EnglishLangAbbrev)

	tests := []struct {
		name         string
		facets       []*solr.Facet
		diagnostics  *solr.Diagnostics
		solrResponse *solr.QueryResponse
		want         *pb.SearchResponse
		wantErr      bool
	}{
		{
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			name: "Matches: Docs array with only non-text fields is correctly mapped to items (" +
				"ID and entity type are outside value-array)",
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      2,
					Start:         0,
					MaxScore:      100,
					Docs: []solr.GenericObject{
						{
							"id":         "123",
							"entityName": "bookItem",
							"author":     "John Doe",
							"category":   "book",
							"price":      23.5,
						},
						{
							"id":         "foo",
							"entityName": "studyItem",
							"author":     "Jane Doe",
							"category":   "study",
							"price":      5,
						},
					},
				},
			},
			want: &pb.SearchResponse{
				NumFound:      2,
				NumFoundExact: true,
				Start:         0,
				MaxScore:      100,
				Items: []*solr.DocItem{
					{
						ItemId:     "123",
						EntityType: "bookItem",
						Values: []*solr.DocValue{
							{FieldName: "author", FieldValue: "John Doe"},
							{FieldName: "category", FieldValue: "book"},
							{FieldName: "price", FieldValue: "23.5"},
						},
					},
					{
						ItemId:     "foo",
						EntityType: "studyItem",
						Values: []*solr.DocValue{
							{FieldName: "author", FieldValue: "Jane Doe"},
							{FieldName: "category", FieldValue: "study"},
							{FieldName: "price", FieldValue: "5"},
						}},
				},
				Facets:     make([]*solr.FacetResult, 0),
				Highlights: make([]*solr.Highlight, 0),
				Diagnostics: &solr.Diagnostics{
					ParsingSucceeded: true,
				},
			},
		},
		{
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			name: "Matches: Multiple values for a non-text field in a Doc is mapped to multiple key-value" +
				" pairs",
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      2,
					Start:         0,
					MaxScore:      100,
					Docs: []solr.GenericObject{
						{
							"id":         "123",
							"entityName": "phonebookItem",
							"author":     []interface{}{"Jones", "Smith"},
						},
					},
				},
			},
			want: &pb.SearchResponse{
				NumFound:      2,
				NumFoundExact: true,
				Start:         0,
				MaxScore:      100,
				Items: []*solr.DocItem{
					{
						ItemId:     "123",
						EntityType: "phonebookItem",
						Values: []*solr.DocValue{
							{FieldName: "author", FieldValue: "Jones"},
							{FieldName: "author", FieldValue: "Smith"},
						},
					},
				},
				Facets:     make([]*solr.FacetResult, 0),
				Highlights: make([]*solr.Highlight, 0),
				Diagnostics: &solr.Diagnostics{
					ParsingSucceeded: true,
				},
			},
		},
		{
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			name: "Matches: For text fields, " +
				"the underlying Language-specific fields are mapped to separate entries with the same field but" +
				" with appropriate kanguage tags",
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      3,
					Start:         0,
					MaxScore:      100,
					Docs: []solr.GenericObject{
						{
							"id":              "123",
							"entityName":      "bookItem",
							genericTitleField: "A Profound Book",
						},
						{
							"id":         "foo",
							"entityName": "studyItem",
							deTitleField: "Etwas Deutsches",
						},
						{
							"id":         "flub",
							"entityName": "studyItem",
							enTitleField: "Something in English",
						},
					},
				},
			},
			want: &pb.SearchResponse{
				NumFound:      3,
				NumFoundExact: true,
				Start:         0,
				MaxScore:      100,
				Items: []*solr.DocItem{
					{
						ItemId:     "123",
						EntityType: "bookItem",
						Values: []*solr.DocValue{
							{
								FieldName:  "title",
								FieldValue: "A Profound Book",
							},
						},
					},
					{
						ItemId:     "flub",
						EntityType: "studyItem",
						Values: []*solr.DocValue{
							{
								FieldName:  "title",
								FieldValue: "Something in English",
								Language:   solr.EnglishLangAbbrev,
							},
						},
					},
					{
						ItemId:     "foo",
						EntityType: "studyItem",
						Values: []*solr.DocValue{
							{
								FieldName:  "title",
								FieldValue: "Etwas Deutsches",
								Language:   solr.GermanLangAbbrev,
							},
						},
					},
				},
				Facets:     make([]*solr.FacetResult, 0),
				Highlights: make([]*solr.Highlight, 0),
				Diagnostics: &solr.Diagnostics{
					ParsingSucceeded: true,
				},
			},
		},
		{
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			name: "Matches: For text fields, " +
				"multiple values in the underlying Language-specific fields are all mapped to separate entries with the same field but appropriate language tags",
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      3,
					Start:         0,
					MaxScore:      100,
					Docs: []solr.GenericObject{
						{
							"id":              "123",
							"entityName":      "bookItem",
							genericTitleField: []interface{}{"A Profound Book", "Another Profound Book"},
						},
						{
							"id":         "foo",
							"entityName": "studyItem",
							deTitleField: []interface{}{"Etwas Deutsches", "Noch Etwas Deutsches"},
						},
						{
							"id":         "flub",
							"entityName": "studyItem",
							enTitleField: []interface{}{"Something in English", "More in English"},
						},
					},
				},
			},
			want: &pb.SearchResponse{
				NumFound:      3,
				NumFoundExact: true,
				Start:         0,
				MaxScore:      100,
				Items: []*solr.DocItem{
					{
						ItemId:     "123",
						EntityType: "bookItem",
						Values: []*solr.DocValue{
							{
								FieldName:  "title",
								FieldValue: "A Profound Book",
							},
							{
								FieldName:  "title",
								FieldValue: "Another Profound Book",
							},
						},
					},
					{
						ItemId:     "flub",
						EntityType: "studyItem",
						Values: []*solr.DocValue{
							{
								FieldName:  "title",
								FieldValue: "More in English",
								Language:   solr.EnglishLangAbbrev,
							},
							{
								FieldName:  "title",
								FieldValue: "Something in English",
								Language:   solr.EnglishLangAbbrev,
							},
						},
					},
					{
						ItemId:     "foo",
						EntityType: "studyItem",
						Values: []*solr.DocValue{
							{
								FieldName:  "title",
								FieldValue: "Etwas Deutsches",
								Language:   solr.GermanLangAbbrev,
							},
							{
								FieldName:  "title",
								FieldValue: "Noch Etwas Deutsches",
								Language:   solr.GermanLangAbbrev,
							},
						},
					},
				},
				Facets:     make([]*solr.FacetResult, 0),
				Highlights: make([]*solr.Highlight, 0),
				Diagnostics: &solr.Diagnostics{
					ParsingSucceeded: true,
				},
			},
		},
		{
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			name:        "Matches: Docs without 'id' field are dropped",
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      2,
					Start:         0,
					MaxScore:      100,
					Docs: []solr.GenericObject{
						{
							"entityName": "bookItem",
							"category":   "book",
							"author":     []interface{}{"Jones", "Smith"},
							"price":      23.5,
						},
						{
							"id":         "foo",
							"entityName": "studyItem",
							"category":   "study",
							"price":      5,
						},
					},
				},
			},
			want: &pb.SearchResponse{
				NumFound:      2,
				NumFoundExact: true,
				Start:         0,
				MaxScore:      100,
				Items: []*solr.DocItem{
					{
						ItemId:     "foo",
						EntityType: "studyItem",
						Values: []*solr.DocValue{
							{FieldName: "category", FieldValue: "study"},
							{FieldName: "price", FieldValue: "5"},
						}},
				},
				Facets:     make([]*solr.FacetResult, 0),
				Highlights: make([]*solr.Highlight, 0),
				Diagnostics: &solr.Diagnostics{
					ParsingSucceeded: true,
				},
			},
		},
		{
			diagnostics: &solr.Diagnostics{ParsingSucceeded: true},
			name:        "Matches: Docs without 'entityName' field are kept and assigned an empty entity type",
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      2,
					Start:         0,
					MaxScore:      100,
					Docs: []solr.GenericObject{
						{
							"id":       "foo",
							"category": "study",
							"price":    5,
						},
					},
				},
			},
			want: &pb.SearchResponse{
				NumFound:      2,
				NumFoundExact: true,
				Start:         0,
				MaxScore:      100,
				Items: []*solr.DocItem{
					{
						ItemId:     "foo",
						EntityType: "",
						Values: []*solr.DocValue{
							{FieldName: "category", FieldValue: "study"},
							{FieldName: "price", FieldValue: "5"},
						}},
				},
				Facets:     make([]*solr.FacetResult, 0),
				Highlights: make([]*solr.Highlight, 0),
				Diagnostics: &solr.Diagnostics{
					ParsingSucceeded: true,
				},
			},
		},
	}

	log := &L.NullLogger{}
	postQueryHooks, _ := hooks.NewPostQueryHooks(hooks.PostQueryHooksConfig{})

	createResponseConstraintFieldsRepo := frepo.NewMockedFieldRepo([]fields.BaseFieldDef{
		(&kind_text.KindText{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "title", Kind: "text",
			IndexDef: &sharedFields.IndexDef{MultiValued: true}}),
		(&kind_string.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "category", Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
		(&kind_number.KindNumber{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "price", Kind: "number",
			IndexDef: &sharedFields.IndexDef{}}),
		(&kind_string.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "author", Kind: "string",
			IndexDef: &sharedFields.IndexDef{MultiValued: true}}),
	})

	createResponseConstraintSearchConfigRepo := screpo.NewMockSearchConfigRepo([]*searchconfig.SearchConfigObject{})

	for _, tt := range tests {
		opts := QueryEngineOptions{
			Log:              log,
			FieldRepo:        createResponseConstraintFieldsRepo,
			SearchConfigRepo: createResponseConstraintSearchConfigRepo,
			PostQueryHooks:   postQueryHooks,
		}
		qe, _ := newQueryEngine(&constantConverterNonPhrase, opts)
		t.Run(tt.name, func(t *testing.T) {
			got, err := qe.CreateResponse(context.TODO(), tt.solrResponse, tt.facets, tt.diagnostics)
			if (err != nil) != tt.wantErr {
				t.Errorf("createResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			sort.Slice(tt.want.Items, func(i, j int) bool {
				return tt.want.Items[i].ItemId > tt.want.Items[j].ItemId
			})
			sort.Slice(got.Items, func(i, j int) bool {
				return got.Items[i].ItemId > got.Items[j].ItemId
			})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createResponse() got: %v, want: %v", got, tt.want)
			}
		})
	}
}
