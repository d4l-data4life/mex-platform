package utils

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubsetEq(t *testing.T) {
	require.True(t, SubsetEq([]string{}, []string{}))
	require.True(t, SubsetEq([]string{"a"}, []string{"a"}))
	require.True(t, SubsetEq([]string{"a"}, []string{"a", "a"}))
	require.True(t, SubsetEq([]string{"a"}, []string{"a", "b"}))
	require.True(t, SubsetEq([]string{"a", "a"}, []string{"a", "a"}))
	require.True(t, SubsetEq([]string{"a", "a", "b"}, []string{"a", "a", "b", "b"}))

	require.False(t, SubsetEq([]string{"a"}, []string{}))
	require.False(t, SubsetEq([]string{"a", "a"}, []string{"a"}))
	require.False(t, SubsetEq([]string{"a", "b"}, []string{"a", "c"}))
	require.False(t, SubsetEq([]string{"a", "a", "b", "b"}, []string{"a", "a", "b"}))
}

func TestStringSliceSetDiff(t *testing.T) {
	tests := []struct {
		x    []string
		y    []string
		want []string
	}{
		{
			x:    []string{},
			y:    []string{},
			want: []string{},
		},
		{
			x:    []string{"b", "a"},
			y:    []string{},
			want: []string{"a", "b"},
		},
		{
			x:    []string{"b", "a"},
			y:    []string{"c"},
			want: []string{"a", "b"},
		},
		{
			x:    []string{"a", "b"},
			y:    []string{"c", "b"},
			want: []string{"a"},
		},
		{
			x:    []string{"a", "b", "b", "b"},
			y:    []string{"c", "b"},
			want: []string{"a"},
		},
		{
			x:    []string{"a", "b"},
			y:    []string{"c", "b", "b"},
			want: []string{"a"},
		},
		{
			x:    []string{"a", "b", "a", "b", "b"},
			y:    []string{"a", "b", "b"},
			want: []string{},
		},
		{
			x:    []string{"a", "b", "a", "b", "b"},
			y:    []string{"a", "b", "b", "c", "c", "d"},
			want: []string{},
		},
		{
			x:    []string{"c", "b", "a", "a"},
			y:    []string{},
			want: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		require.Equal(t, tt.want, sortStrings(SetDiff(tt.x, tt.y)))
	}
}

func TestUniqueString(t *testing.T) {
	require.Equal(t, []string{}, sortStrings(Unique([]string{})))
	require.Equal(t, []string{"x"}, sortStrings(Unique([]string{"x"})))
	require.Equal(t, []string{"x"}, sortStrings(Unique([]string{"x", "x"})))
	require.Equal(t, []string{"x", "y"}, sortStrings(Unique([]string{"y", "x", "x"})))
	require.Equal(t, []string{"x", "y", "z"}, sortStrings(Unique([]string{"z", "y", "x", "x", "z"})))
}

func TestUniqueInt(t *testing.T) {
	require.Equal(t, []int{}, sortInts(Unique([]int{})))
	require.Equal(t, []int{1}, sortInts(Unique([]int{1})))
	require.Equal(t, []int{1}, sortInts(Unique([]int{1, 1})))
	require.Equal(t, []int{1, 2}, sortInts(Unique([]int{2, 1, 1})))
	require.Equal(t, []int{1, 2, 3}, sortInts(Unique([]int{3, 2, 1, 1, 3})))
}

func sortStrings(s []string) []string {
	sort.Strings(s)
	return s
}

func sortInts(s []int) []int {
	sort.Ints(s)
	return s
}

func TestNormalizeString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Lowercases",
			input: "gIUGkjbf89hfHGFJ",
			want:  "giugkjbf89hfhgfj",
		},
		{
			name:  "Restricts length to first 1024 runes",
			input: "giugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfj",
			want:  "giugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfjgiugkjbf89hfhgfj",
		},
		{
			name:  "Replaces a range of special characters with their equivalents",
			input: "äàáâçèéêëîïíñöôœßüùúûÿ",
			want:  "aaaaceeeeiiinoooessuuuuy",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeString(tt.input); got != tt.want {
				t.Errorf("NormalizeString() = %v, want %v", got, tt.want)
			}
		})
	}
}
