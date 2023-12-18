//revive:disable:var-naming
package kind_hierarchy

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/types/known/anypb"

	fieldUtils "github.com/d4l-data4life/mex/mex/shared/fields"
	"github.com/d4l-data4life/mex/mex/shared/solr"
	"github.com/d4l-data4life/mex/mex/shared/utils"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/datamodel"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/hierarchy/codings"
	kindLink "github.com/d4l-data4life/mex/mex/services/metadata/business/fields/kinds/link"
)

const KindName = "hierarchy"

type kindHierarchy struct {
	codings codings.Codings
}

type HierarchyFieldDef interface {
	kindLink.LinkFieldDef

	CodeSystemNameOrEntityType() string
	LinkFieldName() string
	DisplayFieldName() string
}

type hierarchyFieldDef struct {
	/*
		One could also have built this type by embedding the corresponding type of the link field type. However,
		unlike the Link interface, the link type is not exported (same in all other types), so we would have to
		expose it to include it here. For that reasons, the interface is simply re-implemented here.
	*/
	fields.BaseFieldDef

	codeSystemNameOrEntityType string
	linkFieldName              string
	displayFieldName           string

	relationType       string
	linkedTargetFields []string
}

func (def *hierarchyFieldDef) CodeSystemNameOrEntityType() string {
	return def.codeSystemNameOrEntityType
}

func (def *hierarchyFieldDef) LinkFieldName() string { return def.linkFieldName }

func (def *hierarchyFieldDef) DisplayFieldName() string { return def.displayFieldName }

func (def *hierarchyFieldDef) RelationType() string { return def.relationType }

func (def *hierarchyFieldDef) LinkedTargetFields() []string { return def.linkedTargetFields }

func (kind *kindHierarchy) ValidateDefinition(_ context.Context, request *fieldUtils.FieldDef) (fields.BaseFieldDef, error) {
	err := fieldUtils.ValidateName(request.Name)
	if err != nil {
		return nil, err
	}

	if request.Kind != KindName {
		return nil, fmt.Errorf("kind is not %s: %s", KindName, request.Kind)
	}

	hierarchyExt, err := fields.GetFirstHierarchyExt(request.IndexDef)
	if err != nil {
		return nil, fmt.Errorf("no hierarchy config extension found for hierarchy field")
	}

	extLink, err := fields.GetFirstLinkExt(request.IndexDef)
	if err != nil {
		return nil, fmt.Errorf("no link config extension found for hierarchy field")
	}

	completeFieldDef := hierarchyFieldDef{
		BaseFieldDef: fields.NewBaseFieldDef(request.Name, request.Kind, request.DisplayId, false, fields.BaseIndexDef{
			MultiValued: request.IndexDef.MultiValued,
		}),

		// Hierarchy properties
		codeSystemNameOrEntityType: hierarchyExt.CodeSystemNameOrNodeEntityType,
		linkFieldName:              hierarchyExt.LinkFieldName,
		displayFieldName:           hierarchyExt.DisplayFieldName,

		// Link properties
		relationType:       extLink.RelationType,
		linkedTargetFields: extLink.LinkedTargetFields,
	}

	return &completeFieldDef, nil
}

func (kind *kindHierarchy) MustValidateDefinition(ctx context.Context, request *fieldUtils.FieldDef) fields.BaseFieldDef {
	fieldDef, err := kind.ValidateDefinition(ctx, request)
	if err != nil {
		panic(err)
	}
	return fieldDef
}

func (kind *kindHierarchy) MarshalToProtobufFormat(_ context.Context, fieldDef fields.BaseFieldDef) (*fieldUtils.FieldDef, error) {
	hFieldDef, ok := (fieldDef).(HierarchyFieldDef)
	if !ok {
		return nil, fmt.Errorf("field definition object is not a HierarchyFieldDef")
	}

	// Set hierarchy properties
	hierarchyExt, err := anypb.New(&fieldUtils.IndexDefExtHierarchy{
		CodeSystemNameOrNodeEntityType: hFieldDef.CodeSystemNameOrEntityType(),
		LinkFieldName:                  hFieldDef.LinkFieldName(),
		DisplayFieldName:               hFieldDef.DisplayFieldName(),
	})
	if err != nil {
		return nil, err
	}

	// Set link properties
	linkExt, err := anypb.New(&fieldUtils.IndexDefExtLink{
		RelationType:       hFieldDef.RelationType(),
		LinkedTargetFields: []string{},
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
			Ext: []*anypb.Any{
				hierarchyExt,
				linkExt,
			},
		},
	}, nil
}

func (kind *kindHierarchy) ValidateFieldValue(_ context.Context, _ fields.BaseFieldDef, _ string) error {
	return nil
}

// revive:disable:unexported-return
func NewKindHierarchy(db *pgxpool.Pool) (*kindHierarchy, error) {
	postgresCodings := codings.NewPostgresCodings(db)

	return &kindHierarchy{
		codings: &postgresCodings,
	}, nil
}

type KindHierarchy struct{}

func (kind *KindHierarchy) GenerateSolrFields(_ context.Context, fieldDef fields.BaseFieldDef) (solr.FieldCategoryToSolrFieldDefsMap, error) {
	if fieldDef == nil {
		return nil, fmt.Errorf("cannot generate backing fields from empty field definition")
	}
	baseTransitiveHullFieldName := solr.GetTransitiveHullFieldName(fieldDef.Name())
	baseTransitiveHullDisplayFieldName := solr.GetTransitiveHullDisplayFieldName(fieldDef.Name())
	baseDisplayFieldNameDe, err := solr.GetLangSpecificFieldName(baseTransitiveHullDisplayFieldName, solr.GermanLangAbbrev)
	if err != nil {
		return nil, fmt.Errorf("could not generate name of German language field for code labels for the"+
			" hierarchy code field %s", fieldDef.Name())
	}
	baseDisplayFieldNameEn, err := solr.GetLangSpecificFieldName(baseTransitiveHullDisplayFieldName, solr.EnglishLangAbbrev)
	if err != nil {
		return nil, fmt.Errorf("could not generate name of English language field for code labels for the"+
			" hierarchy code field %s", fieldDef.Name())
	}

	/*
		NOTE: For the hierarchy kind, the support for German and English is hard-coded for the time being.
	*/
	solrFields := make(solr.FieldCategoryToSolrFieldDefsMap)
	solrFields[solr.GenericLangBaseFieldCategory] = solr.GetStandardPrimaryBackingField(fieldDef.Name(), solr.DefaultSolrStringFieldType, fieldDef.MultiValued())
	solrFields[solr.ParentCodesBaseFieldCategory] = solr.GetStandardPrimaryBackingField(baseTransitiveHullFieldName, solr.DefaultSolrStringFieldType, true)
	solrFields[solr.GermanLangBaseFieldCategory] = solr.GetStandardPrimaryBackingField(baseDisplayFieldNameDe, solr.DefaultDeSolrTextFieldType, true)
	solrFields[solr.EnglishLangBaseFieldCategory] = solr.GetStandardPrimaryBackingField(baseDisplayFieldNameEn, solr.DefaultEnSolrTextFieldType, true)

	return solrFields, nil
}

func (kind *kindHierarchy) GenerateXMLFieldTags(ctx context.Context, fieldDef fields.BaseFieldDef, itemValue datamodel.CurrentItemValue) ([]string, error) {
	var hFieldDef HierarchyFieldDef
	var ok bool
	if hFieldDef, ok = (fieldDef).(HierarchyFieldDef); !ok {
		return nil, fmt.Errorf("field definition is not a HierarchyFieldDef, but a %t", fieldDef)
	}

	var tags []string
	// Write the code itself into the field with the same name as the MEx field it backs
	cleanedValue := utils.SanitizeXML(itemValue.FieldValue)
	tags = append(tags, fmt.Sprintf("<field name=\"%s\">%s</field>", itemValue.FieldName, cleanedValue))

	// Write the fields for the transitive hull
	baseTransitiveHullFieldName := solr.GetTransitiveHullFieldName(itemValue.FieldName)
	baseTransitiveHullDisplayFieldName := solr.GetTransitiveHullDisplayFieldName(itemValue.FieldName)
	hulls, err := kind.codings.TransitiveClosure(ctx, hFieldDef.CodeSystemNameOrEntityType(), hFieldDef.LinkFieldName(), hFieldDef.DisplayFieldName(), itemValue.FieldValue)
	if err != nil {
		return nil, err
	}
	codeSeen := make(map[string]bool)
	for language, hull := range hulls {
		languageSpecificTransitiveHullDisplayFieldName, err := solr.GetLangSpecificFieldName(baseTransitiveHullDisplayFieldName, language)
		if err != nil {
			return nil, fmt.Errorf("the language code '%s' used for the hierarchy field '%s' is not supported", language,
				fieldDef.Name())
		}
		for _, coding := range hull {
			// Only add code itself the first time we see it, not for every language
			if _, ok := codeSeen[coding.Code]; !ok {
				tags = append(tags, fmt.Sprintf("<field name=\"%s\">%s</field>", baseTransitiveHullFieldName, coding.Code))
				codeSeen[coding.Code] = true
			}
			cleanedDisplayValue := utils.SanitizeXML(coding.Display)
			tags = append(tags, fmt.Sprintf("<field name=\"%s\">%s</field>", languageSpecificTransitiveHullDisplayFieldName, cleanedDisplayValue))
		}
	}

	return tags, nil
}

func (kind *kindHierarchy) EnrichFacetBucket(_ context.Context, bucket *solr.FacetBucket, fieldDef fields.BaseFieldDef) (*solr.FacetBucket, error) {
	// TODO: Move this logic to the handling of ordinal axes
	hFieldDef, ok := (fieldDef).(HierarchyFieldDef)
	if !ok {
		return nil, fmt.Errorf("field definition is not a HierarchyFieldDef, but a %t", fieldDef)
	}
	info, err := makeHierarchyInfo(
		kind.codings,
		hFieldDef.CodeSystemNameOrEntityType(),
		hFieldDef.LinkFieldName(),
		hFieldDef.DisplayFieldName(),
		bucket.Value)
	if err != nil {
		return nil, fmt.Errorf("error getting hierarchy info: %s", err.Error())
	}
	anyInfo, err := anypb.New(info)
	if err != nil {
		return nil, fmt.Errorf("could not set hierarchy info: %s", err.Error())
	}

	return &solr.FacetBucket{
		Value:         bucket.Value,
		Count:         bucket.Count,
		HierarchyInfo: anyInfo,
	}, nil
}

func makeHierarchyInfo(codings codings.Codings, codeSystemName string, linkType string, displayFieldName string, fieldValue string) (*solr.HierarchyInfo, error) {
	if codings == nil {
		return nil, fmt.Errorf("codings are nil")
	}

	codes, err := codings.Resolve(context.Background(), codeSystemName, linkType, displayFieldName, fieldValue)
	if err != nil {
		return nil, err
	}

	if len(codes) != 1 {
		return nil, fmt.Errorf("(%s) number of codes: %d", fieldValue, len(codes))
	}

	return &solr.HierarchyInfo{
		ParentValue: codes[0].ParentCode,
		Display:     codes[0].Display,
		Depth:       uint32(codes[0].Depth),
	}, nil
}

func (kind *kindHierarchy) ResetCaches() {
	kind.codings.Reset()
}
