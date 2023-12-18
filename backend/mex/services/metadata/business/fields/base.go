package fields

import (
	"fmt"

	"github.com/d4l-data4life/mex/mex/shared/fields"
	"github.com/d4l-data4life/mex/mex/shared/solr"
)

type baseFieldDef struct {
	name          string
	kind          string
	displayID     string
	isLinkedField bool

	multiValued bool
}

func (def *baseFieldDef) Name() string        { return def.name }
func (def *baseFieldDef) Kind() string        { return def.kind }
func (def *baseFieldDef) DisplayID() string   { return def.displayID }
func (def *baseFieldDef) IsLinkedField() bool { return def.isLinkedField }
func (def *baseFieldDef) MultiValued() bool   { return def.multiValued }

type BaseIndexDef struct {
	MultiValued bool
}

func NewBaseFieldDef(fieldName string, kind string, displayID string, isLinkedField bool, baseIndexDef BaseIndexDef) BaseFieldDef {
	return &baseFieldDef{
		name:          fieldName,
		kind:          kind,
		displayID:     displayID,
		isLinkedField: isLinkedField,
		multiValued:   baseIndexDef.MultiValued,
	}
}

// GetFirstLinkExt returns the first linked field extension in a given field definition (error if none found)
func GetFirstLinkExt(indexDef *fields.IndexDef) (*fields.IndexDefExtLink, error) {
	for _, ext := range indexDef.Ext {
		if ext.MessageName() == solr.LinkExtID {
			var linkExt fields.IndexDefExtLink
			err := ext.UnmarshalTo(&linkExt)
			if err != nil {
				return nil, err
			}
			return &linkExt, nil
		}
	}
	return nil, fmt.Errorf("no IndexDefExtLink in extension")
}

// GetFirstHierarchyExt returns the first hierarchy extension in a given field definition (error if none found)
func GetFirstHierarchyExt(indexDef *fields.IndexDef) (*fields.IndexDefExtHierarchy, error) {
	for _, ext := range indexDef.Ext {
		if ext.MessageName() == solr.HierarchyExtID {
			var hierarchyExt fields.IndexDefExtHierarchy
			err := ext.UnmarshalTo(&hierarchyExt)
			if err != nil {
				return nil, err
			}
			return &hierarchyExt, nil
		}
	}
	return nil, fmt.Errorf("no IndexDefExtHierarchy in extension")
}
