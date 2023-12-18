package solr

import (
	"encoding/json"
	"fmt"

	"github.com/d4l-data4life/mex/mex/shared/errstat"
)

const (
	/*
		These are types used to distinguish how a given facet should be serialized for Solr.
		It distinguishes not only between the facet types Solr defines, but also (where needed)
		between different data types for a given facet types  (e.g. a Solr range facet can have
		its range given as strings or numbers).
	*/
	SolrTermsFacetType       = "exact"
	SolrStringRangeFacetType = "stringRange"
	SolrStringStatFacetType  = "stringStat"
)

// GenericObject is shorthand for a generic JSON object
type GenericObject map[string]interface{}

// StringFieldRanges is used for mapping field names to corresponding ranges
type StringFieldRanges map[string]*StringRange

/*
SolrFacet represents the different facet objects that Solr offers. Because Solr allows
specifying different types of facets with different properties this is in effects a "union type"
in that in covers all properties that may be present although each specific facet type will only
use a subset thereof. Since the Solr API allows some properties to contain values of different types
(e.g. 'start', 'end' and 'gap' can be strings or numbers), this Go type further splits such properties
into one property per type.
*/
type SolrFacet struct {
	DetailedType string // More detailed type  than the Solr classification as it includes the technical
	// type as well (string vs. number)
	Field       string
	NumBuckets  bool
	Limit       uint32
	Offset      uint32
	StartString string
	EndString   string
	GapString   string
	StatOp      string
	ExcludeTags []string
}

/*
SolrFacetSet represents a set of facets of the different facet objects that Solr offers.

We split it out into a separate type to specify custom JSON marshalling log. This cannot be done at the
level of individual facet because stat facets are just strings, not objects, and hence to not serialize to
valid JSON by themselves.

NOTE THAT THIS TYPE THEREFORE HAS CUSTOM JSON MARSHALLING LOGIC (SEE BELOW) TO ENSURE THAT EACH TYPE
OF FACET IS SERIALIZED IN THE CORRECT WAY.
*/
type SolrFacetSet map[string]SolrFacet

/*
MarshalJSON does custom JSON marshalling for the SolrFacetSet type.

This is required because:

 1. the Solr JSON facet API uses the same property name for different data types - e.g.
    for range facets the 'start' property can be a number or a date string.
 2. THe facet information can be an object or just a string (for stat facets)

Serialization is chosen based on the 'DetailedType' struct field.
*/
func (sfs SolrFacetSet) MarshalJSON() ([]byte, error) {
	rawObj := make(GenericObject)
	for facetName, facet := range sfs {
		switch facet.DetailedType {
		case SolrTermsFacetType:
			entry := map[string]interface{}{
				"type":  "terms",
				"field": facet.Field,
				"limit": facet.Limit,
			}
			if facet.NumBuckets {
				entry["numBuckets"] = facet.NumBuckets
			}
			if facet.Offset != 0 {
				entry["offset"] = facet.Offset
			}
			if len(facet.ExcludeTags) > 0 {
				entry["domain"] = map[string][]string{"excludeTags": facet.ExcludeTags}
			}
			rawObj[facetName] = entry
		case SolrStringRangeFacetType:
			entry := map[string]interface{}{
				"type":  "range",
				"field": facet.Field,
				"start": facet.StartString,
				"end":   facet.EndString,
				"gap":   facet.GapString,
			}
			if len(facet.ExcludeTags) > 0 {
				entry["domain"] = map[string][]string{"excludeTags": facet.ExcludeTags}
			}
			rawObj[facetName] = entry
		case SolrStringStatFacetType:
			expr, err := CreateStatExpression(facet.Field, facet.StatOp)
			if err != nil {
				return nil, errstat.MakeMexStatus(errstat.QueryCreationFailedInternal, fmt.Sprintf("failed to compute stat expression: %s", err.Error())).Err()
			}
			rawObj[facetName] = expr
		default:
			return nil, fmt.Errorf("facet type cannot be serialized to JSON through marshalling")
		}
	}
	return json.Marshal(rawObj)
}

// ParamObj is a key-value object that holds parameters that otherwise need to be put in the Solr query URL
type ParamObj struct {
	DefType    string `json:"defType,omitempty"`
	QOp        string `json:"q.op,omitempty"`
	OmitHeader bool   `json:"omitHeader,omitempty"`
	Hl         bool   `json:"hl,omitempty"`
	HlFl       string `json:"hl.fl,omitempty"`
	HlQ        string `json:"hl.q,omitempty"`
	HlMethod   string `json:"hl.method,omitempty"`
	HlSnippets uint32 `json:"hl.snippets,omitempty"`
	HlFragsize uint32 `json:"hl.fragsize,omitempty"`
	HlTagPre   string `json:"hl.tag.pre,omitempty"`
	HlTagPost  string `json:"hl.tag.post,omitempty"`
}

// QueryBody represents the body of a query for the Solr JSON query API
type QueryBody struct {
	Sort   string       `json:"sort,omitempty"`
	Query  string       `json:"query"` // Empty query string is not the same as absent --> no omitempty
	Filter []string     `json:"filter,omitempty"`
	Limit  uint32       `json:"limit"` // The value 0 is not the same as absent --> no omitempty
	Offset uint32       `json:"offset,omitempty"`
	Fields []string     `json:"fields,omitempty"`
	Facet  SolrFacetSet `json:"facet,omitempty"` // facet is a JSON object with freely chosen labels as the top-level properties
	Params ParamObj     `json:"params,omitempty"`
}

// QueryResult represents the result returned by the Solr JSON query API
type QueryResult struct {
	NumFound      uint32          `json:"numFound"`
	NumFoundExact bool            `json:"numFoundExact"`
	Start         uint32          `json:"start"`
	MaxScore      float64         `json:"maxScore"`
	Docs          []GenericObject `json:"docs"` // docs is an array of the (selected fields of) the matching documents
}

// QueryResponse represents the Solr response to a query
type QueryResponse struct {
	Response     QueryResult            `json:"response"`
	Facets       map[string]interface{} `json:"facets"`       // facet is a JSON object with freely chosen labels as the top-level properties
	Highlighting map[string]interface{} `json:"highlighting"` // highlighting is a JSON object with the document IDs as the top-level properties
	Error        map[string]interface{} `json:"error"`        // error is a JSON object with freely chosen labels as the top-level properties
}

// CopyFieldResponse represents the SOlr information for a single copy field returned by Solr
type CopyFieldResponse struct {
	Source      string `json:"source"`
	Destination string `json:"dest"`
}
