package fields

import (
	"fmt"
	"regexp"
	"strings"
)

var pattern = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]+$")

func ValidateInternalName(fieldName string) error {
	if !pattern.MatchString(fieldName) {
		return fmt.Errorf("invalid field (definition) name")
	}

	return nil
}

func QuoteAndSanitize(sl []string) ([]string, error) {
	x := make([]string, len(sl))
	for i, s := range sl {
		if err := ValidateInternalName(s); err != nil {
			return nil, fmt.Errorf("malformed: '%s', (%w)", s, err)
		}

		x[i] = fmt.Sprintf("'%s'", s)
	}
	return x, nil
}

func ValidateName(fieldName string) error {
	if !pattern.MatchString(fieldName) {
		return fmt.Errorf("field names can only contain alphanumeric characters and underscores and cannot start with a digit")
	}

	if strings.Contains(fieldName, "__") {
		return fmt.Errorf("field names are not allowed to contain double underscores ('__')")
	}

	return nil
}
