package kind_timestamp

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
			name: "non-timestamp field causes an error",
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
			kind := &KindTimestamp{}
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
			name: "generate both a timestamp field and string field for the raw value (with correct propoerties)",
			input: &fieldUtils.FieldDef{
				Name:     "test",
				Kind:     KindName,
				IndexDef: &fieldUtils.IndexDef{},
			},
			want: solr.FieldCategoryToSolrFieldDefsMap{
				solr.GenericLangBaseFieldCategory: solr.FieldDef{
					Name:         "test",
					Type:         solr.DefaultSolrTimestampFieldType,
					Stored:       true,
					Indexed:      false,
					MultiValued:  false,
					DocValues:    false,
					Uninvertible: false,
				},
				solr.RawContentBaseFieldCategory: solr.FieldDef{
					Name:         solr.GetRawValTimestampName("test"),
					Type:         solr.DefaultSolrStringFieldType,
					Stored:       true,
					Indexed:      false,
					MultiValued:  false,
					DocValues:    false,
					Uninvertible: false,
				},
			},
		},
	}

	for _, tt := range tests {
		kind := KindTimestamp{}
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

func Test_GenerateXMLFieldTags(t *testing.T) {
	tests := []struct {
		name     string
		fieldDef *fieldUtils.FieldDef
		itemVal  datamodel.CurrentItemValue
		want     []string
	}{
		{
			name: "fills the index field with the base name with the timestamp and the corresponding raw value field" +
				" with the string pass stored in DB - full precision datetime string",
			fieldDef: &fieldUtils.FieldDef{
				Name:     "test",
				Kind:     KindName,
				IndexDef: &fieldUtils.IndexDef{},
			},
			itemVal: datamodel.CurrentItemValue{
				ID:         "abc",
				ItemID:     "123",
				FieldName:  "test",
				FieldValue: "2003-04-30T12:32:13Z",
				Place:      0,
				Revision:   0,
				Language:   pgtype.Text{},
			},
			want: []string{
				"<field name=\"test\">2003-04-30T12:32:13Z</field>",
				"<field name=\"test_raw_value\">2003-04-30T12:32:13Z</field>",
			},
		},
		{
			name: "fills the index field with the base name with the timestamp and the corresponding raw value field" +
				" with the string pass stored in DB - day-precision datetime string",
			fieldDef: &fieldUtils.FieldDef{
				Name:     "test",
				Kind:     KindName,
				IndexDef: &fieldUtils.IndexDef{},
			},
			itemVal: datamodel.CurrentItemValue{
				ID:         "abc",
				ItemID:     "123",
				FieldName:  "test",
				FieldValue: "2003-04-30",
				Place:      0,
				Revision:   0,
				Language:   pgtype.Text{},
			},
			want: []string{
				"<field name=\"test\">2003-04-30T00:00:00Z</field>",
				"<field name=\"test_raw_value\">2003-04-30</field>",
			},
		},
		{
			name: "fills the index field with the base name with the timestamp and the corresponding raw value field" +
				" with the string pass stored in DB - month-precision datetime string",
			fieldDef: &fieldUtils.FieldDef{
				Name:     "test",
				Kind:     KindName,
				IndexDef: &fieldUtils.IndexDef{},
			},
			itemVal: datamodel.CurrentItemValue{
				ID:         "abc",
				ItemID:     "123",
				FieldName:  "test",
				FieldValue: "2003-04",
				Place:      0,
				Revision:   0,
				Language:   pgtype.Text{},
			},
			want: []string{
				"<field name=\"test\">2003-04-01T00:00:00Z</field>",
				"<field name=\"test_raw_value\">2003-04</field>",
			},
		},
		{
			name: "fills the index field with the base name with the timestamp and the corresponding raw value field" +
				" with the string pass stored in DB - year-precision datetime string",
			fieldDef: &fieldUtils.FieldDef{
				Name:     "test",
				Kind:     KindName,
				IndexDef: &fieldUtils.IndexDef{},
			},
			itemVal: datamodel.CurrentItemValue{
				ID:         "abc",
				ItemID:     "123",
				FieldName:  "test",
				FieldValue: "2003",
				Place:      0,
				Revision:   0,
				Language:   pgtype.Text{},
			},
			want: []string{
				"<field name=\"test\">2003-01-01T00:00:00Z</field>",
				"<field name=\"test_raw_value\">2003</field>",
			},
		},
	}

	for _, tt := range tests {
		kind := KindTimestamp{}
		t.Run(tt.name, func(t *testing.T) {
			fieldDef, err := kind.ValidateDefinition(context.TODO(), tt.fieldDef)
			if err != nil {
				t.Error(err)
			}

			x, err := kind.GenerateXMLFieldTags(context.TODO(), fieldDef, tt.itemVal)
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(x, tt.want) {
				t.Errorf("solrFieldDefs: wanted %v but got %v", tt.want, x)
			}
		})
	}
}
