package parser

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/d4l-data4life/mex/mex/shared/solr"
)

const (
	AndSeparator = " AND "
	OrSeparator  = " OR "
	NotOperator  = "NOT "

	CleanAndSeparator       = " + "
	CleanOrSeparator        = " | "
	CleanNotOperator        = "-"
	CleanStatementSeparator = " "

	/*
		NonSupportedSolrSpecialChars contains all Solr special characters except '*'
		(which we want Solr to interpret as a wildcard)
	*/
	NonSupportedSolrSpecialChars = `"+\-!(){}\[\]^~?:/&|`
)

// Bracket a term
func bracket(s string) string {
	return fmt.Sprintf("(%s)", s)
}

/*
SanitizeTerm escapes or otherwise disables character and terms with special meaning in Solr.
*/
func SanitizeTerm(s string) string {
	/*
		Remove all escape characters - since this method should only be called when the parser has run,
		all escape characters have served their purpose (by guiding parsing).

		The one exception is the wildcard symbol ('*') which is a MEx symbol that still needs to be treated as
		escaped even after MEx parsing since we implement it directly with the Solr *-wildcard. That is,
		if * is escaped, that escape needs to be passed through to Solr; if it is not escaped, it should be kept unescaped
	*/
	reRemoveMexEscapes := regexp.MustCompile(`(\\)($|[^*])`)
	es := reRemoveMexEscapes.ReplaceAllString(s, `${2}`)
	// Lower-case AND, OR, and NOT (if occurring as whole words) to avoid them being interpreted as Boolean operators
	reBooleanWords := regexp.MustCompile(`(?:^|\s)(AND|OR|NOT)(?:$|\s)`)
	es = reBooleanWords.ReplaceAllStringFunc(es, strings.ToLower)
	// Escape all Solr special characters except '*' and '\'
	reSingle := regexp.MustCompile(fmt.Sprintf(`([%s])`, NonSupportedSolrSpecialChars))
	es = reSingle.ReplaceAllString(es, `\${1}`)
	return es
}

/*
addEditDistance appends the maximum edit distance per term for fuzzy search in Solr
*/
func addEditDistance(s string, maxDistance uint32) string {
	if maxDistance == 0 {
		return s
	}
	// We calculate the effective edit distance d as follows:
	// L <= Lmin --> d = 0
	// Lmin < L < Lmax --> linearly interpolated between 1 and d_max, rounding up
	// L >= Lmax --> d = d_max
	rawEffectiveDistance := float64(maxDistance) * float64(len(s)-solr.EditLowerCutoff) / float64(solr.EditUpperCutoff-solr.EditLowerCutoff)
	effectiveDistance := int(math.Min(float64(maxDistance), math.Ceil(rawEffectiveDistance)))
	if effectiveDistance > 0 {
		/*
			Fuzzy search in Solr is allowed by appending a tilde ('~') followed by a distance parameter
			specifying the maximum number of edits allowed
		*/
		s = fmt.Sprintf("%s~%d", s, effectiveDistance)
	}
	return s
}

// joinNonempty joins a list of strings with the given separator, dropping empty strings from the list
func joinNonempty(strList []string, separator string) string {
	filteredStrList := []string{}
	for _, s := range strList {
		if s != "" {
			filteredStrList = append(filteredStrList, s)
		}
	}
	if len(filteredStrList) == 0 {
		return ""
	}
	return strings.Join(filteredStrList, separator)
}

// subTreeSolrQuery represents the compiled query corresponding to a particular sub-tree of the parse tree.
type subTreeSolrQuery struct {
	solrQuery    string // Solr query corresponding to a given sub-tree (tree below a given node)
	cleanedQuery string // Query string meant to illustrate how the parser interpreted the query
	isComplex    bool   // Whether sub-tree query is complex i.e. involves multiple elements --> requires brackets
}

// queryStack provides a simple stack implementation
type queryStack struct {
	elem []subTreeSolrQuery
}

func (s *queryStack) isEmpty() bool {
	return len(s.elem) == 0
}

func (s *queryStack) push(e subTreeSolrQuery) {
	s.elem = append(s.elem, e)
}

func (s *queryStack) pop() (subTreeSolrQuery, error) {
	if s.isEmpty() {
		return subTreeSolrQuery{}, fmt.Errorf("queryStack is empty")
	}
	lastIndex := len(s.elem) - 1
	e := s.elem[lastIndex]
	s.elem = s.elem[:lastIndex]
	return e, nil
}

type MatchingFieldConfig struct {
	FieldName       string
	BoostFactor     string
	MaxEditDistance uint32
}

type MatchingOpsConfig struct {
	Term   []MatchingFieldConfig
	Phrase []MatchingFieldConfig
}

// SolrQueryBuilderListener is a listener used for building a Solr query while walking a MEx parse tree
type SolrQueryBuilderListener struct {
	*BaseMexQueryGrammarListener

	matchingOpsConfig MatchingOpsConfig
	queryStack        queryStack
	queryCleaned      bool
	errs              []error

	phrasesOnlyQuery bool
}

// NewListener creates a new SolrQueryBuilderListener struct - ONLY create new listeners with this method to ensure
// validation & correct initialization
func NewListener(matchingOpsConfig MatchingOpsConfig) (*SolrQueryBuilderListener, error) {
	for _, matchConfig := range matchingOpsConfig.Term {
		if matchConfig.FieldName == "" {
			return nil, fmt.Errorf("term matching operator: empty field name not allowed")
		}
		if matchConfig.MaxEditDistance > solr.MaxEditDistance {
			return nil, fmt.Errorf("term matching operator: max edit distance cannot be set to more that %d", solr.MaxEditDistance)
		}
	}
	for _, matchConfig := range matchingOpsConfig.Phrase {
		if matchConfig.MaxEditDistance > 0 {
			return nil, fmt.Errorf("phrase matching operator: max edit distance must be zero but was %d", matchConfig.MaxEditDistance)
		}
		if matchConfig.FieldName == "" {
			return nil, fmt.Errorf("term matching operator: empty field name not allowed")
		}
	}

	return &SolrQueryBuilderListener{
		matchingOpsConfig: matchingOpsConfig,
		// We assume the query contains only phrases and switch when we see a term
		phrasesOnlyQuery: true,
		queryCleaned:     false,
	}, nil
}

// stackEmpty returns true if the stack is empty
func (listener *SolrQueryBuilderListener) stackEmpty() bool {
	return listener.queryStack.isEmpty()
}

// multiPop returns multiple strings from stack (last popped first), bracketing complex terms if bracketComplex is true
func (listener *SolrQueryBuilderListener) multiPop(count int, bracketComplex bool) ([]string, []string, error) {
	solrTerms := make([]string, count)
	cleanedTerms := make([]string, count)
	// We fill the slice from the end to recover the original term order
	for i := count - 1; i >= 0; i-- {
		q, err := listener.queryStack.pop()
		if err != nil {
			return []string{}, []string{},
				fmt.Errorf("requested %d elements but stack empty after popping %d elements", count, count-1-i)
		}
		// We bracket complex sub-queries if required
		if bracketComplex && q.isComplex {
			solrTerms[i] = bracket(q.solrQuery)
			cleanedTerms[i] = bracket(q.cleanedQuery)
		} else {
			solrTerms[i] = q.solrQuery
			cleanedTerms[i] = q.cleanedQuery
		}
	}
	return solrTerms, cleanedTerms, nil
}

type QueryParseResult struct {
	SolrQuery        string
	CleanedQuery     string
	QueryWasCleaned  bool
	PhrasesOnlyQuery bool
}

/*
GetResult returns the top stack element checking that it is the only element left.
After successful conversion, the listener should contain the final Solr query as the only element
*/
func (listener *SolrQueryBuilderListener) GetResult() (*QueryParseResult, error) {
	var err error
	queryList, cleanedQueryList, err := listener.multiPop(1, false)
	if err != nil {
		return nil, fmt.Errorf("stack does not contain final query: %s", err.Error())
	} else if !listener.stackEmpty() {
		return nil, fmt.Errorf("listener stack contained multiple elements at end of parsing")
	}
	return &QueryParseResult{
		SolrQuery:        queryList[0],
		CleanedQuery:     cleanedQueryList[0],
		QueryWasCleaned:  listener.queryCleaned,
		PhrasesOnlyQuery: listener.phrasesOnlyQuery,
	}, nil
}

// handleMultiTermNode processes a node that may have multiple children (e.g. AND or OR nodes)
func (listener *SolrQueryBuilderListener) handleMultiTermNode(childNo int, termSeparator string, cleanTermSeparator string) {
	// If we only have a single term stack can be left as it is
	if childNo == 1 {
		return
	}

	solrTerms, cleanedTerms, err := listener.multiPop(childNo, true)
	if err != nil {
		listener.errs = append(listener.errs, err)
		return
	}
	res := subTreeSolrQuery{
		solrQuery:    joinNonempty(solrTerms, termSeparator),
		cleanedQuery: joinNonempty(cleanedTerms, cleanTermSeparator),
		isComplex:    true,
	}
	listener.queryStack.push(res)
}

// ExitQuery is called when production statement is exited - all non-empty terms are ANDed together.
func (listener *SolrQueryBuilderListener) ExitQuery(ctx *QueryContext) {
	listener.handleMultiTermNode(len(ctx.AllStatement()), AndSeparator, CleanStatementSeparator)
}

// ExitStatement is called when production statement is exited - all non-empty terms are ANDed together.
func (listener *SolrQueryBuilderListener) ExitStatement(ctx *StatementContext) {
	listener.handleMultiTermNode(len(ctx.AllOr_expr()), AndSeparator, CleanAndSeparator)
}

// ExitOr_expr is called when production or_expr is exited - all non-empty terms are ORed together.
func (listener *SolrQueryBuilderListener) ExitOr_expr(ctx *Or_exprContext) {
	listener.handleMultiTermNode(len(ctx.AllAnd_expr()), OrSeparator, CleanOrSeparator)
}

// ExitAnd_expr is called when production and_expr is exited - all non-empty terms are ANDed together.
func (listener *SolrQueryBuilderListener) ExitAnd_expr(ctx *And_exprContext) {
	listener.handleMultiTermNode(len(ctx.AllOperand_expr()), AndSeparator, CleanAndSeparator)
}

// ExitNegated_block_expr is called when production negated_block_expr is exited - the Solr NOT-operator is added.
func (listener *SolrQueryBuilderListener) ExitNegated_block_expr(_ *Negated_block_exprContext) {
	solrQuery, cleanedQuery, err := listener.multiPop(1, false)
	if err != nil {
		listener.errs = append(listener.errs, err)
		return
	}
	res := subTreeSolrQuery{
		solrQuery:    bracket(NotOperator + solrQuery[0]),
		cleanedQuery: bracket(CleanNotOperator + cleanedQuery[0]),
	}
	listener.queryStack.push(res)
}

// ExitTerm is called when production term is exited - term is pushed onto the stack.
func (listener *SolrQueryBuilderListener) ExitTerm(ctx *TermContext) {
	listener.phrasesOnlyQuery = false
	rawTerm, sanitizedTerm, err := listener.constructTermQuery(ctx.GetText())
	if err != nil {
		listener.errs = append(listener.errs, fmt.Errorf("could not create matching query for term: %s", err.Error()))
		return
	}
	res := subTreeSolrQuery{
		solrQuery:    sanitizedTerm,
		cleanedQuery: rawTerm,
	}
	listener.queryStack.push(res)
}

// ExitPhrase is called when production phrase is exited - phrase is pushed onto the stack.
func (listener *SolrQueryBuilderListener) ExitPhrase(ctx *PhraseContext) {
	rawTerm, sanitizedTerm, err := listener.constructPhraseQuery(ctx.GetText())
	if err != nil {
		listener.errs = append(listener.errs, fmt.Errorf("could not create matching query for phrase: %s", err.Error()))
		return
	}
	res := subTreeSolrQuery{
		solrQuery:    sanitizedTerm,
		cleanedQuery: rawTerm,
	}
	listener.queryStack.push(res)
}

// ExitBracketed_statement is called when production bracketed_query is exited - brackets are added to top stack entry.
func (listener *SolrQueryBuilderListener) ExitBracketed_statement(_ *Bracketed_statementContext) {
	solrQuery, cleanedQuery, err := listener.multiPop(1, false)
	if err != nil {
		listener.errs = append(listener.errs, fmt.Errorf("stack was empty inside bracket query node"))
		return
	}
	res := subTreeSolrQuery{
		solrQuery:    bracket(solrQuery[0]),
		cleanedQuery: bracket(cleanedQuery[0]),
	}
	listener.queryStack.push(res)
}

// ExitDangling_op is called when exiting the dangling_op production - empty node is pushed onto the stack.
func (listener *SolrQueryBuilderListener) ExitDangling_op(_ *Dangling_opContext) {
	// We ignore dangling operators and just flag that such an operator was encountered
	listener.queryCleaned = true
}

// ExitLeftover is called when production leftover is exited (should not happen) - causes an error.
func (listener *SolrQueryBuilderListener) ExitLeftover(_ *LeftoverContext) {
	// Backstop: If we encounter a left-over, we have not succeeded in parsing the statement --> add error
	err := fmt.Errorf("input not matching any token rule found")
	listener.errs = append(listener.errs, err)
}

// constructPhraseQuery constructs the sub-query for matching a given term
func (listener *SolrQueryBuilderListener) constructTermQuery(extractedText string) (string, string, error) {
	rawTerm := extractedText
	processedTerm := SanitizeTerm(rawTerm)
	matchingQuery, err := constructMatchQuery(processedTerm, listener.matchingOpsConfig.Term)
	if err != nil {
		return "", "", err
	}
	return rawTerm, matchingQuery, nil
}

// constructPhraseQuery constructs the sub-query for matching a given phrase
func (listener *SolrQueryBuilderListener) constructPhraseQuery(extractedText string) (string, string, error) {
	rawTerm := extractedText
	processedTerm := extractedText // No need to sanitize inside phrases
	matchingQuery, err := constructMatchQuery(processedTerm, listener.matchingOpsConfig.Phrase)
	if err != nil {
		return "", "", err
	}
	return rawTerm, matchingQuery, nil
}

/*
constructMatchQuery constructs a sub-query for matching a single term. The query takes the form or
an OR-conjunction of individual queries that can differ by the field to search, the boost factor,
and the edit distance for fuzzy search.
*/
func constructMatchQuery(processedTerm string, matchingFieldConfigs []MatchingFieldConfig) (string, error) {
	var collectedTerms []string
	for _, termMatchField := range matchingFieldConfigs {
		if termMatchField.FieldName == "" {
			return "", fmt.Errorf("matching operator config with empty field not allowed")
		}
		queryPart := addEditDistance(processedTerm, termMatchField.MaxEditDistance)
		queryPart = fmt.Sprintf("%s:%s", termMatchField.FieldName, queryPart)
		if termMatchField.BoostFactor != "" {
			queryPart = fmt.Sprintf("(%s)^%s", queryPart, termMatchField.BoostFactor)
		}
		collectedTerms = append(collectedTerms, queryPart)
	}
	result := strings.Join(collectedTerms, OrSeparator)
	if len(collectedTerms) > 1 {
		result = bracket(result)
	}
	return result, nil
}
