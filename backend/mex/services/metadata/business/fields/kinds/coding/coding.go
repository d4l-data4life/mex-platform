package kind_coding // revive:disable

import (
	"context"
	"fmt"
	"regexp"

	"google.golang.org/protobuf/types/known/anypb"

	"github.com/d4l-data4life/mex/mex/shared/codings/csrepo"
	fieldUtils "github.com/d4l-data4life/mex/mex/shared/fields"
	"github.com/d4l-data4life/mex/mex/shared/solr"
	"github.com/d4l-data4life/mex/mex/shared/utils"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/datamodel"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
)

const KindName = "coding"

type KindCoding struct {
	CodingsetRepo csrepo.CodingsetRepo
}

type CodingFieldDef interface {
	fields.BaseFieldDef

	CodingsetNames() []string
}

type codingFieldDef struct {
	fields.BaseFieldDef

	codingsetNames []string
}

func (def *codingFieldDef) CodingsetNames() []string {
	return def.codingsetNames
}

func (kind *KindCoding) ValidateDefinition(_ context.Context, request *fieldUtils.FieldDef) (fields.BaseFieldDef, error) {
	err := fieldUtils.ValidateName(request.Name)
	if err != nil {
		return nil, err
	}

	if request.Kind != KindName {
		return nil, fmt.Errorf("kind is not %s: %s", KindName, request.Kind)
	}

	if len(request.IndexDef.Ext) == 0 {
		return nil, fmt.Errorf("index definition extension is empty")
	}

	var codingExt fieldUtils.IndexDefExtCoding
	err = request.IndexDef.Ext[0].UnmarshalTo(&codingExt)
	if err != nil {
		return nil, fmt.Errorf("malformed index definition extension: %s", err.Error())
	}

	return &codingFieldDef{
		BaseFieldDef: fields.NewBaseFieldDef(request.Name, request.Kind, request.DisplayId, false, fields.BaseIndexDef{
			MultiValued: request.IndexDef.MultiValued,
		}),
		codingsetNames: codingExt.CodingsetNames,
	}, nil
}

func (kind *KindCoding) MustValidateDefinition(ctx context.Context, request *fieldUtils.FieldDef) fields.BaseFieldDef {
	fieldDef, err := kind.ValidateDefinition(ctx, request)
	if err != nil {
		panic(err)
	}
	return fieldDef
}

func (kind *KindCoding) MarshalToProtobufFormat(_ context.Context, fieldDef fields.BaseFieldDef) (*fieldUtils.FieldDef, error) {
	if hFieldDef, ok := (fieldDef).(CodingFieldDef); ok {
		ext, err := anypb.New(&fieldUtils.IndexDefExtCoding{
			CodingsetNames: hFieldDef.CodingsetNames(),
		})
		if err != nil {
			return nil, err
		}

		return &fieldUtils.FieldDef{
			Name:      fieldDef.Name(),
			Kind:      fieldDef.Kind(),
			DisplayId: fieldDef.DisplayID(),
			IndexDef: &fieldUtils.IndexDef{
				MultiValued: fieldDef.MultiValued(),
				Ext:         []*anypb.Any{ext},
			},
		}, nil
	}

	return nil, fmt.Errorf("field definition object is not a CodingFieldDef")
}

func (kind *KindCoding) ValidateFieldValue(context.Context, fields.BaseFieldDef, string) error {
	return nil
}

func (kind *KindCoding) GenerateSolrFields(_ context.Context, fieldDef fields.BaseFieldDef) (solr.FieldCategoryToSolrFieldDefsMap, error) {
	if fieldDef == nil {
		return nil, fmt.Errorf("cannot generate backing fields from empty field definition")
	}
	solrFields := solr.FieldCategoryToSolrFieldDefsMap{
		solr.GenericLangBaseFieldCategory: solr.GetStandardPrimaryBackingField(fieldDef.Name(), solr.DefaultSolrStringFieldType, fieldDef.MultiValued()),
		solr.GermanLangBaseFieldCategory:  solr.GetStandardPrimaryBackingField(solr.GetDisplayFieldName(fieldDef.Name(), solr.GermanLangAbbrev), solr.DefaultDeSolrTextFieldType, true),
		solr.EnglishLangBaseFieldCategory: solr.GetStandardPrimaryBackingField(solr.GetDisplayFieldName(fieldDef.Name(), solr.EnglishLangAbbrev), solr.DefaultEnSolrTextFieldType, true),
	}
	return solrFields, nil
}

var headingCode = regexp.MustCompile(`.*(D\d+)$`)

// Return the trailing 'Ddddddddd' part of a string, if exists.
func extractHeadingCode(s string) string {
	sm := headingCode.FindStringSubmatch(s)
	if len(sm) == 2 {
		return sm[1]
	}
	return ""
}

func (kind *KindCoding) GenerateXMLFieldTags(_ context.Context, fieldDef fields.BaseFieldDef, itemValue datamodel.CurrentItemValue) ([]string, error) {
	ret := []string{fmt.Sprintf("<field name=\"%s\">%s</field>", itemValue.FieldName, utils.SanitizeXML(itemValue.FieldValue))}

	var hFieldDef CodingFieldDef
	var ok bool
	if hFieldDef, ok = (fieldDef).(CodingFieldDef); !ok {
		return nil, fmt.Errorf("field definition is not a CodingFieldDef, but a %t", fieldDef)
	}

	for _, name := range hFieldDef.CodingsetNames() {
		codingset, err := kind.CodingsetRepo.GetCodingset(name)
		if err != nil {
			return nil, err
		}

		dnumber := extractHeadingCode(itemValue.FieldValue)
		displays, err := codingset.ResolveMainHeadings([]string{dnumber}, "de")
		if err != nil {
			return nil, err
		}

		for _, d := range displays {
			ret = append(ret, fmt.Sprintf("<field name=\"%s\">%s</field>", solr.GetDisplayFieldName(itemValue.FieldName, "de"), utils.SanitizeXML(d)))
		}

		displays, err = codingset.ResolveMainHeadings([]string{dnumber}, "en")
		if err != nil {
			return nil, err
		}

		for _, d := range displays {
			ret = append(ret, fmt.Sprintf("<field name=\"%s\">%s</field>", solr.GetDisplayFieldName(itemValue.FieldName, "en"), utils.SanitizeXML(d)))
		}
	}

	return ret, nil
}

func (*KindCoding) GetSortAndFacetFieldName(_ context.Context, fieldDef fields.BaseFieldDef) string {
	return fieldDef.Name()
}

func (kind *KindCoding) EnrichFacetBucket(_ context.Context, bucket *solr.FacetBucket, _ fields.BaseFieldDef) (*solr.FacetBucket, error) {
	return bucket, nil
}

func (kind *KindCoding) ResetCaches() {
}
