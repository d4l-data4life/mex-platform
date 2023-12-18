// Code generated from MexQueryGrammar.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // MexQueryGrammar

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type MexQueryGrammarParser struct {
	*antlr.BaseParser
}

var mexquerygrammarParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func mexquerygrammarParserInit() {
	staticData := &mexquerygrammarParserStaticData
	staticData.literalNames = []string{
		"", "", "", "", "'+'", "'|'", "'-'", "'('", "')'", "'\"'",
	}
	staticData.symbolicNames = []string{
		"", "QUOTED_TERM", "TERM", "FREE_DASH", "AND", "OR", "NOT", "LPAR",
		"RPAR", "QUOTE", "WS", "LEFTOVER",
	}
	staticData.ruleNames = []string{
		"query", "statement", "or_expr", "and_expr", "operand_expr", "unary_expr",
		"dangling_op", "leftover",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 11, 77, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 1, 0, 1, 0, 1, 0, 5, 0, 20, 8, 0,
		10, 0, 12, 0, 23, 9, 0, 1, 0, 1, 0, 1, 1, 1, 1, 5, 1, 29, 8, 1, 10, 1,
		12, 1, 32, 9, 1, 1, 2, 1, 2, 1, 2, 5, 2, 37, 8, 2, 10, 2, 12, 2, 40, 9,
		2, 1, 3, 1, 3, 1, 3, 5, 3, 45, 8, 3, 10, 3, 12, 3, 48, 9, 3, 1, 4, 5, 4,
		51, 8, 4, 10, 4, 12, 4, 54, 9, 4, 1, 4, 1, 4, 5, 4, 58, 8, 4, 10, 4, 12,
		4, 61, 9, 4, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 3, 5, 71,
		8, 5, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 3, 30, 52, 59, 0, 8, 0, 2, 4, 6, 8,
		10, 12, 14, 0, 1, 1, 0, 3, 9, 79, 0, 21, 1, 0, 0, 0, 2, 26, 1, 0, 0, 0,
		4, 33, 1, 0, 0, 0, 6, 41, 1, 0, 0, 0, 8, 52, 1, 0, 0, 0, 10, 70, 1, 0,
		0, 0, 12, 72, 1, 0, 0, 0, 14, 74, 1, 0, 0, 0, 16, 20, 3, 2, 1, 0, 17, 20,
		3, 12, 6, 0, 18, 20, 3, 14, 7, 0, 19, 16, 1, 0, 0, 0, 19, 17, 1, 0, 0,
		0, 19, 18, 1, 0, 0, 0, 20, 23, 1, 0, 0, 0, 21, 19, 1, 0, 0, 0, 21, 22,
		1, 0, 0, 0, 22, 24, 1, 0, 0, 0, 23, 21, 1, 0, 0, 0, 24, 25, 5, 0, 0, 1,
		25, 1, 1, 0, 0, 0, 26, 30, 3, 4, 2, 0, 27, 29, 3, 4, 2, 0, 28, 27, 1, 0,
		0, 0, 29, 32, 1, 0, 0, 0, 30, 31, 1, 0, 0, 0, 30, 28, 1, 0, 0, 0, 31, 3,
		1, 0, 0, 0, 32, 30, 1, 0, 0, 0, 33, 38, 3, 6, 3, 0, 34, 35, 5, 5, 0, 0,
		35, 37, 3, 6, 3, 0, 36, 34, 1, 0, 0, 0, 37, 40, 1, 0, 0, 0, 38, 36, 1,
		0, 0, 0, 38, 39, 1, 0, 0, 0, 39, 5, 1, 0, 0, 0, 40, 38, 1, 0, 0, 0, 41,
		46, 3, 8, 4, 0, 42, 43, 5, 4, 0, 0, 43, 45, 3, 8, 4, 0, 44, 42, 1, 0, 0,
		0, 45, 48, 1, 0, 0, 0, 46, 44, 1, 0, 0, 0, 46, 47, 1, 0, 0, 0, 47, 7, 1,
		0, 0, 0, 48, 46, 1, 0, 0, 0, 49, 51, 3, 12, 6, 0, 50, 49, 1, 0, 0, 0, 51,
		54, 1, 0, 0, 0, 52, 53, 1, 0, 0, 0, 52, 50, 1, 0, 0, 0, 53, 55, 1, 0, 0,
		0, 54, 52, 1, 0, 0, 0, 55, 59, 3, 10, 5, 0, 56, 58, 3, 12, 6, 0, 57, 56,
		1, 0, 0, 0, 58, 61, 1, 0, 0, 0, 59, 60, 1, 0, 0, 0, 59, 57, 1, 0, 0, 0,
		60, 9, 1, 0, 0, 0, 61, 59, 1, 0, 0, 0, 62, 63, 5, 6, 0, 0, 63, 71, 3, 8,
		4, 0, 64, 65, 5, 7, 0, 0, 65, 66, 3, 2, 1, 0, 66, 67, 5, 8, 0, 0, 67, 71,
		1, 0, 0, 0, 68, 71, 5, 2, 0, 0, 69, 71, 5, 1, 0, 0, 70, 62, 1, 0, 0, 0,
		70, 64, 1, 0, 0, 0, 70, 68, 1, 0, 0, 0, 70, 69, 1, 0, 0, 0, 71, 11, 1,
		0, 0, 0, 72, 73, 7, 0, 0, 0, 73, 13, 1, 0, 0, 0, 74, 75, 5, 11, 0, 0, 75,
		15, 1, 0, 0, 0, 8, 19, 21, 30, 38, 46, 52, 59, 70,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// MexQueryGrammarParserInit initializes any static state used to implement MexQueryGrammarParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewMexQueryGrammarParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func MexQueryGrammarParserInit() {
	staticData := &mexquerygrammarParserStaticData
	staticData.once.Do(mexquerygrammarParserInit)
}

// NewMexQueryGrammarParser produces a new parser instance for the optional input antlr.TokenStream.
func NewMexQueryGrammarParser(input antlr.TokenStream) *MexQueryGrammarParser {
	MexQueryGrammarParserInit()
	this := new(MexQueryGrammarParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &mexquerygrammarParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	this.RuleNames = staticData.ruleNames
	this.LiteralNames = staticData.literalNames
	this.SymbolicNames = staticData.symbolicNames
	this.GrammarFileName = "MexQueryGrammar.g4"

	return this
}

// MexQueryGrammarParser tokens.
const (
	MexQueryGrammarParserEOF         = antlr.TokenEOF
	MexQueryGrammarParserQUOTED_TERM = 1
	MexQueryGrammarParserTERM        = 2
	MexQueryGrammarParserFREE_DASH   = 3
	MexQueryGrammarParserAND         = 4
	MexQueryGrammarParserOR          = 5
	MexQueryGrammarParserNOT         = 6
	MexQueryGrammarParserLPAR        = 7
	MexQueryGrammarParserRPAR        = 8
	MexQueryGrammarParserQUOTE       = 9
	MexQueryGrammarParserWS          = 10
	MexQueryGrammarParserLEFTOVER    = 11
)

// MexQueryGrammarParser rules.
const (
	MexQueryGrammarParserRULE_query        = 0
	MexQueryGrammarParserRULE_statement    = 1
	MexQueryGrammarParserRULE_or_expr      = 2
	MexQueryGrammarParserRULE_and_expr     = 3
	MexQueryGrammarParserRULE_operand_expr = 4
	MexQueryGrammarParserRULE_unary_expr   = 5
	MexQueryGrammarParserRULE_dangling_op  = 6
	MexQueryGrammarParserRULE_leftover     = 7
)

// IQueryContext is an interface to support dynamic dispatch.
type IQueryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsQueryContext differentiates from other interfaces.
	IsQueryContext()
}

type QueryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQueryContext() *QueryContext {
	var p = new(QueryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MexQueryGrammarParserRULE_query
	return p
}

func (*QueryContext) IsQueryContext() {}

func NewQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QueryContext {
	var p = new(QueryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MexQueryGrammarParserRULE_query

	return p
}

func (s *QueryContext) GetParser() antlr.Parser { return s.parser }

func (s *QueryContext) EOF() antlr.TerminalNode {
	return s.GetToken(MexQueryGrammarParserEOF, 0)
}

func (s *QueryContext) AllStatement() []IStatementContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IStatementContext); ok {
			len++
		}
	}

	tst := make([]IStatementContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IStatementContext); ok {
			tst[i] = t.(IStatementContext)
			i++
		}
	}

	return tst
}

func (s *QueryContext) Statement(i int) IStatementContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStatementContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStatementContext)
}

func (s *QueryContext) AllDangling_op() []IDangling_opContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IDangling_opContext); ok {
			len++
		}
	}

	tst := make([]IDangling_opContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IDangling_opContext); ok {
			tst[i] = t.(IDangling_opContext)
			i++
		}
	}

	return tst
}

func (s *QueryContext) Dangling_op(i int) IDangling_opContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDangling_opContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDangling_opContext)
}

func (s *QueryContext) AllLeftover() []ILeftoverContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ILeftoverContext); ok {
			len++
		}
	}

	tst := make([]ILeftoverContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ILeftoverContext); ok {
			tst[i] = t.(ILeftoverContext)
			i++
		}
	}

	return tst
}

func (s *QueryContext) Leftover(i int) ILeftoverContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILeftoverContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILeftoverContext)
}

func (s *QueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QueryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *QueryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MexQueryGrammarListener); ok {
		listenerT.EnterQuery(s)
	}
}

func (s *QueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MexQueryGrammarListener); ok {
		listenerT.ExitQuery(s)
	}
}

func (p *MexQueryGrammarParser) Query() (localctx IQueryContext) {
	this := p
	_ = this

	localctx = NewQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, MexQueryGrammarParserRULE_query)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(21)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<MexQueryGrammarParserQUOTED_TERM)|(1<<MexQueryGrammarParserTERM)|(1<<MexQueryGrammarParserFREE_DASH)|(1<<MexQueryGrammarParserAND)|(1<<MexQueryGrammarParserOR)|(1<<MexQueryGrammarParserNOT)|(1<<MexQueryGrammarParserLPAR)|(1<<MexQueryGrammarParserRPAR)|(1<<MexQueryGrammarParserQUOTE)|(1<<MexQueryGrammarParserLEFTOVER))) != 0 {
		p.SetState(19)
		p.GetErrorHandler().Sync(p)
		switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext()) {
		case 1:
			{
				p.SetState(16)
				p.Statement()
			}

		case 2:
			{
				p.SetState(17)
				p.Dangling_op()
			}

		case 3:
			{
				p.SetState(18)
				p.Leftover()
			}

		}

		p.SetState(23)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(24)
		p.Match(MexQueryGrammarParserEOF)
	}

	return localctx
}

// IStatementContext is an interface to support dynamic dispatch.
type IStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStatementContext differentiates from other interfaces.
	IsStatementContext()
}

type StatementContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStatementContext() *StatementContext {
	var p = new(StatementContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MexQueryGrammarParserRULE_statement
	return p
}

func (*StatementContext) IsStatementContext() {}

func NewStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatementContext {
	var p = new(StatementContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MexQueryGrammarParserRULE_statement

	return p
}

func (s *StatementContext) GetParser() antlr.Parser { return s.parser }

func (s *StatementContext) AllOr_expr() []IOr_exprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IOr_exprContext); ok {
			len++
		}
	}

	tst := make([]IOr_exprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IOr_exprContext); ok {
			tst[i] = t.(IOr_exprContext)
			i++
		}
	}

	return tst
}

func (s *StatementContext) Or_expr(i int) IOr_exprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOr_exprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOr_exprContext)
}

func (s *StatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MexQueryGrammarListener); ok {
		listenerT.EnterStatement(s)
	}
}

func (s *StatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MexQueryGrammarListener); ok {
		listenerT.ExitStatement(s)
	}
}

func (p *MexQueryGrammarParser) Statement() (localctx IStatementContext) {
	this := p
	_ = this

	localctx = NewStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, MexQueryGrammarParserRULE_statement)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(26)
		p.Or_expr()
	}
	p.SetState(30)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext())

	for _alt != 1 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1+1 {
			{
				p.SetState(27)
				p.Or_expr()
			}

		}
		p.SetState(32)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext())
	}

	return localctx
}

// IOr_exprContext is an interface to support dynamic dispatch.
type IOr_exprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsOr_exprContext differentiates from other interfaces.
	IsOr_exprContext()
}

type Or_exprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOr_exprContext() *Or_exprContext {
	var p = new(Or_exprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MexQueryGrammarParserRULE_or_expr
	return p
}

func (*Or_exprContext) IsOr_exprContext() {}

func NewOr_exprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Or_exprContext {
	var p = new(Or_exprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MexQueryGrammarParserRULE_or_expr

	return p
}

func (s *Or_exprContext) GetParser() antlr.Parser { return s.parser }

func (s *Or_exprContext) AllAnd_expr() []IAnd_exprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAnd_exprContext); ok {
			len++
		}
	}

	tst := make([]IAnd_exprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAnd_exprContext); ok {
			tst[i] = t.(IAnd_exprContext)
			i++
		}
	}

	return tst
}

func (s *Or_exprContext) And_expr(i int) IAnd_exprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAnd_exprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAnd_exprContext)
}

func (s *Or_exprContext) AllOR() []antlr.TerminalNode {
	return s.GetTokens(MexQueryGrammarParserOR)
}

func (s *Or_exprContext) OR(i int) antlr.TerminalNode {
	return s.GetToken(MexQueryGrammarParserOR, i)
}

func (s *Or_exprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Or_exprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Or_exprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MexQueryGrammarListener); ok {
		listenerT.EnterOr_expr(s)
	}
}

func (s *Or_exprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MexQueryGrammarListener); ok {
		listenerT.ExitOr_expr(s)
	}
}

func (p *MexQueryGrammarParser) Or_expr() (localctx IOr_exprContext) {
	this := p
	_ = this

	localctx = NewOr_exprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, MexQueryGrammarParserRULE_or_expr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(33)
		p.And_expr()
	}
	p.SetState(38)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(34)
				p.Match(MexQueryGrammarParserOR)
			}
			{
				p.SetState(35)
				p.And_expr()
			}

		}
		p.SetState(40)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext())
	}

	return localctx
}

// IAnd_exprContext is an interface to support dynamic dispatch.
type IAnd_exprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAnd_exprContext differentiates from other interfaces.
	IsAnd_exprContext()
}

type And_exprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAnd_exprContext() *And_exprContext {
	var p = new(And_exprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MexQueryGrammarParserRULE_and_expr
	return p
}

func (*And_exprContext) IsAnd_exprContext() {}

func NewAnd_exprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *And_exprContext {
	var p = new(And_exprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MexQueryGrammarParserRULE_and_expr

	return p
}

func (s *And_exprContext) GetParser() antlr.Parser { return s.parser }

func (s *And_exprContext) AllOperand_expr() []IOperand_exprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IOperand_exprContext); ok {
			len++
		}
	}

	tst := make([]IOperand_exprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IOperand_exprContext); ok {
			tst[i] = t.(IOperand_exprContext)
			i++
		}
	}

	return tst
}

func (s *And_exprContext) Operand_expr(i int) IOperand_exprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOperand_exprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOperand_exprContext)
}

func (s *And_exprContext) AllAND() []antlr.TerminalNode {
	return s.GetTokens(MexQueryGrammarParserAND)
}

func (s *And_exprContext) AND(i int) antlr.TerminalNode {
	return s.GetToken(MexQueryGrammarParserAND, i)
}

func (s *And_exprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *And_exprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *And_exprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MexQueryGrammarListener); ok {
		listenerT.EnterAnd_expr(s)
	}
}

func (s *And_exprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MexQueryGrammarListener); ok {
		listenerT.ExitAnd_expr(s)
	}
}

func (p *MexQueryGrammarParser) And_expr() (localctx IAnd_exprContext) {
	this := p
	_ = this

	localctx = NewAnd_exprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, MexQueryGrammarParserRULE_and_expr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(41)
		p.Operand_expr()
	}
	p.SetState(46)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 4, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(42)
				p.Match(MexQueryGrammarParserAND)
			}
			{
				p.SetState(43)
				p.Operand_expr()
			}

		}
		p.SetState(48)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 4, p.GetParserRuleContext())
	}

	return localctx
}

// IOperand_exprContext is an interface to support dynamic dispatch.
type IOperand_exprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsOperand_exprContext differentiates from other interfaces.
	IsOperand_exprContext()
}

type Operand_exprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOperand_exprContext() *Operand_exprContext {
	var p = new(Operand_exprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MexQueryGrammarParserRULE_operand_expr
	return p
}

func (*Operand_exprContext) IsOperand_exprContext() {}

func NewOperand_exprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Operand_exprContext {
	var p = new(Operand_exprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MexQueryGrammarParserRULE_operand_expr

	return p
}

func (s *Operand_exprContext) GetParser() antlr.Parser { return s.parser }

func (s *Operand_exprContext) Unary_expr() IUnary_exprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnary_exprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnary_exprContext)
}

func (s *Operand_exprContext) AllDangling_op() []IDangling_opContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IDangling_opContext); ok {
			len++
		}
	}

	tst := make([]IDangling_opContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IDangling_opContext); ok {
			tst[i] = t.(IDangling_opContext)
			i++
		}
	}

	return tst
}

func (s *Operand_exprContext) Dangling_op(i int) IDangling_opContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDangling_opContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDangling_opContext)
}

func (s *Operand_exprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Operand_exprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Operand_exprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MexQueryGrammarListener); ok {
		listenerT.EnterOperand_expr(s)
	}
}

func (s *Operand_exprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MexQueryGrammarListener); ok {
		listenerT.ExitOperand_expr(s)
	}
}

func (p *MexQueryGrammarParser) Operand_expr() (localctx IOperand_exprContext) {
	this := p
	_ = this

	localctx = NewOperand_exprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, MexQueryGrammarParserRULE_operand_expr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(52)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 5, p.GetParserRuleContext())

	for _alt != 1 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1+1 {
			{
				p.SetState(49)
				p.Dangling_op()
			}

		}
		p.SetState(54)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 5, p.GetParserRuleContext())
	}
	{
		p.SetState(55)
		p.Unary_expr()
	}
	p.SetState(59)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 6, p.GetParserRuleContext())

	for _alt != 1 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1+1 {
			{
				p.SetState(56)
				p.Dangling_op()
			}

		}
		p.SetState(61)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 6, p.GetParserRuleContext())
	}

	return localctx
}

// IUnary_exprContext is an interface to support dynamic dispatch.
type IUnary_exprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsUnary_exprContext differentiates from other interfaces.
	IsUnary_exprContext()
}

type Unary_exprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUnary_exprContext() *Unary_exprContext {
	var p = new(Unary_exprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MexQueryGrammarParserRULE_unary_expr
	return p
}

func (*Unary_exprContext) IsUnary_exprContext() {}

func NewUnary_exprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Unary_exprContext {
	var p = new(Unary_exprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MexQueryGrammarParserRULE_unary_expr

	return p
}

func (s *Unary_exprContext) GetParser() antlr.Parser { return s.parser }

func (s *Unary_exprContext) CopyFrom(ctx *Unary_exprContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *Unary_exprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Unary_exprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type PhraseContext struct {
	*Unary_exprContext
}

func NewPhraseContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PhraseContext {
	var p = new(PhraseContext)

	p.Unary_exprContext = NewEmptyUnary_exprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*Unary_exprContext))

	return p
}

func (s *PhraseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PhraseContext) QUOTED_TERM() antlr.TerminalNode {
	return s.GetToken(MexQueryGrammarParserQUOTED_TERM, 0)
}

func (s *PhraseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MexQueryGrammarListener); ok {
		listenerT.EnterPhrase(s)
	}
}

func (s *PhraseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MexQueryGrammarListener); ok {
		listenerT.ExitPhrase(s)
	}
}

type Negated_block_exprContext struct {
	*Unary_exprContext
}

func NewNegated_block_exprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *Negated_block_exprContext {
	var p = new(Negated_block_exprContext)

	p.Unary_exprContext = NewEmptyUnary_exprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*Unary_exprContext))

	return p
}

func (s *Negated_block_exprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Negated_block_exprContext) NOT() antlr.TerminalNode {
	return s.GetToken(MexQueryGrammarParserNOT, 0)
}

func (s *Negated_block_exprContext) Operand_expr() IOperand_exprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOperand_exprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOperand_exprContext)
}

func (s *Negated_block_exprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MexQueryGrammarListener); ok {
		listenerT.EnterNegated_block_expr(s)
	}
}

func (s *Negated_block_exprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MexQueryGrammarListener); ok {
		listenerT.ExitNegated_block_expr(s)
	}
}

type Bracketed_statementContext struct {
	*Unary_exprContext
}

func NewBracketed_statementContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *Bracketed_statementContext {
	var p = new(Bracketed_statementContext)

	p.Unary_exprContext = NewEmptyUnary_exprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*Unary_exprContext))

	return p
}

func (s *Bracketed_statementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Bracketed_statementContext) LPAR() antlr.TerminalNode {
	return s.GetToken(MexQueryGrammarParserLPAR, 0)
}

func (s *Bracketed_statementContext) Statement() IStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStatementContext)
}

func (s *Bracketed_statementContext) RPAR() antlr.TerminalNode {
	return s.GetToken(MexQueryGrammarParserRPAR, 0)
}

func (s *Bracketed_statementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MexQueryGrammarListener); ok {
		listenerT.EnterBracketed_statement(s)
	}
}

func (s *Bracketed_statementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MexQueryGrammarListener); ok {
		listenerT.ExitBracketed_statement(s)
	}
}

type TermContext struct {
	*Unary_exprContext
}

func NewTermContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TermContext {
	var p = new(TermContext)

	p.Unary_exprContext = NewEmptyUnary_exprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*Unary_exprContext))

	return p
}

func (s *TermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TermContext) TERM() antlr.TerminalNode {
	return s.GetToken(MexQueryGrammarParserTERM, 0)
}

func (s *TermContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MexQueryGrammarListener); ok {
		listenerT.EnterTerm(s)
	}
}

func (s *TermContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MexQueryGrammarListener); ok {
		listenerT.ExitTerm(s)
	}
}

func (p *MexQueryGrammarParser) Unary_expr() (localctx IUnary_exprContext) {
	this := p
	_ = this

	localctx = NewUnary_exprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, MexQueryGrammarParserRULE_unary_expr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(70)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case MexQueryGrammarParserNOT:
		localctx = NewNegated_block_exprContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(62)
			p.Match(MexQueryGrammarParserNOT)
		}
		{
			p.SetState(63)
			p.Operand_expr()
		}

	case MexQueryGrammarParserLPAR:
		localctx = NewBracketed_statementContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(64)
			p.Match(MexQueryGrammarParserLPAR)
		}
		{
			p.SetState(65)
			p.Statement()
		}
		{
			p.SetState(66)
			p.Match(MexQueryGrammarParserRPAR)
		}

	case MexQueryGrammarParserTERM:
		localctx = NewTermContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(68)
			p.Match(MexQueryGrammarParserTERM)
		}

	case MexQueryGrammarParserQUOTED_TERM:
		localctx = NewPhraseContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(69)
			p.Match(MexQueryGrammarParserQUOTED_TERM)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IDangling_opContext is an interface to support dynamic dispatch.
type IDangling_opContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDangling_opContext differentiates from other interfaces.
	IsDangling_opContext()
}

type Dangling_opContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDangling_opContext() *Dangling_opContext {
	var p = new(Dangling_opContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MexQueryGrammarParserRULE_dangling_op
	return p
}

func (*Dangling_opContext) IsDangling_opContext() {}

func NewDangling_opContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Dangling_opContext {
	var p = new(Dangling_opContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MexQueryGrammarParserRULE_dangling_op

	return p
}

func (s *Dangling_opContext) GetParser() antlr.Parser { return s.parser }

func (s *Dangling_opContext) AND() antlr.TerminalNode {
	return s.GetToken(MexQueryGrammarParserAND, 0)
}

func (s *Dangling_opContext) OR() antlr.TerminalNode {
	return s.GetToken(MexQueryGrammarParserOR, 0)
}

func (s *Dangling_opContext) LPAR() antlr.TerminalNode {
	return s.GetToken(MexQueryGrammarParserLPAR, 0)
}

func (s *Dangling_opContext) RPAR() antlr.TerminalNode {
	return s.GetToken(MexQueryGrammarParserRPAR, 0)
}

func (s *Dangling_opContext) QUOTE() antlr.TerminalNode {
	return s.GetToken(MexQueryGrammarParserQUOTE, 0)
}

func (s *Dangling_opContext) NOT() antlr.TerminalNode {
	return s.GetToken(MexQueryGrammarParserNOT, 0)
}

func (s *Dangling_opContext) FREE_DASH() antlr.TerminalNode {
	return s.GetToken(MexQueryGrammarParserFREE_DASH, 0)
}

func (s *Dangling_opContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Dangling_opContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Dangling_opContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MexQueryGrammarListener); ok {
		listenerT.EnterDangling_op(s)
	}
}

func (s *Dangling_opContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MexQueryGrammarListener); ok {
		listenerT.ExitDangling_op(s)
	}
}

func (p *MexQueryGrammarParser) Dangling_op() (localctx IDangling_opContext) {
	this := p
	_ = this

	localctx = NewDangling_opContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, MexQueryGrammarParserRULE_dangling_op)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(72)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<MexQueryGrammarParserFREE_DASH)|(1<<MexQueryGrammarParserAND)|(1<<MexQueryGrammarParserOR)|(1<<MexQueryGrammarParserNOT)|(1<<MexQueryGrammarParserLPAR)|(1<<MexQueryGrammarParserRPAR)|(1<<MexQueryGrammarParserQUOTE))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// ILeftoverContext is an interface to support dynamic dispatch.
type ILeftoverContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLeftoverContext differentiates from other interfaces.
	IsLeftoverContext()
}

type LeftoverContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLeftoverContext() *LeftoverContext {
	var p = new(LeftoverContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = MexQueryGrammarParserRULE_leftover
	return p
}

func (*LeftoverContext) IsLeftoverContext() {}

func NewLeftoverContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LeftoverContext {
	var p = new(LeftoverContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = MexQueryGrammarParserRULE_leftover

	return p
}

func (s *LeftoverContext) GetParser() antlr.Parser { return s.parser }

func (s *LeftoverContext) LEFTOVER() antlr.TerminalNode {
	return s.GetToken(MexQueryGrammarParserLEFTOVER, 0)
}

func (s *LeftoverContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LeftoverContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LeftoverContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MexQueryGrammarListener); ok {
		listenerT.EnterLeftover(s)
	}
}

func (s *LeftoverContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MexQueryGrammarListener); ok {
		listenerT.ExitLeftover(s)
	}
}

func (p *MexQueryGrammarParser) Leftover() (localctx ILeftoverContext) {
	this := p
	_ = this

	localctx = NewLeftoverContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, MexQueryGrammarParserRULE_leftover)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(74)
		p.Match(MexQueryGrammarParserLEFTOVER)
	}

	return localctx
}
