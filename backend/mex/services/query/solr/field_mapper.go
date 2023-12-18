package solr

import (
	"context"

	"github.com/d4l-data4life/mex/mex/shared/solr"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
	kindText "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/text"
)

/*
fieldMapper is a helper object that maps MEx field requested by clients (return fields & highlighting)
to the underlying Solr fields. The mapping has goes in both directions:

* MEx --> Solr: Each MEx field is mapped to a map from base field category to Solr backing field
* Solr --> MEx: Every Solr field backing a MEx field is mapped to the MEx field it backs
*/
type fieldMapper struct {
	mexToSolr map[string]map[string]string
	solrToMex map[string]solr.FieldWithLanguage
}

// newFieldMapper creates a mapper with correct mapping for all fields
func newFieldMapper(ctx context.Context, fieldRepo fields.FieldRepo) (*fieldMapper, error) {
	mexToSolr := make(map[string]map[string]string)
	solrToMex := make(map[string]solr.FieldWithLanguage)

	// We only make explicit entries for text fields since everything else maps trivially
	textFieldDefs, err := fieldRepo.GetFieldDefsByKind(ctx, kindText.KindName)
	if err != nil {
		return nil, err
	}
	for _, tfDef := range textFieldDefs {
		mexToSolr[tfDef.Name()] = make(map[string]string)
		mexToSolrCoreMap := solr.GetTextCoreBackingFieldInfo(tfDef.Name(), false)

		for category, coreInfo := range mexToSolrCoreMap {
			if category == solr.NormalizedBaseFieldCategory {
				// Normalized field is only used for sorting and never search, so it has no mapping
				continue
			}
			mexToSolr[tfDef.Name()][category] = coreInfo.SolrName
			solrToMex[coreInfo.SolrName] = solr.FieldWithLanguage{
				Name:     tfDef.Name(),
				Language: coreInfo.Language,
			}
		}
	}

	return &fieldMapper{
		solrToMex: solrToMex,
		mexToSolr: mexToSolr,
	}, nil
}

// getMexNameFromBackingFieldName returns the (unique) base MEx field name for a given Solr backing field
func (fm *fieldMapper) getMexNameFromBackingFieldName(solrName string) solr.FieldWithLanguage {
	mexName, ok := fm.solrToMex[solrName]
	if !ok {
		return solr.FieldWithLanguage{Name: solrName}
	}
	return mexName
}

// getHighlightBackingFieldNamesFromMexName returns the underlying Solr fields to be used for highlighting
func (fm *fieldMapper) getHighlightBackingFieldNamesFromMexName(fieldName string, isPhraseOnlyQuery bool) ([]string, error) {
	solrFieldNames, ok := fm.mexToSolr[fieldName]
	// If we have no explicit mapping information, the mapped name is assumed to be the base one
	if !ok {
		return []string{fieldName}, nil
	}
	// For a phrase query, only the unanalyzed field is queried, so we highlight there
	if isPhraseOnlyQuery {
		return []string{solrFieldNames[solr.RawContentBaseFieldCategory]}, nil
	}
	// Otherwise, highlight in all stored backing fields
	var returnFields []string
	for _, rf := range solrFieldNames {
		returnFields = append(returnFields, rf)
	}
	return returnFields, nil
}

// getReturnBackingFieldNamesFromMexName returns the underlying Solr fields to be used for return values
func (fm *fieldMapper) getReturnBackingFieldNamesFromMexName(fieldName string) ([]string, error) {
	solrFieldNames, ok := fm.mexToSolr[fieldName]
	// If we have no explicit mapping information, the mapped name is assumed to be the base one
	if !ok {
		return []string{fieldName}, nil
	}
	// Otherwise return the language-specific backing fields stored
	var requestedFieldNames []string
	for _, langCat := range solr.KnownLanguagesCategoryMap {
		if fName, langOk := solrFieldNames[langCat]; langOk {
			requestedFieldNames = append(requestedFieldNames, fName)
		}
	}
	return requestedFieldNames, nil
}
