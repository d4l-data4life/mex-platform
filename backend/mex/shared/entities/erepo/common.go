package erepo

import (
	"fmt"
	"regexp"
)

var pattern = regexp.MustCompile("^[a-zA-Z][a-zA-z_0-9]*[a-zA-z0-9]$")

func ValidateName(fieldName string) error {
	if !pattern.MatchString(fieldName) {
		return fmt.Errorf("entity type names can only contain alphanumeric characters, underscores, must start with a letter and end alphanumerically")
	}

	return nil
}

func QuoteAndSanitize(sl []string) ([]string, error) {
	x := make([]string, len(sl))
	for i, s := range sl {
		if err := ValidateName(s); err != nil {
			return nil, fmt.Errorf("malformed: '%s', (%w)", s, err)
		}

		x[i] = fmt.Sprintf("'%s'", s)
	}
	return x, nil
}
