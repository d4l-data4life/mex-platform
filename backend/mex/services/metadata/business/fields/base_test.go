package fields

import (
	"reflect"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/d4l-data4life/mex/mex/shared/fields"
)

func TestValidateName(t *testing.T) {
	tests := []struct {
		name       string
		fieldNames []string
		wantErr    bool
	}{
		{
			name:       "Name containing underscores is invalid",
			fieldNames: []string{"_this_IS_1_valid_1234_Name"},
			wantErr:    true,
		},
		{
			name:       "Field name starting with a digit is invalid",
			fieldNames: []string{"1name"},
			wantErr:    true,
		},
		{
			name:       "Field name containing characters other than alphanumeric and underscore is invalid",
			fieldNames: []string{"some-name"},
			wantErr:    true,
		},
		{
			name:       "Field name containing double underscores is invalid",
			fieldNames: []string{"some__name"},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, n := range tt.fieldNames {
				if err := fields.ValidateName(n); (err != nil) != tt.wantErr {
					t.Errorf("ValidateName() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestValidateInternalName(t *testing.T) {
	tests := []struct {
		name       string
		fieldNames []string
		wantErr    bool
	}{
		{
			name:       "Name starting with an underscore is invalid",
			fieldNames: []string{"_thisIS1valid1234Name"},
			wantErr:    true,
		},
		{
			name:       "Name containing an underscore (not first character) is valid",
			fieldNames: []string{"this_IS_1_valid_1234_Name"},
			wantErr:    false,
		},
		{
			name:       "Field name starting with a digit is invalid",
			fieldNames: []string{"1name"},
			wantErr:    true,
		},
		{
			name:       "Field name containing characters other than alphanumeric and underscore is invalid",
			fieldNames: []string{"some-name"},
			wantErr:    true,
		},
		{
			name:       "Field name containing double underscores is valid",
			fieldNames: []string{"some__name"},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, n := range tt.fieldNames {
				if err := fields.ValidateInternalName(n); (err != nil) != tt.wantErr {
					t.Errorf("ValidateName() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

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

func TestGetFirstLinkExt(t *testing.T) {
	tests := []struct {
		name     string
		indexDef *fields.IndexDef
		want     *fields.IndexDefExtLink
		wantErr  bool
	}{
		{
			name: "Return error if no link extension is found",
			indexDef: &fields.IndexDef{
				Ext: toAnySlice(
					&fields.IndexDefExtCoding{
						CodingsetNames: []string{"mesh"},
					},
				),
			},
			wantErr: true,
		},
		{
			name: "If single link extension is present, it is returned",
			indexDef: &fields.IndexDef{
				Ext: toAnySlice(
					&fields.IndexDefExtLink{RelationType: "someType"},
				),
			},
			want: &fields.IndexDefExtLink{
				RelationType:       "someType",
				LinkedTargetFields: nil,
			},
		},
		{
			name: "Other extensions are ignored",
			indexDef: &fields.IndexDef{
				Ext: toAnySlice(
					&fields.IndexDefExtHierarchy{CodeSystemNameOrNodeEntityType: "orgUnits"},
					&fields.IndexDefExtCoding{
						CodingsetNames: []string{"mesh"},
					},
					&fields.IndexDefExtLink{RelationType: "someType"},
				),
			},
			want: &fields.IndexDefExtLink{
				RelationType:       "someType",
				LinkedTargetFields: nil,
			},
		},
		{
			name: "If multiple link extensions are present, the first is returned",
			indexDef: &fields.IndexDef{
				Ext: toAnySlice(
					&fields.IndexDefExtHierarchy{CodeSystemNameOrNodeEntityType: "orgUnits"},
					&fields.IndexDefExtLink{RelationType: "firstType"},
					&fields.IndexDefExtCoding{
						CodingsetNames: []string{"mesh"},
					},
					&fields.IndexDefExtLink{RelationType: "secondType"},
				),
			},
			want: &fields.IndexDefExtLink{
				RelationType:       "firstType",
				LinkedTargetFields: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFirstLinkExt(tt.indexDef)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFirstLinkExt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && (!reflect.DeepEqual(got.LinkedTargetFields, tt.want.LinkedTargetFields) || got.RelationType != tt.want.RelationType) {
				t.Errorf("GetFirstLinkExt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFirstHierarchyExt(t *testing.T) {
	tests := []struct {
		name     string
		indexDef *fields.IndexDef
		want     *fields.IndexDefExtHierarchy
		wantErr  bool
	}{
		{
			name: "Return error if no hierarchy extension is found",
			indexDef: &fields.IndexDef{
				Ext: toAnySlice(
					&fields.IndexDefExtLink{RelationType: "firstType"},
					&fields.IndexDefExtCoding{
						CodingsetNames: []string{"mesh"},
					},
				),
			},
			wantErr: true,
		},
		{
			name: "If single hierarchy extension is present, it is returned",
			indexDef: &fields.IndexDef{
				Ext: toAnySlice(
					&fields.IndexDefExtHierarchy{CodeSystemNameOrNodeEntityType: "orgUnits"},
				),
			},
			want: &fields.IndexDefExtHierarchy{
				CodeSystemNameOrNodeEntityType: "orgUnits",
				LinkFieldName:                  "",
				DisplayFieldName:               "",
			},
		},
		{
			name: "Other extensions are ignored",
			indexDef: &fields.IndexDef{
				Ext: toAnySlice(
					&fields.IndexDefExtLink{RelationType: "firstType"},
					&fields.IndexDefExtCoding{
						CodingsetNames: []string{"mesh"},
					},
					&fields.IndexDefExtHierarchy{CodeSystemNameOrNodeEntityType: "orgUnits"},
				),
			},
			want: &fields.IndexDefExtHierarchy{
				CodeSystemNameOrNodeEntityType: "orgUnits",
				LinkFieldName:                  "",
				DisplayFieldName:               "",
			},
		},
		{
			name: "If multiple hierarchy extensions are present, the first is returned",
			indexDef: &fields.IndexDef{
				Ext: toAnySlice(
					&fields.IndexDefExtHierarchy{CodeSystemNameOrNodeEntityType: "orgUnits"},
					&fields.IndexDefExtLink{RelationType: "firstType"},
					&fields.IndexDefExtCoding{
						CodingsetNames: []string{"mesh"},
					},
					&fields.IndexDefExtHierarchy{CodeSystemNameOrNodeEntityType: "otherOrgUnits"},
				),
			},
			want: &fields.IndexDefExtHierarchy{
				CodeSystemNameOrNodeEntityType: "orgUnits",
				LinkFieldName:                  "",
				DisplayFieldName:               "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFirstHierarchyExt(tt.indexDef)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFirstHierarchyExt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && (got.DisplayFieldName != tt.want.DisplayFieldName || got.LinkFieldName != tt.want.LinkFieldName || got.CodeSystemNameOrNodeEntityType != tt.want.CodeSystemNameOrNodeEntityType) {
				t.Errorf("GetFirstLinkExt() got = %v, want %v", got, tt.want)
			}
		})
	}
}
