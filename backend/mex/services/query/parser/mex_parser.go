package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"github.com/d4l-data4life/mex/mex/shared/uuid"
)

const (
	ParserErrorType            = "PARSER_ERROR"
	QueryConstructionErrorType = "QUERY_CONSTRUCTION_ERROR"
	InvalidArgumentErrorType   = "INVALID_ARGUMENT_ERROR"
	EmptyQuery                 = ""
	GetAllQuery                = "*"
)

// TypedParserError is an error than can indicate its type
type TypedParserError interface {
	error
	GetType() string
	GetMessages() []string
}

//nolint:revive
type ParserError struct {
	msg     string
	errType string
	errMsgs []string
}

func (p ParserError) Error() string {
	return p.msg
}

func (p ParserError) GetType() string {
	return p.errType
}

func (p ParserError) GetMessages() []string {
	return p.errMsgs
}

func NewParserError(errType string, msg string, errMsgs []string) ParserError {
	return ParserError{
		msg:     msg,
		errType: errType,
		errMsgs: errMsgs,
	}
}

type MexErrorListener struct {
	*antlr.DefaultErrorListener
	errs []string
}

func NewMexErrorListener() *MexErrorListener {
	return new(MexErrorListener)
}

func (d *MexErrorListener) SyntaxError(_ antlr.Recognizer, offendingSymbol interface{}, _, _ int,
	msg string, _ antlr.RecognitionException) {
	d.errs = append(d.errs, fmt.Sprintf("Msg: %s - symbol: %s", msg, offendingSymbol))
}

// QueryConverter is the interface for an object that can build Solr queries from MEx queries
type QueryConverter interface {
	ConvertToSolrQuery(rawQuery string) (*QueryParseResult, TypedParserError)
}

type MexParser struct {
	listener *SolrQueryBuilderListener
}

func NewMexParser(listener *SolrQueryBuilderListener) *MexParser {
	return &MexParser{listener: listener}
}

// ConvertToSolrQuery turns a MEx search query string into a Solr query string
func (p *MexParser) ConvertToSolrQuery(rawQuery string) (*QueryParseResult, TypedParserError) {
	query := cleanInput(rawQuery)
	if query == EmptyQuery {
		query = GetAllQuery
	}

	errListener := NewMexErrorListener()
	parseTree := getMexQueryTree(query, errListener)
	walkMexQueryTreeWithListeners(parseTree, []antlr.ParseTreeListener{p.listener})

	// Check for errors
	if len(errListener.errs) > 0 {
		// Parsing error
		return nil, NewParserError(ParserErrorType, "parsing error(s) occurred", errListener.errs)
	} else if len(p.listener.errs) > 0 {
		msgs := make([]string, len(p.listener.errs))
		for i, s := range p.listener.errs {
			msgs[i] = s.Error()
		}
		return nil, NewParserError(QueryConstructionErrorType, "Error(s) occurred while constructing the query", msgs)
	}

	queryParseResult, err := p.listener.GetResult()
	if err != nil {
		return nil, NewParserError(QueryConstructionErrorType, "Error(s) occurred while constructing the query",
			[]string{err.Error()})
	}

	return queryParseResult, nil
}

// cleanInput rewrites the raw user query to prepare it for transformation to a Solr query.
func cleanInput(rawQuery string) string {
	trimmedQuery := strings.TrimSpace(rawQuery)
	query := removeInternalHyphens(trimmedQuery)
	return query
}

// MapParseTreeToObject turns a MEx search query into a generic object (to be serialized as JSON)
func MapParseTreeToObject(rawQuery string, showFull bool) (map[string]interface{}, TypedParserError) {
	query := strings.TrimSpace(rawQuery)
	if query == "" {
		return nil, nil
	}

	debugListener := NewDebugListener(showFull)
	errListener := NewMexErrorListener()
	parseTree := getMexQueryTree(rawQuery, errListener)
	walkMexQueryTreeWithListeners(parseTree, []antlr.ParseTreeListener{debugListener})

	// Check for errors
	if len(errListener.errs) > 0 {
		// Parsing error
		return nil, NewParserError(ParserErrorType, "parsing error(s) occurred", errListener.errs)
	} else if len(debugListener.errs) > 0 {
		msgs := make([]string, len(debugListener.errs))
		for i, s := range debugListener.errs {
			msgs[i] = s.Error()
		}
		return nil, NewParserError(QueryConstructionErrorType, "Error(s) occurred while constructing the query", msgs)
	}

	obj, err := debugListener.GetResult()
	if err != nil {
		return nil, NewParserError(QueryConstructionErrorType,
			"Error(s) occurred while constructing the parse tree object",
			[]string{err.Error()})
	}

	return obj, nil
}

func getMexQueryTree(query string, errListener antlr.ErrorListener) antlr.Tree {
	// Lex
	input := antlr.NewInputStream(query)
	lexer := NewMexQueryGrammarLexer(input)
	lexer.AddErrorListener(errListener)
	stream := antlr.NewCommonTokenStream(lexer, 0)

	// Parse
	p := NewMexQueryGrammarParser(stream)
	p.AddErrorListener(errListener)
	return p.Query()
}

func walkMexQueryTreeWithListeners(tree antlr.Tree, listeners []antlr.ParseTreeListener) {
	for _, l := range listeners {
		antlr.ParseTreeWalkerDefault.Walk(l, tree)
	}
}

/*
removeInternalHyphens removes any internal hyphen in a string, unless it is in a phrase.
The individual parts are joined the MEx AND-operator and surrounded with a bracket.

This helps prevent problematic interactions between hyphenated terms and the Solr fuzzy search operator.
*/
func removeInternalHyphens(s string) string {
	if !strings.Contains(s, "-") {
		return s
	}
	/*
		To replace internal hyphens without touching phrases, the code below does the following
		(1) Extract all phrases (substrings surrounded by double quotes), replacing them with UUIDs
		(2) Do the hyphen-replacement without worrying about phrases
		(3) Re-substitute the phrases for the UUIDs
	*/

	// (1) Split out and replace phrases
	cleanedString, phraseMap := pullOutPhrases(s)

	// (2) Remove internal hyphens
	// Check if we have any words with internal hyphens
	internalHyphenatedWordRegExp := regexp.MustCompile(`(?P<front>^|\s)(?P<str>(?:[^\t\n\f\r ]*[[:alnum:]]-+[[:alnum:]][^\t\n\f\r ]*)+)(?P<back>$|\s)`)
	match := internalHyphenatedWordRegExp.Match([]byte(s))
	if !match {
		// No words with internal hyphens
		return s
	}
	// Put a bracket around words with internal hyphens to ensure that it is treated as unit
	bracketedString := internalHyphenatedWordRegExp.ReplaceAllString(cleanedString, "$front($str)$back")
	// ... then split on hyphens and join with the AND-operator (with a regexp that matches
	// the bracketed terms created above)
	internalHyphenRegExp := regexp.MustCompile(`(?:^|\s)\((?:[^\t\n\f\r ]*[[:alnum:]]-+[[:alnum:]][^\t\n\f\r ]*)+\)(?:$|\s)`)
	replacedString := internalHyphenRegExp.ReplaceAllStringFunc(bracketedString, replaceFunc)

	// (3) Substitute back the phrases
	for uuidVal, phrase := range phraseMap {
		replacedString = strings.ReplaceAll(replacedString, uuidVal, phrase)
	}

	return replacedString
}

/*
pullOutPhrases replaces all phrases (string in double-quotes) in the string with uuids without hyphen, returning
the resulting string and a map from uuids to the phrase that it replaced.

For instance the string 'this "hyphe-nated: word' would return a string like `this a4b987f word` and the
map {"a4b987f": "\"hyphe-nated\""}
*/
func pullOutPhrases(s string) (string, map[string]string) {
	// Since double quotes must be matched pairwise, we do some minimal parsing
	inPhrase := false
	phraseMap := map[string]string{}
	var completeArr []rune
	var phraseArr []rune
	for _, runeVal := range s {
		if runeVal == '"' {
			if inPhrase {
				// Insert uuidVal (without hyphens) in complete string
				uuidVal := strings.ReplaceAll(uuid.MustNewV4(), "-", "")
				completeArr = append(completeArr, []rune(uuidVal)...)
				// Store phrase (including closing double quote) under this uuidVal
				phraseArr = append(phraseArr, runeVal)
				phraseMap[uuidVal] = string(phraseArr)
				// Prepare to store new phrase
				phraseArr = []rune{}
			} else {
				// Store the double quote
				phraseArr = append(phraseArr, runeVal)
			}
			// Flip in-phrase state
			inPhrase = !inPhrase
			// Go to next character (this character has still been dealt with)
			continue
		}
		// Accumulate character to the relevant array
		if inPhrase {
			phraseArr = append(phraseArr, runeVal)
		} else {
			completeArr = append(completeArr, runeVal)
		}
	}
	// If we have an unterminated phrase, add the characters back to the full string
	if len(phraseArr) > 0 {
		completeArr = append(completeArr, phraseArr...)
	}
	// Convert back to string
	cleanedString := string(completeArr)
	return cleanedString, phraseMap
}

// replaceFunc takes a single bracket term containing a hyphenated word and replaces the hyphenated word with the
// hyphenated parts joined with the MEx AND-operator.
func replaceFunc(s string) string {
	// The input here is word with internal hyphens, no double-quotes and either a space or a sentence end/start and a bracket at either end
	removeInternalHyphensRegExp := regexp.MustCompile(`(?P<left>[[:alnum:]])-+(?P<right>[[:alnum:]])`)
	return removeInternalHyphensRegExp.ReplaceAllString(s, "$left"+CleanAndSeparator+"$right")
}
