package solr

import (
	"context"
	"reflect"
	"testing"

	"google.golang.org/protobuf/types/known/anypb"

	sharedFields "github.com/d4l-data4life/mex/mex/shared/fields"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/searchconfig"
	"github.com/d4l-data4life/mex/mex/shared/searchconfig/screpo"
	"github.com/d4l-data4life/mex/mex/shared/solr"
	"github.com/d4l-data4life/mex/mex/shared/testutils"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/frepo"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/hooks"
	kindhierarchy "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/hierarchy"
	kindstring "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/string"
	kindtext "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/text"
	kindtimestamp "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/timestamp"

	"github.com/d4l-data4life/mex/mex/services/query/endpoints/search/pb"
	"github.com/d4l-data4life/mex/mex/services/query/parser"
)

type QueryTestInfo struct {
	name                 string
	searchRequest        *pb.SearchRequest
	checks               *[]testutils.BodyCheck
	wantErr              bool
	diagChecks           *[]testutils.DiagnosticCheck
	converter            parser.QueryConverter
	dateRanges           *solr.StringFieldRanges
	yearRangeFacets      []*solr.Facet
	constraintIgnoreAxis string
	tolerateBrokenConfig bool
}

func runQueryBodyChecks(got *solr.QueryBody, diag *solr.Diagnostics, err error, t *testing.T, tt QueryTestInfo) {
	if (err != nil) != tt.wantErr {
		t.Errorf("createSolrQueryBody() error = %v, wantErr %v", err, tt.wantErr)
		return
	}
	if tt.checks != nil {
		for _, check := range *tt.checks {
			check(got, t)
		}
	}
	if tt.diagChecks != nil {
		for _, dCheck := range *tt.diagChecks {
			dCheck(diag, t)
		}
	}
}

/*
mapToLanguageSpecificFieldNames returns all language-variants (including
that for generic/unknown language) of a given field name.
*/
func mapToLanguageSpecificFieldNames(fieldName string) ([]string, error) {
	var langSpecificNames []string
	for lc := range solr.KnownLanguagesFieldTypeMap {
		fn, err := solr.GetLangSpecificFieldName(fieldName, lc)
		if err != nil {
			return nil, err
		}
		langSpecificNames = append(langSpecificNames, fn)
	}
	return langSpecificNames, nil
}

// Mocked MEx --> Solr query converter
type MockConverter struct {
	returnQuery            string
	returnCleanedQuery     string
	returnQueryWasCleaned  bool
	returnPhrasesOnlyQuery bool
	err                    parser.TypedParserError
}

func NewMockConverter(query string, cleanedQuery string, queryWasCleaned bool, phrasesOnlyQuery bool, err parser.TypedParserError) MockConverter {
	return MockConverter{
		returnQuery:            query,
		returnCleanedQuery:     cleanedQuery,
		returnQueryWasCleaned:  queryWasCleaned,
		returnPhrasesOnlyQuery: phrasesOnlyQuery,
		err:                    err,
	}
}

func (c *MockConverter) ConvertToSolrQuery(_ string) (*parser.QueryParseResult, parser.TypedParserError) {
	if c.err != nil {
		return nil, c.err
	}
	return &parser.QueryParseResult{
		SolrQuery:        c.returnQuery,
		CleanedQuery:     c.returnCleanedQuery,
		QueryWasCleaned:  c.returnQueryWasCleaned,
		PhrasesOnlyQuery: c.returnPhrasesOnlyQuery,
	}, nil
}

func getStandardSearchConfigRepo(defaultSearchFields []string) searchconfig.SearchConfigRepo {
	return screpo.NewMockSearchConfigRepo([]*searchconfig.SearchConfigObject{
		{
			Type:   solr.MexSearchFocusType,
			Name:   solr.MexDefaultSearchFocusName,
			Fields: defaultSearchFields,
		},
	})
}

// Standard converters used for testing
const DefaultMockedReturnQuery = "X:abc"
const DefaultMockedReturnCleanedQuery = "xyz"
const DefaultMockedReturnQueryWasCleaned = false

var (
	constantConverterNonPhrase = NewMockConverter(DefaultMockedReturnQuery, DefaultMockedReturnCleanedQuery, DefaultMockedReturnQueryWasCleaned, false, nil)
	constantConverterPhrase    = NewMockConverter(DefaultMockedReturnQuery, DefaultMockedReturnCleanedQuery, DefaultMockedReturnQueryWasCleaned, true, nil)
	emptyConverter             = NewMockConverter("", "", false, false, nil)
	parserErrorConverter       = NewMockConverter("", "", false, false, parser.NewParserError(parser.ParserErrorType, "problem!",
		[]string{"error1", "error2"}))
	converterErrorConverter = NewMockConverter("", "", false, false, parser.NewParserError(parser.QueryConstructionErrorType,
		"problem!", nil))
	unknownErrorConverter = NewMockConverter("", "", false, false, parser.NewParserError("unknown", "problem!", nil))
)

func Test_queryEngineFactory(t *testing.T) {
	type testType struct {
		name            string
		searchFocusName string
		maxEditDistance uint32
		useNgramField   bool
		wantErr         bool
	}

	log := &L.NullLogger{}
	postQueryHooks, _ := hooks.NewPostQueryHooks(hooks.PostQueryHooksConfig{})

	orgLinkExt, _ := anypb.New(&sharedFields.IndexDefExtLink{
		RelationType: "relationParent",
	})
	icdLinkExt, _ := anypb.New(&sharedFields.IndexDefExtLink{
		RelationType: "relationIcdParent",
	})
	hierarchy, _ := kindhierarchy.NewKindHierarchy(nil)
	orgHierarchyExt, _ := anypb.New(&sharedFields.IndexDefExtHierarchy{
		CodeSystemNameOrNodeEntityType: "unit",
		LinkFieldName:                  "parentUnit",
		DisplayFieldName:               "label",
	})
	icdHierarchyExt, _ := anypb.New(&sharedFields.IndexDefExtHierarchy{
		CodeSystemNameOrNodeEntityType: "icd",
		LinkFieldName:                  "parentCode",
		DisplayFieldName:               "label",
	})

	axisConstraintFieldsRepo := frepo.NewMockedFieldRepo([]fields.BaseFieldDef{
		(&kindstring.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "category", Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
		(&kindstring.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "type", Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
		(&kindstring.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "country",
			Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
		(&kindtext.KindText{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "title", Kind: "text", IndexDef: &sharedFields.IndexDef{}}),
		hierarchy.MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "orgUnit", Kind: "hierarchy", IndexDef: &sharedFields.IndexDef{
			MultiValued: true,
			Ext:         []*anypb.Any{orgHierarchyExt, orgLinkExt},
		}}),
		hierarchy.MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "icdCode", Kind: "hierarchy", IndexDef: &sharedFields.IndexDef{
			MultiValued: true,
			Ext:         []*anypb.Any{icdHierarchyExt, icdLinkExt},
		}}),
	})
	axisConstraintSearchConfigRepo := screpo.NewMockSearchConfigRepo([]*searchconfig.SearchConfigObject{
		{
			Type:   solr.MexSearchFocusType,
			Name:   solr.MexDefaultSearchFocusName,
			Fields: []string{"type"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "typeAxis",
			Fields: []string{"type"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "countryAxis",
			Fields: []string{"countryDate"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "categoryAxis",
			Fields: []string{"category"},
		},
		{
			Type:   solr.MexHierarchyAxisType,
			Name:   "orgHierarchyAxis",
			Fields: []string{"orgUnitCode"},
		},
		{
			Type:   solr.MexHierarchyAxisType,
			Name:   "icdHierarchyAxis",
			Fields: []string{"icdCode"},
		},
	})

	tests := []testType{
		{
			name:            "QE factory: constructing QE with an edit distance higher than the allowed max value fails",
			searchFocusName: solr.MexDefaultSearchFocusName,
			maxEditDistance: solr.MaxEditDistance + 1,
			useNgramField:   false,
			wantErr:         true,
		},
		{
			name:            "QE factory: constructing QE with an unknown search focus value fails",
			searchFocusName: "notarealfocus",
			maxEditDistance: 0,
			useNgramField:   false,
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		engineOpts := QueryEngineOptions{
			Log:              log,
			FieldRepo:        axisConstraintFieldsRepo,
			SearchConfigRepo: axisConstraintSearchConfigRepo,
			PostQueryHooks:   postQueryHooks,
		}
		queryOpts := QueryOptions{
			SearchFocusName: tt.searchFocusName,
			MaxEditDistance: tt.maxEditDistance,
			UseNgramField:   tt.useNgramField,
		}
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := QueryEngineFactory(context.Background(), queryOpts, engineOpts)
			if tt.wantErr != (gotErr != nil) {
				t.Errorf(": wanted error to occur: %v - error returned: %v", tt.wantErr, gotErr)
			}
		})
	}
}

func Test_createSolrQueryBody_paging(t *testing.T) {
	testQuery := "sunday"

	tests := []QueryTestInfo{
		{
			name: "Paging: Limit defaults to 0",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckLimit(0),
			},
		},
		{
			name: "Paging: Limit > 0 is correctly mapped",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Limit: 25,
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckLimit(25),
			},
		},
		{
			name: "Paging: Limit set to 0 is kept",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Limit: 0,
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckLimit(0),
			},
		},
		{
			name: "Paging: Limit above allowed max is reset to max",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Limit: solr.MaxDocLimit + 100,
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckLimit(solr.MaxDocLimit),
			},
		},
		{
			name: "Paging: Offset defaults to 0",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckOffset(0),
			},
		},
		{
			name: "Paging: Offset == 0 is correctly mapped",
			searchRequest: &pb.SearchRequest{
				Query:  testQuery,
				Offset: 0,
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckOffset(0),
			},
		},
		{
			name: "Paging: Offset > 0 is correctly mapped",
			searchRequest: &pb.SearchRequest{
				Query:  testQuery,
				Offset: 5,
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckOffset(5),
			},
		},
	}

	log := &L.NullLogger{}
	postQueryHooks, _ := hooks.NewPostQueryHooks(hooks.PostQueryHooksConfig{})
	opts := QueryEngineOptions{
		Log:              log,
		FieldRepo:        frepo.NewMockedFieldRepo([]fields.BaseFieldDef{}),
		SearchConfigRepo: getStandardSearchConfigRepo([]string{}),
		PostQueryHooks:   postQueryHooks,
	}
	qe, _ := newQueryEngine(&constantConverterNonPhrase, opts)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, diag, err := qe.CreateSolrQuery(context.TODO(), tt.searchRequest, nil)
			runQueryBodyChecks(body, diag, err, t, tt)
		})
	}
}

func Test_createSolrQueryBody_sorting(t *testing.T) {
	testQuery := "sunday"

	tests := []QueryTestInfo{
		{
			name: "Sorting: Defaults to empty string if no sort information is given",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckSorting("", ""),
			},
		},
		{
			name: "Sorting: For a configured sort axis, " +
				"the correct backing field is used and order defaults to ascending",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Sorting: &solr.Sorting{
					Axis: "general",
				},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckSorting(solr.GetOrdinalAxisSortFieldName("general"), "asc"),
			},
		},
		{
			name: "Sorting: using an invalid ordinal axis name causes an error when used for sorting",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Sorting: &solr.Sorting{
					Axis: "notDefined",
				},
			},
			wantErr: true,
		},
		{
			name: "Sorting: If both field and order is set, they are correctly mapped",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Sorting: &solr.Sorting{
					Axis:  "general",
					Order: "desc",
				},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckSorting(solr.GetOrdinalAxisSortFieldName("general"), "desc"),
			},
		},
		{
			name: "Sorting: An invalid sort order causes an error",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Sorting: &solr.Sorting{
					Axis:  "general",
					Order: "notAnOrder",
				},
			},
			checks:  nil,
			wantErr: true,
		},
		{
			name: "Sorting: If no field but sort order is given, an error occurs",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Sorting: &solr.Sorting{
					Order: "desc",
				},
			},
			checks:  nil,
			wantErr: true,
		},
	}

	log := &L.NullLogger{}
	postQueryHooks, _ := hooks.NewPostQueryHooks(hooks.PostQueryHooksConfig{})
	fieldRepo := frepo.NewMockedFieldRepo([]fields.BaseFieldDef{
		(&kindstring.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "keyword", Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
		(&kindtext.KindText{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "author", Kind: "text", IndexDef: &sharedFields.IndexDef{}}),
	})
	searchConfigRepo := screpo.NewMockSearchConfigRepo([]*searchconfig.SearchConfigObject{
		{
			Type:   solr.MexSearchFocusType,
			Name:   "default",
			Fields: []string{"author"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "general",
			Fields: []string{"keyword", "author"},
		},
	})
	opts := QueryEngineOptions{
		Log:              log,
		FieldRepo:        fieldRepo,
		SearchConfigRepo: searchConfigRepo,
		PostQueryHooks:   postQueryHooks,
	}
	qe, _ := newQueryEngine(&constantConverterNonPhrase, opts)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, diag, err := qe.CreateSolrQuery(context.TODO(), tt.searchRequest, nil)
			runQueryBodyChecks(body, diag, err, t, tt)
		})
	}
}

func Test_createSolrQueryBody_params(t *testing.T) {
	tests := []QueryTestInfo{
		{
			name: "Params: Omits the response header",
			searchRequest: &pb.SearchRequest{
				Query: "",
			},
			checks: &[]testutils.BodyCheck{testutils.CheckOmitHeader(true)},
		},
		{
			name: "Params: Instructs Solr to combine terms with AND",
			searchRequest: &pb.SearchRequest{
				Query: "",
			},
			checks: &[]testutils.BodyCheck{testutils.CheckQOp("AND")},
		},
	}

	log := &L.NullLogger{}
	postQueryHooks, _ := hooks.NewPostQueryHooks(hooks.PostQueryHooksConfig{})
	opts := QueryEngineOptions{
		Log:              log,
		FieldRepo:        frepo.NewMockedFieldRepo([]fields.BaseFieldDef{}),
		SearchConfigRepo: getStandardSearchConfigRepo([]string{}),
		PostQueryHooks:   postQueryHooks,
	}
	qe, _ := newQueryEngine(&constantConverterNonPhrase, opts)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, diag, err := qe.CreateSolrQuery(context.TODO(), tt.searchRequest, nil)
			runQueryBodyChecks(body, diag, err, t, tt)
		})
	}
}

func Test_createSolrQueryBody_query(t *testing.T) {
	tests := []QueryTestInfo{
		{
			name: "Query string: If converter judges search query to be valid (" +
				"grammatical), ParsingSucceeded flag is true",
			searchRequest: &pb.SearchRequest{
				Query: "not relevant since we mock return",
			},
			diagChecks: &[]testutils.DiagnosticCheck{testutils.CheckParsingSuccess(true)},
			converter:  &constantConverterNonPhrase,
		},
		{
			name: "Query string: If converter judges search query to be valid (" +
				"grammatical), no parse errors are returned",
			searchRequest: &pb.SearchRequest{
				Query: "not relevant since we mock return",
			},
			diagChecks: &[]testutils.DiagnosticCheck{testutils.CheckParsingErrors(nil)},
			converter:  &constantConverterNonPhrase,
		},
		{
			name: "Query string: If converter judges search query to be valid (" +
				"grammatical), the cleaned query is returned",
			searchRequest: &pb.SearchRequest{
				Query: "not relevant since we mock return",
			},
			diagChecks: &[]testutils.DiagnosticCheck{testutils.CheckCleanedQuery(DefaultMockedReturnCleanedQuery)},
			converter:  &constantConverterNonPhrase,
		},
		{
			// This is also implicitly tested in the tests below, but it's good to have it as a separate statement
			name: "Query string: If converter judges search query to be invalid (" +
				"ungrammatical), an error is caused",
			searchRequest: &pb.SearchRequest{
				Query: "not relevant since we mock return",
			},
			converter: &parserErrorConverter,
			wantErr:   true,
		},
		{
			name: "Query string: If converter judges search query to be invalid (" +
				"ungrammatical), ParsingSucceeded flag is false",
			searchRequest: &pb.SearchRequest{
				Query: "not relevant since we mock return",
			},
			diagChecks: &[]testutils.DiagnosticCheck{testutils.CheckParsingSuccess(false)},
			converter:  &parserErrorConverter,
			wantErr:    true,
		},
		{
			name: "Query string: If converter judges search query to be invalid (" +
				"ungrammatical) parse errors messages are returned",
			searchRequest: &pb.SearchRequest{
				Query: "not relevant since we mock return",
			},
			diagChecks: &[]testutils.DiagnosticCheck{testutils.CheckParsingErrors([]string{"error1", "error2"})},
			converter:  &parserErrorConverter,
			wantErr:    true,
		},
		{
			name: "Query string: If query is valid but we fail to create a Solr query, an error is returned",
			searchRequest: &pb.SearchRequest{
				Query: "not relevant since we mock return",
			},
			wantErr:   true,
			converter: &converterErrorConverter,
		},
		{
			name: "Query string: If query is valid but an unknown error occurs, an error is returned",
			searchRequest: &pb.SearchRequest{
				Query: "not relevant since we mock return",
			},
			wantErr:   true,
			converter: &unknownErrorConverter,
		},
	}

	log := &L.NullLogger{}
	postQueryHooks, _ := hooks.NewPostQueryHooks(hooks.PostQueryHooksConfig{})

	for _, tt := range tests {
		opts := QueryEngineOptions{
			Log:              log,
			FieldRepo:        frepo.NewMockedFieldRepo([]fields.BaseFieldDef{}),
			SearchConfigRepo: getStandardSearchConfigRepo([]string{}),
			PostQueryHooks:   postQueryHooks,
		}
		qe, _ := newQueryEngine(tt.converter, opts)
		t.Run(tt.name, func(t *testing.T) {
			body, diag, err := qe.CreateSolrQuery(context.TODO(), tt.searchRequest, nil)
			runQueryBodyChecks(body, diag, err, t, tt)
		})
	}
}

func Test_createSolrQueryBody_exact_faceting(t *testing.T) {
	testQuery := "sunday"

	categoryFacet := []*solr.Facet{
		{
			Type: solr.MexExactFacetType,
			Axis: "categoryAxis",
		},
	}
	unconfiguredFacet := []*solr.Facet{
		{
			Type: solr.MexExactFacetType,
			Axis: "journalAxis",
		},
	}
	unknownFieldFacet := []*solr.Facet{
		{
			Type: solr.MexExactFacetType,
			Axis: "missingFieldAxis",
		},
	}
	multipleFacet := []*solr.Facet{
		{
			Type: solr.MexExactFacetType,
			Axis: "categoryAxis",
		},
		{
			Type: solr.MexExactFacetType,
			Axis: "titleAxis",
		},
	}
	multipleFacetResponse := map[string]solr.SolrFacet{}
	multipleFacetResponse[createFacetName("categoryAxis", "")] = solr.SolrFacet{DetailedType: solr.SolrTermsFacetType, Field: "category"}
	multipleFacetResponse[createFacetName("titleAxis", "")] = solr.SolrFacet{DetailedType: solr.SolrTermsFacetType, Field: "title"}

	tests := []QueryTestInfo{
		{
			name: "Faceting - terms facets: Default is no faceting",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(nil),
			},
		},
		{
			name: "Exact (terms) faceting: For single facet axis, " +
				"the correct backing field is used and numBucket is on",
			searchRequest: &pb.SearchRequest{
				Query:  testQuery,
				Facets: categoryFacet,
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(map[string]solr.SolrFacet{
					createFacetName("categoryAxis", ""): {
						DetailedType: solr.SolrTermsFacetType,
						NumBuckets:   true,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("categoryAxis"),
					},
				}),
			},
		},
		{
			name: "Exact (terms) faceting: facets using an unknown ordinal axis cause error if using strict config error policy",
			searchRequest: &pb.SearchRequest{
				Query:  testQuery,
				Facets: unconfiguredFacet,
			},
			tolerateBrokenConfig: false,
			wantErr:              true,
		},
		{
			name: "Exact (terms) faceting: facets using an unknown ordinal axis are skipped (and an ignored error added to diagnostics) if using relaxed config error policy",
			searchRequest: &pb.SearchRequest{
				Query:  testQuery,
				Facets: unconfiguredFacet,
			},
			tolerateBrokenConfig: true,
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(nil),
			},
			diagChecks: &[]testutils.DiagnosticCheck{
				testutils.CheckIgnoredErrors([]string{SkippedFacetWarning}),
			},
		},
		{
			name: "Exact (terms) faceting: a configured axis using unknown fields is still accepted for faceting",
			searchRequest: &pb.SearchRequest{
				Query:  testQuery,
				Facets: unknownFieldFacet,
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(map[string]solr.SolrFacet{
					createFacetName("missingFieldAxis", ""): {
						DetailedType: solr.SolrTermsFacetType,
						NumBuckets:   true,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("missingFieldAxis"),
					},
				}),
			},
		},
		{
			name: "Exact (terms) faceting: For multiple facets, names and back fields are correctly set, " +
				"numBucket is on ",
			searchRequest: &pb.SearchRequest{
				Query:  testQuery,
				Facets: multipleFacet,
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(map[string]solr.SolrFacet{
					createFacetName("categoryAxis", ""): {
						DetailedType: solr.SolrTermsFacetType,
						NumBuckets:   true,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("categoryAxis"),
					},
					createFacetName("titleAxis", ""): {
						DetailedType: solr.SolrTermsFacetType,
						NumBuckets:   true,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("titleAxis"),
					},
				}),
			},
		},
		{
			name: "Exact (terms) faceting: Facet limit > 0 is correctly set",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Facets: []*solr.Facet{
					{
						Type: solr.MexExactFacetType,
						Axis: "categoryAxis",
					},
					{
						Type:  solr.MexExactFacetType,
						Axis:  "titleAxis",
						Limit: 15,
					},
				},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(map[string]solr.SolrFacet{
					createFacetName("categoryAxis", ""): {
						DetailedType: solr.SolrTermsFacetType,
						NumBuckets:   true,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("categoryAxis"),
					},
					createFacetName("titleAxis", ""): {
						DetailedType: solr.SolrTermsFacetType,
						NumBuckets:   true,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("titleAxis"),
						Limit:        15,
					},
				}),
			},
		},
		{
			name: "Exact (terms) faceting: Setting facet limit == 0 is the same as not setting it at all",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Facets: []*solr.Facet{
					{
						Type: solr.MexExactFacetType,
						Axis: "categoryAxis",
					},
					{
						Type:  solr.MexExactFacetType,
						Axis:  "titleAxis",
						Limit: 0,
					},
				},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(map[string]solr.SolrFacet{
					createFacetName("categoryAxis", ""): {
						DetailedType: solr.SolrTermsFacetType,
						NumBuckets:   true,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("categoryAxis"),
					},
					createFacetName("titleAxis", ""): {
						DetailedType: solr.SolrTermsFacetType,
						NumBuckets:   true,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("titleAxis"),
					},
				}),
			},
		},
		{
			name: "Exact (terms) faceting: Facet limits above allowed limit is reset to max",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Facets: []*solr.Facet{
					{
						Type: solr.MexExactFacetType,
						Axis: "categoryAxis",
					},
					{
						Type:  solr.MexExactFacetType,
						Axis:  "titleAxis",
						Limit: solr.MaxFacetLimit + 100,
					},
				},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(map[string]solr.SolrFacet{
					createFacetName("categoryAxis", ""): {
						DetailedType: solr.SolrTermsFacetType,
						NumBuckets:   true,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("categoryAxis"),
					},
					createFacetName("titleAxis", ""): {
						DetailedType: solr.SolrTermsFacetType,
						NumBuckets:   true,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("titleAxis"),
						Limit:        solr.MaxFacetLimit,
					},
				}),
			},
		},
		{
			name: "Exact (terms) faceting: Facet offset is correctly mapped",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Facets: []*solr.Facet{
					{
						Type: solr.MexExactFacetType,
						Axis: "categoryAxis",
					},
					{
						Type:   solr.MexExactFacetType,
						Axis:   "titleAxis",
						Offset: 15,
					},
				},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(map[string]solr.SolrFacet{
					createFacetName("categoryAxis", ""): {
						DetailedType: solr.SolrTermsFacetType,
						NumBuckets:   true,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("categoryAxis"),
					},
					createFacetName("titleAxis", ""): {
						DetailedType: solr.SolrTermsFacetType,
						NumBuckets:   true,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("titleAxis"),
						Offset:       15,
					},
				}),
			},
		},
		{
			name: "Exact (terms) faceting: An axis constraint on the facet ordinal axis is explicitly excluded in" +
				" faceting",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Facets: []*solr.Facet{
					{
						Type: solr.MexExactFacetType,
						Axis: "categoryAxis",
					},
				},
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:   solr.MexExactAxisConstraint,
						Axis:   "categoryAxis",
						Values: []string{"article"},
					},
				},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(map[string]solr.SolrFacet{
					createFacetName("categoryAxis", ""): {
						DetailedType: solr.SolrTermsFacetType,
						NumBuckets:   true,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("categoryAxis"),
						ExcludeTags:  []string{solr.GenerateTagName("categoryAxis")},
					},
				}),
			},
		},
		{
			name: "Exact (terms) faceting: An axis constraint on an ordinal axis OTHER THAN the facet axis is NOT" +
				" excluded in faceting",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Facets: []*solr.Facet{
					{
						Type: solr.MexExactFacetType,
						Axis: "categoryAxis",
					},
				},
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:   solr.MexExactAxisConstraint,
						Axis:   "titleAxis",
						Values: []string{"hand"},
					},
				},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(map[string]solr.SolrFacet{
					createFacetName("categoryAxis", ""): {
						DetailedType: solr.SolrTermsFacetType,
						NumBuckets:   true,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("categoryAxis"),
					},
				}),
			},
		},
	}

	log := &L.NullLogger{}
	postQueryHooks, _ := hooks.NewPostQueryHooks(hooks.PostQueryHooksConfig{})
	facetQueryFieldsRepo := frepo.NewMockedFieldRepo([]fields.BaseFieldDef{
		(&kindstring.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "category", Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
		(&kindtext.KindText{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "title", Kind: "text", IndexDef: &sharedFields.IndexDef{}}),
	})
	facetSearchConfigRepo := screpo.NewMockSearchConfigRepo([]*searchconfig.SearchConfigObject{
		{
			Type:   solr.MexSearchFocusType,
			Name:   solr.MexDefaultSearchFocusName,
			Fields: []string{"title"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "titleAxis",
			Fields: []string{"title"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "categoryAxis",
			Fields: []string{"category"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "missingFieldAxis",
			Fields: []string{"unknownField"},
		},
	})
	for _, tt := range tests {
		opts := QueryEngineOptions{
			Log:                   log,
			FieldRepo:             facetQueryFieldsRepo,
			SearchConfigRepo:      facetSearchConfigRepo,
			PostQueryHooks:        postQueryHooks,
			TolerantErrorHandling: tt.tolerateBrokenConfig,
		}
		qe, _ := newQueryEngine(&constantConverterNonPhrase, opts)
		t.Run(tt.name, func(t *testing.T) {
			body, diag, err := qe.CreateSolrQuery(context.TODO(), tt.searchRequest, nil)
			runQueryBodyChecks(body, diag, err, t, tt)
		})
	}
}

func Test_createSolrQueryBody_year_range_faceting(t *testing.T) {
	testQuery := "sunday"

	publishingDateFacet := []*solr.Facet{
		{
			Type: solr.MexYearRangeFacetType,
			Axis: "publishingDateAxis",
		},
	}
	unconfiguredDateFacet := []*solr.Facet{
		{
			Type: solr.MexYearRangeFacetType,
			Axis: "secretDateAxis",
		},
	}
	unknownFieldDateFacet := []*solr.Facet{
		{
			Type: solr.MexYearRangeFacetType,
			Axis: "missingDateFieldAxis",
		},
	}
	multipleDateFacet := []*solr.Facet{
		{
			Type: solr.MexYearRangeFacetType,
			Axis: "publishingDateAxis",
		},
		{
			Type: solr.MexYearRangeFacetType,
			Axis: "createdAtAxis",
		},
	}
	multipleFacetResponse := map[string]solr.SolrFacet{}
	multipleFacetResponse[createFacetName("categoryAxis", "")] = solr.SolrFacet{DetailedType: solr.SolrTermsFacetType, Field: "category"}
	multipleFacetResponse[createFacetName("titleAxis", "")] = solr.SolrFacet{DetailedType: solr.SolrTermsFacetType, Field: "title"}

	tests := []QueryTestInfo{
		{
			name: "Year range faceting: if no data ranges are available for facet field, the facet is dropped",
			searchRequest: &pb.SearchRequest{
				Query:  testQuery,
				Facets: publishingDateFacet,
			},
			dateRanges: &solr.StringFieldRanges{
				// Only configure a field not used in requested facet
				"createdAtAxis": &solr.StringRange{
					Min: "2012-05-20T17:33:18",
					Max: "2022-05-20T17:33:18Z",
				},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(map[string]solr.SolrFacet{}),
			},
		},
		{
			name: "Year range faceting: if the date range given does not have an identifiable year, the facet is skipped (and an ignored error added to diagnostics) if using relaxed config error policy",
			searchRequest: &pb.SearchRequest{
				Query:  testQuery,
				Facets: publishingDateFacet,
			},
			dateRanges: &solr.StringFieldRanges{
				// Configure with invalid date strings
				"publishingDateAxis": &solr.StringRange{
					Min: "201205-20T17:33:18Z",
					Max: "202205-20T17:33:18Z",
				},
			},
			tolerateBrokenConfig: true,
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(nil),
			},
			diagChecks: &[]testutils.DiagnosticCheck{
				testutils.CheckIgnoredErrors([]string{SkippedFacetWarning}),
			},
		},
		{
			name: "Year range faceting: if the date range given does not have an identifiable year, an error is raised  if using strict config error policy",
			searchRequest: &pb.SearchRequest{
				Query:  testQuery,
				Facets: publishingDateFacet,
			},
			dateRanges: &solr.StringFieldRanges{
				// Configure with invalid date strings
				"publishingDateAxis": &solr.StringRange{
					Min: "201205-20T17:33:18Z",
					Max: "202205-20T17:33:18Z",
				},
			},
			tolerateBrokenConfig: false,
			wantErr:              true,
		},
		{
			name: "Year range faceting: if a facet axis missing in the query engine configuration, the facet is skipped (and an ignored error added to diagnostics) if using relaxed config error policy",
			searchRequest: &pb.SearchRequest{
				Query:  testQuery,
				Facets: unconfiguredDateFacet,
			},
			tolerateBrokenConfig: true,
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(nil),
			},
			diagChecks: &[]testutils.DiagnosticCheck{
				testutils.CheckIgnoredErrors([]string{SkippedFacetWarning}),
			},
		},
		{
			name: "Year range faceting: if a facet axis missing in the query engine configuration, an error is cause if using strict config error policy",
			searchRequest: &pb.SearchRequest{
				Query:  testQuery,
				Facets: unconfiguredDateFacet,
			},
			tolerateBrokenConfig: false,
			wantErr:              true,
		},
		{
			name: "Year range faceting: faceting by an ordinal axis configured with an unknown field causes the facet to be skipped (and an ignored error added to diagnostics) if using relaxed config error policy",
			searchRequest: &pb.SearchRequest{
				Query:  testQuery,
				Facets: unknownFieldDateFacet,
			},
			tolerateBrokenConfig: true,
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(nil),
			},
		},
		{
			name: "Year range faceting: faceting by an ordinal axis configured with an unknown field causes an error if using strict config error policy",
			searchRequest: &pb.SearchRequest{
				Query:  testQuery,
				Facets: unknownFieldDateFacet,
			},
			tolerateBrokenConfig: false,
			wantErr:              true,
		},
		{
			name: "Year range faceting: For single facet axis, " +
				"the range is set to the passed range for the axis (rounded to whole years) and the gap to 1 year",
			searchRequest: &pb.SearchRequest{
				Query:  testQuery,
				Facets: publishingDateFacet,
			},
			dateRanges: &solr.StringFieldRanges{
				// Only configure a field not used in requested facet
				"publishingDateAxis": &solr.StringRange{
					Min: "2012-05-20T17:33:18.77Z",
					Max: "2022-05-20T17:33:18.77Z",
				},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(map[string]solr.SolrFacet{
					createFacetName("publishingDateAxis", ""): {
						DetailedType: solr.SolrStringRangeFacetType,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("publishingDateAxis"),
						StartString:  "2012-01-01T00:00:00.000Z",
						EndString:    "2022-12-31T23:59:59.999Z",
						GapString:    "+1YEARS",
					},
				}),
			},
		},
		{
			name: "Year range faceting: For multiple facets, names and ranges are correctly set",
			searchRequest: &pb.SearchRequest{
				Query:  testQuery,
				Facets: multipleDateFacet,
			},
			dateRanges: &solr.StringFieldRanges{
				// Only configure a field not used in requested facet
				"publishingDateAxis": &solr.StringRange{
					Min: "2012-05-20T17:33:18Z",
					Max: "2022-05-20T17:33:18Z",
				},
				"createdAtAxis": &solr.StringRange{
					Min: "2002-06-02T17:33:18Z",
					Max: "2018-12-20T17:33:18Z",
				},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(map[string]solr.SolrFacet{
					createFacetName("publishingDateAxis", ""): {
						DetailedType: solr.SolrStringRangeFacetType,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("publishingDateAxis"),
						StartString:  "2012-01-01T00:00:00.000Z",
						EndString:    "2022-12-31T23:59:59.999Z",
						GapString:    "+1YEARS",
					},
					createFacetName("createdAtAxis", ""): {
						DetailedType: solr.SolrStringRangeFacetType,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("createdAtAxis"),
						StartString:  "2002-01-01T00:00:00.000Z",
						EndString:    "2018-12-31T23:59:59.999Z",
						GapString:    "+1YEARS",
					},
				}),
			},
		},
		{
			name: "Year range faceting: An axis constraint on the facet field is excluded in faceting",
			searchRequest: &pb.SearchRequest{
				Query:  testQuery,
				Facets: publishingDateFacet,
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type: solr.MexStringRangeConstraint,
						Axis: "publishingDateAxis",
						StringRanges: []*solr.StringRange{
							{
								Min: "2016-05-20T17:33:18.77Z",
								Max: "2018-05-20T17:33:18.77Z",
							},
						},
					},
				},
			},
			dateRanges: &solr.StringFieldRanges{
				// Only configure a field not used in requested facet
				"publishingDateAxis": &solr.StringRange{
					Min: "2012-05-20T17:33:18.77Z",
					Max: "2022-05-20T17:33:18.77Z",
				},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(map[string]solr.SolrFacet{
					createFacetName("publishingDateAxis", ""): {
						DetailedType: solr.SolrStringRangeFacetType,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("publishingDateAxis"),
						StartString:  "2012-01-01T00:00:00.000Z",
						EndString:    "2022-12-31T23:59:59.999Z",
						GapString:    "+1YEARS",
						ExcludeTags:  []string{solr.GenerateTagName("publishingDateAxis")},
					},
				}),
			},
		},
		{
			name: "Year range faceting: An axis constraint on a field OTHER THAN the facet field is NOT excluded in" +
				" faceting",
			searchRequest: &pb.SearchRequest{
				Query:  testQuery,
				Facets: publishingDateFacet,
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type: solr.MexStringRangeConstraint,
						Axis: "createdAtAxis",
						StringRanges: []*solr.StringRange{
							{
								Min: "2016-05-20T17:33:18.77Z",
								Max: "2018-05-20T17:33:18.77Z",
							},
						},
					},
				},
			},
			dateRanges: &solr.StringFieldRanges{
				// Only configure a field not used in requested facet
				"publishingDateAxis": &solr.StringRange{
					Min: "2012-05-20T17:33:18.77Z",
					Max: "2022-05-20T17:33:18.77Z",
				},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(map[string]solr.SolrFacet{
					createFacetName("publishingDateAxis", ""): {
						DetailedType: solr.SolrStringRangeFacetType,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("publishingDateAxis"),
						StartString:  "2012-01-01T00:00:00.000Z",
						EndString:    "2022-12-31T23:59:59.999Z",
						GapString:    "+1YEARS",
					},
				}),
			},
		},
	}

	log := &L.NullLogger{}
	postQueryHooks, _ := hooks.NewPostQueryHooks(hooks.PostQueryHooksConfig{})
	dateRangeFacetQueryFieldsRepo := frepo.NewMockedFieldRepo([]fields.BaseFieldDef{
		(&kindtimestamp.KindTimestamp{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "publishingDate",
			Kind: "timestamp", IndexDef: &sharedFields.IndexDef{}}),
		(&kindstring.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "category", Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
	})
	dateRangeSearchConfigRepo := screpo.NewMockSearchConfigRepo([]*searchconfig.SearchConfigObject{
		{
			Type:   solr.MexSearchFocusType,
			Name:   solr.MexDefaultSearchFocusName,
			Fields: []string{"category"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "createdAtAxis",
			Fields: []string{"createdAt"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "publishingDateAxis",
			Fields: []string{"publishingDate"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "missingDateFieldAxis",
			Fields: []string{"unknownDateField"},
		},
	})

	for _, tt := range tests {
		opts := QueryEngineOptions{
			Log:                   log,
			FieldRepo:             dateRangeFacetQueryFieldsRepo,
			SearchConfigRepo:      dateRangeSearchConfigRepo,
			PostQueryHooks:        postQueryHooks,
			TolerantErrorHandling: tt.tolerateBrokenConfig,
		}
		qe, _ := newQueryEngine(&constantConverterNonPhrase, opts)
		t.Run(tt.name, func(t *testing.T) {
			body, diag, err := qe.CreateSolrQuery(context.TODO(), tt.searchRequest, tt.dateRanges)
			runQueryBodyChecks(body, diag, err, t, tt)
		})
	}
}

func Test_createSolrQueryBody_string_stat_faceting(t *testing.T) {
	testQuery := "sunday"

	dateRangeFacetQueryFieldsRepo := frepo.NewMockedFieldRepo([]fields.BaseFieldDef{
		(&kindtimestamp.KindTimestamp{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "publishingDate",
			Kind: "timestamp", IndexDef: &sharedFields.IndexDef{}}),
		(&kindstring.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "category", Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
		(&kindtext.KindText{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "title", Kind: "text",
			IndexDef: &sharedFields.IndexDef{}}),
	})

	tests := []QueryTestInfo{
		{
			name: "String stat faceting: the facet is skipped (and an ignored error added to diagnostics) if the axis fields are not of the string or timestamp kind and we use the relaxed config error policy",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Facets: []*solr.Facet{
					{
						Type:   solr.MexStringStatFacetType,
						Axis:   "titleAxis",
						StatOp: "min",
					},
				},
			},
			tolerateBrokenConfig: true,
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(nil),
			},
			diagChecks: &[]testutils.DiagnosticCheck{
				testutils.CheckIgnoredErrors([]string{SkippedFacetWarning}),
			},
		},
		{
			name: "String stat faceting: an error is caused if the axis fields are not of the string or timestamp kind and we use the strict config error policy",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Facets: []*solr.Facet{
					{
						Type:   solr.MexStringStatFacetType,
						Axis:   "titleAxis",
						StatOp: "min",
					},
				},
			},
			tolerateBrokenConfig: false,
			wantErr:              true,
		},
		{
			name: "String stat faceting: the facet is skipped (and an ignored error added to diagnostics) if the operator is unknown and we use the relaxed config error policy",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Facets: []*solr.Facet{
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "publishingDateAxis",
						StatOp:   "magic_op",
						StatName: "min_facet",
					},
				},
			},
			tolerateBrokenConfig: true,
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(nil),
			},
			diagChecks: &[]testutils.DiagnosticCheck{
				testutils.CheckIgnoredErrors([]string{SkippedFacetWarning}),
			},
		},
		{
			name: "String stat faceting: an error is cause if the operator is unknown and we use the strict config error policy",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Facets: []*solr.Facet{
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "publishingDateAxis",
						StatOp:   "magic_op",
						StatName: "min_facet",
					},
				},
			},
			tolerateBrokenConfig: false,
			wantErr:              true,
		},
		{
			name: "String stat faceting: the facet is skipped (and an ignored error added to diagnostics) if stat is missing and we use the relaxed config error policy",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Facets: []*solr.Facet{
					{
						Type:   solr.MexStringStatFacetType,
						Axis:   "publishingDateAxis",
						StatOp: solr.MinOperator,
					},
				},
			},
			tolerateBrokenConfig: true,
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(nil),
			},
			diagChecks: &[]testutils.DiagnosticCheck{
				testutils.CheckIgnoredErrors([]string{SkippedFacetWarning}),
			},
		},
		{
			name: "String stat faceting: an error is cause if stat is missing and we use the strict config error policy",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Facets: []*solr.Facet{
					{
						Type:   solr.MexStringStatFacetType,
						Axis:   "publishingDateAxis",
						StatOp: solr.MinOperator,
					},
				},
			},
			tolerateBrokenConfig: false,
			wantErr:              true,
		},
		{
			name: "String stat faceting: the facet is skipped (and an ignored error added to diagnostics) if the axis is not configured and we use the relaxed config error policy",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Facets: []*solr.Facet{
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "nonExistingAxis",
						StatOp:   solr.MinOperator,
						StatName: "min_facet",
					},
				},
			},
			tolerateBrokenConfig: true,
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(nil),
			},
			diagChecks: &[]testutils.DiagnosticCheck{
				testutils.CheckIgnoredErrors([]string{SkippedFacetWarning}),
			},
		},
		{
			name: "String stat faceting: an error is caused if the axis is not configured and we use the strict config error policy",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Facets: []*solr.Facet{
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "nonExistingAxis",
						StatOp:   solr.MinOperator,
						StatName: "min_facet",
					},
				},
			},
			tolerateBrokenConfig: false,
			wantErr:              true,
		},
		{
			name: "String stat faceting: the facet is skipped (and an ignored error added to diagnostics) if the axis contain an unknown field if using relaxed config error policy",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Facets: []*solr.Facet{
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "missingFieldRangeAxis",
						StatOp:   solr.MinOperator,
						StatName: "min_facet",
					},
				},
			},
			tolerateBrokenConfig: true,
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(nil),
			},
			diagChecks: &[]testutils.DiagnosticCheck{
				testutils.CheckIgnoredErrors([]string{SkippedFacetWarning}),
			},
		},
		{
			name: "String stat faceting: an error is caused if the axis contain an unknown field and we use the strict config error policy",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Facets: []*solr.Facet{
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "missingFieldRangeAxis",
						StatOp:   solr.MinOperator,
						StatName: "min_facet",
					},
				},
			},
			tolerateBrokenConfig: false,
			wantErr:              true,
		},
		{
			name: "String stat faceting: for a single facet, operator and field name are correctly set",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Facets: []*solr.Facet{
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "publishingDateAxis",
						StatOp:   solr.MinOperator,
						StatName: "min_facet",
					},
				},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(map[string]solr.SolrFacet{
					createFacetName("publishingDateAxis", "min_facet"): {
						DetailedType: solr.SolrStringStatFacetType,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("publishingDateAxis"),
						StatOp:       solr.MinOperator,
					},
				}),
			},
		},
		{
			name: "String stat faceting: for multiple facets on the same field, " +
				"operator and field name are correctly set",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Facets: []*solr.Facet{
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "publishingDateAxis",
						StatOp:   solr.MinOperator,
						StatName: "min_facet",
					},
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "publishingDateAxis",
						StatOp:   solr.MaxOperator,
						StatName: "max_facet",
					},
				},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(map[string]solr.SolrFacet{
					createFacetName("publishingDateAxis", "min_facet"): {
						DetailedType: solr.SolrStringStatFacetType,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("publishingDateAxis"),
						StatOp:       solr.MinOperator,
					},
					createFacetName("publishingDateAxis", "max_facet"): {
						DetailedType: solr.SolrStringStatFacetType,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("publishingDateAxis"),
						StatOp:       solr.MaxOperator,
					},
				}),
			},
		},
		{
			name: "String stat faceting: for multiple facets on the same field, " +
				"operator and field name are correctly set",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Facets: []*solr.Facet{
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "publishingDateAxis",
						StatOp:   solr.MinOperator,
						StatName: "min_facet",
					},
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "publishingDateAxis",
						StatOp:   solr.MaxOperator,
						StatName: "max_facet",
					},
				},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(map[string]solr.SolrFacet{
					createFacetName("publishingDateAxis", "min_facet"): {
						DetailedType: solr.SolrStringStatFacetType,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("publishingDateAxis"),
						StatOp:       solr.MinOperator,
					},
					createFacetName("publishingDateAxis", "max_facet"): {
						DetailedType: solr.SolrStringStatFacetType,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("publishingDateAxis"),
						StatOp:       solr.MaxOperator,
					},
				}),
			},
		},
		{
			name: "String stat faceting: for multiple facets on different fields, " +
				"if the stat names are the same for several facets, the second and further occurrences are skipped (and an ignored error added to diagnostics) if we use the relaxed config error policy",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Facets: []*solr.Facet{
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "publishingDateAxis",
						StatOp:   solr.MinOperator,
						StatName: "min_facet",
					},
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "createdAtAxis",
						StatOp:   solr.MinOperator,
						StatName: "min_facet",
					},
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "createdAtAxis",
						StatOp:   solr.MaxOperator,
						StatName: "min_facet",
					},
				},
			},
			tolerateBrokenConfig: true,
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(map[string]solr.SolrFacet{
					createFacetName("publishingDateAxis", "min_facet"): {
						DetailedType: solr.SolrStringStatFacetType,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("publishingDateAxis"),
						StatOp:       solr.MinOperator,
					},
				}),
			},
			diagChecks: &[]testutils.DiagnosticCheck{
				// Two warnings accumulated since the same facet name is used three timers
				testutils.CheckIgnoredErrors([]string{SkippedFacetWarning, SkippedFacetWarning}),
			},
		},
		{
			name: "String stat faceting: for multiple facets on different fields, " +
				"if the stat names are the same for several facets, an error is caused if we use the strict config error policy",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Facets: []*solr.Facet{
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "publishingDateAxis",
						StatOp:   solr.MinOperator,
						StatName: "min_facet",
					},
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "createdAtAxis",
						StatOp:   solr.MinOperator,
						StatName: "min_facet",
					},
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "createdAtAxis",
						StatOp:   solr.MaxOperator,
						StatName: "min_facet",
					},
				},
			},
			tolerateBrokenConfig: false,
			wantErr:              true,
		},
		{
			name: "String stat faceting: axis constraints are NOT excluded when faceting, " +
				"even if set on the facet axis",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
				Facets: []*solr.Facet{
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "publishingDateAxis",
						StatOp:   solr.MinOperator,
						StatName: "min_facet",
					},
				},
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type: solr.MexStringRangeConstraint,
						Axis: "publishingDateAxis",
						StringRanges: []*solr.StringRange{
							{
								Min: "2016-05-20T17:33:18.77Z",
								Max: "2018-05-20T17:33:18.77Z",
							},
						},
					},
					{
						Type: solr.MexStringRangeConstraint,
						Axis: "createdAtAxis",
						StringRanges: []*solr.StringRange{
							{
								Min: "2016-05-20T17:33:18.77Z",
								Max: "2018-05-20T17:33:18.77Z",
							},
						},
					},
				},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(map[string]solr.SolrFacet{
					createFacetName("publishingDateAxis", "min_facet"): {
						DetailedType: solr.SolrStringStatFacetType,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("publishingDateAxis"),
						StatOp:       solr.MinOperator,
					},
				}),
			},
		},
	}

	log := &L.NullLogger{}
	postQueryHooks, _ := hooks.NewPostQueryHooks(hooks.PostQueryHooksConfig{})
	statAxisConstraintSearchConfigRepo := screpo.NewMockSearchConfigRepo([]*searchconfig.SearchConfigObject{
		{
			Type:   solr.MexSearchFocusType,
			Name:   solr.MexDefaultSearchFocusName,
			Fields: []string{"publishingDate"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "publishingDateAxis",
			Fields: []string{"publishingDate"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "createdAtAxis",
			Fields: []string{"createdAt"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "missingFieldRangeAxis",
			Fields: []string{"unknownField"},
		},
	})

	for _, tt := range tests {
		opts := QueryEngineOptions{
			Log:                   log,
			FieldRepo:             dateRangeFacetQueryFieldsRepo,
			SearchConfigRepo:      statAxisConstraintSearchConfigRepo,
			PostQueryHooks:        postQueryHooks,
			TolerantErrorHandling: tt.tolerateBrokenConfig,
		}
		qe, _ := newQueryEngine(&constantConverterNonPhrase, opts)
		t.Run(tt.name, func(t *testing.T) {
			body, diag, err := qe.CreateSolrQuery(context.TODO(), tt.searchRequest, tt.dateRanges)
			runQueryBodyChecks(body, diag, err, t, tt)
		})
	}
}

func Test_createSolrQueryBody_exact_axis_constraints(t *testing.T) {
	tests := []QueryTestInfo{
		{
			name: "Axis constraints: A constraint with no type causes error",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Axis:   "categoryAxis",
						Values: []string{"article"},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			wantErr:   true,
		},
		{
			name: "Axis constraints: A constraint with unknown type causes error",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:   "notAKnownType",
						Axis:   "categoryAxis",
						Values: []string{"article"},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			wantErr:   true,
		},
		{
			name: "Axis constraints: A constraint using a search focus as an axis causes error",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:   solr.MexExactAxisConstraint,
						Axis:   solr.MexDefaultSearchFocusName,
						Values: []string{"article"},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			wantErr:   true,
		},
		{
			name: "Axis constraints: Axis constraints do not affect the main query",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:   solr.MexExactAxisConstraint,
						Axis:   "categoryAxis",
						Values: []string{"article and more"},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks:    &[]testutils.BodyCheck{testutils.CheckQuery(DefaultMockedReturnQuery)},
		},
		{
			name: "Axis constraints: Constraining an axis which is absent from the query engine configuration causes" +
				" an error",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:   solr.MexExactAxisConstraint,
						Axis:   "absentAxis",
						Values: []string{"article"},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			wantErr:   true,
		},
		{
			name: "Axis constraints: Constraining an ordinal axis to a single, " +
				"one-word value maps to single, double-quoted constraint on that field",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:   solr.MexExactAxisConstraint,
						Axis:   "categoryAxis",
						Values: []string{"article"},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetOrdinalAxisFacetAndFilterFieldName("categoryAxis") + ":\"article\""},
				[]string{"categoryAxis"},
			)},
		},
		{
			name: "Axis constraints: For an (non-hierarchical) axis, single-node constraints are ignored",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:             solr.MexExactAxisConstraint,
						Axis:             "categoryAxis",
						Values:           []string{"article"},
						SingleNodeValues: []string{"ignored"},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetOrdinalAxisFacetAndFilterFieldName("categoryAxis") + ":\"article\""},
				[]string{"categoryAxis"},
			)},
		},
		{
			name: "Axis constraints: Constraining a non-hierarchical axis to a multi-word string maps" +
				" to single, double-quoted constraint on that field",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:   solr.MexExactAxisConstraint,
						Axis:   "categoryAxis",
						Values: []string{"article and more"},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetOrdinalAxisFacetAndFilterFieldName("categoryAxis") + ":\"article and more\""},
				[]string{"categoryAxis"},
			)},
		},
		{
			name: "Axis constraints: non-hierarchical axis value is sanitized",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:   solr.MexExactAxisConstraint,
						Axis:   "typeAxis",
						Values: []string{"art&&icle:"},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetOrdinalAxisFacetAndFilterFieldName("typeAxis") + ":\"art\\&\\&icle\\:\""},
				[]string{"typeAxis"},
			)},
		},
		{
			name: "Axis constraints: non-hierarchical multi-line axis value is sanitized",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:   solr.MexExactAxisConstraint,
						Axis:   "typeAxis",
						Values: []string{"first\n art&&icle: AND"},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetOrdinalAxisFacetAndFilterFieldName("typeAxis") + ":\"first\n art\\&\\&icle\\: and\""},
				[]string{"typeAxis"},
			)},
		},
		{
			name: "Axis constraints: in constraints on non-hierarchical axes, non-supported control characters are escaped",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:   solr.MexExactAxisConstraint,
						Axis:   "typeAxis",
						Values: []string{"article:type!"},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetOrdinalAxisFacetAndFilterFieldName("typeAxis") + ":\"article\\:type\\!\""},
				[]string{"typeAxis"},
			)},
		},
		{
			name: "Axis constraints: constraints for several possible values for the same non-hierarchical axis are ORed together by" +
				" default",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:   solr.MexExactAxisConstraint,
						Axis:   "typeAxis",
						Values: []string{"article", "book"},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetOrdinalAxisFacetAndFilterFieldName("typeAxis") + ":\"article\" || " + solr.
					GetOrdinalAxisFacetAndFilterFieldName("typeAxis") + ":\"book\""},
				[]string{"typeAxis"},
			)},
		},
		{
			name: "Axis constraints: constraints for several possible values for the same non-hierarchical axis can be explicitly" +
				" set to be ORed together",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:            solr.MexExactAxisConstraint,
						Axis:            "typeAxis",
						Values:          []string{"article", "book"},
						CombineOperator: "or",
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetOrdinalAxisFacetAndFilterFieldName("typeAxis") + ":\"article\" || " + solr.
					GetOrdinalAxisFacetAndFilterFieldName("typeAxis") + ":\"book\""},
				[]string{"typeAxis"},
			)},
		},
		{
			name: "Axis constraints: constraints for several possible values for the same non-hierarchical axis can explicitly" +
				" be set to be ANDed together",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:            solr.MexExactAxisConstraint,
						Axis:            "typeAxis",
						Values:          []string{"article", "book"},
						CombineOperator: "and",
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetOrdinalAxisFacetAndFilterFieldName("typeAxis") + ":\"article\" && " + solr.
					GetOrdinalAxisFacetAndFilterFieldName("typeAxis") + ":\"book\""},
				[]string{"typeAxis"},
			)},
		},
		{
			name: "Axis constraints: unsupported Boolean combination operator for value constraints" +
				" causes error",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:            solr.MexExactAxisConstraint,
						Axis:            "type",
						Values:          []string{"article", "book"},
						CombineOperator: "xor",
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks:    nil,
			wantErr:   true,
		},
		{
			name: "Axis constraints: constraints on separate non-hierarchical axes lead to separate filters",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:   solr.MexExactAxisConstraint,
						Axis:   "countryAxis",
						Values: []string{"dk"},
					},
					{
						Type:   solr.MexExactAxisConstraint,
						Axis:   "typeAxis",
						Values: []string{"article", "book"},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetOrdinalAxisFacetAndFilterFieldName("countryAxis") + ":\"dk\"",
					solr.GetOrdinalAxisFacetAndFilterFieldName("typeAxis") + ":\"article\" || " + solr.
						GetOrdinalAxisFacetAndFilterFieldName("typeAxis") + ":\"book\""},
				[]string{"countryAxis", "typeAxis"},
			)},
		},
		{
			name: "Axis constraints: For a hierarchical axis, a single-node constraint to a single, " +
				"value maps to single, double-quoted constraint on that field",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:             solr.MexExactAxisConstraint,
						Axis:             "orgHierarchyAxis",
						SingleNodeValues: []string{"U1"},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetSingleNodeAxisFieldName(solr.GetOrdinalAxisFacetAndFilterFieldName("orgHierarchyAxis")) + ":\"U1\""},
				[]string{"orgHierarchyAxis"},
			)},
		},
		{
			name: "Axis constraints: Single-value constraint on a  hierarchical axis to a multi-word string maps" +
				" to single, double-quoted constraint on that field",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:             solr.MexExactAxisConstraint,
						Axis:             "orgHierarchyAxis",
						SingleNodeValues: []string{"article and more"},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetSingleNodeAxisFieldName(solr.GetOrdinalAxisFacetAndFilterFieldName("orgHierarchyAxis")) + ":\"article and more\""},
				[]string{"orgHierarchyAxis"},
			)},
		},
		{
			name: "Axis constraints: hierarchical single-value axis constraint value is sanitized",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:             solr.MexExactAxisConstraint,
						Axis:             "orgHierarchyAxis",
						SingleNodeValues: []string{"art&&icle:"},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetSingleNodeAxisFieldName(solr.GetOrdinalAxisFacetAndFilterFieldName("orgHierarchyAxis")) + ":\"art\\&\\&icle\\:\""},
				[]string{"orgHierarchyAxis"},
			)},
		},
		{
			name: "Axis constraints: multi-line, single-value hierarchical axis value constraint is sanitized",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:             solr.MexExactAxisConstraint,
						Axis:             "orgHierarchyAxis",
						SingleNodeValues: []string{"first\n art&&icle: AND"},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetSingleNodeAxisFieldName(solr.GetOrdinalAxisFacetAndFilterFieldName("orgHierarchyAxis")) + ":\"first\n art\\&\\&icle\\: and\""},
				[]string{"orgHierarchyAxis"},
			)},
		},
		{
			name: "Axis constraints: non-supported control characters are escaped for single-value constraints on hierarchical axis",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:             solr.MexExactAxisConstraint,
						Axis:             "orgHierarchyAxis",
						SingleNodeValues: []string{"article:type!"},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetSingleNodeAxisFieldName(solr.GetOrdinalAxisFacetAndFilterFieldName("orgHierarchyAxis")) + ":\"article\\:type\\!\""},
				[]string{"orgHierarchyAxis"},
			)},
		},
		{
			name: "Axis constraints: single-value constraints for several possible values for the same hierarchical axis are ORed together by" +
				" default",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:             solr.MexExactAxisConstraint,
						Axis:             "orgHierarchyAxis",
						SingleNodeValues: []string{"article", "book"},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetSingleNodeAxisFieldName(solr.GetOrdinalAxisFacetAndFilterFieldName("orgHierarchyAxis")) + ":\"article\" || " + solr.GetSingleNodeAxisFieldName(solr.GetOrdinalAxisFacetAndFilterFieldName("orgHierarchyAxis")) + ":\"book\""},
				[]string{"orgHierarchyAxis"},
			)},
		},
		{
			name: "Axis constraints: single-value constraints for several possible values for the same hierarchical axis can be explicitly" +
				" set to be ORed together",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:             solr.MexExactAxisConstraint,
						Axis:             "orgHierarchyAxis",
						SingleNodeValues: []string{"article", "book"},
						CombineOperator:  "or",
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetSingleNodeAxisFieldName(solr.GetOrdinalAxisFacetAndFilterFieldName("orgHierarchyAxis")) + ":\"article\" || " + solr.GetSingleNodeAxisFieldName(solr.GetOrdinalAxisFacetAndFilterFieldName("orgHierarchyAxis")) + ":\"book\""},
				[]string{"orgHierarchyAxis"},
			)},
		},
		{
			name: "Axis constraints: single-value constraints for several possible values for the same hierarchical axis are can explicitly" +
				" be set to be ANDed together",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:             solr.MexExactAxisConstraint,
						Axis:             "orgHierarchyAxis",
						SingleNodeValues: []string{"article", "book"},
						CombineOperator:  "and",
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetSingleNodeAxisFieldName(solr.GetOrdinalAxisFacetAndFilterFieldName("orgHierarchyAxis")) + ":\"article\" && " + solr.GetSingleNodeAxisFieldName(solr.GetOrdinalAxisFacetAndFilterFieldName("orgHierarchyAxis")) + ":\"book\""},
				[]string{"orgHierarchyAxis"},
			)},
		},
		{
			name: "Axis constraints: mixed single-value and normal (sub-tree) constraints for the same hierarchical axis by default ORed together",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:             solr.MexExactAxisConstraint,
						Axis:             "orgHierarchyAxis",
						Values:           []string{"article"},
						SingleNodeValues: []string{"book"},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetOrdinalAxisFacetAndFilterFieldName("orgHierarchyAxis") + ":\"article\" || " + solr.GetSingleNodeAxisFieldName(solr.GetOrdinalAxisFacetAndFilterFieldName("orgHierarchyAxis")) + ":\"book\""},
				[]string{"orgHierarchyAxis"},
			)},
		},
		{
			name: "Axis constraints: mixed single-value and normal (sub-tree) constraints for the same hierarchical axis be explicitly set to be ORed together",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:             solr.MexExactAxisConstraint,
						Axis:             "orgHierarchyAxis",
						Values:           []string{"article"},
						SingleNodeValues: []string{"book"},
						CombineOperator:  "or",
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetOrdinalAxisFacetAndFilterFieldName("orgHierarchyAxis") + ":\"article\" || " + solr.GetSingleNodeAxisFieldName(solr.GetOrdinalAxisFacetAndFilterFieldName("orgHierarchyAxis")) + ":\"book\""},
				[]string{"orgHierarchyAxis"},
			)},
		},
		{
			name: "Axis constraints: mixed single-value and normal (sub-tree) constraints for the same hierarchical axis be explicitly set to be ORed together",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:             solr.MexExactAxisConstraint,
						Axis:             "orgHierarchyAxis",
						Values:           []string{"article"},
						SingleNodeValues: []string{"book"},
						CombineOperator:  "and",
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetOrdinalAxisFacetAndFilterFieldName("orgHierarchyAxis") + ":\"article\" && " + solr.GetSingleNodeAxisFieldName(solr.GetOrdinalAxisFacetAndFilterFieldName("orgHierarchyAxis")) + ":\"book\""},
				[]string{"orgHierarchyAxis"},
			)},
		},
		{
			name: "Axis constraints: unsupported Boolean combination operator for hierarchical single-value constraints" +
				" causes error",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:             solr.MexExactAxisConstraint,
						Axis:             "orgHierarchyAxis",
						SingleNodeValues: []string{"article", "book"},
						CombineOperator:  "xor",
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks:    nil,
			wantErr:   true,
		},
		{
			name: "Axis constraints: single-value constraints on separate hierarchical axes lead to separate filters",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:             solr.MexExactAxisConstraint,
						Axis:             "orgHierarchyAxis",
						SingleNodeValues: []string{"dk"},
					},
					{
						Type:             solr.MexExactAxisConstraint,
						Axis:             "icdHierarchyAxis",
						SingleNodeValues: []string{"article", "book"},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetSingleNodeAxisFieldName(solr.GetOrdinalAxisFacetAndFilterFieldName("orgHierarchyAxis")) + ":\"dk\"",
					solr.GetSingleNodeAxisFieldName(solr.GetOrdinalAxisFacetAndFilterFieldName("icdHierarchyAxis")) + ":\"article\" || " + solr.GetSingleNodeAxisFieldName(solr.GetOrdinalAxisFacetAndFilterFieldName("icdHierarchyAxis")) + ":\"book\""},
				[]string{"orgHierarchyAxis", "icdHierarchyAxis"},
			)},
		},
		{
			name: "Axis constraints: single-value constraints are ignored for a non-hierarchical axes",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:             solr.MexExactAxisConstraint,
						Axis:             "typeAxis",
						Values:           []string{"article"},
						SingleNodeValues: []string{"book"},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetOrdinalAxisFacetAndFilterFieldName("typeAxis") + ":\"article\""},
				[]string{"typeAxis"},
			)},
		},
	}

	log := &L.NullLogger{}
	postQueryHooks, _ := hooks.NewPostQueryHooks(hooks.PostQueryHooksConfig{})

	orgLinkExt, _ := anypb.New(&sharedFields.IndexDefExtLink{
		RelationType: "relationParent",
	})
	icdLinkExt, _ := anypb.New(&sharedFields.IndexDefExtLink{
		RelationType: "relationIcdParent",
	})
	hierarchy, _ := kindhierarchy.NewKindHierarchy(nil)
	orgHierarchyExt, _ := anypb.New(&sharedFields.IndexDefExtHierarchy{
		CodeSystemNameOrNodeEntityType: "unit",
		LinkFieldName:                  "parentUnit",
		DisplayFieldName:               "label",
	})
	icdHierarchyExt, _ := anypb.New(&sharedFields.IndexDefExtHierarchy{
		CodeSystemNameOrNodeEntityType: "icd",
		LinkFieldName:                  "parentCode",
		DisplayFieldName:               "label",
	})

	axisConstraintFieldsRepo := frepo.NewMockedFieldRepo([]fields.BaseFieldDef{
		(&kindstring.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "category", Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
		(&kindstring.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "type", Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
		(&kindstring.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "country",
			Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
		(&kindtext.KindText{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "title", Kind: "text", IndexDef: &sharedFields.IndexDef{}}),
		hierarchy.MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "parentUnit", Kind: "hierarchy", IndexDef: &sharedFields.IndexDef{
			MultiValued: true,
			Ext:         []*anypb.Any{orgHierarchyExt, orgLinkExt},
		}}),
		hierarchy.MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "icdCode", Kind: "hierarchy", IndexDef: &sharedFields.IndexDef{
			MultiValued: true,
			Ext:         []*anypb.Any{icdHierarchyExt, icdLinkExt},
		}}),
	})
	axisConstraintSearchConfigRepo := screpo.NewMockSearchConfigRepo([]*searchconfig.SearchConfigObject{
		{
			Type:   solr.MexSearchFocusType,
			Name:   solr.MexDefaultSearchFocusName,
			Fields: []string{"type"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "typeAxis",
			Fields: []string{"type"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "countryAxis",
			Fields: []string{"countryDate"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "categoryAxis",
			Fields: []string{"category"},
		},
		{
			Type:   solr.MexHierarchyAxisType,
			Name:   "orgHierarchyAxis",
			Fields: []string{"orgUnitCode"},
		},
		{
			Type:   solr.MexHierarchyAxisType,
			Name:   "icdHierarchyAxis",
			Fields: []string{"icdCode"},
		},
	})

	for _, tt := range tests {
		opts := QueryEngineOptions{
			Log:              log,
			FieldRepo:        axisConstraintFieldsRepo,
			SearchConfigRepo: axisConstraintSearchConfigRepo,
			PostQueryHooks:   postQueryHooks,
		}
		qe, _ := newQueryEngine(tt.converter, opts)
		t.Run(tt.name, func(t *testing.T) {
			body, diag, err := qe.CreateSolrQuery(context.TODO(), tt.searchRequest, nil)
			runQueryBodyChecks(body, diag, err, t, tt)
		})
	}
}

func Test_createSolrQueryBody_string_range_axis_constraints(t *testing.T) {
	tests := []QueryTestInfo{
		{
			name: "String range axis constraints: Constraining an axis to a range with both upper and lower value" +
				" leads to a range constraint with both end inclusive & values double-quoted",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type: solr.MexStringRangeConstraint,
						Axis: "categoryAxis",
						StringRanges: []*solr.StringRange{
							{
								Min: "foot locker",
								Max: "head",
							},
						},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetOrdinalAxisFacetAndFilterFieldName("categoryAxis") + ":[\"foot locker\" TO \"head\"]"},
				[]string{"categoryAxis"},
			)},
		},
		{
			name: "String range axis constraints: Constraining an axis to a range with only lower value" +
				" leads to a range constraint no upper constraint, inclusive lower border & values double-quoted",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type: solr.MexStringRangeConstraint,
						Axis: "categoryAxis",
						StringRanges: []*solr.StringRange{
							{
								Min: "foot locker",
							},
						},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetOrdinalAxisFacetAndFilterFieldName("categoryAxis") + ":[\"foot locker\" TO *]"},
				[]string{"categoryAxis"},
			)},
		},
		{
			name: "String range axis constraints: Constraining an axis to a range with only upper value" +
				" leads to a range constraint no lower constraint, inclusive upper border & values double-quoted",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type: solr.MexStringRangeConstraint,
						Axis: "categoryAxis",
						StringRanges: []*solr.StringRange{
							{
								Max: "head",
							},
						},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetOrdinalAxisFacetAndFilterFieldName("categoryAxis") + ":[* TO \"head\"]"},
				[]string{"categoryAxis"},
			)},
		},
		{
			name: "String range axis constraints: Constraining an axis which is absent from the query engine" +
				" configuration causes" +
				" an error",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type: solr.MexStringRangeConstraint,
						Axis: "journalAxis",
						StringRanges: []*solr.StringRange{
							{
								Min: "foot locker",
								Max: "head",
							},
						},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			wantErr:   true,
		},
		{
			name: "String range axis constraints: constraints for several possible values for the same axis are ORed" +
				" together by" +
				" default",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type: solr.MexStringRangeConstraint,
						Axis: "categoryAxis",
						StringRanges: []*solr.StringRange{
							{
								Min: "foot locker",
								Max: "head",
							},
							{
								Min: "nose",
								Max: "toe",
							},
						},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetOrdinalAxisFacetAndFilterFieldName(
					"categoryAxis") + ":[\"foot locker\" TO \"head\"] || " + solr.GetOrdinalAxisFacetAndFilterFieldName(
					"categoryAxis") + ":[\"nose\" TO \"toe\"]"},
				[]string{"categoryAxis"},
			)},
		},
		{
			name: "String range axis constraints: constraints for several possible values for the same axis can be" +
				" explicitly" +
				" set to be ORed together",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type: solr.MexStringRangeConstraint,
						Axis: "categoryAxis",
						StringRanges: []*solr.StringRange{
							{
								Min: "foot locker",
								Max: "head",
							},
							{
								Min: "nose",
								Max: "toe",
							},
						},
						CombineOperator: "or",
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetOrdinalAxisFacetAndFilterFieldName(
					"categoryAxis") + ":[\"foot locker\" TO \"head\"] || " + solr.GetOrdinalAxisFacetAndFilterFieldName(
					"categoryAxis") + ":[\"nose\" TO \"toe\"]"},
				[]string{"categoryAxis"},
			)},
		},
		{
			name: "String range axis constraints: constraints for several possible values for the same axis are can" +
				" explicitly" +
				" be set to be ANDed together",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type: solr.MexStringRangeConstraint,
						Axis: "categoryAxis",
						StringRanges: []*solr.StringRange{
							{
								Min: "foot locker",
								Max: "head",
							},
							{
								Min: "nose",
								Max: "toe",
							},
						},
						CombineOperator: "and",
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetOrdinalAxisFacetAndFilterFieldName(
					"categoryAxis") + ":[\"foot locker\" TO \"head\"] && " + solr.GetOrdinalAxisFacetAndFilterFieldName(
					"categoryAxis") + ":[\"nose\" TO \"toe\"]"},
				[]string{"categoryAxis"},
			)},
		},
		{
			name: "String range axis constraints: unsupported Boolean combination operator for value constraints" +
				" causes error",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type: solr.MexStringRangeConstraint,
						Axis: "categoryAxis",
						StringRanges: []*solr.StringRange{
							{
								Min: "foot locker",
								Max: "head",
							},
							{
								Min: "nose",
								Max: "toe",
							},
						},
						CombineOperator: "xor",
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks:    nil,
			wantErr:   true,
		},
		{
			name: "String range axis constraints: constraints on separate axes are ANDed together",
			searchRequest: &pb.SearchRequest{
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:   solr.MexExactAxisConstraint,
						Axis:   "typeAxis",
						Values: []string{"article", "book"},
					},
					{
						Type: solr.MexStringRangeConstraint,
						Axis: "categoryAxis",
						StringRanges: []*solr.StringRange{
							{
								Min: "foot locker",
								Max: "head",
							},
						},
					},
				},
			},
			converter: &constantConverterNonPhrase,
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetOrdinalAxisFacetAndFilterFieldName("typeAxis") + ":\"article\" || " + solr.
					GetOrdinalAxisFacetAndFilterFieldName("typeAxis") + ":\"book\"",
					solr.GetOrdinalAxisFacetAndFilterFieldName("categoryAxis") + ":[\"foot locker\" TO \"head\"]"},
				[]string{"typeAxis", "categoryAxis"},
			)},
		},
	}

	log := &L.NullLogger{}
	postQueryHooks, _ := hooks.NewPostQueryHooks(hooks.PostQueryHooksConfig{})

	axisConstraintFieldsRepo := frepo.NewMockedFieldRepo([]fields.BaseFieldDef{
		(&kindstring.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "category", Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
		(&kindstring.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "type", Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
		(&kindstring.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "country",
			Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
		(&kindtext.KindText{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "title", Kind: "text", IndexDef: &sharedFields.IndexDef{}}),
	})
	axisConstraintSearchConfigRepo := screpo.NewMockSearchConfigRepo([]*searchconfig.SearchConfigObject{
		{
			Type:   solr.MexSearchFocusType,
			Name:   solr.MexDefaultSearchFocusName,
			Fields: []string{"category"},
		},
		{
			Type:   solr.MexSearchFocusType,
			Name:   "news",
			Fields: []string{"category"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "typeAxis",
			Fields: []string{"type"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "categoryAxis",
			Fields: []string{"category"},
		},
	})

	for _, tt := range tests {
		opts := QueryEngineOptions{
			Log:              log,
			FieldRepo:        axisConstraintFieldsRepo,
			SearchConfigRepo: axisConstraintSearchConfigRepo,
			PostQueryHooks:   postQueryHooks,
		}
		qe, _ := newQueryEngine(tt.converter, opts)
		t.Run(tt.name, func(t *testing.T) {
			body, diag, err := qe.CreateSolrQuery(context.TODO(), tt.searchRequest, nil)
			runQueryBodyChecks(body, diag, err, t, tt)
		})
	}
}

func Test_createSolrQueryBody_search_focus_non_phrase(t *testing.T) {

	tests := []QueryTestInfo{
		{
			name: "Search focus for non-phrase queries: If no search focus is set and N-gram prefixes are NOT used, the boosted default search focus fields WITHOUT fuzzy field are used",
			searchRequest: &pb.SearchRequest{
				Query:         "not relevant since we mock return",
				UseNgramField: false,
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckQuery(DefaultMockedReturnQuery),
			},
			converter: &constantConverterNonPhrase,
		},
		{
			name: "Search focus for non-phrase queries: If no search focus is set and N-gram prefixes ARE used, the boosted default search focus fields INCLUDING boosted prefix field are used",
			searchRequest: &pb.SearchRequest{
				Query:         "not relevant since we mock return",
				UseNgramField: true,
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckQuery(DefaultMockedReturnQuery),
			},
			converter: &constantConverterNonPhrase,
		},
		{
			name: "Search focus for non-phrase queries: If a search focus is set and N-gram prefixes are NOT used, " +
				"the user query run again all Language-specific versions of the boosted search focus field (but NOT the prefix field)",
			searchRequest: &pb.SearchRequest{
				Query:         "not relevant since we mock return",
				SearchFocus:   "news",
				UseNgramField: false,
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckQuery(DefaultMockedReturnQuery),
			},
			converter: &constantConverterNonPhrase,
		},
		{
			name: "Search focus for non-phrase queries: If a search focus is set and N-gram prefixes ARE used, " +
				"the user query run again all Language-specific versions of the boosted search focus field AND the boosted prefix field",
			searchRequest: &pb.SearchRequest{
				Query:         "not relevant since we mock return",
				SearchFocus:   "news",
				UseNgramField: true,
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckQuery(DefaultMockedReturnQuery),
			},
			converter: &constantConverterNonPhrase,
		},
		{
			name: "Search focus for non-phrase queries: If the query is empty, " +
				"the default search focus is used even if another focus is requested",
			searchRequest: &pb.SearchRequest{
				Query:         "not relevant since we mock return",
				SearchFocus:   "news",
				UseNgramField: false,
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckQuery(""),
			},
			converter: &emptyConverter,
		},
		{
			name: "Search focus for non-phrase queries: If both a search focus and axis constraints are set, " +
				"the restriction to the search focus is not applied to the axis constraints",
			searchRequest: &pb.SearchRequest{
				Query:         "not relevant since we mock return",
				SearchFocus:   "news",
				UseNgramField: false,
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:   solr.MexExactFacetType,
						Axis:   "typeAxis",
						Values: []string{"article"},
					},
				},
			},
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetOrdinalAxisFacetAndFilterFieldName("typeAxis") + ":\"article\""},
				[]string{"typeAxis"},
			)},
			converter: &constantConverterNonPhrase,
		},
	}

	log := &L.NullLogger{}
	postQueryHooks, _ := hooks.NewPostQueryHooks(hooks.PostQueryHooksConfig{})

	searchFocusFieldsRepo := frepo.NewMockedFieldRepo([]fields.BaseFieldDef{
		(&kindstring.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "category",
			Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
		(&kindtext.KindText{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "title", Kind: "text", IndexDef: &sharedFields.IndexDef{}}),
		(&kindstring.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "type", Kind: "string",
			IndexDef: &sharedFields.IndexDef{}}),
	})
	focusSearchConfigRepo := screpo.NewMockSearchConfigRepo([]*searchconfig.SearchConfigObject{
		{
			Type:   solr.MexSearchFocusType,
			Name:   solr.MexDefaultSearchFocusName,
			Fields: []string{"category"},
		},
		{
			Type:   solr.MexSearchFocusType,
			Name:   "news",
			Fields: []string{"category"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "typeAxis",
			Fields: []string{"type"},
		},
	})

	for _, tt := range tests {
		opts := QueryEngineOptions{
			Log:              log,
			FieldRepo:        searchFocusFieldsRepo,
			SearchConfigRepo: focusSearchConfigRepo,
			PostQueryHooks:   postQueryHooks,
		}
		qe, _ := newQueryEngine(tt.converter, opts)

		t.Run(tt.name, func(t *testing.T) {
			body, diag, err := qe.CreateSolrQuery(context.TODO(), tt.searchRequest, nil)
			runQueryBodyChecks(body, diag, err, t, tt)
		})
	}
}

func Test_createSolrQueryBody_search_focus_phrase(t *testing.T) {
	tests := []QueryTestInfo{
		{
			name: "Search focus for phrase queries: If no search focus is set and N-gram prefixes are NOT used, the default unanalyzed search focus field is used",
			searchRequest: &pb.SearchRequest{
				Query:         "not relevant since we mock return",
				UseNgramField: false,
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckQuery(DefaultMockedReturnQuery),
			},
			converter: &constantConverterPhrase,
		},
		{
			name: "Search focus for phrase queries: If no search focus is set and N-gram prefixes ARE used, the default unanalyzed search focus field is used",
			searchRequest: &pb.SearchRequest{
				Query:         "not relevant since we mock return",
				UseNgramField: true,
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckQuery(DefaultMockedReturnQuery),
			},
			converter: &constantConverterPhrase,
		},
		{
			name: "Search focus for phrase queries: If a search focus is set and N-gram prefixes are NOT used, the unanalyzed search focus field is used",
			searchRequest: &pb.SearchRequest{
				Query:         "not relevant since we mock return",
				SearchFocus:   "news",
				UseNgramField: false,
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckQuery(DefaultMockedReturnQuery),
			},
			converter: &constantConverterPhrase,
		},
		{
			name: "Search focus for phrase queries: If a search focus is set and N-gram prefixes ARE used, the unanalyzed search focus field is used",
			searchRequest: &pb.SearchRequest{
				Query:         "not relevant since we mock return",
				SearchFocus:   "news",
				UseNgramField: true,
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckQuery(DefaultMockedReturnQuery),
			},
			converter: &constantConverterPhrase,
		},
		{
			name: "Search focus for phrase queries: If both a search focus and axis constraints are set, " +
				"the restriction to the search focus is not applied to the axis constraints",
			searchRequest: &pb.SearchRequest{
				Query:         "not relevant since we mock return",
				SearchFocus:   "news",
				UseNgramField: false,
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:   solr.MexExactFacetType,
						Axis:   "typeAxis",
						Values: []string{"article"},
					},
				},
			},
			checks: &[]testutils.BodyCheck{testutils.CheckConstraints(
				[]string{solr.GetOrdinalAxisFacetAndFilterFieldName("typeAxis") + ":\"article\""},
				[]string{"typeAxis"},
			)},
			converter: &constantConverterPhrase,
		},
	}

	log := &L.NullLogger{}
	postQueryHooks, _ := hooks.NewPostQueryHooks(hooks.PostQueryHooksConfig{})

	searchFocusFieldsRepo := frepo.NewMockedFieldRepo([]fields.BaseFieldDef{
		(&kindstring.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "category",
			Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
		(&kindtext.KindText{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "title", Kind: "text", IndexDef: &sharedFields.IndexDef{}}),
		(&kindstring.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "type", Kind: "string",
			IndexDef: &sharedFields.IndexDef{}}),
	})
	focusSearchConfigRepo := screpo.NewMockSearchConfigRepo([]*searchconfig.SearchConfigObject{
		{
			Type:   solr.MexSearchFocusType,
			Name:   solr.MexDefaultSearchFocusName,
			Fields: []string{"category"},
		},
		{
			Type:   solr.MexSearchFocusType,
			Name:   "news",
			Fields: []string{"category"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "typeAxis",
			Fields: []string{"type"},
		},
	})

	for _, tt := range tests {
		opts := QueryEngineOptions{
			Log:              log,
			FieldRepo:        searchFocusFieldsRepo,
			SearchConfigRepo: focusSearchConfigRepo,
			PostQueryHooks:   postQueryHooks,
		}
		qe, _ := newQueryEngine(tt.converter, opts)

		t.Run(tt.name, func(t *testing.T) {
			body, diag, err := qe.CreateSolrQuery(context.TODO(), tt.searchRequest, nil)
			runQueryBodyChecks(body, diag, err, t, tt)
		})
	}
}

func Test_createSolrQueryBody_fields(t *testing.T) {
	testQuery := "sunday"

	var genericTitleFieldName, _ = solr.GetLangSpecificFieldName(titleName, solr.GenericLangAbbrev)
	var deTitleFieldName, _ = solr.GetLangSpecificFieldName(titleName, solr.GermanLangAbbrev)
	var enTitleFieldName, _ = solr.GetLangSpecificFieldName(titleName, solr.EnglishLangAbbrev)
	textFieldExample := []string{"id", "entityName", genericTitleFieldName, deTitleFieldName, enTitleFieldName}

	tests := []QueryTestInfo{
		{
			name: "Fields: If the 'fields' property is missing, , only 'id' and `'entityName' is returned",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFields([]string{"id", "entityName"}),
			},
		},
		{
			name: "Fields: If the array of requested fields is empty, only 'id' and `'entityName'  is returned",
			searchRequest: &pb.SearchRequest{
				Query:  testQuery,
				Fields: []string{},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFields([]string{"id", "entityName"}),
			},
		},
		{
			name: "Fields: If the fields 'id' or 'entityName' are not requested, they are automatically added",
			searchRequest: &pb.SearchRequest{
				Query:  testQuery,
				Fields: []string{"category"},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFields([]string{"category", "id", "entityName"}),
			},
		},
		{
			name: "Fields: All requested fields are added",
			searchRequest: &pb.SearchRequest{
				Query:  testQuery,
				Fields: []string{"tag", "category"},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFields([]string{"tag", "category", "id", "entityName"}),
			},
		},
		{
			name: "Fields: If the field 'id' is explicitly added, it is not added again",
			searchRequest: &pb.SearchRequest{
				Query:  testQuery,
				Fields: []string{"id", "category"},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFields([]string{"id", "category", "entityName"}),
			},
		},
		{
			name: "Fields: If the field 'entityName' is explicitly added, it is not added again",
			searchRequest: &pb.SearchRequest{
				Query:  testQuery,
				Fields: []string{"entityName", "category"},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFields([]string{"entityName", "category", "id"}),
			},
		},
		{
			name: "Fields: For text fields, the language-specific backing fields are added",
			searchRequest: &pb.SearchRequest{
				Query:  testQuery,
				Fields: []string{"title"},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFields(textFieldExample),
			},
		},
	}

	log := &L.NullLogger{}
	postQueryHooks, _ := hooks.NewPostQueryHooks(hooks.PostQueryHooksConfig{})

	mockedRepo := frepo.NewMockedFieldRepo([]fields.BaseFieldDef{
		(&kindstring.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "category",
			Kind:     "string",
			IndexDef: &sharedFields.IndexDef{}}),
		(&kindstring.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "tag",
			Kind:     "string",
			IndexDef: &sharedFields.IndexDef{}}),
		(&kindtext.KindText{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "title",
			Kind:     "text",
			IndexDef: &sharedFields.IndexDef{}}),
	})
	opts := QueryEngineOptions{
		Log:              log,
		FieldRepo:        mockedRepo,
		SearchConfigRepo: getStandardSearchConfigRepo([]string{"category"}),
		PostQueryHooks:   postQueryHooks,
	}
	qe, _ := newQueryEngine(&constantConverterNonPhrase, opts)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, diag, err := qe.CreateSolrQuery(context.TODO(), tt.searchRequest, nil)
			runQueryBodyChecks(body, diag, err, t, tt)
		})
	}
}

func Test_createSolrQueryBody_highlighting(t *testing.T) {
	testQuery := "sunday"
	fieldsToHighlightForAbstract, _ := mapToLanguageSpecificFieldNames("abstract")
	fieldsToHighlightForAbstract = append(fieldsToHighlightForAbstract, solr.GetPrefixBackingFieldName("abstract"), solr.GetRawBackingFieldName("abstract"))
	fieldsToDefaultAutoHighlight := []string{"category"}
	fieldsToHighlightForTitle, _ := mapToLanguageSpecificFieldNames("title")
	fieldsToHighlightForTitle = append(fieldsToHighlightForTitle, solr.GetPrefixBackingFieldName("title"), solr.GetRawBackingFieldName("title"))
	fieldsToDefaultAutoHighlight = append(fieldsToDefaultAutoHighlight, fieldsToHighlightForTitle...)
	fieldsToAutoHighlightSearchFocus := []string{"type"}
	fieldsToAutoHighlightSearchFocus = append(fieldsToAutoHighlightSearchFocus, fieldsToHighlightForAbstract...)

	tests := []QueryTestInfo{
		{
			name: "Highlighting: Highlight fields are not mapped into query string",
			searchRequest: &pb.SearchRequest{
				Query:           testQuery,
				HighlightFields: []string{"type"},
			},
			checks: &[]testutils.BodyCheck{testutils.CheckQuery(DefaultMockedReturnQuery)},
		},
		{
			name: "Highlighting: Highlighting is off by default",
			searchRequest: &pb.SearchRequest{
				Query: testQuery,
			},
			checks: &[]testutils.BodyCheck{testutils.CheckHighlightingFields(nil)},
		},
		{
			name: "Highlighting: Default highlighting parameters are set",
			searchRequest: &pb.SearchRequest{
				Query:           testQuery,
				HighlightFields: []string{"type"},
			},
			checks: &[]testutils.BodyCheck{testutils.CheckHighlightingParameters("unified", 10, 100, "\ue000", "\ue001")},
		},
		{
			name: "Highlighting: Single highlight string field is mapped to params",
			searchRequest: &pb.SearchRequest{
				Query:           testQuery,
				HighlightFields: []string{"type"},
			},
			checks: &[]testutils.BodyCheck{testutils.CheckHighlightingFields([]string{"type"})},
		},
		{
			name: "Highlighting: Multiple highlight string fields are mapped to params",
			searchRequest: &pb.SearchRequest{
				Query:           testQuery,
				HighlightFields: []string{"type", "category"},
			},
			checks: &[]testutils.BodyCheck{testutils.CheckHighlightingFields([]string{"type", "category"})},
		},
		{
			name: "Highlighting: Highlight on text field leads to highlighting on all backing fields",
			searchRequest: &pb.SearchRequest{
				Query:           testQuery,
				HighlightFields: []string{"abstract"},
			},
			checks: &[]testutils.BodyCheck{testutils.CheckHighlightingFields(fieldsToHighlightForAbstract)},
		},
		{
			name: "Highlighting: Auto-highlight with no search focus set leads to highlighting of all" +
				" default searchable fields (with text field mapped to the generic-language field)",
			searchRequest: &pb.SearchRequest{
				Query:         testQuery,
				AutoHighlight: true,
			},
			checks: &[]testutils.BodyCheck{testutils.CheckHighlightingFields(fieldsToDefaultAutoHighlight)},
		},
		{
			name: "Highlighting: Auto-highlight with a search focus set leads to highlighting of all" +
				" search focus fields (with text field mapped to the generic-language field)",
			searchRequest: &pb.SearchRequest{
				Query:         testQuery,
				SearchFocus:   "testFocus",
				AutoHighlight: true,
			},
			checks: &[]testutils.BodyCheck{testutils.CheckHighlightingFields(fieldsToAutoHighlightSearchFocus)},
		},
	}

	log := &L.NullLogger{}
	postQueryHooks, _ := hooks.NewPostQueryHooks(hooks.PostQueryHooksConfig{})

	mockedRepo := frepo.NewMockedFieldRepo([]fields.BaseFieldDef{
		(&kindstring.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "type", Kind: "string",
			IndexDef: &sharedFields.IndexDef{}}),
		(&kindstring.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "category",
			Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
		(&kindtext.KindText{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "abstract",
			Kind:     "text",
			IndexDef: &sharedFields.IndexDef{}}),
		(&kindtext.KindText{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "title",
			Kind:     "text",
			IndexDef: &sharedFields.IndexDef{}}),
	})
	highlightSearchConfigRepo := screpo.NewMockSearchConfigRepo([]*searchconfig.SearchConfigObject{
		{
			Type: solr.MexSearchFocusType,
			Name: "default",
			Fields: []string{
				"category",
				"title",
			},
		},
		{
			Type: solr.MexSearchFocusType,
			Name: "testFocus",
			Fields: []string{
				"type",
				"abstract",
			},
		},
	})
	opts := QueryEngineOptions{
		Log:              log,
		FieldRepo:        mockedRepo,
		SearchConfigRepo: highlightSearchConfigRepo,
		PostQueryHooks:   postQueryHooks,
	}
	qe, _ := newQueryEngine(&constantConverterNonPhrase, opts)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, diag, err := qe.CreateSolrQuery(context.TODO(), tt.searchRequest, nil)
			runQueryBodyChecks(body, diag, err, t, tt)
		})
	}
}

func Test_CreateFacetName(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		statName  string
		want      string
	}{
		{
			name:      "Adds a pre-fix to the input if no stat name is given",
			fieldName: "hello",
			want:      solr.FacetPrefix + "_hello",
		},
		{
			name:      "Uses stat name directly if given",
			fieldName: "hello",
			statName:  "statHello",
			want:      "statHello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createFacetName(tt.fieldName, tt.statName)
			if got != tt.want {
				t.Errorf("createFacetName(): got '%s' wanted '%s'", got, tt.want)
			}
		})
	}
}

func Test_CreateStatExpression(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		op        string
		want      string
		wantErr   bool
	}{
		{
			name:      "Work for min operator",
			fieldName: "hello",
			op:        solr.MinOperator,
			want:      "min(hello)",
		},
		{
			name:      "Work for max operator",
			fieldName: "hello",
			op:        solr.MaxOperator,
			want:      "max(hello)",
		},
		{
			name:      "Throws error for unknown operator",
			fieldName: "hello",
			op:        "magicOp",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := solr.CreateStatExpression(tt.fieldName, tt.op)
			if tt.wantErr != (err != nil) {
				t.Errorf("createFacetName(): expected error but got none")
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("createFacetName(): got '%s' wanted '%s'", got, tt.want)
			}
		})
	}
}

func Test_getStatNameForFieldAndOp(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		opName    string
		want      string
		wantErr   bool
	}{
		{
			name:      "Returns error if the op contains a low-dash",
			fieldName: "author",
			opName:    "some_op",
			wantErr:   true,
		},
		{
			name:      "Returns error if the field name is empty",
			fieldName: "",
			opName:    "min",
			wantErr:   true,
		},
		{
			name:      "Returns error if the op name is empty",
			fieldName: "author",
			opName:    "",
			wantErr:   true,
		},
		{
			name:      "Return the field name prefixed with the op name and a low-dash",
			fieldName: "author",
			opName:    "min",
			want:      "min_author",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getStatNameForAxisAndOp(tt.fieldName, tt.opName)
			if tt.wantErr != (err != nil) {
				t.Errorf("Unexpected error response - wanted error = %v, error returned = %v", tt.wantErr, err != nil)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("createFacetName(): got '%s' wanted '%s'", got, tt.want)
			}
		})
	}
}

func Test_getFieldAndOpFromStatName(t *testing.T) {
	tests := []struct {
		name      string
		statName  string
		wantOp    string
		wantField string
		wantErr   bool
	}{
		{
			name:     "Returns error if stat name does not contain a low-dash",
			statName: "nameWithNoDash",
			wantErr:  true,
		},
		{
			name:      "splits off the prefixed operator name",
			statName:  "opName_fieldName",
			wantField: "fieldName",
			wantOp:    "opName",
		},
		{
			name:      "if there are multiple low-dashes, all but the first are treated as part of the field name",
			statName:  "opName_fieldName_with_dashes",
			wantField: "fieldName_with_dashes",
			wantOp:    "opName",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotField, gotOp, err := getSolrFieldAndOpFromStatName(tt.statName)
			if tt.wantErr != (err != nil) {
				t.Errorf("Unexpected error response - wanted error = %v, error returned = %v", tt.wantErr, err != nil)
			}
			if !tt.wantErr && (gotField != tt.wantField || gotOp != tt.wantOp) {
				t.Errorf("createFacetName(): (field name, op) - wanted (%s, %s) but got (%s, %s)", tt.wantField,
					tt.wantOp, gotField, gotOp)
			}
		})
	}
}

func Test_GetRangeStatRequestFacets(t *testing.T) {
	pubMinFacetName, _ := getStatNameForAxisAndOp("publishingDateAxis", "min")
	pubMaxFacetName, _ := getStatNameForAxisAndOp("publishingDateAxis", "max")
	creMinFacetName, _ := getStatNameForAxisAndOp("createdAtAxis", "min")
	creMaxFacetName, _ := getStatNameForAxisAndOp("createdAtAxis", "max")

	tests := []struct {
		name                   string
		searchRequest          *pb.SearchRequest
		wantNoConstraintFacets map[string][]*solr.Facet
		wantConstrainedFacets  map[string][]*solr.Facet
		wantErr                bool
	}{
		{
			name: "Returns empty slices if no facets are set at all",
			searchRequest: &pb.SearchRequest{
				Query: "*",
			},
			wantNoConstraintFacets: make(map[string][]*solr.Facet),
			wantConstrainedFacets:  make(map[string][]*solr.Facet),
		},
		{
			name: "Returns empty slices if the requested facets are not year-range facets",
			searchRequest: &pb.SearchRequest{
				Query: "*",
				Facets: []*solr.Facet{
					{
						Type: solr.MexExactFacetType,
						Axis: "category",
					},
					{
						Type:   solr.MexExactFacetType,
						Axis:   "titleAxis",
						Offset: 15,
					},
				},
			},
			wantNoConstraintFacets: make(map[string][]*solr.Facet),
			wantConstrainedFacets:  make(map[string][]*solr.Facet),
		},
		{
			name: "Returns min & max facets for the fields on which a year-range facet is set",
			searchRequest: &pb.SearchRequest{
				Query: "*",
				Facets: []*solr.Facet{
					{
						Type: solr.MexYearRangeFacetType,
						Axis: "publishingDateAxis",
					},
					{
						Type:   solr.MexExactFacetType,
						Axis:   "titleAxis",
						Offset: 15,
					}, {
						Type: solr.MexYearRangeFacetType,
						Axis: "createdAtAxis",
					},
				},
			},
			wantNoConstraintFacets: map[string][]*solr.Facet{
				"publishingDateAxis": {
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "publishingDateAxis",
						StatName: pubMinFacetName,
						StatOp:   "min",
					},
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "publishingDateAxis",
						StatName: pubMaxFacetName,
						StatOp:   "max",
					},
				},
				"createdAtAxis": {
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "createdAtAxis",
						StatName: creMinFacetName,
						StatOp:   "min",
					},
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "createdAtAxis",
						StatName: creMaxFacetName,
						StatOp:   "max",
					},
				},
			},
			wantConstrainedFacets: make(map[string][]*solr.Facet),
		},
		{
			name: "For fields that have an axis constraint, the min & max facets for the fields are returned separately",
			searchRequest: &pb.SearchRequest{
				Query: "*",
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type: solr.MexStringRangeConstraint,
						Axis: "publishingDateAxis",
						StringRanges: []*solr.StringRange{
							{
								Min: "2018-05-20T17:33:18",
								Max: "2020-05-20T17:33:18",
							},
						},
					},
				},
				Facets: []*solr.Facet{
					{
						Type: solr.MexYearRangeFacetType,
						Axis: "publishingDateAxis",
					},
					{
						Type:   solr.MexExactFacetType,
						Axis:   "titleAxis",
						Offset: 15,
					}, {
						Type: solr.MexYearRangeFacetType,
						Axis: "createdAtAxis",
					},
				},
			},
			wantNoConstraintFacets: map[string][]*solr.Facet{
				"createdAtAxis": {
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "createdAtAxis",
						StatName: creMinFacetName,
						StatOp:   "min",
					},
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "createdAtAxis",
						StatName: creMaxFacetName,
						StatOp:   "max",
					},
				},
			},
			wantConstrainedFacets: map[string][]*solr.Facet{
				"publishingDateAxis": {
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "publishingDateAxis",
						StatName: pubMinFacetName,
						StatOp:   "min",
					},
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "publishingDateAxis",
						StatName: pubMaxFacetName,
						StatOp:   "max",
					},
				},
			},
		},
		{
			name: "Axis constraints on fields on which no year-range facet is requested do not change the result",
			searchRequest: &pb.SearchRequest{
				Query: "*",
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type: solr.MexStringRangeConstraint,
						Axis: "someOtherDate",
						StringRanges: []*solr.StringRange{
							{
								Min: "2018-05-20T17:33:18",
								Max: "2020-05-20T17:33:18",
							},
						},
					},
				},
				Facets: []*solr.Facet{
					{
						Type: solr.MexYearRangeFacetType,
						Axis: "publishingDateAxis",
					},
					{
						Type:   solr.MexExactFacetType,
						Axis:   "titleAxis",
						Offset: 15,
					}, {
						Type: solr.MexYearRangeFacetType,
						Axis: "createdAtAxis",
					},
				},
			},
			wantNoConstraintFacets: map[string][]*solr.Facet{
				"publishingDateAxis": {
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "publishingDateAxis",
						StatName: pubMinFacetName,
						StatOp:   "min",
					},
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "publishingDateAxis",
						StatName: pubMaxFacetName,
						StatOp:   "max",
					},
				},
				"createdAtAxis": {
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "createdAtAxis",
						StatName: creMinFacetName,
						StatOp:   "min",
					},
					{
						Type:     solr.MexStringStatFacetType,
						Axis:     "createdAtAxis",
						StatName: creMaxFacetName,
						StatOp:   "max",
					},
				},
			},
			wantConstrainedFacets: make(map[string][]*solr.Facet),
		},
		{
			name: "Returns an error if there is a year-range facet without specified field",
			searchRequest: &pb.SearchRequest{
				Query: "*",
				Facets: []*solr.Facet{
					{
						Type: solr.MexYearRangeFacetType,
						Axis: "publishingDateAxis",
					},
					{
						Type:   solr.MexExactFacetType,
						Axis:   "titleAxis",
						Offset: 15,
					}, {
						Type: solr.MexYearRangeFacetType,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNoConstraint, gotConstrained, err := GetRangeStatRequestFacets(tt.searchRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetYearRangeFacetFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(gotNoConstraint, tt.wantNoConstraintFacets) {
				t.Errorf("GetYearRangeFacetFields() - facets with no matching axis constraint: got = %v, want %v",
					gotNoConstraint, tt.wantNoConstraintFacets)
			}
			if !tt.wantErr && !reflect.DeepEqual(gotConstrained, tt.wantConstrainedFacets) {
				t.Errorf("GetYearRangeFacetFields() - facet with matching axis constraint: got = %v, want %v", gotConstrained,
					tt.wantConstrainedFacets)
			}
		})
	}
}

func TestQueryEngine_CreateYearRangeQuery(t *testing.T) {
	pubMinFacetName, _ := getStatNameForAxisAndOp("publishingDateAxis", "min")
	pubMaxFacetName, _ := getStatNameForAxisAndOp("publishingDateAxis", "max")
	creMinFacetName, _ := getStatNameForAxisAndOp("createdAtAxis", "min")
	creMaxFacetName, _ := getStatNameForAxisAndOp("createdAtAxis", "max")

	tests := []QueryTestInfo{
		{
			name: "Causes error if there are no year-range facets",
			searchRequest: &pb.SearchRequest{
				Query: "IRRELEVANT SINCE MOCKED",
				Facets: []*solr.Facet{
					{
						Type:   solr.MexExactFacetType,
						Axis:   "titleAxis",
						Offset: 15,
					},
				},
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:   solr.MexExactAxisConstraint,
						Axis:   "titleAxis",
						Values: []string{"article and more"},
					},
				},
			},
			yearRangeFacets: nil,
			wantErr:         true,
		},
		{
			name: "Applies query and and all axis constraints if no (non-empty) ignore field is passed",
			searchRequest: &pb.SearchRequest{
				Query: "IRRELEVANT SINCE MOCKED",
				Facets: []*solr.Facet{
					{
						Type: solr.MexYearRangeFacetType,
						Axis: "publishingDateAxis",
					},
					{
						Type:   solr.MexExactFacetType,
						Axis:   "titleAxis",
						Offset: 15,
					},
				},
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:   solr.MexExactAxisConstraint,
						Axis:   "titleAxis",
						Values: []string{"article and more"},
					},
				},
			},
			yearRangeFacets: []*solr.Facet{
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMinFacetName,
					StatOp:   "min",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMaxFacetName,
					StatOp:   "max",
				},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckQuery(DefaultMockedReturnQuery),
				testutils.CheckConstraints(
					[]string{solr.GetOrdinalAxisFacetAndFilterFieldName("titleAxis") + ":\"article and more\""},
					[]string{"titleAxis"},
				),
			},
		},
		{
			name: "Also applies single-value constraint to hierarchy axes if no (non-empty) ignore field is passed",
			searchRequest: &pb.SearchRequest{
				Query: "IRRELEVANT SINCE MOCKED",
				Facets: []*solr.Facet{
					{
						Type: solr.MexYearRangeFacetType,
						Axis: "publishingDateAxis",
					},
					{
						Type:   solr.MexExactFacetType,
						Axis:   "titleAxis",
						Offset: 15,
					},
				},
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:             solr.MexExactAxisConstraint,
						Axis:             "orgHierarchyAxis",
						SingleNodeValues: []string{"article and more"},
					},
				},
			},
			yearRangeFacets: []*solr.Facet{
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMinFacetName,
					StatOp:   "min",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMaxFacetName,
					StatOp:   "max",
				},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckQuery(DefaultMockedReturnQuery),
				testutils.CheckConstraints(
					[]string{solr.GetSingleNodeAxisFieldName(solr.GetOrdinalAxisFacetAndFilterFieldName("orgHierarchyAxis")) + ":\"article and more\""},
					[]string{"orgHierarchyAxis"},
				),
			},
		},
		{
			name: "Applies the facets passed if no (non-empty) ignore field is passed",
			searchRequest: &pb.SearchRequest{
				Query: "IRRELEVANT SINCE MOCKED",
				Facets: []*solr.Facet{
					{
						Type: solr.MexYearRangeFacetType,
						Axis: "publishingDateAxis",
					},
					{
						Type:   solr.MexExactFacetType,
						Axis:   "titleAxis",
						Offset: 15,
					}, {
						Type: solr.MexYearRangeFacetType,
						Axis: "createdAtAxis",
					},
				},
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:   solr.MexExactAxisConstraint,
						Axis:   "titleAxis",
						Values: []string{"article and more"},
					},
				},
			},
			yearRangeFacets: []*solr.Facet{
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMinFacetName,
					StatOp:   "min",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMaxFacetName,
					StatOp:   "max",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "createdAtAxis",
					Limit:    0,
					StatName: creMinFacetName,
					StatOp:   "min",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "createdAtAxis",
					Limit:    0,
					StatName: creMaxFacetName,
					StatOp:   "max",
				},
			},
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(map[string]solr.SolrFacet{
					pubMinFacetName: {
						DetailedType: solr.SolrStringStatFacetType,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("publishingDateAxis"),
						StatOp:       "min",
					},
					pubMaxFacetName: {
						DetailedType: solr.SolrStringStatFacetType,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("publishingDateAxis"),
						StatOp:       "max",
					},
					creMinFacetName: {
						DetailedType: solr.SolrStringStatFacetType,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("createdAtAxis"),
						StatOp:       "min",
					},
					creMaxFacetName: {
						DetailedType: solr.SolrStringStatFacetType,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("createdAtAxis"),
						StatOp:       "max",
					},
				}),
			},
		},
		{
			name: "Applies the facets passed if a non-empty ignore field is passed",
			searchRequest: &pb.SearchRequest{
				Query: "IRRELEVANT SINCE MOCKED",
				Facets: []*solr.Facet{
					{
						Type: solr.MexYearRangeFacetType,
						Axis: "publishingDateAxis",
					},
					{
						Type:   solr.MexExactFacetType,
						Axis:   "titleAxis",
						Offset: 15,
					}, {
						Type: solr.MexYearRangeFacetType,
						Axis: "createdAtAxis",
					},
				},
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:   solr.MexExactAxisConstraint,
						Axis:   "titleAxis",
						Values: []string{"article and more"},
					},
				},
			},
			yearRangeFacets: []*solr.Facet{
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMinFacetName,
					StatOp:   "min",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMaxFacetName,
					StatOp:   "max",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "createdAtAxis",
					Limit:    0,
					StatName: creMinFacetName,
					StatOp:   "min",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "createdAtAxis",
					Limit:    0,
					StatName: creMaxFacetName,
					StatOp:   "max",
				},
			},
			constraintIgnoreAxis: "createdAtAxis",
			checks: &[]testutils.BodyCheck{
				testutils.CheckFacets(map[string]solr.SolrFacet{
					pubMinFacetName: {
						DetailedType: solr.SolrStringStatFacetType,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("publishingDateAxis"),
						StatOp:       "min",
					},
					pubMaxFacetName: {
						DetailedType: solr.SolrStringStatFacetType,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("publishingDateAxis"),
						StatOp:       "max",
					},
					creMinFacetName: {
						DetailedType: solr.SolrStringStatFacetType,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("createdAtAxis"),
						StatOp:       "min",
					},
					creMaxFacetName: {
						DetailedType: solr.SolrStringStatFacetType,
						Field:        solr.GetOrdinalAxisFacetAndFilterFieldName("createdAtAxis"),
						StatOp:       "max",
					},
				}),
			},
		},
		{
			name: "Applies the query constraint if a non-empty ignore field is passed",
			searchRequest: &pb.SearchRequest{
				Query: "IRRELEVANT SINCE MOCKED",
				Facets: []*solr.Facet{
					{
						Type: solr.MexYearRangeFacetType,
						Axis: "publishingDateAxis",
					},
					{
						Type:   solr.MexExactFacetType,
						Axis:   "titleAxis",
						Offset: 15,
					}, {
						Type: solr.MexYearRangeFacetType,
						Axis: "createdAtAxis",
					},
				},
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type:   solr.MexExactAxisConstraint,
						Axis:   "titleAxis",
						Values: []string{"article and more"},
					},
				},
			},
			yearRangeFacets: []*solr.Facet{
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMinFacetName,
					StatOp:   "min",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMaxFacetName,
					StatOp:   "max",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "createdAtAxis",
					Limit:    0,
					StatName: creMinFacetName,
					StatOp:   "min",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "createdAtAxis",
					Limit:    0,
					StatName: creMaxFacetName,
					StatOp:   "max",
				},
			},
			constraintIgnoreAxis: "createdAtAxis",
			checks:               &[]testutils.BodyCheck{testutils.CheckQuery(DefaultMockedReturnQuery)},
		},
		{
			name: "Applies the facets constraints only for non-ignored axes if a non-empty ignore axis" +
				" is passed",
			searchRequest: &pb.SearchRequest{
				Query: "IRRELEVANT SINCE MOCKED",
				Facets: []*solr.Facet{
					{
						Type: solr.MexYearRangeFacetType,
						Axis: "publishingDateAxis",
					},
					{
						Type:   solr.MexExactFacetType,
						Axis:   "titleAxis",
						Offset: 15,
					}, {
						Type: solr.MexYearRangeFacetType,
						Axis: "createdAtAxis",
					},
				},
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type: solr.MexYearRangeFacetType,
						Axis: "createdAtAxis",
						StringRanges: []*solr.StringRange{
							{
								Min: "2011-01-01T00:00:00Z",
								Max: "2019-01-01T00:00:00Z",
							},
						},
					},
					{
						Type:   solr.MexExactAxisConstraint,
						Axis:   "titleAxis",
						Values: []string{"article and more"},
					},
				},
			},
			yearRangeFacets: []*solr.Facet{
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMinFacetName,
					StatOp:   "min",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMaxFacetName,
					StatOp:   "max",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "createdAtAxis",
					Limit:    0,
					StatName: creMinFacetName,
					StatOp:   "min",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "createdAtAxis",
					Limit:    0,
					StatName: creMaxFacetName,
					StatOp:   "max",
				},
			},
			constraintIgnoreAxis: "createdAtAxis",
			checks: &[]testutils.BodyCheck{
				testutils.CheckConstraints(
					[]string{solr.GetOrdinalAxisFacetAndFilterFieldName("titleAxis") + ":\"article and more\""},
					[]string{"titleAxis"},
				),
			},
		},
		{
			name: "Also applies single-value hierary facets constraints only for non-ignored axes if a non-empty ignore axis" +
				" is passed",
			searchRequest: &pb.SearchRequest{
				Query: "IRRELEVANT SINCE MOCKED",
				Facets: []*solr.Facet{
					{
						Type: solr.MexYearRangeFacetType,
						Axis: "publishingDateAxis",
					},
					{
						Type:   solr.MexExactFacetType,
						Axis:   "titleAxis",
						Offset: 15,
					},
					{
						Type: solr.MexYearRangeFacetType,
						Axis: "createdAtAxis",
					},
				},
				AxisConstraints: []*solr.AxisConstraint{
					{
						Type: solr.MexYearRangeFacetType,
						Axis: "createdAtAxis",
						StringRanges: []*solr.StringRange{
							{
								Min: "2011-01-01T00:00:00Z",
								Max: "2019-01-01T00:00:00Z",
							},
						},
					},
					{
						Type:             solr.MexExactAxisConstraint,
						Axis:             "orgHierarchyAxis",
						SingleNodeValues: []string{"article and more"},
					},
				},
			},
			yearRangeFacets: []*solr.Facet{
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMinFacetName,
					StatOp:   "min",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMaxFacetName,
					StatOp:   "max",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "createdAtAxis",
					Limit:    0,
					StatName: creMinFacetName,
					StatOp:   "min",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "createdAtAxis",
					Limit:    0,
					StatName: creMaxFacetName,
					StatOp:   "max",
				},
			},
			constraintIgnoreAxis: "createdAtAxis",
			checks: &[]testutils.BodyCheck{
				testutils.CheckConstraints(
					[]string{solr.GetSingleNodeAxisFieldName(solr.GetOrdinalAxisFacetAndFilterFieldName("orgHierarchyAxis")) + ":\"article and more\""},
					[]string{"orgHierarchyAxis"},
				),
			},
		},
	}

	orgLinkExt, _ := anypb.New(&sharedFields.IndexDefExtLink{
		RelationType: "relationParent",
	})
	hierarchy, _ := kindhierarchy.NewKindHierarchy(nil)
	orgHierarchyExt, _ := anypb.New(&sharedFields.IndexDefExtHierarchy{
		CodeSystemNameOrNodeEntityType: "unit",
		LinkFieldName:                  "parentUnit",
		DisplayFieldName:               "label",
	})

	facetQueryFieldsRepo := frepo.NewMockedFieldRepo([]fields.BaseFieldDef{
		(&kindstring.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "category", Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
		(&kindtext.KindText{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "title", Kind: "text", IndexDef: &sharedFields.IndexDef{}}),
		(&kindtimestamp.KindTimestamp{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "publishingDate",
			Kind: "timestamp", IndexDef: &sharedFields.IndexDef{}}),
		hierarchy.MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "orgUnit", Kind: "hierarchy", IndexDef: &sharedFields.IndexDef{
			MultiValued: true,
			Ext:         []*anypb.Any{orgHierarchyExt, orgLinkExt},
		}}),
	})
	searchConfigRepo := screpo.NewMockSearchConfigRepo([]*searchconfig.SearchConfigObject{
		{
			Type:   solr.MexSearchFocusType,
			Name:   solr.MexDefaultSearchFocusName,
			Fields: []string{"title"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "titleAxis",
			Fields: []string{"title"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "publishingDateAxis",
			Fields: []string{"publishingDate"},
		},
		{
			Type:   solr.MexOrdinalAxisType,
			Name:   "createdAtAxis",
			Fields: []string{"createdAt"},
		},
		{
			Type:   solr.MexHierarchyAxisType,
			Name:   "orgHierarchyAxis",
			Fields: []string{"orgUnitCode"},
		},
	})

	log := &L.NullLogger{}
	postQueryHooks, _ := hooks.NewPostQueryHooks(hooks.PostQueryHooksConfig{})
	opts := QueryEngineOptions{
		Log:              log,
		FieldRepo:        facetQueryFieldsRepo,
		SearchConfigRepo: searchConfigRepo,
		PostQueryHooks:   postQueryHooks,
	}
	qe, _ := newQueryEngine(&constantConverterNonPhrase, opts)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _, err := qe.CreateYearRangeQuery(context.TODO(), tt.searchRequest, tt.yearRangeFacets, tt.constraintIgnoreAxis)
			runQueryBodyChecks(body, nil, err, t, tt)
		})
	}
}

func Test_GetDateFieldRangesFromResponse(t *testing.T) {
	pubMinFacetName, _ := getStatNameForAxisAndOp("publishingDateAxis", "min")
	pubMaxFacetName, _ := getStatNameForAxisAndOp("publishingDateAxis", "max")
	creMinFacetName, _ := getStatNameForAxisAndOp("createdAtAxis", "min")
	creMaxFacetName, _ := getStatNameForAxisAndOp("createdAtAxis", "max")

	tests := []struct {
		name         string
		solrResponse *solr.QueryResponse
		facets       []*solr.Facet
		want         solr.StringFieldRanges
		wantErr      bool
	}{
		{
			name: "maps all returned facets to min-max range for the corresponding axes",
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      14,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Facets: map[string]interface{}{
					pubMinFacetName: "2010-01-05T00:00:00Z",
					pubMaxFacetName: "2020-02-06T00:00:00Z",
					creMinFacetName: "2011-03-07T00:00:00Z",
					creMaxFacetName: "2021-04-08T00:00:00Z",
				},
			},
			facets: []*solr.Facet{
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMinFacetName,
					StatOp:   "min",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMaxFacetName,
					StatOp:   "max",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "createdAtAxis",
					Limit:    0,
					StatName: creMinFacetName,
					StatOp:   "min",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "createdAtAxis",
					Limit:    0,
					StatName: creMaxFacetName,
					StatOp:   "max",
				},
			},
			want: solr.StringFieldRanges{
				"publishingDateAxis": &solr.StringRange{
					Min: "2010-01-05T00:00:00Z",
					Max: "2020-02-06T00:00:00Z",
				},
				"createdAtAxis": &solr.StringRange{
					Min: "2011-03-07T00:00:00Z",
					Max: "2021-04-08T00:00:00Z",
				},
			},
		},
		{
			name: "returns an empty string if the facet is not found in the Solr response",
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      0,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Facets: map[string]interface{}{
					// We assume the constraints are such that there are NO matches
				},
			},
			facets: []*solr.Facet{
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMinFacetName,
					StatOp:   "min",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMaxFacetName,
					StatOp:   "max",
				},
			},
			want: solr.StringFieldRanges{},
		},
		{
			name: "returns error if one of the requested & correctly returned facets cannot be parsed as a min or max" +
				" facet",
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      14,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Facets: map[string]interface{}{
					pubMinFacetName:     "2010-01-05T00:00:00Z",
					pubMaxFacetName:     "2020-02-06T00:00:00Z",
					"magicOp_someField": "2011-03-07T00:00:00Z",
				},
			},
			facets: []*solr.Facet{
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMinFacetName,
					StatOp:   "min",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMaxFacetName,
					StatOp:   "max",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "createdAtAxis",
					Limit:    0,
					StatName: "magicOp_someField",
					StatOp:   "min",
				},
			},
			wantErr: true,
		},
		{
			name: "returns error if multiple copies of a given requested min facet are returned",
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      14,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Facets: map[string]interface{}{
					pubMinFacetName: "2010-01-05T00:00:00Z",
					pubMaxFacetName: "2020-02-06T00:00:00Z",
				},
			},
			facets: []*solr.Facet{
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMinFacetName,
					StatOp:   "min",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMinFacetName,
					StatOp:   "min",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMaxFacetName,
					StatOp:   "max",
				},
			},
			wantErr: true,
		},
		{
			name: "returns error if multiple copies of a given requested max facet are returned",
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      14,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Facets: map[string]interface{}{
					pubMinFacetName: "2010-01-05T00:00:00Z",
					pubMaxFacetName: "2020-02-06T00:00:00Z",
				},
			},
			facets: []*solr.Facet{
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMinFacetName,
					StatOp:   "min",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMaxFacetName,
					StatOp:   "max",
				},
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMaxFacetName,
					StatOp:   "max",
				},
			},
			wantErr: true,
		},
		{
			name: "returns an error if min value is missing from a range (even if not requested)",
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      14,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Facets: map[string]interface{}{
					pubMaxFacetName: "2020-02-06T00:00:00Z",
				},
			},
			facets: []*solr.Facet{
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMaxFacetName,
					StatOp:   "max",
				},
			},
			wantErr: true,
		},
		{
			name: "returns an error if max value is missing from a range (even if not requested)",
			solrResponse: &solr.QueryResponse{
				Response: solr.QueryResult{
					NumFoundExact: true,
					NumFound:      14,
					Start:         0,
					MaxScore:      100,
				},
				// This only works if every level of the object is typed as in the generic go JSON representation!
				Facets: map[string]interface{}{
					pubMinFacetName: "2020-02-06T00:00:00Z",
				},
			},
			facets: []*solr.Facet{
				{
					Type:     solr.MexStringStatFacetType,
					Axis:     "publishingDateAxis",
					Limit:    0,
					StatName: pubMinFacetName,
					StatOp:   "min",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDateFieldRangesFromResponse(tt.solrResponse, tt.facets)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDateFieldRangesFromResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDateFieldRangesFromResponse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCombineRanges(t *testing.T) {
	tests := []struct {
		name         string
		stringRanges []solr.StringFieldRanges
		want         solr.StringFieldRanges
		wantErr      bool
	}{
		{
			name: "combines the entries for all array members",
			stringRanges: []solr.StringFieldRanges{
				{
					"key1": &solr.StringRange{
						Min: "A",
						Max: "Q",
					},
					"key2": &solr.StringRange{
						Min: "R",
						Max: "T",
					},
				},
				{
					"key3": &solr.StringRange{
						Min: "min",
						Max: "max",
					},
				},
			},
			want: solr.StringFieldRanges{
				"key1": &solr.StringRange{
					Min: "A",
					Max: "Q",
				},
				"key2": &solr.StringRange{
					Min: "R",
					Max: "T",
				},
				"key3": &solr.StringRange{
					Min: "min",
					Max: "max",
				},
			},
		},
		{
			name: "ignores nil entries",
			stringRanges: []solr.StringFieldRanges{
				{
					"key1": &solr.StringRange{
						Min: "A",
						Max: "Q",
					},
					"key2": &solr.StringRange{
						Min: "R",
						Max: "T",
					},
				},
				nil,
				{
					"key3": &solr.StringRange{
						Min: "min",
						Max: "max",
					},
				},
			},
			want: solr.StringFieldRanges{
				"key1": &solr.StringRange{
					Min: "A",
					Max: "Q",
				},
				"key2": &solr.StringRange{
					Min: "R",
					Max: "T",
				},
				"key3": &solr.StringRange{
					Min: "min",
					Max: "max",
				},
			},
		},
		{
			name: "combines empty range to a single empty range",
			stringRanges: []solr.StringFieldRanges{
				{},
				{},
			},
			want: solr.StringFieldRanges{},
		},
		{
			name: "throws an error if a key occurs more than once",
			stringRanges: []solr.StringFieldRanges{
				{
					"key1": &solr.StringRange{
						Min: "A",
						Max: "Q",
					},
					"repeatedKey": &solr.StringRange{
						Min: "R",
						Max: "T",
					},
				},
				{
					"repeatedKey": &solr.StringRange{
						Min: "min",
						Max: "max",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CombineRanges(tt.stringRanges)
			if (err != nil) != tt.wantErr {
				t.Errorf("CombineRanges() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CombineRanges() got = %v, want %v", got, tt.want)
			}
		})
	}
}
