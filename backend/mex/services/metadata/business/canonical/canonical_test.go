package canonical

import (
	"testing"

	"github.com/d4l-data4life/mex/mex/shared/items"
)

type testCase struct {
	values   []*items.ItemValue
	expected string
}

var testCases = []testCase{
	{
		values:   []*items.ItemValue{},
		expected: "()",
	},
	{
		values:   []*items.ItemValue{{FieldName: "name", FieldValue: "Guybrush   Threepwood"}},
		expected: "((name:(('Guybrush   Threepwood',''))))",
	},
	{
		values: []*items.ItemValue{
			{
				FieldName:  "flag",
				FieldValue: "🏴󠁧󠁢󠁳󠁣󠁴󠁿",
				Language:   "sco",
			},
			{
				FieldName:  "flag",
				FieldValue: "🇪🇺",
				Language:   "eu",
			},
		},
		expected: "((flag:(('🏴󠁧󠁢󠁳󠁣󠁴󠁿','sco'),('🇪🇺','eu'))))",
	},
	{
		values: []*items.ItemValue{
			{
				FieldName:  "keyword",
				FieldValue: "foo",
			},
			{
				FieldName:  "title",
				FieldValue: "This is an integration test item",
				Language:   "en",
			},
			{
				FieldName:  "author",
				FieldValue: "Mr A",
			},
			{
				FieldName:  "author",
				FieldValue: "Mr B",
			},
			{
				FieldName:  "abstract",
				FieldValue: "Lorem ipsum",
			},
			{
				FieldName:  "keyword",
				FieldValue: "bar",
			},
			{
				FieldName:  "author",
				FieldValue: "Mr C",
			},
		},
		expected: "((abstract:(('Lorem ipsum',''))),(author:(('Mr A',''),('Mr B',''),('Mr C',''))),(keyword:(('foo',''),('bar',''))),(title:(('This is an integration test item','en'))))",
	},
}

func TestCanonicalize(t *testing.T) {
	for _, c := range testCases {
		actual := canonicalizedValues(c.values)
		if actual != c.expected {
			t.Errorf("wanted: %s, got: %s", c.expected, actual)
		}
	}
}
