package solr

import (
	"context"
	"reflect"
	"sort"
	"testing"

	sharedFields "github.com/d4l-data4life/mex/mex/shared/fields"
	"github.com/d4l-data4life/mex/mex/shared/solr"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/frepo"
	kindnumber "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/number"
	kindstring "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/string"
	kindtext "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/text"
)

var testRepo = frepo.NewMockedFieldRepo([]fields.BaseFieldDef{
	(&kindnumber.KindNumber{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "count",
		Kind: "number", IndexDef: &sharedFields.IndexDef{}}),
	(&kindstring.KindString{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "category",
		Kind: "string", IndexDef: &sharedFields.IndexDef{}}),
	(&kindtext.KindText{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "title", Kind: "text", IndexDef: &sharedFields.IndexDef{}}),
	(&kindtext.KindText{}).MustValidateDefinition(context.TODO(), &sharedFields.FieldDef{Name: "abstract", Kind: "text",
		IndexDef: &sharedFields.IndexDef{}}),
})

const titleName = "title"

func TestFieldMapper_GetHighlightBackingFieldNamesFromMexName(t *testing.T) {
	var genericTitleFieldName, _ = solr.GetLangSpecificFieldName(titleName, solr.GenericLangAbbrev)
	var deTitleFieldName, _ = solr.GetLangSpecificFieldName(titleName, solr.GermanLangAbbrev)
	var enTitleFieldName, _ = solr.GetLangSpecificFieldName(titleName, solr.EnglishLangAbbrev)
	var prefixTitleFieldName = solr.GetPrefixBackingFieldName(titleName)
	var unanalyzedTitleFieldName = solr.GetRawBackingFieldName(titleName)

	var tests = []struct {
		name              string
		repo              fields.FieldRepo
		mexFieldName      string
		isPhraseOnlyQuery bool
		want              []string
		wantErr           bool
	}{
		{
			name: "For an unknown non-text MEx field, " +
				"returns a single Solr field (with same name as MEx field)",
			repo:         testRepo,
			mexFieldName: "unknown",
			want:         []string{"unknown"},
		},
		{
			name: "For a configured non-text, non-string MEx field, " +
				"returns a single Solr field (with same name as MEx)",
			repo:         testRepo,
			mexFieldName: "count",
			want:         []string{"count"},
		},
		{
			name: "For MEx string field, " +
				"returns the basic backing field for a phrase query",
			repo:              testRepo,
			mexFieldName:      "category",
			isPhraseOnlyQuery: true,
			want:              []string{"category"},
		},
		{
			name: "For MEx text field, " +
				"returns the unanalyzed backing field for a phrase query",
			repo:              testRepo,
			mexFieldName:      "title",
			isPhraseOnlyQuery: true,
			want:              []string{unanalyzedTitleFieldName},
		},
		{
			name: "For MEx string field, " +
				"returns the base backing fields for a non-phrase query",
			repo:              testRepo,
			mexFieldName:      "category",
			isPhraseOnlyQuery: false,
			want:              []string{"category"},
		},
		{
			name: "For MEx text field, " +
				"returns all backing fields for a non-phrase query",
			repo:              testRepo,
			mexFieldName:      "title",
			isPhraseOnlyQuery: false,
			want:              []string{genericTitleFieldName, deTitleFieldName, enTitleFieldName, prefixTitleFieldName, unanalyzedTitleFieldName},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, _ := newFieldMapper(context.Background(), tt.repo)
			got, gotErr := fm.getHighlightBackingFieldNamesFromMexName(tt.mexFieldName, tt.isPhraseOnlyQuery)
			sort.Strings(got)
			sort.Strings(tt.want)
			if tt.wantErr != (gotErr != nil) {
				t.Errorf("GetBackingFieldNamesFromMexName(): wantError = %v but got error = %v", tt.wantErr, gotErr)
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBackingFieldNamesFromMexName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldMapper_GetReturnBackingFieldNamesFromMexName(t *testing.T) {
	var genericTitleFieldName, _ = solr.GetLangSpecificFieldName(titleName, solr.GenericLangAbbrev)
	var deTitleFieldName, _ = solr.GetLangSpecificFieldName(titleName, solr.GermanLangAbbrev)
	var enTitleFieldName, _ = solr.GetLangSpecificFieldName(titleName, solr.EnglishLangAbbrev)

	var tests = []struct {
		name         string
		repo         fields.FieldRepo
		mexFieldName string
		want         []string
		wantErr      bool
	}{
		{
			name: "For an unknown non-text MEx field, " +
				"returns a single Solr field (with same name as MEx field)",
			repo:         testRepo,
			mexFieldName: "unknown",
			want:         []string{"unknown"},
		},
		{
			name: "For a configured non-text MEx field, " +
				"returns a single Solr field (with same name as MEx)",
			repo:         testRepo,
			mexFieldName: "category",
			want:         []string{"category"},
		},
		{
			name: "For MEx text field, " +
				"returns all the language-specific fields",
			repo:         testRepo,
			mexFieldName: "title",
			want:         []string{genericTitleFieldName, deTitleFieldName, enTitleFieldName},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, _ := newFieldMapper(context.Background(), tt.repo)
			got, gotErr := fm.getReturnBackingFieldNamesFromMexName(tt.mexFieldName)
			sort.Strings(got)
			sort.Strings(tt.want)
			if tt.wantErr != (gotErr != nil) {
				t.Errorf("GetBackingFieldNamesFromMexName(): wantError = %v but got error = %v", tt.wantErr, gotErr)
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBackingFieldNamesFromMexName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldMapper_GetMexNameFromBackingFieldName(t *testing.T) {
	genericTitleFieldName, _ := solr.GetLangSpecificFieldName(titleName, solr.GenericLangAbbrev)
	deTitleFieldName, _ := solr.GetLangSpecificFieldName(titleName, solr.GermanLangAbbrev)
	enTitleFieldName, _ := solr.GetLangSpecificFieldName(titleName, solr.EnglishLangAbbrev)
	unsupportedTitleFieldName, _ := solr.GetLangSpecificFieldName(titleName, "dk")

	tests := []struct {
		name             string
		repo             fields.FieldRepo
		indexedFieldName string
		want             solr.FieldWithLanguage
	}{
		{
			name:             "Returns the Solr field unchanged & with no Language if not matching a known (field, Language) combination",
			repo:             testRepo,
			indexedFieldName: unsupportedTitleFieldName,
			want: solr.FieldWithLanguage{
				Name: unsupportedTitleFieldName,
			},
		},
		{
			name: "For a non-text field, " +
				"the Solr field is mapped to a field with the same name and no Language set",
			repo:             testRepo,
			indexedFieldName: "category",
			want: solr.FieldWithLanguage{
				Name: "category",
			},
		},
		{
			name: "For a text field, " +
				"the generic-language Solr fields is mapped to the base Solr field with no Language set",
			repo:             testRepo,
			indexedFieldName: genericTitleFieldName,
			want: solr.FieldWithLanguage{
				Name: "title",
			},
		},
		{
			name: "For a text field, " +
				"the German-language Solr fields is mapped to the MEx field with language set to 'de'",
			repo:             testRepo,
			indexedFieldName: deTitleFieldName,
			want: solr.FieldWithLanguage{
				Name:     "title",
				Language: solr.GermanLangAbbrev,
			},
		},
		{
			name: "For a text field, " +
				"the English-language Solr fields is mapped to the MEx field with language set to 'en'",
			repo:             testRepo,
			indexedFieldName: enTitleFieldName,
			want: solr.FieldWithLanguage{
				Name:     "title",
				Language: solr.EnglishLangAbbrev,
			},
		},
		{
			name: "For a text field, " +
				"a Solr field matching that of a non-supported language in format returns error",
			repo:             testRepo,
			indexedFieldName: unsupportedTitleFieldName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, _ := newFieldMapper(context.Background(), tt.repo)
			got := fm.getMexNameFromBackingFieldName(tt.indexedFieldName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMexNameFromBackingFieldName() = %v, want %v", got, tt.want)
			}
		})
	}
}
