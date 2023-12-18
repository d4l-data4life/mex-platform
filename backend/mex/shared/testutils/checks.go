package testutils

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/d4l-data4life/mex/mex/services/query/endpoints/search/pb"
	"github.com/d4l-data4life/mex/mex/services/query/parser"
	"github.com/d4l-data4life/mex/mex/shared/solr"
	"github.com/d4l-data4life/mex/mex/shared/utils"
)

type FieldCheck func(fields solr.FieldCategoryToSolrFieldDefsMap, t *testing.T)

type BodyCheck func(body *solr.QueryBody, t *testing.T)

func CheckOmitHeader(expectedVal bool) BodyCheck {
	return func(body *solr.QueryBody, t *testing.T) {
		if body.Params.OmitHeader != expectedVal {
			t.Errorf("incorrect omit header status: wanted '%v', got '%v'", expectedVal, body.Params.OmitHeader)
		}
	}
}

func CheckQOp(expectedOp string) BodyCheck {
	return func(body *solr.QueryBody, t *testing.T) {
		if body.Params.QOp != expectedOp {
			t.Errorf("incorrect operator for combining terms: wanted '%s', got '%s'", expectedOp,
				body.Params.QOp)
		}
	}
}

func CheckSorting(expectedField string, expectedOrder string) BodyCheck {
	return func(body *solr.QueryBody, t *testing.T) {
		expectedSortStr := ""
		if expectedField != "" || expectedOrder != "" {
			expectedSortStr = fmt.Sprintf("%s %s", expectedField, expectedOrder)
		}
		if body.Sort != expectedSortStr {
			t.Errorf("incorrect sorting: wanted '%s', got '%s'", expectedSortStr,
				body.Sort)
		}
	}
}

func CheckHighlightingFields(expectedFields []string) BodyCheck {
	return func(body *solr.QueryBody, t *testing.T) {
		if expectedFields == nil {
			// Check for absence
			if body.Params.Hl || body.Params.HlFl != "" {
				t.Errorf("highlighting set where not was expected: highlighting on '%v', highlight fields '%s'",
					body.Params.Hl, body.Params.HlFl)
			}
		} else {
			if !body.Params.Hl {
				t.Errorf("highlighting was not activated though it was expected to be")
			}
			hFieldArray := strings.Split(body.Params.HlFl, ",")
			sort.Strings(hFieldArray)
			sort.Strings(expectedFields)
			if !reflect.DeepEqual(hFieldArray, expectedFields) {
				t.Errorf("incorrect highlight fields set: wanted '%v', got '%v'", expectedFields, hFieldArray)
			}
		}
	}
}

func CheckHighlightingParameters(expectedMethod string, expectedSnippets uint32, expectedFragsize uint32, expectedPreTag string, expectedPostTag string) BodyCheck {
	return func(body *solr.QueryBody, t *testing.T) {
		if body.Params.HlMethod != expectedMethod {
			t.Errorf("highlighting method not set to expected value: method was '%s', expected '%s'",
				body.Params.HlMethod, expectedMethod)
		}
		if body.Params.HlSnippets != expectedSnippets {
			t.Errorf("number of highlighting snippets not set to expected value: method was '%v', expected '%v'",
				body.Params.HlSnippets, expectedSnippets)
		}
		if body.Params.HlFragsize != expectedFragsize {
			t.Errorf("highlighting snippet fragsize not set to expected value: method was '%v', expected '%v'",
				body.Params.HlFragsize, expectedFragsize)
		}
		if (body.Params.HlTagPre != expectedPreTag) || (body.Params.HlTagPost != expectedPostTag) {
			t.Errorf("highlighting tags not set to expected value: pre tag was '%s', expected '%s'; post tag was '%s', expected '%s'",
				body.Params.HlTagPre, expectedPreTag, body.Params.HlTagPost, expectedPostTag)
		}
	}
}

func CheckQuery(expectedQuery string) BodyCheck {
	return func(body *solr.QueryBody, t *testing.T) {
		if body.Query != expectedQuery {
			t.Errorf("incorrect query: wanted '%s', got '%s'", expectedQuery, body.Query)
		}
	}
}

func CheckLimit(expectedLimit uint32) BodyCheck {
	return func(body *solr.QueryBody, t *testing.T) {
		if body.Limit != expectedLimit {
			t.Errorf("incorrect limit: wanted '%d', got '%d'", expectedLimit, body.Limit)
		}
	}
}

func CheckOffset(expectedOffset uint32) BodyCheck {
	return func(body *solr.QueryBody, t *testing.T) {
		if body.Offset != expectedOffset {
			t.Errorf("incorrect offset: wanted '%d', got '%d'", expectedOffset, body.Offset)
		}
	}
}

func CheckFields(expectedFields []string) BodyCheck {
	return func(body *solr.QueryBody, t *testing.T) {
		sort.Strings(expectedFields)
		sort.Strings(body.Fields)
		if !reflect.DeepEqual(body.Fields, expectedFields) {
			t.Errorf("incorrect fields: wanted '%v', got '%v'", expectedFields, body.Fields)
		}
	}
}

func CheckConstraints(expectedFilters []string, expectedAxes []string) BodyCheck {
	if len(expectedFilters) != len(expectedAxes) {
		panic(fmt.Sprintf("Filter check not possible no. of filter & no. of axes do not match: #filters = %d, "+
			"#axes = %d",
			len(expectedFilters), len(expectedAxes)))
	}
	var expectedFullFilters []string
	for i, filter := range expectedFilters {
		fullFilter, _ := solr.TagExpr(filter, expectedAxes[i])
		expectedFullFilters = append(expectedFullFilters, fullFilter)
	}
	return func(body *solr.QueryBody, t *testing.T) {
		if !reflect.DeepEqual(body.Filter, expectedFullFilters) {
			t.Errorf("incorrect filters: wanted '%v', got '%v'", expectedFullFilters, body.Filter)
		}
	}
}

func CheckFacets(expectedFacets solr.SolrFacetSet) BodyCheck {
	return func(body *solr.QueryBody, t *testing.T) {
		// Separate check for empty maps since DeepEqual does not work for them
		if len(body.Facet) == 0 && len(expectedFacets) == 0 {
			return
		}
		if !reflect.DeepEqual(body.Facet, expectedFacets) {
			t.Errorf("faceting incorrect: wanted '%v', got '%v'", expectedFacets, body.Facet)
		}
	}
}

type DiagnosticCheck func(diag *solr.Diagnostics, t *testing.T)

func CheckParsingSuccess(expectedVal bool) DiagnosticCheck {
	return func(diag *solr.Diagnostics, t *testing.T) {
		if diag.ParsingSucceeded != expectedVal {
			t.Errorf("diagnostics incorrect: wanted parsing success flag '%v', got '%v'", expectedVal,
				diag.ParsingSucceeded)
		}
	}
}

func CheckCleanedQuery(expectedQuery string) DiagnosticCheck {
	return func(diag *solr.Diagnostics, t *testing.T) {
		if diag.CleanedQuery != expectedQuery {
			t.Errorf("diagnostics incorrect: wrong cleaned query - wanted: %s, got: %s",
				expectedQuery, diag.CleanedQuery)
		}
	}
}

func CheckParsingErrors(expectedErrors []string) DiagnosticCheck {
	return func(diag *solr.Diagnostics, t *testing.T) {
		if expectedErrors == nil {
			if diag.ParsingErrors != nil {
				t.Errorf("diagnostics incorrect: did not expect parsing errors but got '%s'",
					strings.Join(diag.ParsingErrors, ", "))
			}
		} else if !reflect.DeepEqual(diag.ParsingErrors, expectedErrors) {
			t.Errorf("diagnostics incorrect: wanted parsing errors '%s', got '%s'", strings.Join(expectedErrors, ", "),
				strings.Join(diag.ParsingErrors, ", "))
		}
	}
}

func CheckIgnoredErrors(expectedIgnoredErrors []string) DiagnosticCheck {
	return func(diag *solr.Diagnostics, t *testing.T) {
		if expectedIgnoredErrors == nil {
			if diag.IgnoredErrors != nil {
				t.Errorf("diagnostics incorrect: did not expect ignored errors but got '%s'",
					strings.Join(diag.IgnoredErrors, ", "))
			}
		} else if !reflect.DeepEqual(diag.IgnoredErrors, expectedIgnoredErrors) {
			t.Errorf("diagnostics incorrect: wanted ignored errors '%s', got '%s'", strings.Join(expectedIgnoredErrors, ", "),
				strings.Join(diag.IgnoredErrors, ", "))
		}
	}
}

type ResponseCheck func(response *pb.SearchResponse, t *testing.T)

func CheckHighlightResponse(expectedHighlights []*solr.Highlight) ResponseCheck {
	return func(response *pb.SearchResponse, t *testing.T) {
		gotHighlights := response.Highlights
		sort.Slice(expectedHighlights, func(i, j int) bool {
			return expectedHighlights[i].ItemId > expectedHighlights[j].ItemId
		})
		sort.Slice(gotHighlights, func(i, j int) bool {
			return gotHighlights[i].ItemId > gotHighlights[j].ItemId
		})
		if len(expectedHighlights) != len(gotHighlights) {
			t.Errorf("Highlight: wrong no. items with highlights returned - wanted %d, got %d", len(expectedHighlights), len(gotHighlights))
			return
		}
		for k, wantH := range expectedHighlights {
			gotH := gotHighlights[k]
			if wantH.ItemId != gotH.ItemId {
				t.Errorf("Highlight: item ID for highlight with index %d does not match - wanted %s, got %s", k, wantH.ItemId, gotH.ItemId)
				continue
			}
			wantMatches := wantH.Matches
			gotMatches := gotH.Matches
			if len(wantMatches) != len(gotMatches) {
				t.Errorf("Highlight: wrong no. matches returned for item ID %s - wanted %d, got %d", wantH.ItemId, len(wantMatches), len(gotMatches))
				return
			}
			sort.Slice(wantMatches, func(i, j int) bool {
				return wantMatches[i].FieldName > wantMatches[j].FieldName
			})
			sort.Slice(gotMatches, func(i, j int) bool {
				return gotMatches[i].FieldName > gotMatches[j].FieldName
			})
			for l, wantM := range wantMatches {
				gotM := gotMatches[l]
				if wantM.FieldName != gotM.FieldName {
					t.Errorf("Highlight: field name mismatch for match with index %d returned for item ID %s - wanted %s, got %s", l, wantH.ItemId, wantM.FieldName, gotM.FieldName)
				}
				if wantM.Language != gotM.Language {
					t.Errorf("Highlight: language mismatch for match with index %d returned for item ID %s - wanted %s, got %s", l, wantH.ItemId, wantM.Language, gotM.Language)
				}
				sort.Strings(wantM.Snippets)
				sort.Strings(gotM.Snippets)
				if !reflect.DeepEqual(wantM.Snippets, gotM.Snippets) {
					t.Errorf("Highlight: snippet mismatch for match with index %d returned for item ID %s - wanted %v, got %v", l, wantH.ItemId, wantM.Snippets, gotM.Snippets)
				}
			}
		}
	}
}

type ClientCheck func(solr.MockClient, *testing.T)

func CheckUniqueIDFetchedViaClient(shouldBeFetched bool) ClientCheck {
	return func(client solr.MockClient, t *testing.T) {
		wasFetched := utils.Contains(client.CallQueue, "GetSchemaUniqueKey")
		if wasFetched != shouldBeFetched {
			t.Errorf("unexpected call/missing call to unique ID: wanted %v, got %v", shouldBeFetched,
				wasFetched)
		}
	}
}

func CheckFieldsFetchedViaClient(shouldBeFetched bool) ClientCheck {
	return func(client solr.MockClient, t *testing.T) {
		wasFetched := utils.Contains(client.CallQueue, "GetSchemaFields")
		if wasFetched != shouldBeFetched {
			t.Errorf("unexpected call/missing call to fetch fields: wanted %v, got %v", shouldBeFetched,
				wasFetched)
		}
	}
}

func CheckCopyFieldsFetchedViaClient(shouldBeFetched bool) ClientCheck {
	return func(client solr.MockClient, t *testing.T) {
		wasFetched := utils.Contains(client.CallQueue, "GetSchemaCopyFields")
		if wasFetched != shouldBeFetched {
			t.Errorf("unexpected call/missing call to fetch copy fields: wanted %v, got %v", shouldBeFetched,
				wasFetched)
		}
	}
}

func CheckDynamicFieldsFetchedViaClient(shouldBeFetched bool) ClientCheck {
	return func(client solr.MockClient, t *testing.T) {
		wasFetched := utils.Contains(client.CallQueue, "GetSchemaDynamicFields")
		if wasFetched != shouldBeFetched {
			t.Errorf("unexpected call/missing call to fetch dynamic fields: wanted %v, got %v", shouldBeFetched,
				wasFetched)
		}
	}
}

func CheckFieldsRemovedViaClient(removedFields []string) ClientCheck {
	return func(client solr.MockClient, t *testing.T) {
		if len(removedFields) > len(client.FieldsRemoved) {
			t.Errorf("given %d field removal constraints, but only %d fields were removed", len(removedFields),
				len(client.FieldsRemoved))
			return
		}
		sort.Strings(removedFields)
		sort.Strings(client.FieldsRemoved)
		for i, fieldStr := range removedFields {
			if fieldStr != client.FieldsRemoved[i] {
				t.Errorf("removed fields do not match on index %d: wanted %s, got %s", i, fieldStr,
					client.FieldsRemoved[i])
			}
		}
	}
}

func CheckFieldsAddedViaClient(addedFields []string) ClientCheck {
	return func(client solr.MockClient, t *testing.T) {
		if len(addedFields) > len(client.FieldsAdded) {
			t.Errorf("given %d field addition constraints, but only %d fields were added", len(addedFields),
				len(client.FieldsAdded))
			return
		}
		for i, fieldStr := range addedFields {
			if fieldStr != client.FieldsAdded[i] {
				t.Errorf("added fields do not match on index %d: wanted %s, got %s", i, fieldStr,
					client.FieldsAdded[i])
			}
		}
	}
}

func CheckCopyFieldsRemovedViaClient(removedFields []string) ClientCheck {
	return func(client solr.MockClient, t *testing.T) {
		if len(removedFields) > len(client.CopyFieldsRemoved) {
			t.Errorf("given %d copy field removal constraints, but only %d copy fields were removed",
				len(removedFields),
				len(client.CopyFieldsRemoved))
			return
		}
		for i, fieldStr := range removedFields {
			if fieldStr != client.CopyFieldsRemoved[i] {
				t.Errorf("removed copy fields do not match on index %d: wanted %s, got %s", i, fieldStr,
					client.CopyFieldsRemoved[i])
			}
		}
	}
}

func CheckCopyFieldsAddedViaClient(addedFields []string) ClientCheck {
	return func(client solr.MockClient, t *testing.T) {
		if len(addedFields) > len(client.CopyFieldsAdded) {
			t.Errorf("given %d copy field addition constraints, but only %d copy fields were added", len(addedFields),
				len(client.CopyFieldsAdded))
			return
		}
		for i, fieldStr := range addedFields {
			if fieldStr != client.CopyFieldsAdded[i] {
				t.Errorf("added copy fields do not match on index %d: wanted %s, got %s", i, fieldStr,
					client.CopyFieldsAdded[i])
			}
		}
	}
}

func CheckDynamicFieldsRemovedViaClient(removedFields []string) ClientCheck {
	return func(client solr.MockClient, t *testing.T) {
		if len(removedFields) > len(client.DynamicFieldsRemoved) {
			t.Errorf("given %d dynamic field removal constraints, but only %d copy fields were removed",
				len(removedFields),
				len(client.DynamicFieldsRemoved))
			return
		}
		for i, fieldStr := range removedFields {
			if fieldStr != client.DynamicFieldsRemoved[i] {
				t.Errorf("removed dynamic fields do not match on index %d: wanted %s, got %s", i, fieldStr,
					client.DynamicFieldsRemoved[i])
			}
		}
	}
}

func CheckDynamicFieldsAddedViaClient(addedFields []string) ClientCheck {
	return func(client solr.MockClient, t *testing.T) {
		if len(addedFields) > len(client.DynamicFieldsAdded) {
			t.Errorf("given %d dynamic field addition constraints, but only %d dynamic fields were added",
				len(addedFields),
				len(client.DynamicFieldsAdded))
			return
		}
		for i, fieldStr := range addedFields {
			if fieldStr != client.DynamicFieldsAdded[i] {
				t.Errorf("added dynamic fields do not match on index %d: wanted %s, got %s", i, fieldStr,
					client.DynamicFieldsAdded[i])
			}
		}
	}
}

// Checks that all calls to firstMethod come before the first call to secondMethod
func CheckClientCallOrder(firstMethod string, secondMethod string) ClientCheck {
	return func(client solr.MockClient, t *testing.T) {
		firstIndex := -1
		secondIndex := -1
		for i, m := range client.CallQueue {
			if m == firstMethod {
				// Update every time we see this since we want the last occurrence
				firstIndex = i
			} else if m == secondMethod && secondIndex == -1 {
				// Only need to set on first occurrence
				secondIndex = i
			}
		}
		if firstIndex == -1 || secondIndex == -1 {
			t.Errorf("one or both of the two required methods weere not called: %s, %s", firstMethod, secondMethod)
		} else if firstIndex > secondIndex {
			t.Errorf("methods called in wrong order: some calls to %s came before the last call to %s",
				secondMethod, firstMethod)
		}
	}
}

type MatchingOpCheck func(conf parser.MatchingOpsConfig, t *testing.T)

func CheckSearchedFields(expectedTermFields []string, expectedPhraseFields []string) MatchingOpCheck {
	return func(conf parser.MatchingOpsConfig, t *testing.T) {
		gotTermFields := make([]string, len(conf.Term))
		for i1, tf := range conf.Term {
			gotTermFields[i1] = tf.FieldName
		}
		if len(gotTermFields) != len(expectedTermFields) {
			t.Errorf("expected %d fields to be searched for terms but got %d", len(expectedTermFields), len(gotTermFields))
		}
		gotPhraseFields := make([]string, len(conf.Phrase))
		for i2, tf := range conf.Phrase {
			gotPhraseFields[i2] = tf.FieldName
		}
		if len(gotPhraseFields) != len(expectedPhraseFields) {
			t.Errorf("expected %d fields to be searched for phrases but got %d", len(expectedPhraseFields), len(gotPhraseFields))
		}
		sort.Strings(gotTermFields)
		sort.Strings(gotPhraseFields)
		sort.Strings(expectedTermFields)
		sort.Strings(expectedPhraseFields)
		for j1, fieldName := range expectedTermFields {
			if fieldName != gotTermFields[j1] {
				t.Errorf("term search fields differ on index %d: wanted %s, got %s", j1, fieldName, gotTermFields[j1])
			}
		}
		for j2, fieldName := range expectedPhraseFields {
			if fieldName != gotPhraseFields[j2] {
				t.Errorf("phrase search fields differ on index %d: wanted %s, got %s", j2, fieldName, gotPhraseFields[j2])
			}
		}
	}
}

func CheckBoostFactor(expectedTermBoostFactors map[string]string, expectedPhraseBoostFactors map[string]string) MatchingOpCheck {
	return func(conf parser.MatchingOpsConfig, t *testing.T) {
		gotTermBoostFactors := make(map[string]string)
		for _, tf := range conf.Term {
			gotTermBoostFactors[tf.FieldName] = tf.BoostFactor
		}
		if len(gotTermBoostFactors) != len(expectedTermBoostFactors) {
			t.Errorf("expected %d fields to be searched for terms but got %d", len(expectedTermBoostFactors), len(gotTermBoostFactors))
		}
		gotPhraseBoostFactors := make(map[string]string)
		for _, tf := range conf.Phrase {
			gotPhraseBoostFactors[tf.FieldName] = tf.BoostFactor
		}
		if len(gotPhraseBoostFactors) != len(expectedPhraseBoostFactors) {
			t.Errorf("expected %d fields to be searched for phrases but got %d", len(expectedPhraseBoostFactors), len(gotPhraseBoostFactors))
		}
		for key, val := range expectedTermBoostFactors {
			if gotTermBoostFactors[key] != val {
				t.Errorf("term search boost factors for field '%s': wanted %s, got %s", key, val, gotTermBoostFactors[key])
			}
		}
		for key, val := range expectedPhraseBoostFactors {
			if gotPhraseBoostFactors[key] != val {
				t.Errorf("phrase search boost factors for field '%s': wanted %s, got %s", key, val, gotPhraseBoostFactors[key])
			}
		}
	}
}

func CheckEditDistance(expectedTermEditDistances map[string]uint32, expectedPhraseEditDistances map[string]uint32) MatchingOpCheck {
	return func(conf parser.MatchingOpsConfig, t *testing.T) {
		gotTermEditDistances := make(map[string]uint32)
		for _, tf := range conf.Term {
			gotTermEditDistances[tf.FieldName] = tf.MaxEditDistance
		}
		if len(gotTermEditDistances) != len(expectedTermEditDistances) {
			t.Errorf("expected %d fields to be searched for terms but got %d", len(expectedTermEditDistances), len(gotTermEditDistances))
		}
		gotPhraseEditDistances := make(map[string]uint32)
		for _, tf := range conf.Phrase {
			gotPhraseEditDistances[tf.FieldName] = tf.MaxEditDistance
		}
		if len(gotPhraseEditDistances) != len(expectedPhraseEditDistances) {
			t.Errorf("expected %d fields to be searched for phrases but got %d", len(expectedPhraseEditDistances), len(gotPhraseEditDistances))
		}
		for key, val := range expectedTermEditDistances {
			if gotTermEditDistances[key] != val {
				t.Errorf("term search edit distances for field '%s': wanted %d, got %d", key, val, gotTermEditDistances[key])
			}
		}
		for key, val := range expectedPhraseEditDistances {
			if gotPhraseEditDistances[key] != val {
				t.Errorf("phrase search edit distances for field '%s': wanted %d, got %d", key, val, gotPhraseEditDistances[key])
			}
		}
	}
}
