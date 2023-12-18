package parser

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

const DefaultMaxTermEditDistance uint32 = 0

type TestDef struct {
	name                     string
	mexQuery                 string
	maxTermEditDistance      uint32
	useNgramFuzz             bool
	expectedSolrQuery        string
	expectedCleanedQuery     string
	expectedQueryWasCleaned  bool
	expectedErrorType        string
	expectedPhrasesOnlyQuery bool
}

func checkPhraseOnlyStatus(t *testing.T, tt TestDef, gotParseResult *QueryParseResult, err TypedParserError) {
	if err != nil && err.GetType() != tt.expectedErrorType {
		t.Errorf("mexQuery mapping: mexQuery error of type %s and msg. '%s', "+
			"with sub-msgs. '%s' - expected error type %v", err.GetType(),
			err.Error(), strings.Join(err.GetMessages(), ", "),
			tt.expectedErrorType)
		return
	}
	if gotParseResult.PhrasesOnlyQuery != tt.expectedPhrasesOnlyQuery {
		t.Errorf(`mexQuery mapping: wanted phrase-only flag '%v' but got '%v'`, tt.expectedPhrasesOnlyQuery, gotParseResult.PhrasesOnlyQuery)
	}
}

func checkQueryGenerationOutput(t *testing.T, tt TestDef, gotParseResult *QueryParseResult, err TypedParserError) {
	if err == nil && tt.expectedErrorType != "" {
		t.Errorf("mexQuery mapping: expected error type %v but did not get any errors", tt.expectedErrorType)
		return
	} else if err != nil && err.GetType() != tt.expectedErrorType {
		t.Errorf("mexQuery mapping: mexQuery error of type %s and msg. '%s', "+
			"with sub-msgs. '%s' - expected error type %v", err.GetType(),
			err.Error(), strings.Join(err.GetMessages(), ", "),
			tt.expectedErrorType)
		return
	}
	if gotParseResult.SolrQuery != tt.expectedSolrQuery {
		t.Errorf(`mexQuery mapping: wanted mexQuery '%s' but got '%s'`, tt.expectedSolrQuery, gotParseResult.SolrQuery)
	}
	if tt.expectedCleanedQuery != "" && gotParseResult.CleanedQuery != tt.expectedCleanedQuery {
		t.Errorf(`mexQuery mapping: wanted cleaned mexQuery '%s' but got '%s'`, tt.expectedCleanedQuery, gotParseResult.CleanedQuery)
	}
	if tt.expectedQueryWasCleaned != gotParseResult.QueryWasCleaned {
		t.Errorf(`mexQuery mapping: wanted mexQuery to be cleaned '%v' but got '%v'`,
			tt.expectedQueryWasCleaned, gotParseResult.QueryWasCleaned)
	}
}

var standardMatchingOpsConfig = MatchingOpsConfig{
	Term: []MatchingFieldConfig{
		{
			FieldName:       "fieldA",
			BoostFactor:     "3",
			MaxEditDistance: 0,
		},
		{
			FieldName:       "fieldB",
			BoostFactor:     "0.5",
			MaxEditDistance: 1,
		},
	},
	Phrase: []MatchingFieldConfig{
		{
			FieldName:       "fieldC",
			BoostFactor:     "2",
			MaxEditDistance: 0,
		},
		{
			FieldName:       "fieldD",
			BoostFactor:     "",
			MaxEditDistance: 0,
		},
	},
}

func termMatcher(term string) string {
	if len(term) <= 4 {
		return fmt.Sprintf("((fieldA:%s)^3 OR (fieldB:%s)^0.5)", term, term)
	}
	return fmt.Sprintf("((fieldA:%s)^3 OR (fieldB:%s~1)^0.5)", term, term)
}
func phraseMatcher(phrase string) string {
	quotedPhrase := fmt.Sprintf("\"%s\"", phrase)
	return fmt.Sprintf("((fieldC:%s)^2 OR fieldD:%s)", quotedPhrase, quotedPhrase)
}

func sortMarchingConfig(matchingConfig MatchingOpsConfig) {
	sort.Slice(matchingConfig.Term, func(i, j int) bool {
		return matchingConfig.Term[i].FieldName < matchingConfig.Term[j].FieldName
	})
	sort.Slice(matchingConfig.Phrase, func(i, j int) bool {
		return matchingConfig.Phrase[i].FieldName < matchingConfig.Phrase[j].FieldName
	})
}

func Test_MatchingOp_ParseSearchQuery_Basics(t *testing.T) {

	tests := []TestDef{
		{
			name:              "Completely empty mexQuery string request is mapped to an get-all term query",
			mexQuery:          ``,
			expectedSolrQuery: termMatcher("*"),
		},
		{
			name:              "Query string containing only spaces is mapped to an get-all term query",
			mexQuery:          `    `,
			expectedSolrQuery: termMatcher("*"),
		},
		{
			name:              "Query string '*' is mapped to  to an get-all term query",
			mexQuery:          `*`,
			expectedSolrQuery: termMatcher("*"),
		},
		{
			name:              "Single term is left unchanged",
			mexQuery:          `hello`,
			expectedSolrQuery: termMatcher("hello"),
		},
		{
			name:              "Multiple search terms separated by spaces are ANDed together",
			mexQuery:          `hello there again`,
			expectedSolrQuery: fmt.Sprintf("%s AND %s AND %s", termMatcher("hello"), termMatcher("there"), termMatcher("again")),
		},
		{
			name: "Tabs, carriage returns, " +
				"and newlines are interpreted as space characters and hence separate terms",
			mexQuery:          "hello\tthere\ragain\ntest",
			expectedSolrQuery: fmt.Sprintf("%s AND %s AND %s AND %s", termMatcher("hello"), termMatcher("there"), termMatcher("again"), termMatcher("test")),
		},
		{
			name: "Escaped tabs, carriage returns, " +
				"and newlines and not interpreted as spaces but kept unchanged",
			mexQuery:          `hello\\tthere\\ragain\\ntest`,
			expectedSolrQuery: termMatcher(`hello\tthere\ragain\ntest`),
		},
		{
			name: "Mixtures of terms, " +
				"quoted terms and Boolean expression are converted and the results ANDed together - Boolean" +
				" expressions are bracketed",
			mexQuery:          `hello "you there" he + said | joyfully`,
			expectedSolrQuery: fmt.Sprintf(`%s AND %s AND ((%s AND %s) OR %s)`, termMatcher("hello"), phraseMatcher("you there"), termMatcher("he"), termMatcher("said"), termMatcher("joyfully")),
		},
		{
			name: "Escape characters in front of characters that are neither MEx nor Solr symbols are" +
				" removed",
			mexQuery:          `first\second\ third\@`,
			expectedSolrQuery: fmt.Sprintf("%s AND %s", termMatcher("firstsecond"), termMatcher("third@")),
		},
	}
	for _, tt := range tests {
		sortMarchingConfig(standardMatchingOpsConfig)
		listener, _ := NewListener(standardMatchingOpsConfig)
		parser := NewMexParser(listener)
		t.Run(tt.name, func(t *testing.T) {
			gotParseResult, err := parser.ConvertToSolrQuery(tt.mexQuery)
			checkQueryGenerationOutput(t, tt, gotParseResult, err)
		})
	}
}

func Test_MatchingOp_PhraseQueryClassification(t *testing.T) {
	tests := []TestDef{
		{
			name:                     "A single-word phrase is classified as a phrase-only query",
			mexQuery:                 `"quoted"`,
			expectedPhrasesOnlyQuery: true,
		},
		{
			name:                     "A multiple-word phrase is classified as a phrase-only query",
			mexQuery:                 `"a quoted phrase"`,
			expectedPhrasesOnlyQuery: true,
		},
		{
			name:                     "Several multiple-word phrases is classified as a phrase-only query",
			mexQuery:                 `"a quoted phrase" "and another one"`,
			expectedPhrasesOnlyQuery: true,
		},
		{
			name:                     "Several phrases combined with Boolean operators is classified as a phrase-only query",
			mexQuery:                 `("a quoted phrase" + "and another one") | -"last chance"`,
			expectedPhrasesOnlyQuery: true,
		},
		{
			name:                     "Several words combined with Boolean operators is NOT classified as a phrase-only query",
			mexQuery:                 `(word + more) | -another`,
			expectedPhrasesOnlyQuery: false,
		},
		{
			name:                     "A mixture of words and phrases combined with Boolean operators is NOT classified as a phrase-only query",
			mexQuery:                 `("a quoted phrase" + word) | -"last chance"`,
			expectedPhrasesOnlyQuery: false,
		},
		{
			name:                     "A single unquoted word is NOT classified as a phrase-only query",
			mexQuery:                 `word`,
			expectedPhrasesOnlyQuery: false,
		},
		{
			name:                     "Multiple unquoted words is NOT classified as a phrase-only query",
			mexQuery:                 `a word and another`,
			expectedPhrasesOnlyQuery: false,
		},
		{
			name:                     "A mixture of quoted and unquoted single words is NOT classified as a phrase-only query",
			mexQuery:                 `"word" and "another" as well"`,
			expectedPhrasesOnlyQuery: false,
		},
		{
			name:                     "A mixture of quoted and unquoted words is NOT classified as a phrase-only query",
			mexQuery:                 `"a short phrase" and another "as well"`,
			expectedPhrasesOnlyQuery: false,
		},
	}
	for _, tt := range tests {
		sortMarchingConfig(standardMatchingOpsConfig)
		listener, _ := NewListener(standardMatchingOpsConfig)
		parser := NewMexParser(listener)
		t.Run(tt.name, func(t *testing.T) {
			gotParseResult, err := parser.ConvertToSolrQuery(tt.mexQuery)
			checkPhraseOnlyStatus(t, tt, gotParseResult, err)
		})
	}
}

func Test_MatchingOp_ParseSearchQuery_Phrases(t *testing.T) {
	tests := []TestDef{
		{
			name:                    "Query consisting of just a single double quote is mapped to the empty search",
			mexQuery:                `"`,
			expectedSolrQuery:       ``,
			expectedQueryWasCleaned: true,
		},
		{
			name:              "Query consisting of just a single escaped double quote is kept unchanged",
			mexQuery:          `\"`,
			expectedSolrQuery: termMatcher(`\"`),
		},
		{
			name:              "Empty quote is kept unchanged",
			mexQuery:          `""`,
			expectedSolrQuery: phraseMatcher(""),
		},
		{
			name:              "Quoted single-word phrases are treated as single terms",
			mexQuery:          `"hello"`,
			expectedSolrQuery: phraseMatcher("hello"),
		},
		{
			name:              "Quoted phrases containing spaces are treated as single terms",
			mexQuery:          `"hello there"`,
			expectedSolrQuery: phraseMatcher("hello there"),
		},
		{
			name:              "Separate quoted phrases separated by spaces are ANDed with other search terms",
			mexQuery:          `"hello there" precedes "goodbye here"`,
			expectedSolrQuery: fmt.Sprintf("%s AND %s AND %s", phraseMatcher("hello there"), termMatcher("precedes"), phraseMatcher("goodbye here")),
		},
		{
			name:              "Separate quoted phrases not separated by spaces are ANDed together",
			mexQuery:          `"hello there""goodbye here"`,
			expectedSolrQuery: fmt.Sprintf("%s AND %s", phraseMatcher("hello there"), phraseMatcher("goodbye here")),
		},
		{
			name: "Quoted phrase followed by normal search term with an intervening space results in" +
				" ANDing the phrase and the term",
			mexQuery:          `"hello there" goodbye`,
			expectedSolrQuery: fmt.Sprintf("%s AND %s", phraseMatcher("hello there"), termMatcher("goodbye")),
		},
		{
			name: "Quoted phrase followed by normal search term without intervening space results in" +
				" ANDing the phrase and the term",
			mexQuery:          `"hello there"goodbye`,
			expectedSolrQuery: fmt.Sprintf("%s AND %s", phraseMatcher("hello there"), termMatcher("goodbye")),
		},
		{
			name:              "Empty quoted phrase combined with other terms are kept and treated as separate terms",
			mexQuery:          `"" hello there""`,
			expectedSolrQuery: fmt.Sprintf("%s AND %s AND %s AND %s", phraseMatcher(""), termMatcher("hello"), termMatcher("there"), phraseMatcher("")),
		},
		{
			name: "Double quotation at start of word not matched by double quote at end of" +
				" word is discarded",
			mexQuery:                `"hello there`,
			expectedSolrQuery:       fmt.Sprintf("%s AND %s", termMatcher("hello"), termMatcher("there")),
			expectedQueryWasCleaned: true,
		},
		{
			name: "Double quotation at end of word not matched by double quote at start of word is" +
				" discarded",
			mexQuery:                `hello there"`,
			expectedSolrQuery:       fmt.Sprintf("%s AND %s", termMatcher("hello"), termMatcher("there")),
			expectedQueryWasCleaned: true,
		},
		{
			name:              "Escaped double quotation symbol inside word is left unchanged",
			mexQuery:          `hello\"there`,
			expectedSolrQuery: termMatcher(`hello\"there`),
		},
		{
			name:              "Escaped double quotes do not start or end phrases",
			mexQuery:          `\"hello there\"`,
			expectedSolrQuery: fmt.Sprintf("%s AND %s", termMatcher(`\"hello`), termMatcher(`there\"`)),
		},
		{
			name:              "NOT-operator inside a word that is inside a double-quoted phrase is left unchanged",
			mexQuery:          `"hello-there"`,
			expectedSolrQuery: phraseMatcher("hello-there"),
		},
		{
			name:              "NOT-operator at the beginning of a word that is inside a double-quoted phrase is left unchanged",
			mexQuery:          `"hello -there"`,
			expectedSolrQuery: phraseMatcher("hello -there"),
		},
		{
			name:              "NOT-operator surrounded by space inside a double-quoted phrase is left unchanged",
			mexQuery:          `"hello - there"`,
			expectedSolrQuery: phraseMatcher("hello - there"),
		},
		{
			name:              "Escaped NOT-operator inside a word that is inside a double-quoted phrase is left unchanged",
			mexQuery:          `"hello\-there"`,
			expectedSolrQuery: phraseMatcher(`hello\-there`),
		},
		{
			name: "Escaped  NOT-operator at the beginning of a word that is inside a double-quoted" +
				" phrase is left unchanged",
			mexQuery:          `"hello \-there"`,
			expectedSolrQuery: phraseMatcher(`hello \-there`),
		},
		{
			name:              "Escaped NOT-operator sounded by space inside a double-quoted phrase is left unchanged",
			mexQuery:          `"hello \- there"`,
			expectedSolrQuery: phraseMatcher(`hello \- there`),
		},
		{
			name:              "OR-operator surrounded by space inside double-quoted term is left unchanged",
			mexQuery:          `"hello | there"`,
			expectedSolrQuery: phraseMatcher("hello | there"),
		},
		{
			name:              "OR-operator not surrounded by space inside double-quoted term is left unchanged",
			mexQuery:          `"hello|there"`,
			expectedSolrQuery: phraseMatcher("hello|there"),
		},
		{
			name:              "Escaped OR-operator not surrounded by space inside double-quoted term is left unchanged",
			mexQuery:          `"hello\|there"`,
			expectedSolrQuery: phraseMatcher(`hello\|there`),
		},
		{
			name:              "Escaped OR-operator surrounded by space inside double-quoted term is left unchanged",
			mexQuery:          `"hello \| there"`,
			expectedSolrQuery: phraseMatcher(`hello \| there`),
		},
		{
			name:              "AND-operator not surrounded by space inside double-quoted term is left unchanged",
			mexQuery:          `"hello+there"`,
			expectedSolrQuery: phraseMatcher("hello+there"),
		},
		{
			name:              "AND-operator not surrounded by space inside double-quoted term is left unchanged",
			mexQuery:          `"hello + there"`,
			expectedSolrQuery: phraseMatcher("hello + there"),
		},
		{
			name:              "Escaped AND-operator not surrounded by space inside double-quoted term is left unchanged",
			mexQuery:          `"hello\+there"`,
			expectedSolrQuery: phraseMatcher(`hello\+there`),
		},
		{
			name:              "Escaped AND-operator surrounded by space inside double-quoted term is left unchanged",
			mexQuery:          `"hello \+ there"`,
			expectedSolrQuery: phraseMatcher(`hello \+ there`),
		},
		{
			name:              "Brackets not surrounded by space inside double-quoted term are left unchanged",
			mexQuery:          `"hel)lo(there)"`,
			expectedSolrQuery: phraseMatcher("hel)lo(there)"),
		},
		{
			name:              "Brackets on term surrounded by space inside double-quoted term are left unchanged",
			mexQuery:          `"hello )there("`,
			expectedSolrQuery: phraseMatcher("hello )there("),
		},
		{
			name:              "Escaped brackets not surrounded by space inside double-quoted term are left unchanged",
			mexQuery:          `"hel\)lo\(there\)"`,
			expectedSolrQuery: phraseMatcher(`hel\)lo\(there\)`),
		},
		{
			name:              "Escaped brackets on term surrounded by space inside double-quoted term are left unchanged",
			mexQuery:          `"hello \)there\("`,
			expectedSolrQuery: phraseMatcher(`hello \)there\(`),
		},
	}
	for _, tt := range tests {
		sortMarchingConfig(standardMatchingOpsConfig)
		listener, _ := NewListener(standardMatchingOpsConfig)
		parser := NewMexParser(listener)
		t.Run(tt.name, func(t *testing.T) {
			gotParseResult, err := parser.ConvertToSolrQuery(tt.mexQuery)
			checkQueryGenerationOutput(t, tt, gotParseResult, err)
		})
	}
}

func Test_MatchingOp_ParseSearchQuery_Wildcard(t *testing.T) {
	tests := []TestDef{
		{
			name: "A mexQuery consisting of a single wildcard character is left unchanged and" +
				" not escaped",
			mexQuery:          `*`,
			expectedSolrQuery: "X:*",
		},
		{
			name:              "Wildcard character in word is left unchanged and not escaped",
			mexQuery:          `he*lo the*e again*`,
			expectedSolrQuery: `X:he*lo~1 AND X:the*e~1 AND X:again*~1`,
		},
		{
			name:              "An isolated wildcard character is kept and combined with other terms with AND, but is not fuzzied",
			mexQuery:          `* hello * again *`,
			expectedSolrQuery: `X:* AND X:hello~1 AND X:* AND X:again~1 AND X:*`,
		},
		{
			name:              "Escaped wildcard character is left unchanged",
			mexQuery:          `he\*lo the\*e`,
			expectedSolrQuery: `X:he\*lo~1 AND X:the\*e~1`,
		},
	}
	for _, tt := range tests {
		matchingOpConfig := MatchingOpsConfig{
			Term: []MatchingFieldConfig{
				{
					FieldName:       "X",
					MaxEditDistance: 1,
				},
			},
			Phrase: []MatchingFieldConfig{
				{
					FieldName:       "Y",
					MaxEditDistance: 0,
				},
			},
		}
		listener, _ := NewListener(matchingOpConfig)
		parser := NewMexParser(listener)
		t.Run(tt.name, func(t *testing.T) {
			gotParseResult, err := parser.ConvertToSolrQuery(tt.mexQuery)
			if err == nil && tt.expectedErrorType != "" {
				t.Errorf("mexQuery mapping: expected error type %v but did not get any errors", tt.expectedErrorType)
				return
			} else if err != nil && err.GetType() != tt.expectedErrorType {
				t.Errorf("mexQuery mapping: mexQuery error of type %s and msg. '%s', "+
					"with sub-msgs. '%s' - expected error type %v", err.GetType(),
					err.Error(), strings.Join(err.GetMessages(), ", "),
					tt.expectedErrorType)
				return
			}
			if gotParseResult.SolrQuery != tt.expectedSolrQuery {
				t.Errorf(`mexQuery mapping: wanted mexQuery '%s' but got '%s'`, tt.expectedSolrQuery, gotParseResult.SolrQuery)
			}
		})
	}
}

func Test_MatchingOp_ParseSearchQuery_Boolean_Not(t *testing.T) {
	tests := []TestDef{
		{
			name:              "Query consisting of only an escaped NOT-operator is left unchanged",
			mexQuery:          `\-`,
			expectedSolrQuery: termMatcher(`\-`),
		},
		{
			name:              "NOT-operator right in front of a word leads to negation",
			mexQuery:          `-hello`,
			expectedSolrQuery: fmt.Sprintf("(NOT %s)", termMatcher("hello")),
		},
		{
			name: "Escaped NOT-operator at beginning of word is left unchanged and does not lead to" +
				" negation",
			mexQuery:          `\-hello`,
			expectedSolrQuery: termMatcher(`\-hello`),
		},
		{
			name:              "NOT-operator inside a word does not lead to negation but leads to a bracketed AND-expression of the hyphenated parts",
			mexQuery:          `hello-there`,
			expectedSolrQuery: fmt.Sprintf("(%s AND %s)", termMatcher("hello"), termMatcher("there")),
		},
		{
			name:              "A NOT-operator that ends a word is treated as a part of the word",
			mexQuery:          `hello-`,
			expectedSolrQuery: termMatcher(`hello\-`),
		},
		{
			name:                    "A separated NOT operator at the very beginning of a mexQuery is discarded",
			mexQuery:                `- hello`,
			expectedSolrQuery:       termMatcher(`hello`),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "A single NOT operator surrounded by spaces inside larger mexQuery is discarded",
			mexQuery:                `hello - there`,
			expectedSolrQuery:       fmt.Sprintf("%s AND %s", termMatcher("hello"), termMatcher("there")),
			expectedQueryWasCleaned: true,
		},
		{
			name:              "Negated bracketed expressions are mapped to negated bracketed expression",
			mexQuery:          `-(hello there)`,
			expectedSolrQuery: fmt.Sprintf("(NOT (%s AND %s))", termMatcher("hello"), termMatcher("there")),
		},
		{
			name:                    "NOT-operator separated from bracketed expressions is discarded",
			mexQuery:                `- (hello there)`,
			expectedSolrQuery:       fmt.Sprintf("(%s AND %s)", termMatcher("hello"), termMatcher("there")),
			expectedQueryWasCleaned: true,
		},
	}

	for _, tt := range tests {
		sortMarchingConfig(standardMatchingOpsConfig)
		listener, _ := NewListener(standardMatchingOpsConfig)
		parser := NewMexParser(listener)
		t.Run(tt.name, func(t *testing.T) {
			gotParseResult, err := parser.ConvertToSolrQuery(tt.mexQuery)
			checkQueryGenerationOutput(t, tt, gotParseResult, err)
		})
	}
}

func Test_MatchingOp_ParseSearchQuery_Boolean_And(t *testing.T) {
	tests := []TestDef{
		{
			name:              "A single escaped AND operator as search term is kept",
			mexQuery:          `\+`,
			expectedSolrQuery: termMatcher(`\+`),
		},
		{
			name:              "AND-operator surrounded by spaces is mapped to AND-operator",
			mexQuery:          `hello + there + again`,
			expectedSolrQuery: fmt.Sprintf("%s AND %s AND %s", termMatcher("hello"), termMatcher("there"), termMatcher("again")),
		},
		{
			name:              "AND-operator not (or only partially) surrounded by spaces is mapped to AND-operator",
			mexQuery:          `hello+there+ again + today`,
			expectedSolrQuery: fmt.Sprintf("%s AND %s AND %s AND %s", termMatcher("hello"), termMatcher("there"), termMatcher("again"), termMatcher("today")),
		},
	}
	for _, tt := range tests {
		sortMarchingConfig(standardMatchingOpsConfig)
		listener, _ := NewListener(standardMatchingOpsConfig)
		parser := NewMexParser(listener)
		t.Run(tt.name, func(t *testing.T) {
			gotParseResult, err := parser.ConvertToSolrQuery(tt.mexQuery)
			checkQueryGenerationOutput(t, tt, gotParseResult, err)
		})
	}
}

func Test_MatchingOp_ParseSearchQuery_Boolean_Or(t *testing.T) {
	tests := []TestDef{
		{
			name:              "A single escaped OR operator as search term is kept",
			mexQuery:          `\|`,
			expectedSolrQuery: termMatcher(`\|`),
		},
		{
			name:              "OR-operator surrounded by space is mapped to OR-operator",
			mexQuery:          `hello | there | again`,
			expectedSolrQuery: fmt.Sprintf("%s OR %s OR %s", termMatcher("hello"), termMatcher("there"), termMatcher("again")),
		},
		{
			name:              "OR-operator not (or only partially) surrounded by space is mapped to OR-operator",
			mexQuery:          `hello|there| again |today `,
			expectedSolrQuery: fmt.Sprintf("%s OR %s OR %s OR %s", termMatcher("hello"), termMatcher("there"), termMatcher("again"), termMatcher("today")),
		},
	}
	for _, tt := range tests {
		sortMarchingConfig(standardMatchingOpsConfig)
		listener, _ := NewListener(standardMatchingOpsConfig)
		parser := NewMexParser(listener)
		t.Run(tt.name, func(t *testing.T) {
			gotParseResult, err := parser.ConvertToSolrQuery(tt.mexQuery)
			checkQueryGenerationOutput(t, tt, gotParseResult, err)
		})
	}
}

func Test_MatchingOp_ParseSearchQuery_Boolean_Combinations(t *testing.T) {
	tests := []TestDef{
		{
			name:              "AND-operator surrounded by spaces is mapped to AND-operator",
			mexQuery:          `hello + there + again`,
			expectedSolrQuery: fmt.Sprintf("%s AND %s AND %s", termMatcher("hello"), termMatcher("there"), termMatcher("again")),
		},
		{
			name: "Precedence of explicit operators is NOT > AND > OR, " +
				"and is made explicit by brackets",
			mexQuery:          `-hello | there + again`,
			expectedSolrQuery: fmt.Sprintf("(NOT %s) OR (%s AND %s)", termMatcher("hello"), termMatcher("there"), termMatcher("again")),
		},
		{
			name:              "Whitespace acting as AND has lower precedence than an explicit AND",
			mexQuery:          `hello there + again`,
			expectedSolrQuery: fmt.Sprintf("%s AND (%s AND %s)", termMatcher("hello"), termMatcher("there"), termMatcher("again")),
		},
		{
			name:              "Whitespace acting as AND has lower precedence than an explicit OR",
			mexQuery:          `-hello | there again`,
			expectedSolrQuery: fmt.Sprintf("((NOT %s) OR %s) AND %s", termMatcher("hello"), termMatcher("there"), termMatcher("again")),
		},
		{
			name: "In mixtures of words/phrases and Boolean expression, " +
				"the Boolean expressions are converted before being combined with terms/phrases (" +
				"follows from precedence rules)",
			mexQuery:          `ear hand + foot | nose "eyes open"`,
			expectedSolrQuery: fmt.Sprintf("%s AND ((%s AND %s) OR %s) AND %s", termMatcher("ear"), termMatcher("hand"), termMatcher("foot"), termMatcher("nose"), phraseMatcher("eyes open")),
		},
		{
			name: "In combinations of valid Boolean expression and dangling operators, " +
				"the Boolean expression is still correctly parsed",
			mexQuery:                `hand | -foot +`,
			expectedSolrQuery:       fmt.Sprintf("%s OR (NOT %s)", termMatcher("hand"), termMatcher("foot")),
			expectedQueryWasCleaned: true,
		},
		{
			name:              "Whitespace acts as AND even inside complex Boolean expression",
			mexQuery:          `-hello + there | here never`,
			expectedSolrQuery: fmt.Sprintf("(((NOT %s) AND %s) OR %s) AND %s", termMatcher("hello"), termMatcher("there"), termMatcher("here"), termMatcher("never")),
		},
		{
			name:              "Explicit brackets override normal operator precedence",
			mexQuery:          `-((hello | there) + here)`,
			expectedSolrQuery: fmt.Sprintf("(NOT ((%s OR %s) AND %s))", termMatcher("hello"), termMatcher("there"), termMatcher("here")),
		},
		{
			name:              "Very complicated combinations are correctly handled",
			mexQuery:          `(-(hand| foot) nose) + (("ear eye" + -mo*)) leg`,
			expectedSolrQuery: fmt.Sprintf(`(((NOT (%s OR %s)) AND %s) AND ((%s AND (NOT %s)))) AND %s`, termMatcher("hand"), termMatcher("foot"), termMatcher("nose"), phraseMatcher("ear eye"), termMatcher("mo*"), termMatcher("leg")),
		},
	}
	for _, tt := range tests {
		sortMarchingConfig(standardMatchingOpsConfig)
		listener, _ := NewListener(standardMatchingOpsConfig)
		parser := NewMexParser(listener)
		t.Run(tt.name, func(t *testing.T) {
			gotParseResult, err := parser.ConvertToSolrQuery(tt.mexQuery)
			checkQueryGenerationOutput(t, tt, gotParseResult, err)
		})
	}
}

func Test_MatchingOp_ParseSearchQuery_Dangling_Operators(t *testing.T) {
	tests := []TestDef{
		{
			name: "Double quotation at start of word not matched by double quote at end of" +
				" word is discarded",
			mexQuery:                `"hello there`,
			expectedSolrQuery:       fmt.Sprintf("%s AND %s", termMatcher("hello"), termMatcher("there")),
			expectedQueryWasCleaned: true,
		},
		{
			name: "Double quotation at end of word not matched by double quote at start of word is" +
				" discarded",
			mexQuery:                `hello there"`,
			expectedSolrQuery:       fmt.Sprintf("%s AND %s", termMatcher("hello"), termMatcher("there")),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "Query consisting of only a NOT-operator maps to the empty query",
			mexQuery:                `-`,
			expectedSolrQuery:       ``,
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "A separated NOT operator at the very end of a query is discarded",
			mexQuery:                `hello -`,
			expectedSolrQuery:       termMatcher("hello"),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "A single NOT operator surrounded by spaces inside larger mexQuery is discarded",
			mexQuery:                `hello - there`,
			expectedSolrQuery:       fmt.Sprintf("%s AND %s", termMatcher("hello"), termMatcher("there")),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "A separated NOT operator at the very beginning of a mexQuery is discarded",
			mexQuery:                `- hello`,
			expectedSolrQuery:       termMatcher("hello"),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "A single AND operator as search term is discarded",
			mexQuery:                `+`,
			expectedSolrQuery:       ``,
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "Separated AND-operator with missing first operand is discarded",
			mexQuery:                `+ hello`,
			expectedSolrQuery:       termMatcher("hello"),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "AND-operator that starts a word but has a missing first operand is discarded",
			mexQuery:                `+hello`,
			expectedSolrQuery:       termMatcher("hello"),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "Separated AND-operator with missing second operand is discarded",
			mexQuery:                `hello +`,
			expectedSolrQuery:       termMatcher("hello"),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "AND-operator that ends a word but has a missing second operand is discarded",
			mexQuery:                `hello+`,
			expectedSolrQuery:       termMatcher("hello"),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "A single OR operator as search term is discarded",
			mexQuery:                `|`,
			expectedSolrQuery:       ``,
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "Separated OR-operator with missing first operand is discarded",
			mexQuery:                `| hello`,
			expectedSolrQuery:       termMatcher("hello"),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "OR-operator that starts a word but has a missing first operand is discarded",
			mexQuery:                `|hello`,
			expectedSolrQuery:       termMatcher("hello"),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "OR-operator with missing second operand is discarded",
			mexQuery:                `hello |`,
			expectedSolrQuery:       termMatcher("hello"),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "OR-operator that ends a word but has a missing first operand is discarded",
			mexQuery:                `hello|`,
			expectedSolrQuery:       termMatcher("hello"),
			expectedQueryWasCleaned: true,
		},
		{
			name: "In combinations of valid Boolean expression and dangling operators, " +
				"the Boolean expression is still correctly parsed",
			mexQuery:                `hand | -foot +`,
			expectedSolrQuery:       fmt.Sprintf("%s OR (NOT %s)", termMatcher("hand"), termMatcher("foot")),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "Combinations of valid Boolean expression and dangling operator work",
			mexQuery:                `hand | -foot +`,
			expectedSolrQuery:       fmt.Sprintf("%s OR (NOT %s)", termMatcher("hand"), termMatcher("foot")),
			expectedQueryWasCleaned: true,
		},
		{
			name:              "Very complicated combinations are correctly handled",
			mexQuery:          `(-(hand| foot) nose) + (("ear eye" + -mo*)) leg`,
			expectedSolrQuery: fmt.Sprintf("(((NOT (%s OR %s)) AND %s) AND ((%s AND (NOT %s)))) AND %s", termMatcher("hand"), termMatcher("foot"), termMatcher("nose"), phraseMatcher("ear eye"), termMatcher("mo*"), termMatcher("leg")),
		},
		{
			name:                    "Very complicated combinations with broken parts are correctly handled",
			mexQuery:                `(-(hand| foot) nose) + (("ear eye" + -mo*)) ("leg +`,
			expectedSolrQuery:       fmt.Sprintf("(((NOT (%s OR %s)) AND %s) AND ((%s AND (NOT %s)))) AND %s", termMatcher("hand"), termMatcher("foot"), termMatcher("nose"), phraseMatcher("ear eye"), termMatcher("mo*"), termMatcher("leg")),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "NOT-operator separated from bracketed expressions is discarded",
			mexQuery:                `- (hello there)`,
			expectedSolrQuery:       fmt.Sprintf("(%s AND %s)", termMatcher("hello"), termMatcher("there")),
			expectedQueryWasCleaned: true,
		},
		{
			name: "Unmatched double quote directly following an operator is discarded and the term" +
				" directly following (no intervening space) is taken as the second operand of the operator",
			mexQuery:                `hello | "there`,
			expectedSolrQuery:       fmt.Sprintf("%s OR %s", termMatcher("hello"), termMatcher("there")),
			expectedQueryWasCleaned: true,
		},
		{
			name: "Unmatched double quote directly following an operator is discarded and the term" +
				" following (with intervening space) is taken as the second operand of the operator",
			mexQuery:                `hello | " there`,
			expectedSolrQuery:       fmt.Sprintf("%s OR %s", termMatcher("hello"), termMatcher("there")),
			expectedQueryWasCleaned: true,
		},
		{
			name: "Unmatched bracket directly following an operator is discarded and the" +
				" term" +
				" directly following (no intervening space) is taken as the second operand of the operator",
			mexQuery:                `hello | (there`,
			expectedSolrQuery:       fmt.Sprintf("%s OR %s", termMatcher("hello"), termMatcher("there")),
			expectedQueryWasCleaned: true,
		},
		{
			name: "Unmatched bracket directly following an operator is discarded and the" +
				" term" +
				" following (with intervening space) is taken as the second operand of the operator",
			mexQuery:                `hello | ( there`,
			expectedSolrQuery:       fmt.Sprintf("%s OR %s", termMatcher("hello"), termMatcher("there")),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "Unmatched opening bracket is discarded",
			mexQuery:                `(hello again`,
			expectedSolrQuery:       fmt.Sprintf("%s AND %s", termMatcher("hello"), termMatcher("again")),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "Unmatched closing bracket at end of word is discarded",
			mexQuery:                `hello) again`,
			expectedSolrQuery:       fmt.Sprintf("%s AND %s", termMatcher("hello"), termMatcher("again")),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "Unmatched closing bracket at end of word is discarded",
			mexQuery:                `hello) again`,
			expectedSolrQuery:       fmt.Sprintf("%s AND %s", termMatcher("hello"), termMatcher("again")),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "Separate and unmatched opening bracket is discarded",
			mexQuery:                `hello  ( again`,
			expectedSolrQuery:       fmt.Sprintf("%s AND %s", termMatcher("hello"), termMatcher("again")),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "Separate and unmatched closing bracket is discarded",
			mexQuery:                `hello there ) again`,
			expectedSolrQuery:       fmt.Sprintf("%s AND %s AND %s", termMatcher("hello"), termMatcher("there"), termMatcher("again")),
			expectedQueryWasCleaned: true,
		},
		{
			name: `Regression test: Does not swallow NOT operator if there is a dangling operator between NOT and
actual operand`,
			mexQuery:                `-"paris`,
			expectedSolrQuery:       fmt.Sprintf("(NOT %s)", termMatcher("paris")),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    `Regression test: OR-operator not swallowed in query test -covid [OR] "sars`,
			mexQuery:                `test -covid | "sars`,
			expectedSolrQuery:       fmt.Sprintf("%s AND ((NOT %s) OR %s)", termMatcher("test"), termMatcher("covid"), termMatcher("sars")),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    `Regression test: OR-operator not swallowed in query (paris [OR] (rome + venice)`,
			mexQuery:                `(paris | (rome + venice)`,
			expectedSolrQuery:       fmt.Sprintf("(%s OR (%s AND %s))", termMatcher("paris"), termMatcher("rome"), termMatcher("venice")),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    `Regression test: OR-operator not swallowed in query cov * [OR] "sars" + "test"-"19`,
			mexQuery:                `cov * | "sars" + "test"-"19`,
			expectedSolrQuery:       fmt.Sprintf("%s AND (%s OR (%s AND %s)) AND (NOT %s)", termMatcher("cov"), termMatcher("*"), phraseMatcher("sars"), phraseMatcher("test"), termMatcher("19")),
			expectedQueryWasCleaned: true,
		},
	}
	for _, tt := range tests {
		sortMarchingConfig(standardMatchingOpsConfig)
		listener, _ := NewListener(standardMatchingOpsConfig)
		parser := NewMexParser(listener)
		t.Run(tt.name, func(t *testing.T) {
			gotParseResult, err := parser.ConvertToSolrQuery(tt.mexQuery)
			checkQueryGenerationOutput(t, tt, gotParseResult, err)
		})
	}
}

func Test_MatchingOp_ParseSearchQuery_Boolean_CleanedQueryString(t *testing.T) {
	tests := []TestDef{
		{
			name:                 "Cleaned query is correct for simple combinations",
			mexQuery:             `-hand (foot | nose)`,
			expectedCleanedQuery: `(-hand) (foot | nose)`,
		},
		{
			name:                 "Cleaned query is correct for very complicated combinations",
			mexQuery:             `(-(hand| foot) nose) + (("ear eye" + -mo*)) leg`,
			expectedCleanedQuery: `(((-(hand | foot)) + nose) + (("ear eye" + (-mo*)))) leg`,
		},
		{
			name:                    "Cleaned query is correct very complicated combinations handled",
			mexQuery:                `(-(hand| foot) nose) + (("ear eye" + -mo*)) ("leg +`,
			expectedCleanedQuery:    `(((-(hand | foot)) + nose) + (("ear eye" + (-mo*)))) leg`,
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "Cleaned query is correct for simple combinations",
			mexQuery:                `"hello)`,
			expectedCleanedQuery:    `hello`,
			expectedQueryWasCleaned: true,
		},
	}
	for _, tt := range tests {
		sortMarchingConfig(standardMatchingOpsConfig)
		listener, _ := NewListener(standardMatchingOpsConfig)
		parser := NewMexParser(listener)
		t.Run(tt.name, func(t *testing.T) {
			gotParseResult, err := parser.ConvertToSolrQuery(tt.mexQuery)
			if err == nil && tt.expectedErrorType != "" {
				t.Errorf("mexQuery mapping: expected error type %v but did not get any errors", tt.expectedErrorType)
				return
			} else if err != nil && err.GetType() != tt.expectedErrorType {
				t.Errorf("mapping of mexQuery '%s': mexQuery error of type %s and msg. '%s', "+
					"with sub-msgs. '%s' - expected error type %v", tt.mexQuery, err.GetType(),
					err.Error(), strings.Join(err.GetMessages(), ", "),
					tt.expectedErrorType)
				return
			}
			if gotParseResult.CleanedQuery != tt.expectedCleanedQuery {
				t.Errorf(`mexQuery mapping: wanted cleaned mexQuery '%s' but got '%s'`, tt.expectedCleanedQuery, gotParseResult.CleanedQuery)
			}
			if tt.expectedQueryWasCleaned != gotParseResult.QueryWasCleaned {
				t.Errorf(`mexQuery mapping: wanted mexQuery to be cleaned '%s' but got '%s'`, tt.expectedCleanedQuery, gotParseResult.CleanedQuery)
			}
		})
	}
}

func Test_MatchingOp_ParseSearchQuery_Brackets(t *testing.T) {
	tests := []TestDef{
		{
			name:                    "A single opening bracket as search term is discarded",
			mexQuery:                `(`,
			expectedSolrQuery:       ``,
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "A single closing bracket as search term is discarded",
			mexQuery:                `)`,
			expectedSolrQuery:       ``,
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "Query containing only an empty bracket expression maps to the empty query",
			mexQuery:                `()`,
			expectedSolrQuery:       ``,
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "Query containing only bracket around spaces maps to the empty query",
			mexQuery:                `(  )`,
			expectedSolrQuery:       ``,
			expectedQueryWasCleaned: true,
		},
		{
			name: "Empty bracket expression is discarded when combined with" +
				" other space-separated terms",
			mexQuery:                `() hello`,
			expectedSolrQuery:       termMatcher("hello"),
			expectedQueryWasCleaned: true,
		},
		{
			name: "Empty bracket expression inside a word is discarded but" +
				" split the word into two",
			mexQuery:                `hel()lo`,
			expectedSolrQuery:       fmt.Sprintf("%s AND %s", termMatcher("hel"), termMatcher("lo")),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "Empty bracket expression at front or end of word or inside a word are discarded",
			mexQuery:                `()hello`,
			expectedSolrQuery:       termMatcher("hello"),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "Empty bracket expression at end of word are discarded",
			mexQuery:                `hello()`,
			expectedSolrQuery:       termMatcher("hello"),
			expectedQueryWasCleaned: true,
		},
		{
			name: "If there are empty bracket expressions both before and after a non-bracket string, " +
				"the outer brackets are matched and the inner are discarded",
			mexQuery:                `()hello there()`,
			expectedSolrQuery:       fmt.Sprintf("(%s AND %s)", termMatcher("hello"), termMatcher("there")),
			expectedQueryWasCleaned: true,
		},
		{
			name: "Query containing only an empty bracket expression is discarded when combined with" +
				" other space-separated terms",
			mexQuery:                `() hello`,
			expectedSolrQuery:       termMatcher("hello"),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "Query containing only a bracket expression with spaces maps to the empty query",
			mexQuery:                `(     )`,
			expectedSolrQuery:       ``,
			expectedQueryWasCleaned: true,
		},
		{
			name: "Spaces-only bracket expression is discarded when" +
				" combined with other space-separated terms",
			mexQuery:                `(     ) hello`,
			expectedSolrQuery:       termMatcher("hello"),
			expectedQueryWasCleaned: true,
		}, {
			name: "Spaces-only bracket expression inside a word is discarded but" +
				" separates terms",
			mexQuery:                `hel(  )lo`,
			expectedSolrQuery:       fmt.Sprintf("%s AND %s", termMatcher("hel"), termMatcher("lo")),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "Spaces-only  bracket expression at front or end of word or inside a word are discarded",
			mexQuery:                `(  )hello`,
			expectedSolrQuery:       termMatcher("hello"),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "Spaces-only  bracket expression at end of word are discarded",
			mexQuery:                `hello(  )`,
			expectedSolrQuery:       termMatcher("hello"),
			expectedQueryWasCleaned: true,
		},
		{
			name: "If there are Spaces-only bracket expressions both before and after a non-bracket string, " +
				"the outer brackets are matched and the inner are discarded",
			mexQuery:                `(  )hello there(  )`,
			expectedSolrQuery:       fmt.Sprintf("(%s AND %s)", termMatcher("hello"), termMatcher("there")),
			expectedQueryWasCleaned: true,
		},
		{
			name:              "Brackets surrounding single term are left in place",
			mexQuery:          `(hello)`,
			expectedSolrQuery: fmt.Sprintf("(%s)", termMatcher("hello")),
		},
		{
			name:              "Bracketed terms separated by space are ANDed together",
			mexQuery:          `(hello) (there)`,
			expectedSolrQuery: fmt.Sprintf("(%s) AND (%s)", termMatcher("hello"), termMatcher("there")),
		},
		{
			name:              "Bracketed terms not separated by space are ANDed together",
			mexQuery:          `(hello)(there)`,
			expectedSolrQuery: fmt.Sprintf("(%s) AND (%s)", termMatcher("hello"), termMatcher("there")),
		},
		{
			name:              "Bracketed terms and normal terms separated by space are ANDed together",
			mexQuery:          `(hello) there (again)`,
			expectedSolrQuery: fmt.Sprintf("(%s) AND %s AND (%s)", termMatcher("hello"), termMatcher("there"), termMatcher("again")),
		},
		{
			name:              "Bracketed terms and normal terms not separated by space are ANDed together",
			mexQuery:          `(hello)there(again)`,
			expectedSolrQuery: fmt.Sprintf("(%s) AND %s AND (%s)", termMatcher("hello"), termMatcher("there"), termMatcher("again")),
		},
		{
			name:              "Brackets surrounding multiple terms are left in place",
			mexQuery:          `(hello there)`,
			expectedSolrQuery: fmt.Sprintf("(%s AND %s)", termMatcher("hello"), termMatcher("there")),
		},
		{
			name:              "Nested brackets are left in place",
			mexQuery:          `(hello (there again))`,
			expectedSolrQuery: fmt.Sprintf("(%s AND (%s AND %s))", termMatcher("hello"), termMatcher("there"), termMatcher("again")),
		},
		{
			name:              "Redundant multiple brackets are left in place",
			mexQuery:          `((hello))`,
			expectedSolrQuery: fmt.Sprintf("((%s))", termMatcher("hello")),
		},
		{
			name:              "Escaped brackets are left in place",
			mexQuery:          `he\(lo the\)e`,
			expectedSolrQuery: fmt.Sprintf("%s AND %s", termMatcher(`he\(lo`), termMatcher(`the\)e`)),
		},
		{
			name:              "Brackets around multiple-term mexQuery in brackets are kept in output",
			mexQuery:          `("hello there" I) say`,
			expectedSolrQuery: fmt.Sprintf("(%s AND %s) AND %s", phraseMatcher(`hello there`), termMatcher("I"), termMatcher("say")),
		},
		{
			name:                    "Unmatched opening bracket is discarded",
			mexQuery:                `(hello again`,
			expectedSolrQuery:       fmt.Sprintf("%s AND %s", termMatcher("hello"), termMatcher("again")),
			expectedQueryWasCleaned: true,
		},
		{
			name: "Unmatched escaped opening bracket is treated as part of containing word and" +
				" does not cause parsing error",
			mexQuery:          `\(hello again`,
			expectedSolrQuery: fmt.Sprintf("%s AND %s", termMatcher(`\(hello`), termMatcher("again")),
		},
		{
			name:                    "Unmatched closing bracket at end of word is discarded",
			mexQuery:                `hello) again`,
			expectedSolrQuery:       fmt.Sprintf("%s AND %s", termMatcher("hello"), termMatcher("again")),
			expectedQueryWasCleaned: true,
		},
		{
			name:              "Unmatched escaped closing bracket is treated as part of containing word",
			mexQuery:          `hello\) again`,
			expectedSolrQuery: fmt.Sprintf("%s AND %s", termMatcher(`hello\)`), termMatcher("again")),
		},
		{
			name:                    "Unmatched closing bracket at end of word is discarded",
			mexQuery:                `hello) again`,
			expectedSolrQuery:       fmt.Sprintf("%s AND %s", termMatcher("hello"), termMatcher("again")),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "Separate and unmatched opening bracket is discarded",
			mexQuery:                `hello  ( again`,
			expectedSolrQuery:       fmt.Sprintf("%s AND %s", termMatcher("hello"), termMatcher("again")),
			expectedQueryWasCleaned: true,
		},
		{
			name:                    "Separate and unmatched closing bracket is discarded",
			mexQuery:                `hello there ) again`,
			expectedSolrQuery:       fmt.Sprintf("%s AND %s AND %s", termMatcher("hello"), termMatcher("there"), termMatcher("again")),
			expectedQueryWasCleaned: true,
		},
		{
			name:              "Separate and unmatched escaped opening bracket does not cause parsing error",
			mexQuery:          `hello there \( again`,
			expectedSolrQuery: fmt.Sprintf("%s AND %s AND %s AND %s", termMatcher("hello"), termMatcher("there"), termMatcher(`\(`), termMatcher("again")),
		},
		{
			name:              "Separate and unmatched escaped closing bracket does not cause parsing error",
			mexQuery:          `hello there \) again`,
			expectedSolrQuery: fmt.Sprintf("%s AND %s AND %s AND %s", termMatcher("hello"), termMatcher("there"), termMatcher(`\)`), termMatcher("again")),
		},
	}

	for _, tt := range tests {
		sortMarchingConfig(standardMatchingOpsConfig)
		listener, _ := NewListener(standardMatchingOpsConfig)
		parser := NewMexParser(listener)
		t.Run(tt.name, func(t *testing.T) {
			gotParseResult, err := parser.ConvertToSolrQuery(tt.mexQuery)
			checkQueryGenerationOutput(t, tt, gotParseResult, err)
		})
	}
}

func Test_MatchingOp_ParseSearchQuery_Escaping_Solr_Symbols(t *testing.T) {
	tests := []TestDef{
		{
			name:              "Single Solr control characters which are not MEx operators inside words are escaped (hyphen is missing as it is written out)",
			mexQuery:          `hello bc!def{g}g[i]j^kl~m?n:o/p there`,
			expectedSolrQuery: fmt.Sprintf("%s AND %s AND %s", termMatcher("hello"), termMatcher(`bc\!def\{g\}g\[i\]j\^kl\~m\?n\:o\/p`), termMatcher("there")),
		},
		{
			name:              "Escaped single Solr control characters which are not MEx operators inside terms are left unchanged",
			mexQuery:          `hello a\&b\-c\!def\{g\}g\[i\]j\^kl\~m\?n\:o\/p there`,
			expectedSolrQuery: fmt.Sprintf("%s AND %s AND %s", termMatcher("hello"), termMatcher(`a\&b\-c\!def\{g\}g\[i\]j\^kl\~m\?n\:o\/p`), termMatcher("there")),
		},
		{
			name: "Single Solr control characters (other than double quotes) inside quoted phrases (" +
				"where allowed MEx) are left unchanged",
			mexQuery:          `"hello a+b-c!def{g}g[i]j^kl~m?n:o/p there"`,
			expectedSolrQuery: phraseMatcher(`hello a+b-c!def{g}g[i]j^kl~m?n:o/p there`),
		},
		{
			name:              "A wildcard character (*) in a term is not escaped",
			mexQuery:          `hello*there`,
			expectedSolrQuery: termMatcher("hello*there"),
		},
		{
			name:              "An already escaped wildcard character (*) is not escaped further",
			mexQuery:          `hello\*there`,
			expectedSolrQuery: termMatcher(`hello\*there`),
		},
		{
			name:              "Solr AND operator is escaped",
			mexQuery:          `first&&second`,
			expectedSolrQuery: termMatcher(`first\&\&second`),
		},
		{
			name:              "The words 'AND', 'OR' and 'NOT' are lower-cased in queries (to avoid interpretation as special Solr terms)",
			mexQuery:          `OR first AND second OR third NOT fourth NOT`,
			expectedSolrQuery: fmt.Sprintf("%s AND %s AND %s AND %s AND %s AND %s AND %s AND %s AND %s", termMatcher("or"), termMatcher("first"), termMatcher("and"), termMatcher("second"), termMatcher("or"), termMatcher("third"), termMatcher("not"), termMatcher("fourth"), termMatcher("not")),
		},
		{
			name: "The strings 'AND', " +
				"'OR' and 'NOT' are not lower-cased if they occur as parts of other words",
			mexQuery:          `ANDing stANDs fOR NOThing nOR`,
			expectedSolrQuery: fmt.Sprintf("%s AND %s AND %s AND %s AND %s", termMatcher("ANDing"), termMatcher("stANDs"), termMatcher("fOR"), termMatcher("NOThing"), termMatcher("nOR")),
		},
		{
			name:              "Sanitization also works for multiline strings",
			mexQuery:          "first\n second OR thir?d",
			expectedSolrQuery: fmt.Sprintf("%s AND %s AND %s AND %s", termMatcher("first"), termMatcher("second"), termMatcher("or"), termMatcher(`thir\?d`)),
		},
	}

	for _, tt := range tests {
		sortMarchingConfig(standardMatchingOpsConfig)
		listener, _ := NewListener(standardMatchingOpsConfig)
		parser := NewMexParser(listener)
		t.Run(tt.name, func(t *testing.T) {
			gotParseResult, err := parser.ConvertToSolrQuery(tt.mexQuery)
			checkQueryGenerationOutput(t, tt, gotParseResult, err)
		})
	}
}

func Test_MatchingOp_ParseSearchQuery_FuzzySearch_Normal(t *testing.T) {
	tests := []TestDef{
		{
			name:                "Edit distance of 0 does not affect the Solr Query string",
			mexQuery:            `hello`,
			expectedSolrQuery:   "(X:hello)^3",
			maxTermEditDistance: 0,
			useNgramFuzz:        false,
		},
		{
			name:                "No distance is added to words of length equal to or less than the lower cut-off in the Solr Query string",
			mexQuery:            `tens`,
			expectedSolrQuery:   "(X:tens)^3",
			maxTermEditDistance: 2,
			useNgramFuzz:        false,
		},
		{
			name:                "Edit distance = 1 is added to a single term just above the lower cut-off in the Solr Query string",
			mexQuery:            `joins`,
			expectedSolrQuery:   "(X:joins~1)^3",
			maxTermEditDistance: 2,
			useNgramFuzz:        false,
		},
		{
			name:                "Full edit distance is added to a single term that is as long as the upper cut-off in the Solr Query string",
			mexQuery:            `eventually`,
			expectedSolrQuery:   "(X:eventually~2)^3",
			maxTermEditDistance: 2,
			useNgramFuzz:        false,
		}, {
			name:                "Edit distance > 0 is added to every term individually in the Solr Query string",
			mexQuery:            `hello there again`,
			expectedSolrQuery:   `(X:hello~1)^3 AND (X:there~1)^3 AND (X:again~1)^3`,
			maxTermEditDistance: 1,
			useNgramFuzz:        false,
		},
		{
			name:                "Edit distance is added (for long enough words) to every term individually also in cases of more complicated boolean logic",
			mexQuery:            `(hello there ) | (and -again)`,
			expectedSolrQuery:   `((X:hello~1)^3 AND (X:there~1)^3) OR ((X:and)^3 AND (NOT (X:again~1)^3))`,
			maxTermEditDistance: 1,
			useNgramFuzz:        false,
		},
		{
			name:                "No edit distance is added to phrases when so configured",
			mexQuery:            `"hello there"`,
			expectedSolrQuery:   `(Y:"hello there")^2`,
			maxTermEditDistance: 2,
			useNgramFuzz:        false,
		},
		{
			name:                "No edit distance is added to phrases even when terms follow the phrase",
			mexQuery:            `"hello there" and again`,
			expectedSolrQuery:   `(Y:"hello there")^2 AND (X:and)^3 AND (X:again~1)^3`,
			maxTermEditDistance: 1,
			useNgramFuzz:        false,
		},
		{
			name:                "No edit distance is added to phrases even when terms precede the phrase",
			mexQuery:            `hello there "and again"`,
			expectedSolrQuery:   `(X:hello~1)^3 AND (X:there~1)^3 AND (Y:"and again")^2`,
			maxTermEditDistance: 1,
			useNgramFuzz:        false,
		},
	}

	for _, tt := range tests {
		variableMatchingOpsConfig := MatchingOpsConfig{
			Term: []MatchingFieldConfig{
				{
					FieldName:       "X",
					BoostFactor:     "3",
					MaxEditDistance: tt.maxTermEditDistance,
				},
			},
			Phrase: []MatchingFieldConfig{
				{
					FieldName:       "Y",
					BoostFactor:     "2",
					MaxEditDistance: 0,
				},
			},
		}
		sortMarchingConfig(variableMatchingOpsConfig)
		listener, _ := NewListener(variableMatchingOpsConfig)
		parser := NewMexParser(listener)
		t.Run(tt.name, func(t *testing.T) {
			gotParseResult, err := parser.ConvertToSolrQuery(tt.mexQuery)
			checkQueryGenerationOutput(t, tt, gotParseResult, err)
		})
	}
}

func Test_MatchingOp_ParseSearchQuery_FuzzySearch_Ngram(t *testing.T) {
	tests := []TestDef{
		{
			name:                "If N-gram fuzziness is on, the edit distance is not ignored - simple case",
			mexQuery:            `hello`,
			expectedSolrQuery:   `(X:hello~1)^3`,
			maxTermEditDistance: 1,
			useNgramFuzz:        true,
		},
		{
			name:                "If N-gram fuzziness is on, the edit distance is not ignored - complex case",
			mexQuery:            `hello there "and again"`,
			expectedSolrQuery:   `(X:hello~1)^3 AND (X:there~1)^3 AND (Y:"and again")^2`,
			maxTermEditDistance: 1,
			useNgramFuzz:        true,
		},
	}

	for _, tt := range tests {
		variableMatchingOpsConfig := MatchingOpsConfig{
			Term: []MatchingFieldConfig{
				{
					FieldName:       "X",
					BoostFactor:     "3",
					MaxEditDistance: tt.maxTermEditDistance,
				},
			},
			Phrase: []MatchingFieldConfig{
				{
					FieldName:       "Y",
					BoostFactor:     "2",
					MaxEditDistance: 0,
				},
			},
		}
		sortMarchingConfig(variableMatchingOpsConfig)
		listener, _ := NewListener(variableMatchingOpsConfig)
		parser := NewMexParser(listener)
		t.Run(tt.name, func(t *testing.T) {
			gotParseResult, err := parser.ConvertToSolrQuery(tt.mexQuery)
			checkQueryGenerationOutput(t, tt, gotParseResult, err)
		})
	}
}

func Test_RemoveInternalHyphens(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "A string with no hyphens is returned unchanged",
			input: "something",
			want:  "something",
		},
		{
			name:  "A string with only a prefixed hyphen is returned unchanged",
			input: "-something",
			want:  "-something",
		},
		{
			name:  "A string with only a postfixed hyphen is returned unchanged",
			input: "something-",
			want:  "something-",
		},
		{
			name:  "In string with internal hyphens, the word is surrounded by brackets and the hyphens are replaced by the MEx AND-operator",
			input: "some-thing-else-too",
			want:  "(some + thing + else + too)",
		},
		{
			name:  "A phrase with an internal hyphen is left untouched",
			input: "\"some-thing-else\"",
			want:  "\"some-thing-else\"",
		},
		{
			name:  "A multi-word phrase with an internal hyphen is also converted is left untouched",
			input: "\"when some-thing-else appears\"",
			want:  "\"when some-thing-else appears\"",
		},
		{
			name:  "Multiple consecutive hyphens are also replaced",
			input: "some--thing---else",
			want:  "(some + thing + else)",
		},
		{
			name:  "Pre- and postfixed hyphen are left even if there are also internal hyphens that are replaced",
			input: "-some-thing-else-",
			want:  "(-some + thing + else-)",
		},
		{
			name:  "Does not substitute is the character before or after the hyphen is not alphanumeric",
			input: `some\-thing-.else`,
			want:  `some\-thing-.else`,
		},
		{
			name:  "Combinations of hyphenated phrases and words work",
			input: "\"some-thing-else\" plus this-other-thing",
			want:  "\"some-thing-else\" plus (this + other + thing)",
		},
		{
			name:  "Does the substitution in the hyphenated words without affects terms around it",
			input: "first some-thing-else middle this-other-thing last",
			want:  "first (some + thing + else) middle (this + other + thing) last",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeInternalHyphens(tt.input); got != tt.want {
				t.Errorf("RemoveInternalHyphens() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}
