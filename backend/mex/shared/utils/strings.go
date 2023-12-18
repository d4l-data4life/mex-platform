package utils

import (
	"strings"
)

// Returns A ⊆ B
// Note: Duplicate values are taken into account:
//
// [ foo, bar, foo ] ⊆ [ foo, bar, foo, wom ]
// [ foo, bar ] ⊆ [ foo, bar, foo, wom ]
// but
// [ foo, bar, foo ] ⊈ [ foo, bar, wom ]
func SubsetEq[V comparable](a []V, b []V) bool {
	dictA := make(map[V]int)
	dictB := make(map[V]int)

	// Fill dictA with A
	for _, a := range a {
		if na, ok := dictA[a]; ok {
			dictA[a] = na + 1
		} else {
			dictA[a] = 1
		}
	}

	// Fill dictB with B
	for _, b := range b {
		if nb, ok := dictB[b]; ok {
			dictB[b] = nb + 1
		} else {
			dictB[b] = 1
		}
	}

	for a, na := range dictA {
		if nb, ok := dictB[a]; ok {
			if nb < na {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

func Contains[V comparable](list []V, elem V) bool {
	for _, s := range list {
		if s == elem {
			return true
		}
	}
	return false
}

// Returns sort(set(a) \ set(b))
// (duplicates are removed)
func SetDiff[V comparable](a []V, b []V) []V {
	m := make(map[V]struct{})

	for _, s := range a {
		m[s] = struct{}{}
	}

	for _, s := range b {
		delete(m, s)
	}

	r := make([]V, len(m))
	i := 0
	for k := range m {
		r[i] = k
		i++
	}

	return r
}

func Unique[V comparable](a []V) []V {
	return SetDiff(a, []V{})
}

type str string

func (s str) replace(old string, new string) str {
	return str(strings.ReplaceAll(string(s), old, new))
}

func SanitizeXML(s string) string {
	return string(str(s).
		replace("&", "&amp;"). // replace this first in order not to replace the & of the subsequent patterns
		replace("<", "&lt;").
		replace(">", "&gt;"))
}

/*
NormalizeString normalized a string in the following ways:

1: restricts to the first 1024 runes
2: lowercases
3: maps certain non-standard ASCII-characters to one or more standard ASCII-characters

The aim of step (3) is to improve the sort ordering, in particular giving the correct sorting of German
text. However, it does NOT ensure correct alphabetic ordering in all cases - indeed, no simple mapping can achieve that
since sorting conventions differs between languages and may disagree. However, the mapping is not a pure ASCII folding
either, since some strings (like the german 'ß') are mapped to multiple characters to ensure the correct sorting.
*/
func NormalizeString(s string) string {
	const MaxNormalizedStringLength = 1024
	replaceMap := map[rune][]rune{
		'ä': []rune("a"),
		'à': []rune("a"),
		'á': []rune("a"),
		'â': []rune("a"),
		'ç': []rune("c"),
		'è': []rune("e"),
		'é': []rune("e"),
		'ê': []rune("e"),
		'ë': []rune("e"),
		'î': []rune("i"),
		'ï': []rune("i"),
		'í': []rune("i"),
		'ñ': []rune("n"),
		'ö': []rune("o"),
		'ô': []rune("o"),
		'œ': []rune("oe"),
		'ß': []rune("ss"),
		'ü': []rune("u"),
		'ù': []rune("u"),
		'ú': []rune("u"),
		'û': []rune("u"),
		'ÿ': []rune("y"),
	}
	transformedS := s
	// Restrict to the first MaxNormalizedStringLength characters
	if len(transformedS) > MaxNormalizedStringLength {
		transformedS = transformedS[:MaxNormalizedStringLength]
	}
	// Lowercase
	transformedS = strings.ToLower(transformedS)
	var runeVersion []rune
	for _, curRune := range transformedS {
		if replacementRunes, ok := replaceMap[curRune]; ok {
			runeVersion = append(runeVersion, replacementRunes...)
		} else {
			runeVersion = append(runeVersion, curRune)
		}
	}

	return string(runeVersion)
}
