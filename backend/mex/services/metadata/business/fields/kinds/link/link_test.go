package kind_link

import (
	"context"
	"reflect"
	"testing"

	"github.com/d4l-data4life/mex/mex/shared/solr"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	fieldUtils "github.com/d4l-data4life/mex/mex/shared/fields"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
)

func toAnySlice(msg ...proto.Message) []*anypb.Any {
	anys := make([]*anypb.Any, len(msg))

	for i, m := range msg {
		anyOne, err := anypb.New(m)
		if err != nil {
			panic(err)
		}

		anys[i] = anyOne
	}

	return anys
}

func TestKindLink_MarshalToProtobufFormat(t *testing.T) {

	fieldDef := &fieldUtils.FieldDef{
		Name:      "testField",
		Kind:      KindName,
		DisplayId: "SOMETHING",
		IndexDef: &fieldUtils.IndexDef{
			MultiValued: true,
			Ext: toAnySlice(&fieldUtils.IndexDefExtLink{
				RelationType: "someType",
				LinkedTargetFields: []string{
					"linkedField1",
					"linkedField2",
				},
			}),
		},
	}
	linkField, _ := (&KindLink{}).ValidateDefinition(context.TODO(), fieldDef)

	tests := []struct {
		name     string
		fieldDef fields.BaseFieldDef
		want     *fieldUtils.FieldDef
		wantErr  bool
	}{
		{
			name:     "Correctly includes all fields",
			fieldDef: linkField,
			want:     fieldDef,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kind := &KindLink{}
			got, err := kind.MarshalToProtobufFormat(context.TODO(), tt.fieldDef)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalToProtobufFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalToProtobufFormat() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GenerateSolrFields(t *testing.T) {
	genericTypeName := solr.KnownLanguagesCategoryMap[solr.GenericLangAbbrev]
	tests := []struct {
		name  string
		input *fieldUtils.FieldDef

		wantErr bool
		want    solr.FieldCategoryToSolrFieldDefsMap
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
			name: "generates a string field with the right properties",
			input: &fieldUtils.FieldDef{
				Name: "test",
				Kind: KindName,
				IndexDef: &fieldUtils.IndexDef{
					MultiValued: false,
					Ext: toAnySlice(&fieldUtils.IndexDefExtLink{
						RelationType: "someType",
						LinkedTargetFields: []string{
							"linkedField1",
							"linkedField2",
						},
					}),
				},
			},
			want: solr.FieldCategoryToSolrFieldDefsMap{
				genericTypeName: solr.FieldDef{
					Name:         "test",
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
		kind := KindLink{}
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
