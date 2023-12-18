package linked

import (
	"context"
	"sort"
	"strings"
	"testing"

	"google.golang.org/protobuf/types/known/anypb"

	fieldUtils "github.com/d4l-data4life/mex/mex/shared/fields"
	"github.com/d4l-data4life/mex/mex/shared/solr"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
	kindHierarchy "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/hierarchy"
	kindLink "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/link"
	kindString "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/string"
)

func Test_GetExactFieldName(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Adds a post-fix to the input",
			input: "hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := solr.GetExactFieldName(tt.input)
			if !strings.HasPrefix(got, tt.input) {
				t.Errorf("GetExactFieldName(): output '%s' did not start with the input '%s'", got, tt.input)
			}
		})
	}
}

func Test_GetLinkedFieldName(t *testing.T) {
	tests := []struct {
		name   string
		link   string
		target string
		want   string
	}{
		{
			name:   "Glues together link and target with the double underscore",
			link:   "hello",
			target: "there",
			want:   "hello__there",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := solr.GetLinkedFieldName(tt.link, tt.target)
			if got != tt.want {
				t.Errorf("Wanted the linked field name '%s' but got '%s'", tt.want, got)
			}
		})
	}
}

func Test_GetLangSpecificFieldName(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		langCode  string
		want      string
		wantErr   bool
	}{
		{
			name:      "Accepts 'de' as language code and appends it to name with triple-underscore",
			fieldName: "hello",
			langCode:  solr.GermanLangAbbrev,
			want:      "hello___de",
		},
		{
			name:      "Accepts 'en' as language code and appends it to name with triple-underscore",
			fieldName: "hello",
			langCode:  solr.EnglishLangAbbrev,
			want:      "hello___en",
		},
		{
			name:      "Appends '___generic 'it to name if an empty string is passed as language code",
			fieldName: "hello",
			langCode:  solr.GenericLangAbbrev,
			want:      "hello___" + solr.GenericLanguageSuffix,
		},
		{
			name:      "Returns error when 'fr' is passed as language code",
			fieldName: "hello",
			langCode:  "fr",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := solr.GetLangSpecificFieldName(tt.fieldName, tt.langCode)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("GetLangSpecificFieldName(): wantErr '%v' but got '%s'", tt.wantErr, gotErr)
			} else if got != tt.want {
				t.Errorf("GetLangSpecificFieldName(): wanted '%s' but got '%s'", tt.want, got)
			}
		})
	}
}

func Test_GetLinkedFieldDefsLinkFields(t *testing.T) {
	singleValuedField1 := (&kindString.KindString{}).MustValidateDefinition(
		context.Background(),
		&fieldUtils.FieldDef{
			Name: "singleValuedField1",
			Kind: "string",
			IndexDef: &fieldUtils.IndexDef{
				MultiValued: false,
			},
		})
	singleValuedField2 := (&kindString.KindString{}).MustValidateDefinition(
		context.Background(),
		&fieldUtils.FieldDef{
			Name: "singleValuedField2",
			Kind: "string",
			IndexDef: &fieldUtils.IndexDef{
				MultiValued: false,
			},
		})
	multiValuedField := (&kindString.KindString{}).MustValidateDefinition(
		context.Background(),
		&fieldUtils.FieldDef{
			Name: "multiValuedField",
			Kind: "string",
			IndexDef: &fieldUtils.IndexDef{
				MultiValued: true,
			},
		})
	linkExtLinkingToOneSingleValuedField, _ := anypb.New(&fieldUtils.IndexDefExtLink{
		RelationType:       "link",
		LinkedTargetFields: []string{"singleValuedField1"},
	})
	singleValuedLinkFieldLinkingToOneSingleValuedField := (&kindLink.KindLink{}).MustValidateDefinition(
		context.Background(),
		&fieldUtils.FieldDef{
			Name: "singleValuedLinkFieldLinkingToOneSingleValuedField",
			Kind: "link",
			IndexDef: &fieldUtils.IndexDef{
				MultiValued: false,
				Ext:         []*anypb.Any{linkExtLinkingToOneSingleValuedField},
			}})
	linkExtLinkingToSeveralSingleValuedFields, _ := anypb.New(&fieldUtils.IndexDefExtLink{
		RelationType:       "link",
		LinkedTargetFields: []string{"singleValuedField1", "singleValuedField2"},
	})
	singleValuedLinkFieldLinkingToSeveralSingleValuedFields := (&kindLink.KindLink{}).MustValidateDefinition(
		context.Background(),
		&fieldUtils.FieldDef{
			Name: "singleValuedLinkFieldLinkingToSeveralSingleValuedFields",
			Kind: "link",
			IndexDef: &fieldUtils.IndexDef{
				MultiValued: false,
				Ext:         []*anypb.Any{linkExtLinkingToSeveralSingleValuedFields},
			}})
	multiValuedLinkFieldLinkingToOneSingleValuedField := (&kindLink.KindLink{}).MustValidateDefinition(
		context.Background(),
		&fieldUtils.FieldDef{
			Name: "multiValuedLinkFieldLinkingToOneSingleValuedField",
			Kind: "link",
			IndexDef: &fieldUtils.IndexDef{
				MultiValued: true,
				Ext:         []*anypb.Any{linkExtLinkingToOneSingleValuedField},
			}})
	linkExtLinkingToOneMultiValuedField, _ := anypb.New(&fieldUtils.IndexDefExtLink{
		RelationType:       "link",
		LinkedTargetFields: []string{"multiValuedField"},
	})
	singleValuedLinkFieldLinkingToOneMultiValuedField := (&kindLink.KindLink{}).MustValidateDefinition(
		context.Background(),
		&fieldUtils.FieldDef{
			Name: "singleValuedLinkFieldLinkingToOneMultiValuedField",
			Kind: "link",
			IndexDef: &fieldUtils.IndexDef{
				MultiValued: false,
				Ext:         []*anypb.Any{linkExtLinkingToOneMultiValuedField},
			}})
	multiValuedLinkFieldLinkingToSeveralSingleValuedFields := (&kindLink.KindLink{}).MustValidateDefinition(
		context.Background(),
		&fieldUtils.FieldDef{
			Name: "multiValuedLinkFieldLinkingToSeveralSingleValuedFields",
			Kind: "link",
			IndexDef: &fieldUtils.IndexDef{
				MultiValued: true,
				Ext:         []*anypb.Any{linkExtLinkingToSeveralSingleValuedFields},
			}})
	multiValuedLinkFieldLinkingToOneMultiValuedField := (&kindLink.KindLink{}).MustValidateDefinition(
		context.Background(),
		&fieldUtils.FieldDef{
			Name: "multiValuedLinkFieldLinkingToOneMultiValuedField",
			Kind: "link",
			IndexDef: &fieldUtils.IndexDef{
				MultiValued: true,
				Ext:         []*anypb.Any{linkExtLinkingToOneMultiValuedField},
			}})
	invalidLinkExt, _ := anypb.New(
		&fieldUtils.IndexDefExtLink{
			RelationType:       "link",
			LinkedTargetFields: []string{"singleValuedField1", "notAField"},
		})
	invalidLinkField := (&kindLink.KindLink{}).MustValidateDefinition(
		context.Background(),
		&fieldUtils.FieldDef{
			Name: "invalidLink",
			Kind: "link",
			IndexDef: &fieldUtils.IndexDef{
				MultiValued: false,
				Ext:         []*anypb.Any{invalidLinkExt},
			}})

	// Expected outputs
	singleValuedLinkFieldLinkingToOneSingleValuedField_generated := fields.NewBaseFieldDef(
		"singleValuedLinkFieldLinkingToOneSingleValuedField__singleValuedField1",
		singleValuedField1.Kind(),
		solr.GetLinkedDisplayID("singleValuedLinkFieldLinkingToOneSingleValuedField__singleValuedField1"),
		true,
		fields.BaseIndexDef{
			MultiValued: false,
		})
	fieldWithOneInvalidTargetLinkFieldLinkingToOneSingleValuedField_generated := fields.NewBaseFieldDef(
		"invalidLink__singleValuedField1",
		singleValuedField1.Kind(),
		solr.GetLinkedDisplayID("invalidLink__singleValuedField1"),
		true,
		fields.BaseIndexDef{
			MultiValued: false,
		})
	singleValuedLinkFieldLinkingToSeveralSingleValuedFields_generated1 := fields.NewBaseFieldDef(
		"singleValuedLinkFieldLinkingToSeveralSingleValuedFields__singleValuedField1",
		singleValuedField1.Kind(),
		solr.GetLinkedDisplayID("singleValuedLinkFieldLinkingToSeveralSingleValuedFields__singleValuedField1"),
		true,
		fields.BaseIndexDef{
			MultiValued: false,
		})
	singleValuedLinkFieldLinkingToSeveralSingleValuedFields_generated2 := fields.NewBaseFieldDef(
		"singleValuedLinkFieldLinkingToSeveralSingleValuedFields__singleValuedField2",
		singleValuedField2.Kind(),
		solr.GetLinkedDisplayID("singleValuedLinkFieldLinkingToSeveralSingleValuedFields__singleValuedField2"),
		true,
		fields.BaseIndexDef{
			MultiValued: false,
		})
	multiValuedLinkFieldLinkingToOneSingleValuedField_generated := fields.NewBaseFieldDef(
		"multiValuedLinkFieldLinkingToOneSingleValuedField__singleValuedField1",
		singleValuedField1.Kind(),
		solr.GetLinkedDisplayID("multiValuedLinkFieldLinkingToOneSingleValuedField__singleValuedField1"),
		true,
		fields.BaseIndexDef{
			MultiValued: true,
		})
	singleValuedLinkFieldLinkingToOneMultiValuedField_generated := fields.NewBaseFieldDef(
		"singleValuedLinkFieldLinkingToOneMultiValuedField__multiValuedField",
		multiValuedField.Kind(),
		solr.GetLinkedDisplayID("singleValuedLinkFieldLinkingToOneMultiValuedField__multiValuedField"),
		true,
		fields.BaseIndexDef{
			MultiValued: true,
		})
	multiValuedLinkFieldLinkingToSeveralSingleValuedFields_generated1 := fields.NewBaseFieldDef(
		"multiValuedLinkFieldLinkingToSeveralSingleValuedFields__singleValuedField1",
		singleValuedField1.Kind(),
		solr.GetLinkedDisplayID("multiValuedLinkFieldLinkingToSeveralSingleValuedFields__singleValuedField1"),
		true,
		fields.BaseIndexDef{
			MultiValued: true,
		})
	multiValuedLinkFieldLinkingToSeveralSingleValuedFields_generated2 := fields.NewBaseFieldDef(
		"multiValuedLinkFieldLinkingToSeveralSingleValuedFields__singleValuedField2",
		singleValuedField2.Kind(),
		solr.GetLinkedDisplayID("multiValuedLinkFieldLinkingToSeveralSingleValuedFields__singleValuedField2"),
		true,
		fields.BaseIndexDef{
			MultiValued: true,
		})
	multiValuedLinkFieldLinkingToOneMultiValuedField_generated := fields.NewBaseFieldDef(
		"multiValuedLinkFieldLinkingToOneMultiValuedField__multiValuedField",
		multiValuedField.Kind(),
		solr.GetLinkedDisplayID("multiValuedLinkFieldLinkingToOneMultiValuedField__multiValuedField"),
		true,
		fields.BaseIndexDef{
			MultiValued: true,
		})

	tests := []struct {
		name      string
		fieldDefs []fields.BaseFieldDef
		want      []fields.BaseFieldDef
		wantErr   bool
	}{
		{
			name: "Single-valued link field with a single, " +
				"single-valued target field generates a single-valued linked field",
			fieldDefs: []fields.BaseFieldDef{
				singleValuedField1,
				singleValuedField2,
				multiValuedField,
				singleValuedLinkFieldLinkingToOneSingleValuedField,
			},
			want: []fields.BaseFieldDef{
				singleValuedLinkFieldLinkingToOneSingleValuedField_generated,
			},
		},
		{
			name: "Single-valued link field with multiple, " +
				"single-valued target fields generates one single-valued linked field for each target",
			fieldDefs: []fields.BaseFieldDef{
				singleValuedField1,
				singleValuedField2,
				multiValuedField,
				singleValuedLinkFieldLinkingToSeveralSingleValuedFields,
			},
			want: []fields.BaseFieldDef{
				singleValuedLinkFieldLinkingToSeveralSingleValuedFields_generated1,
				singleValuedLinkFieldLinkingToSeveralSingleValuedFields_generated2,
			},
		},
		{
			name: "Multi-valued link field with a single, " +
				"single-valued target field generates a single, single-valued linked field",
			fieldDefs: []fields.BaseFieldDef{
				singleValuedField1,
				singleValuedField2,
				multiValuedField,
				multiValuedLinkFieldLinkingToOneSingleValuedField,
			},
			want: []fields.BaseFieldDef{
				multiValuedLinkFieldLinkingToOneSingleValuedField_generated,
			},
		},
		{
			name: "Single-valued link field with a single, " +
				"multi-valued target field generates a single, multi-valued linked field",
			fieldDefs: []fields.BaseFieldDef{
				singleValuedField1,
				singleValuedField2,
				multiValuedField,
				singleValuedLinkFieldLinkingToOneMultiValuedField,
			},
			want: []fields.BaseFieldDef{
				singleValuedLinkFieldLinkingToOneMultiValuedField_generated,
			},
		},
		{
			name: "Multi-valued link field with several, " +
				"single-valued target fields generates a multi-valued linked field for each target",
			fieldDefs: []fields.BaseFieldDef{
				singleValuedField1,
				singleValuedField2,
				multiValuedField,
				multiValuedLinkFieldLinkingToSeveralSingleValuedFields,
			},
			want: []fields.BaseFieldDef{
				multiValuedLinkFieldLinkingToSeveralSingleValuedFields_generated1,
				multiValuedLinkFieldLinkingToSeveralSingleValuedFields_generated2,
			},
		},
		{
			name: "Multi-valued link field with a single, " +
				"multi-valued target field generates a single, multi-valued linked field",
			fieldDefs: []fields.BaseFieldDef{
				singleValuedField1,
				singleValuedField2,
				multiValuedField,
				multiValuedLinkFieldLinkingToOneMultiValuedField,
			},
			want: []fields.BaseFieldDef{
				multiValuedLinkFieldLinkingToOneMultiValuedField_generated,
			},
		},
		{
			name: "If a link field specifies a target field that is not configured, it is ignored",
			fieldDefs: []fields.BaseFieldDef{
				singleValuedField1,
				singleValuedField2,
				multiValuedField,
				invalidLinkField,
			},
			want: []fields.BaseFieldDef{
				fieldWithOneInvalidTargetLinkFieldLinkingToOneSingleValuedField_generated,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := GetLinkedFieldDefs(tt.fieldDefs)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("getSearchFocusFields(), wanted error was %v, but the error returned was %v", tt.wantErr, gotErr)
			} else {
				if len(got) != len(tt.want) {
					t.Errorf("Expected %d generated linked fields but found %d", len(tt.want), len(got))
				}
				sort.Slice(got, func(i, j int) bool {
					return got[i].Name() > got[j].Name()
				})
				sort.Slice(tt.want, func(i, j int) bool {
					return tt.want[i].Name() > tt.want[j].Name()
				})
				for i := range got {
					if got[i].Name() != tt.want[i].Name() {
						t.Errorf("Names do not match - got %s but expected %s", got[i].Name(), tt.want[i].Name())
					}
					if got[i].Kind() != tt.want[i].Kind() {
						t.Errorf("Kinds do not match - got %s but expected %s", got[i].Kind(), tt.want[i].Kind())
					}
				}
			}
		})
	}
}

func Test_GetLinkedFieldDefsHierarchyFields(t *testing.T) {
	singleValuedField1 := (&kindString.KindString{}).MustValidateDefinition(
		context.Background(),
		&fieldUtils.FieldDef{
			Name: "singleValuedField1",
			Kind: "string",
			IndexDef: &fieldUtils.IndexDef{
				MultiValued: false,
			},
		})
	singleValuedField2 := (&kindString.KindString{}).MustValidateDefinition(
		context.Background(),
		&fieldUtils.FieldDef{
			Name: "singleValuedField2",
			Kind: "string",
			IndexDef: &fieldUtils.IndexDef{
				MultiValued: false,
			},
		})
	multiValuedField := (&kindString.KindString{}).MustValidateDefinition(
		context.Background(),
		&fieldUtils.FieldDef{
			Name: "multiValuedField",
			Kind: "string",
			IndexDef: &fieldUtils.IndexDef{
				MultiValued: true,
			},
		})
	linkExtLinkingToOneSingleValuedField, _ := anypb.New(&fieldUtils.IndexDefExtLink{
		RelationType:       "hierarchy",
		LinkedTargetFields: []string{"singleValuedField1"},
	})
	hierarchyExt, _ := anypb.New(&fieldUtils.IndexDefExtHierarchy{
		CodeSystemNameOrNodeEntityType: "mesh",
	})

	hierarchyInstance, _ := kindHierarchy.NewKindHierarchy(nil)

	singleValuedHierarchyFieldLinkingToOneSingleValuedField := hierarchyInstance.MustValidateDefinition(
		context.Background(),
		&fieldUtils.FieldDef{
			Name: "singleValuedHierarchyFieldLinkingToOneSingleValuedField",
			Kind: "hierarchy",
			IndexDef: &fieldUtils.IndexDef{
				MultiValued: false,
				Ext: []*anypb.Any{
					linkExtLinkingToOneSingleValuedField,
					hierarchyExt,
				},
			}})
	extLinkingToSeveralSingleValuedFields, _ := anypb.New(&fieldUtils.IndexDefExtLink{
		RelationType:       "link",
		LinkedTargetFields: []string{"singleValuedField1", "singleValuedField2"},
	})

	singleValuedHierarchyFieldLinkingToSeveralSingleValuedFields := hierarchyInstance.MustValidateDefinition(
		context.Background(),
		&fieldUtils.FieldDef{
			Name: "singleValuedHierarchyFieldLinkingToSeveralSingleValuedFields",
			Kind: "hierarchy",
			IndexDef: &fieldUtils.IndexDef{
				MultiValued: false,
				Ext: []*anypb.Any{
					extLinkingToSeveralSingleValuedFields,
					hierarchyExt,
				},
			}})
	multiValuedHierarchyFieldLinkingToOneSingleValuedField := hierarchyInstance.MustValidateDefinition(
		context.Background(),
		&fieldUtils.FieldDef{
			Name: "multiValuedHierarchyFieldLinkingToOneSingleValuedField",
			Kind: "hierarchy",
			IndexDef: &fieldUtils.IndexDef{
				MultiValued: true,
				Ext: []*anypb.Any{
					linkExtLinkingToOneSingleValuedField,
					hierarchyExt,
				},
			}})
	extLinkingToOneMultiValuedField, _ := anypb.New(&fieldUtils.IndexDefExtLink{
		RelationType:       "link",
		LinkedTargetFields: []string{"multiValuedField"},
	})
	singleValuedHierarchyFieldLinkingToOneMultiValuedField := hierarchyInstance.MustValidateDefinition(
		context.Background(),
		&fieldUtils.FieldDef{
			Name: "singleValuedHierarchyFieldLinkingToOneMultiValuedField",
			Kind: "hierarchy",
			IndexDef: &fieldUtils.IndexDef{
				MultiValued: false,
				Ext: []*anypb.Any{
					extLinkingToOneMultiValuedField,
					hierarchyExt,
				},
			}})
	multiValuedHierarchyFieldLinkingToSeveralSingleValuedFields := hierarchyInstance.MustValidateDefinition(
		context.Background(),
		&fieldUtils.FieldDef{
			Name: "multiValuedHierarchyFieldLinkingToSeveralSingleValuedFields",
			Kind: "hierarchy",
			IndexDef: &fieldUtils.IndexDef{
				MultiValued: true,
				Ext: []*anypb.Any{
					extLinkingToSeveralSingleValuedFields,
					hierarchyExt,
				},
			}})
	multiValuedHierarchyFieldLinkingToOneMultiValuedField := hierarchyInstance.MustValidateDefinition(
		context.Background(),
		&fieldUtils.FieldDef{
			Name: "multiValuedHierarchyFieldLinkingToOneMultiValuedField",
			Kind: "hierarchy",
			IndexDef: &fieldUtils.IndexDef{
				MultiValued: true,
				Ext: []*anypb.Any{
					extLinkingToOneMultiValuedField,
					hierarchyExt,
				},
			}})
	invalidLinkExt, _ := anypb.New(
		&fieldUtils.IndexDefExtLink{
			RelationType:       "link",
			LinkedTargetFields: []string{"singleValuedField1", "notAField"},
		})
	invalidHierarchyField := hierarchyInstance.MustValidateDefinition(
		context.Background(),
		&fieldUtils.FieldDef{
			Name: "invalidHierarchy",
			Kind: "hierarchy",
			IndexDef: &fieldUtils.IndexDef{
				MultiValued: false,
				Ext: []*anypb.Any{
					invalidLinkExt,
					hierarchyExt,
				},
			}})

	// Expected outputs
	singleValuedHierarchyFieldLinkingToOneSingleValuedField_generated := fields.NewBaseFieldDef(
		"singleValuedHierarchyFieldLinkingToOneSingleValuedField__singleValuedField1",
		singleValuedField1.Kind(),
		solr.GetLinkedDisplayID("singleValuedHierarchyFieldLinkingToOneSingleValuedField__singleValuedField1"),
		true,
		fields.BaseIndexDef{
			MultiValued: false,
		})
	fieldWithOneInvalidTargetHierarchyFieldLinkingToOneSingleValuedField_generated := fields.NewBaseFieldDef(
		"invalidHierarchy__singleValuedField1",
		singleValuedField1.Kind(),
		solr.GetLinkedDisplayID("invalidLink__singleValuedField1"),
		true,
		fields.BaseIndexDef{
			MultiValued: false,
		})
	singleValuedHierarchyFieldLinkingToSeveralSingleValuedFields_generated1 := fields.NewBaseFieldDef(
		"singleValuedHierarchyFieldLinkingToSeveralSingleValuedFields__singleValuedField1",
		singleValuedField1.Kind(),
		solr.GetLinkedDisplayID("singleValuedHierarchyFieldLinkingToSeveralSingleValuedFields__singleValuedField1"),
		true,
		fields.BaseIndexDef{
			MultiValued: false,
		})
	singleValuedHierarchyFieldLinkingToSeveralSingleValuedFields_generated2 := fields.NewBaseFieldDef(
		"singleValuedHierarchyFieldLinkingToSeveralSingleValuedFields__singleValuedField2",
		singleValuedField2.Kind(),
		solr.GetLinkedDisplayID("singleValuedHierarchyFieldLinkingToSeveralSingleValuedFields__singleValuedField2"),
		true,
		fields.BaseIndexDef{
			MultiValued: false,
		})
	multiValuedHierarchyFieldLinkingToOneSingleValuedField_generated := fields.NewBaseFieldDef(
		"multiValuedHierarchyFieldLinkingToOneSingleValuedField__singleValuedField1",
		singleValuedField1.Kind(),
		solr.GetLinkedDisplayID("multiValuedHierarchyFieldLinkingToOneSingleValuedField__singleValuedField1"),
		true,
		fields.BaseIndexDef{
			MultiValued: true,
		})
	singleValuedHierarchyFieldLinkingToOneMultiValuedField_generated := fields.NewBaseFieldDef(
		"singleValuedHierarchyFieldLinkingToOneMultiValuedField__multiValuedField",
		multiValuedField.Kind(),
		solr.GetLinkedDisplayID("singleValuedHierarchyFieldLinkingToOneMultiValuedField__multiValuedField"),
		true,
		fields.BaseIndexDef{
			MultiValued: true,
		})
	multiValuedHierarchyFieldLinkingToSeveralSingleValuedFields_generated1 := fields.NewBaseFieldDef(
		"multiValuedHierarchyFieldLinkingToSeveralSingleValuedFields__singleValuedField1",
		singleValuedField1.Kind(),
		solr.GetLinkedDisplayID("multiValuedHierarchyFieldLinkingToSeveralSingleValuedFields__singleValuedField1"),
		true,
		fields.BaseIndexDef{
			MultiValued: true,
		})
	multiValuedHierarchyFieldLinkingToSeveralSingleValuedFields_generated2 := fields.NewBaseFieldDef(
		"multiValuedHierarchyFieldLinkingToSeveralSingleValuedFields__singleValuedField2",
		singleValuedField2.Kind(),
		solr.GetLinkedDisplayID("multiValuedHierarchyFieldLinkingToSeveralSingleValuedFields__singleValuedField2"),
		true,
		fields.BaseIndexDef{
			MultiValued: true,
		})
	multiValuedHierarchyFieldLinkingToOneMultiValuedField_generated := fields.NewBaseFieldDef(
		"multiValuedHierarchyFieldLinkingToOneMultiValuedField__multiValuedField",
		multiValuedField.Kind(),
		solr.GetLinkedDisplayID("multiValuedHierarchyFieldLinkingToOneMultiValuedField__multiValuedField"),
		true,
		fields.BaseIndexDef{
			MultiValued: true,
		})

	tests := []struct {
		name      string
		fieldDefs []fields.BaseFieldDef
		want      []fields.BaseFieldDef
		wantErr   bool
	}{
		{
			name: "Single-valued hierarchy field with a single, " +
				"single-valued target field generates a single-valued linked field",
			fieldDefs: []fields.BaseFieldDef{
				singleValuedField1,
				singleValuedField2,
				multiValuedField,
				singleValuedHierarchyFieldLinkingToOneSingleValuedField,
			},
			want: []fields.BaseFieldDef{
				singleValuedHierarchyFieldLinkingToOneSingleValuedField_generated,
			},
		},
		{
			name: "Single-valued hierarchy field with multiple, " +
				"single-valued target fields generates one single-valued linked field for each target",
			fieldDefs: []fields.BaseFieldDef{
				singleValuedField1,
				singleValuedField2,
				multiValuedField,
				singleValuedHierarchyFieldLinkingToSeveralSingleValuedFields,
			},
			want: []fields.BaseFieldDef{
				singleValuedHierarchyFieldLinkingToSeveralSingleValuedFields_generated1,
				singleValuedHierarchyFieldLinkingToSeveralSingleValuedFields_generated2,
			},
		},
		{
			name: "Multi-valued hierarchy field with a single, " +
				"single-valued target field generates a single, single-valued linked field",
			fieldDefs: []fields.BaseFieldDef{
				singleValuedField1,
				singleValuedField2,
				multiValuedField,
				multiValuedHierarchyFieldLinkingToOneSingleValuedField,
			},
			want: []fields.BaseFieldDef{
				multiValuedHierarchyFieldLinkingToOneSingleValuedField_generated,
			},
		},
		{
			name: "Single-valued hierarchy field with a single, " +
				"multi-valued target field generates a single, multi-valued linked field",
			fieldDefs: []fields.BaseFieldDef{
				singleValuedField1,
				singleValuedField2,
				multiValuedField,
				singleValuedHierarchyFieldLinkingToOneMultiValuedField,
			},
			want: []fields.BaseFieldDef{
				singleValuedHierarchyFieldLinkingToOneMultiValuedField_generated,
			},
		},
		{
			name: "Multi-valued hierarchy field with several, " +
				"single-valued target fields generates a multi-valued linked field for each target",
			fieldDefs: []fields.BaseFieldDef{
				singleValuedField1,
				singleValuedField2,
				multiValuedField,
				multiValuedHierarchyFieldLinkingToSeveralSingleValuedFields,
			},
			want: []fields.BaseFieldDef{
				multiValuedHierarchyFieldLinkingToSeveralSingleValuedFields_generated1,
				multiValuedHierarchyFieldLinkingToSeveralSingleValuedFields_generated2,
			},
		},
		{
			name: "Multi-valued hierarchy field with a single, " +
				"multi-valued target field generates a single, multi-valued linked field",
			fieldDefs: []fields.BaseFieldDef{
				singleValuedField1,
				singleValuedField2,
				multiValuedField,
				multiValuedHierarchyFieldLinkingToOneMultiValuedField,
			},
			want: []fields.BaseFieldDef{
				multiValuedHierarchyFieldLinkingToOneMultiValuedField_generated,
			},
		},
		{
			name: "If a hierarchy field specifies a target field that is not configured, it is ignored",
			fieldDefs: []fields.BaseFieldDef{
				singleValuedField1,
				singleValuedField2,
				multiValuedField,
				invalidHierarchyField,
			},
			want: []fields.BaseFieldDef{
				fieldWithOneInvalidTargetHierarchyFieldLinkingToOneSingleValuedField_generated,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := GetLinkedFieldDefs(tt.fieldDefs)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("getSearchFocusFields(), wanted error was %v, but the error returned was %v", tt.wantErr, gotErr)
			} else {
				if len(got) != len(tt.want) {
					t.Errorf("Expected %d generated linked fields but found %d", len(tt.want), len(got))
				}
				sort.Slice(got, func(i, j int) bool {
					return got[i].Name() > got[j].Name()
				})
				sort.Slice(tt.want, func(i, j int) bool {
					return tt.want[i].Name() > tt.want[j].Name()
				})
				for i := range got {
					if got[i].Name() != tt.want[i].Name() {
						t.Errorf("Names do not match - got %s but expected %s", got[i].Name(), tt.want[i].Name())
					}
					if got[i].Kind() != tt.want[i].Kind() {
						t.Errorf("Kinds do not match - got %s but expected %s", got[i].Kind(), tt.want[i].Kind())
					}
				}
			}
		})
	}
}
