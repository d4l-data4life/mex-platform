package searchconfig

import (
	"github.com/d4l-data4life/mex/mex/shared/solr"
)

type FunctionBackingFieldInfo struct {
	Def                solr.FieldDef
	FunctionCategoryID string // Functional category of field
}
