package frepo

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	kind_coding "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/coding"
	kind_link "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/link"
	kind_number "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/number"
	kind_string "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/string"
	kind_text "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/text"
	kind_timestamp "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/timestamp"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/hooks"
	sharedFields "github.com/d4l-data4life/mex/mex/shared/fields"
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

func getTestFieldDefinitionHooks() hooks.FieldDefinitionHooks {
	hooksHere := make(hooks.FieldDefinitionHooks)

	hooksHere[kind_number.KindName] = &kind_number.KindNumber{}
	hooksHere[kind_string.KindName] = &kind_string.KindString{}
	hooksHere[kind_text.KindName] = &kind_text.KindText{}
	hooksHere[kind_timestamp.KindName] = &kind_timestamp.KindTimestamp{}
	hooksHere[kind_link.KindName] = &kind_link.KindLink{}

	// No hierarchy kind to avoid having to supply a DB
	hooksHere[kind_coding.KindName] = &kind_coding.KindCoding{}

	return hooksHere
}

var cmsFieldContent = sharedFields.FieldDefList{
	FieldDefs: []*sharedFields.FieldDef{
		{
			Name: "label",
			Kind: "text",
			IndexDef: &sharedFields.IndexDef{
				MultiValued: true,
			},
		},
		{
			Name:     "count",
			Kind:     "number",
			IndexDef: &sharedFields.IndexDef{},
		},
		{
			Name:     "category",
			Kind:     "string",
			IndexDef: &sharedFields.IndexDef{},
		},
		{
			Name: "contact",
			Kind: "link",
			IndexDef: &sharedFields.IndexDef{
				MultiValued: true,
				Ext: toAnySlice(&sharedFields.IndexDefExtLink{
					RelationType: "myRelation",
					LinkedTargetFields: []string{
						"label",
						"category",
					},
				}),
			},
		},
		{
			Name: "ranking",
			Kind: "link",
			IndexDef: &sharedFields.IndexDef{
				Ext: toAnySlice(&sharedFields.IndexDefExtLink{
					RelationType: "myCountRelation",
					LinkedTargetFields: []string{
						"count",
					},
				}),
			},
		},
	},
}

var labelField = fields.NewBaseFieldDef("label", "text", "", false, fields.BaseIndexDef{MultiValued: true})
var categoryField = fields.NewBaseFieldDef("category", "string", "", false, fields.BaseIndexDef{})
var countField = fields.NewBaseFieldDef("count", "number", "", false, fields.BaseIndexDef{})
var contactField = fields.NewBaseFieldDef("contact", "link", "", false, fields.BaseIndexDef{MultiValued: true})
var contactLabelField = fields.NewBaseFieldDef("contact__label", "text", "FIELD_CONTACT__LABEL", true, fields.BaseIndexDef{MultiValued: true})
var contactCategoryField = fields.NewBaseFieldDef("contact__category", "string", "FIELD_CONTACT__CATEGORY", true, fields.BaseIndexDef{MultiValued: true})
var rankingField = fields.NewBaseFieldDef("ranking", "link", "", false, fields.BaseIndexDef{})
var rankingCountField = fields.NewBaseFieldDef("ranking__count", "number", "FIELD_RANKING__COUNT", true, fields.BaseIndexDef{})

func Test_fieldDefsRepoDirectCMS_GetFieldDefByName(t *testing.T) {
	tests := []struct {
		name               string
		statusCodeReturned int
		fieldName          string
		want               fields.BaseFieldDef
		wantErr            bool
	}{
		{
			name:               "Returns a normal field if available",
			statusCodeReturned: http.StatusOK,
			fieldName:          "label",
			want:               labelField,
		},
		{
			name:               "Returns a linked field if available",
			statusCodeReturned: http.StatusOK,
			fieldName:          "contact__label",
			want:               contactLabelField,
		},
		{
			name:               "Throws an error for unknown normal field",
			statusCodeReturned: http.StatusOK,
			fieldName:          "unknown",
			wantErr:            true,
		},
		{
			name:               "Throws an error for unknown linked field",
			statusCodeReturned: http.StatusOK,
			fieldName:          "contact__count",
			wantErr:            true,
		},
		{
			name:               "Throws an error if CMS returns error",
			statusCodeReturned: http.StatusInternalServerError,
			fieldName:          "label",
			wantErr:            true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange - set up mock server to return the required value
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCodeReturned)
				body, _ := protojson.Marshal(&cmsFieldContent)
				fmt.Fprintln(w, string(body))
			}))
			defer ts.Close()
			// Act
			repo := &fieldDefsRepoDirectCMS{
				originCMS: ts.URL,
				parsers:   getTestFieldDefinitionHooks(),
			}
			got, err := repo.GetFieldDefByName(context.TODO(), tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFieldDefByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFieldDefByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fieldDefsRepoDirectCMS_GetFieldDefNames(t *testing.T) {
	tests := []struct {
		name               string
		statusCodeReturned int
		want               []string
		wantErr            bool
	}{
		{
			name:               "Returns all field names",
			statusCodeReturned: http.StatusOK,
			want: []string{
				"label",
				"contact",
				"category",
				"count",
				"ranking",
				"ranking__count",
				"contact__label",
				"contact__category",
			},
		},
		{
			name:               "Returns an error if the CMS returns an error",
			statusCodeReturned: http.StatusInternalServerError,
			wantErr:            true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange - set up mock server to return the required value
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCodeReturned)
				body, _ := protojson.Marshal(&cmsFieldContent)
				fmt.Fprintln(w, string(body))
			}))
			defer ts.Close()
			// Act
			repo := &fieldDefsRepoDirectCMS{
				originCMS: ts.URL,
				parsers:   getTestFieldDefinitionHooks(),
			}
			got, err := repo.GetFieldDefNames(context.TODO())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFieldDefNames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			sort.Strings(tt.want)
			sort.Strings(got)
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFieldDefByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fieldDefsRepoDirectCMS_GetFieldDefsByKind(t *testing.T) {
	tests := []struct {
		name               string
		statusCodeReturned int
		kindName           string
		want               []fields.BaseFieldDef
		wantErr            bool
	}{
		{
			name:               "Returns both normal and linked fields of the given kind",
			statusCodeReturned: http.StatusOK,
			kindName:           "text",
			want: []fields.BaseFieldDef{
				contactLabelField,
				labelField,
			},
		},
		{
			name:               "Returns a nil slice if no field of the requested kind are available",
			statusCodeReturned: http.StatusOK,
			kindName:           "timestamp",
			want:               nil,
		},
		{
			name:               "Throws an error if CMS returns error",
			statusCodeReturned: http.StatusInternalServerError,
			wantErr:            true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange - set up mock server to return the required value
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCodeReturned)
				body, _ := protojson.Marshal(&cmsFieldContent)
				fmt.Fprintln(w, string(body))
			}))
			defer ts.Close()
			// Act
			repo := &fieldDefsRepoDirectCMS{
				originCMS: ts.URL,
				parsers:   getTestFieldDefinitionHooks(),
			}
			got, err := repo.GetFieldDefsByKind(context.TODO(), tt.kindName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFieldDefsByKind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			sort.Slice(got, func(i, j int) bool {
				return got[i].Name() < got[j].Name()
			})
			sort.Slice(tt.want, func(i, j int) bool {
				return tt.want[i].Name() < tt.want[j].Name()
			})
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFieldDefByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fieldDefsRepoDirectCMS_ListFieldDefs(t *testing.T) {
	tests := []struct {
		name               string
		statusCodeReturned int
		kindName           string
		want               []fields.BaseFieldDef
		wantErr            bool
	}{
		{
			name:               "Returns all normal and linked fields",
			statusCodeReturned: http.StatusOK,
			kindName:           "text",
			want: []fields.BaseFieldDef{
				labelField,
				categoryField,
				contactField,
				countField,
				rankingField,
				contactLabelField,
				contactCategoryField,
				rankingCountField,
			},
		},
		{
			name:               "Throws an error if CMS returns error",
			statusCodeReturned: http.StatusInternalServerError,
			wantErr:            true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange - set up mock server to return the required value
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCodeReturned)
				body, _ := protojson.Marshal(&cmsFieldContent)
				fmt.Fprintln(w, string(body))
			}))
			defer ts.Close()
			// Act
			repo := &fieldDefsRepoDirectCMS{
				originCMS: ts.URL,
				parsers:   getTestFieldDefinitionHooks(),
			}
			got, err := repo.ListFieldDefs(context.TODO())
			if (err != nil) != tt.wantErr {
				t.Errorf("ListFieldDefs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			sort.Slice(got, func(i, j int) bool {
				return got[i].Name() < got[j].Name()
			})
			sort.Slice(tt.want, func(i, j int) bool {
				return tt.want[i].Name() < tt.want[j].Name()
			})
			if len(tt.want) != len(got) {
				t.Errorf("ListFieldDefs() got length = %d, want length %v", len(got), len(tt.want))
			}
			for i, gotVal := range got {
				if gotVal.Name() != tt.want[i].Name() {
					t.Errorf("ListFieldDefs() at index %d: got = %v, want %v", i, gotVal, tt.want[i])
				}
				if gotVal.Name() != tt.want[i].Name() || gotVal.Kind() != tt.want[i].Kind() {
					t.Errorf("ListFieldDefs() at index %d: got = %v, want %v", i, gotVal, tt.want[i])
				}
				if gotVal.Name() != tt.want[i].Name() || gotVal.Kind() != tt.want[i].Kind() || gotVal.MultiValued() != tt.want[i].MultiValued() {
					t.Errorf("ListFieldDefs() at index %d: got = %v, want %v", i, gotVal, tt.want[i])
				}
				if gotVal.Name() != tt.want[i].Name() || gotVal.Kind() != tt.want[i].Kind() || gotVal.MultiValued() != tt.want[i].MultiValued() {
					t.Errorf("ListFieldDefs() at index %d: got = %v, want %v", i, gotVal, tt.want[i])
				}
			}
		})
	}
}
