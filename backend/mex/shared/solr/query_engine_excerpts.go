package solr

import (
	"fmt"
)

// createStatExpression return the Solr expression when applying the given operator to the given axis
func CreateStatExpression(axisName string, op string) (string, error) {
	if axisName == "" {
		return "", fmt.Errorf("cannot create stat operator expression with empty axis name")
	}
	var expr string
	switch op {
	case MinOperator:
		expr = fmt.Sprintf("min(%s)", axisName)
	case MaxOperator:
		expr = fmt.Sprintf("max(%s)", axisName)
	default:
		return "", fmt.Errorf("the used stat operator is not allowed by MEx")
	}
	return expr, nil
}
