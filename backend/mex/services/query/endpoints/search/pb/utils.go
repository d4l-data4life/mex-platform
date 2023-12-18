package pb

import "github.com/d4l-data4life/mex/mex/shared/solr"

func (sr *SearchResponse) AddDocItem(doc *solr.DocItem) *SearchResponse {
	sr.Items = append(sr.Items, doc)
	return sr
}

func NewSearchResponse() *SearchResponse {
	return &SearchResponse{
		Highlights: make([]*solr.Highlight, 0),
		Facets:     make([]*solr.FacetResult, 0),
		Items:      make([]*solr.DocItem, 0),
		Diagnostics: &solr.Diagnostics{
			ParsingSucceeded: true,
		},
	}
}
