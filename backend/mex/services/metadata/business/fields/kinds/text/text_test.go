package kind_text

import (
	"context"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"

	fieldUtils "github.com/d4l-data4life/mex/mex/shared/fields"
	"github.com/d4l-data4life/mex/mex/shared/solr"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/datamodel"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
)

func TestKindContact_ValidateDefinition(t *testing.T) {
	tests := []struct {
		name    string
		request *fieldUtils.FieldDef
		want    fields.BaseFieldDef
		wantErr bool
	}{
		{
			name: "non-text field causes an error",
			request: &fieldUtils.FieldDef{
				Name:     "someName",
				Kind:     "number",
				IndexDef: &fieldUtils.IndexDef{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kind := &KindText{}
			got, err := kind.ValidateDefinition(context.TODO(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateDefinition() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GenerateSolrFields(t *testing.T) {
	testFieldName := "text"

	genericName, _ := solr.GetLangSpecificFieldName(testFieldName, "")
	deName, _ := solr.GetLangSpecificFieldName(testFieldName, solr.GermanLangAbbrev)
	enName, _ := solr.GetLangSpecificFieldName(testFieldName, solr.EnglishLangAbbrev)
	prefixName := solr.GetPrefixBackingFieldName(testFieldName)
	rawName := solr.GetRawBackingFieldName(testFieldName)

	tests := []struct {
		name    string
		input   *fieldUtils.FieldDef
		want    solr.FieldCategoryToSolrFieldDefsMap
		wantErr bool
	}{
		{
			name: "invalid kind causes error",
			input: &fieldUtils.FieldDef{
				Name:     "test",
				Kind:     "xxx",
				IndexDef: &fieldUtils.IndexDef{},
			},
			wantErr: true,
		},
		{
			name: "generates all backing fields (all languages and generic, prefix, raw, and normalized) with correct properties (singlevalued field)",
			input: &fieldUtils.FieldDef{
				Name: testFieldName,
				Kind: KindName,
				IndexDef: &fieldUtils.IndexDef{
					MultiValued: false,
				},
			},
			want: solr.FieldCategoryToSolrFieldDefsMap{
				solr.GenericLangBaseFieldCategory: solr.FieldDef{
					Name:         genericName,
					Type:         solr.DefaultSolrTextFieldType,
					Stored:       true,
					Indexed:      false,
					MultiValued:  false,
					DocValues:    false,
					Uninvertible: false,
					Required:     false,
				},
				solr.GermanLangBaseFieldCategory: solr.FieldDef{
					Name:         deName,
					Type:         solr.DefaultDeSolrTextFieldType,
					Stored:       true,
					Indexed:      false,
					MultiValued:  false,
					DocValues:    false,
					Uninvertible: false,
					Required:     false,
				},
				solr.EnglishLangBaseFieldCategory: solr.FieldDef{
					Name:         enName,
					Type:         solr.DefaultEnSolrTextFieldType,
					Stored:       true,
					Indexed:      false,
					MultiValued:  false,
					DocValues:    false,
					Uninvertible: false,
					Required:     false,
				},
				solr.PrefixContentBaseFieldCategory: solr.FieldDef{
					Name:         prefixName,
					Type:         solr.DefaultPrefixSolrTextFieldType,
					Stored:       true,
					Indexed:      false,
					MultiValued:  true,
					DocValues:    false,
					Uninvertible: false,
					Required:     false,
				},
				solr.RawContentBaseFieldCategory: solr.FieldDef{
					Name:         rawName,
					Type:         solr.DefaultRawSolrTextFieldType,
					Stored:       true,
					Indexed:      false,
					MultiValued:  false,
					DocValues:    false,
					Uninvertible: false,
					Required:     false,
				},
				solr.NormalizedBaseFieldCategory: solr.FieldDef{
					Name:         solr.GetNormalizedBackingFieldName(testFieldName),
					Type:         solr.DefaultSolrStringFieldType,
					Stored:       true,
					Indexed:      false,
					MultiValued:  false,
					DocValues:    false,
					Uninvertible: false,
					Required:     false,
				},
			},
		},
		{
			name: "generates all backing fields (all languages and generic, prefix, raw, and normalized) with correct properties (multivalued field)",
			input: &fieldUtils.FieldDef{
				Name: testFieldName,
				Kind: KindName,
				IndexDef: &fieldUtils.IndexDef{
					MultiValued: true,
				},
			},
			want: solr.FieldCategoryToSolrFieldDefsMap{
				solr.GenericLangBaseFieldCategory: solr.FieldDef{
					Name:         genericName,
					Type:         solr.DefaultSolrTextFieldType,
					Stored:       true,
					Indexed:      false,
					MultiValued:  true,
					DocValues:    false,
					Uninvertible: false,
					Required:     false,
				},
				solr.GermanLangBaseFieldCategory: solr.FieldDef{
					Name:         deName,
					Type:         solr.DefaultDeSolrTextFieldType,
					Stored:       true,
					Indexed:      false,
					MultiValued:  true,
					DocValues:    false,
					Uninvertible: false,
					Required:     false,
				},
				solr.EnglishLangBaseFieldCategory: solr.FieldDef{
					Name:         enName,
					Type:         solr.DefaultEnSolrTextFieldType,
					Stored:       true,
					Indexed:      false,
					MultiValued:  true,
					DocValues:    false,
					Uninvertible: false,
					Required:     false,
				},
				solr.PrefixContentBaseFieldCategory: solr.FieldDef{
					Name:         prefixName,
					Type:         solr.DefaultPrefixSolrTextFieldType,
					Stored:       true,
					Indexed:      false,
					MultiValued:  true,
					DocValues:    false,
					Uninvertible: false,
					Required:     false,
				},
				solr.RawContentBaseFieldCategory: solr.FieldDef{
					Name:         rawName,
					Type:         solr.DefaultRawSolrTextFieldType,
					Stored:       true,
					Indexed:      false,
					MultiValued:  true,
					DocValues:    false,
					Uninvertible: false,
					Required:     false,
				},
				solr.NormalizedBaseFieldCategory: solr.FieldDef{
					Name:         solr.GetNormalizedBackingFieldName(testFieldName),
					Type:         solr.DefaultSolrStringFieldType,
					Stored:       true,
					Indexed:      false,
					MultiValued:  true,
					DocValues:    false,
					Uninvertible: false,
					Required:     false,
				},
			},
		},
	}

	for _, tt := range tests {
		kind := KindText{}
		t.Run(tt.name, func(t *testing.T) {
			fieldDef, _ := kind.ValidateDefinition(context.TODO(), tt.input)

			got, err := kind.GenerateSolrFields(context.TODO(), fieldDef)
			if tt.wantErr == (err == nil) {
				t.Errorf("GenerateSolrFields: wanted error (true/false) = %v, got error %v", tt.wantErr, err)
				return
			}
			if err != nil {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("solrFieldDefs: wanted %v but got %v", tt.want, got)
			}
		})
	}
}

func TestKindText_GenerateXMLFieldTags(t *testing.T) {
	testFieldName := "test"
	testFieldNameGeneric, _ := solr.GetLangSpecificFieldName(testFieldName, solr.GenericLangAbbrev)
	testFieldNameDe, _ := solr.GetLangSpecificFieldName(testFieldName, solr.GermanLangAbbrev)
	testFieldNameEn, _ := solr.GetLangSpecificFieldName(testFieldName, solr.EnglishLangAbbrev)
	testFieldNamePrefix := solr.GetPrefixBackingFieldName(testFieldName)
	testFieldNameNormed := solr.GetNormalizedBackingFieldName(testFieldName)
	testFieldNameRaw := solr.GetRawBackingFieldName(testFieldName)

	tests := []struct {
		name      string
		fieldDef  fields.BaseFieldDef
		itemValue datamodel.CurrentItemValue
		want      []string
		wantErr   bool
	}{
		{
			name:     "Places data in the generic language, prefix, normalized, and unanalyzed fields if no language is specified",
			fieldDef: nil,
			itemValue: datamodel.CurrentItemValue{
				ID:         "abc123",
				ItemID:     "def456",
				FieldName:  testFieldName,
				FieldValue: "hello",
			},
			want: []string{
				"<field name=\"" + testFieldNameGeneric + "\">hello</field>",
				"<field name=\"" + testFieldNamePrefix + "\">hello</field>",
				"<field name=\"" + testFieldNameRaw + "\">hello</field>",
				"<field name=\"" + testFieldNameNormed + "\">hello</field>",
			},
		},
		{
			name:     "Places data in generic language, prefix, normalized, and unanalyzed fields if an unknown language is specified",
			fieldDef: nil,
			itemValue: datamodel.CurrentItemValue{
				ID:         "abc123",
				ItemID:     "def456",
				FieldName:  testFieldName,
				FieldValue: "hello",
				Language: pgtype.Text{
					String: "dk",
					Valid:  true,
				},
			},
			want: []string{
				"<field name=\"" + testFieldNameGeneric + "\">hello</field>",
				"<field name=\"" + testFieldNamePrefix + "\">hello</field>",
				"<field name=\"" + testFieldNameRaw + "\">hello</field>",
				"<field name=\"" + testFieldNameNormed + "\">hello</field>",
			},
		},
		{
			name:     "Places data in DE, prefix, normalized, and unanalyzed fields if language is 'de'",
			fieldDef: nil,
			itemValue: datamodel.CurrentItemValue{
				ID:         "abc123",
				ItemID:     "def456",
				FieldName:  testFieldName,
				FieldValue: "hello",
				Language: pgtype.Text{
					String: solr.GermanLangAbbrev,
					Valid:  true,
				},
			},
			want: []string{
				"<field name=\"" + testFieldNameDe + "\">hello</field>",
				"<field name=\"" + testFieldNamePrefix + "\">hello</field>",
				"<field name=\"" + testFieldNameRaw + "\">hello</field>",
				"<field name=\"" + testFieldNameNormed + "\">hello</field>",
			},
		},
		{
			name:     "Places data in EN, prefix, normalized, and unanalyzed fields if language is 'en'",
			fieldDef: nil,
			itemValue: datamodel.CurrentItemValue{
				ID:         "abc123",
				ItemID:     "def456",
				FieldName:  testFieldName,
				FieldValue: "hello",
				Language: pgtype.Text{
					String: solr.EnglishLangAbbrev,
					Valid:  true,
				},
			},
			want: []string{
				"<field name=\"" + testFieldNameEn + "\">hello</field>",
				"<field name=\"" + testFieldNamePrefix + "\">hello</field>",
				"<field name=\"" + testFieldNameRaw + "\">hello</field>",
				"<field name=\"" + testFieldNameNormed + "\">hello</field>",
			},
		},
		{
			name:     "Does not store two copies of the data in the generic language field even if the language is explicitly set to generic",
			fieldDef: nil,
			itemValue: datamodel.CurrentItemValue{
				ID:         "abc123",
				ItemID:     "def456",
				FieldName:  testFieldName,
				FieldValue: "hello",
				Language: pgtype.Text{
					String: solr.GenericLangAbbrev,
					Valid:  true,
				},
			},
			want: []string{
				"<field name=\"" + testFieldNameGeneric + "\">hello</field>",
				"<field name=\"" + testFieldNamePrefix + "\">hello</field>",
				"<field name=\"" + testFieldNameRaw + "\">hello</field>",
				"<field name=\"" + testFieldNameNormed + "\">hello</field>",
			},
		},
		{
			name:     "Lowercases in the text in the normalized fied",
			fieldDef: nil,
			itemValue: datamodel.CurrentItemValue{
				ID:         "abc123",
				ItemID:     "def456",
				FieldName:  testFieldName,
				FieldValue: "REALLY REALLY loud TEXT",
			},
			want: []string{
				"<field name=\"" + testFieldNameGeneric + "\">REALLY REALLY loud TEXT</field>",
				"<field name=\"" + testFieldNamePrefix + "\">REALLY REALLY loud TEXT</field>",
				"<field name=\"" + testFieldNameRaw + "\">REALLY REALLY loud TEXT</field>",
				"<field name=\"" + testFieldNameNormed + "\">really really loud text</field>",
			},
		},
		{
			name:     "Replaces German character with their sort-equivalents in the text in the normalized fied",
			fieldDef: nil,
			itemValue: datamodel.CurrentItemValue{
				ID:         "abc123",
				ItemID:     "def456",
				FieldName:  testFieldName,
				FieldValue: "Örtlich MÜßE ähnlich",
			},
			want: []string{
				"<field name=\"" + testFieldNameGeneric + "\">Örtlich MÜßE ähnlich</field>",
				"<field name=\"" + testFieldNamePrefix + "\">Örtlich MÜßE ähnlich</field>",
				"<field name=\"" + testFieldNameRaw + "\">Örtlich MÜßE ähnlich</field>",
				"<field name=\"" + testFieldNameNormed + "\">ortlich musse ahnlich</field>",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ki := &KindText{}
			got, err := ki.GenerateXMLFieldTags(context.Background(), tt.fieldDef, tt.itemValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateXMLFieldTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateXMLFieldTags() got = %v, want %v", got, tt.want)
			}
		})
	}
}
