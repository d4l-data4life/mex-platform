package frepo

import (
	"context"
	"fmt"

	"github.com/d4l-data4life/mex/mex/shared/utils"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/linked"
)

type mockedFieldRepo struct {
	fieldDefs []fields.BaseFieldDef
}

func NewMockedFieldRepo(fieldDefs []fields.BaseFieldDef) fields.FieldRepo {
	fullFieldList := fieldDefs
	// Add all linked fields to the stored list as well
	linkedFieldDefs, _ := linked.GetLinkedFieldDefs(fieldDefs)
	fullFieldList = append(fullFieldList, linkedFieldDefs...)
	// Add all pre-defined to the stored list as well
	preDefinedFields := getPredefinedFields()
	fullFieldList = append(fullFieldList, preDefinedFields...)
	return &mockedFieldRepo{fieldDefs: fullFieldList}
}

func (repo *mockedFieldRepo) GetFieldDefNames(_ context.Context) ([]string, error) {
	fieldNames := []string{}
	for _, fieldDef := range repo.fieldDefs {
		fieldNames = append(fieldNames, fieldDef.Name())
	}

	fieldNames = utils.Unique(fieldNames)
	if len(fieldNames) != len(repo.fieldDefs) {
		return nil, fmt.Errorf("mock: duplicate field names")
	}

	return fieldNames, nil
}

func (repo *mockedFieldRepo) GetFieldDefByName(_ context.Context, fieldName string) (fields.BaseFieldDef, error) {
	for _, fieldDef := range repo.fieldDefs {
		if fieldDef.Name() == fieldName {
			return fieldDef, nil
		}
	}
	return nil, fmt.Errorf("mock: field def not found: %s", fieldName)
}

func (repo *mockedFieldRepo) GetFieldDefsByKind(_ context.Context, fieldKind string) ([]fields.BaseFieldDef, error) {
	fieldDefs := []fields.BaseFieldDef{}
	for _, fieldDef := range repo.fieldDefs {
		if fieldDef.Kind() == fieldKind {
			fieldDefs = append(fieldDefs, fieldDef)
		}
	}

	return fieldDefs, nil
}

func (repo *mockedFieldRepo) ListFieldDefs(_ context.Context) ([]fields.BaseFieldDef, error) {
	return repo.fieldDefs, nil
}

func (*mockedFieldRepo) Purge(_ context.Context) error { return nil }
