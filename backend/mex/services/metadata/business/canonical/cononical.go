package canonical

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"

	"github.com/d4l-data4life/mex/mex/shared/items"
)

// The fingerprint of an Item is the uppercase hexadecimal SHA-256 hash of its canonical representation.
func Fingerprint(item *items.Item) string {
	s1 := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s", item.EntityType, item.BusinessId, canonicalizedValues(item.Values))))
	return hexa(s1[:], "")
}

// Canonicalization rules:
//   - Collect all values of the same item (preserving order).
//   - Order fields lexicographically.
//   - Turn into tuples
//   - Marshal to string removing any superfluous whitespace.
//
// Example:
//
// Input:
// [
//
//	{
//	    fieldName: "keyword",
//	    fieldValue: "foo"
//	},
//	{
//	    fieldName: "title",
//	    fieldValue: "This is an integration test item",
//	    language: "en"
//	},
//	{
//	    fieldName: "author",
//	    fieldValue: "Mr A"
//	},
//	{
//	    fieldName: "author",
//	    fieldValue: "Mr B"
//	},
//	{
//	    fieldName: "abstract",
//	    fieldValue: "Lorem ipsum"
//	},
//	{
//	    fieldName: "keyword",
//	    fieldValue: "bar"
//	},
//	{
//	    fieldName: "author",
//	    fieldValue: "Mr C"
//	}
//
// ]
//
// Canonicalized string:
//
// ((abstract:(("Lorem ipsum",""))),(author:(("Mr A",""),("Mr B",""),("Mr C",""))),(keyword:(("foo",""),("bar",""))),(title:(("This is an integration test item","en"))))
//
// --
func canonicalizedValues(values []*items.ItemValue) string {
	// Pivot the field values list into a map assembling all values as a slice.
	fmap := make(map[string]*F)
	for _, value := range values {
		if f, ok := fmap[value.FieldName]; ok {
			f.vs = append(f.vs, &V{v: value.FieldValue, l: value.Language})
		} else {
			fmap[value.FieldName] = &F{
				n:  value.FieldName,
				vs: []*V{{v: value.FieldValue, l: value.Language}},
			}
		}
	}

	// Un-pivot back into slice, one entry per field.
	fslice := make(Fs, len(fmap))
	i := 0
	for _, v := range fmap {
		fslice[i] = v
		i++
	}

	sort.Sort(fslice)

	return fslice.String()
}

type V struct {
	v string
	l string
}

type Vs []*V

type F struct {
	n  string
	vs Vs
}

type Fs []*F

func (f F) String() string {
	return fmt.Sprintf("(%s:%s)", f.n, f.vs.String())
}

func (fs Fs) String() string {
	s := make([]string, len(fs))
	for i, f := range fs {
		s[i] = f.String()
	}
	return fmt.Sprintf("(%s)", strings.Join(s, ","))
}

func (v V) String() string {
	// We do not escape the string v.v here as it would not change the canonicalization properties.
	return fmt.Sprintf("('%s','%s')", v.v, v.l)
}

func (vs Vs) String() string {
	s := make([]string, len(vs))
	for i, v := range vs {
		s[i] = v.String()
	}
	return fmt.Sprintf("(%s)", strings.Join(s, ","))
}

func (fs Fs) Len() int { return len(fs) }

func (fs Fs) Less(i, j int) bool {
	return fs[i].n < fs[j].n
}

func (fs Fs) Swap(i, j int) { fs[i], fs[j] = fs[j], fs[i] }

func hexa(hash []byte, sep string) string {
	s := make([]string, len(hash))

	for i, b := range hash {
		s[i] = fmt.Sprintf("%02X", b)
	}

	return strings.Join(s, sep)
}
