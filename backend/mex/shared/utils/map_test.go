package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	input    map[string]string
	expected []string
}

func TestKeysOfMap(t *testing.T) {
	tests := []testCase{
		{
			input:    map[string]string{},
			expected: []string{},
		},
		{
			input:    map[string]string{"foo": "bar"},
			expected: []string{"foo"},
		},
		{
			input:    map[string]string{"xxx": "yyy", "foo": "bar", "1234": ""},
			expected: []string{"1234", "foo", "xxx"},
		},
		{
			input:    map[string]string{"xxx": "yyy", "foo": "bar", "1234": "", "": ""},
			expected: []string{"", "1234", "foo", "xxx"},
		},
	}

	for _, test := range tests {
		require.Equal(t, test.expected, sortStrings(KeysOfMap(test.input)))
	}
}
