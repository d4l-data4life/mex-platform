// Code generated from MexQueryGrammar.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser

import (
	"fmt"
	"sync"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type MexQueryGrammarLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var mexquerygrammarlexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	channelNames           []string
	modeNames              []string
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func mexquerygrammarlexerLexerInit() {
	staticData := &mexquerygrammarlexerLexerStaticData
	staticData.channelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.modeNames = []string{
		"DEFAULT_MODE",
	}
	staticData.literalNames = []string{
		"", "", "", "", "'+'", "'|'", "'-'", "'('", "')'", "'\"'",
	}
	staticData.symbolicNames = []string{
		"", "QUOTED_TERM", "TERM", "FREE_DASH", "AND", "OR", "NOT", "LPAR",
		"RPAR", "QUOTE", "WS", "LEFTOVER",
	}
	staticData.ruleNames = []string{
		"QUOTED_TERM", "TERM", "VALID_TERM_MIDDLE_SYMBOL", "VALID_TERM_START_SYMBOL",
		"NORMAL_TERM_SYMBOL", "ESCAPED_CONTROL_SYMBOL", "FREE_DASH", "AND",
		"OR", "NOT", "LPAR", "RPAR", "QUOTE", "WS", "LEFTOVER",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 11, 83, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 1, 0,
		1, 0, 1, 0, 1, 0, 5, 0, 36, 8, 0, 10, 0, 12, 0, 39, 9, 0, 1, 0, 1, 0, 1,
		1, 1, 1, 5, 1, 45, 8, 1, 10, 1, 12, 1, 48, 9, 1, 1, 2, 1, 2, 3, 2, 52,
		8, 2, 1, 3, 1, 3, 3, 3, 56, 8, 3, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 6, 1,
		6, 1, 6, 1, 7, 1, 7, 1, 8, 1, 8, 1, 9, 1, 9, 1, 10, 1, 10, 1, 11, 1, 11,
		1, 12, 1, 12, 1, 13, 1, 13, 1, 13, 1, 13, 1, 14, 1, 14, 1, 37, 0, 15, 1,
		1, 3, 2, 5, 0, 7, 0, 9, 0, 11, 0, 13, 3, 15, 4, 17, 5, 19, 6, 21, 7, 23,
		8, 25, 9, 27, 10, 29, 11, 1, 0, 3, 8, 0, 9, 10, 12, 13, 32, 32, 34, 34,
		40, 41, 43, 43, 45, 45, 124, 124, 5, 0, 34, 34, 40, 41, 43, 43, 45, 45,
		124, 124, 3, 0, 9, 10, 12, 13, 32, 32, 83, 0, 1, 1, 0, 0, 0, 0, 3, 1, 0,
		0, 0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 0, 17, 1, 0, 0, 0, 0, 19, 1,
		0, 0, 0, 0, 21, 1, 0, 0, 0, 0, 23, 1, 0, 0, 0, 0, 25, 1, 0, 0, 0, 0, 27,
		1, 0, 0, 0, 0, 29, 1, 0, 0, 0, 1, 31, 1, 0, 0, 0, 3, 42, 1, 0, 0, 0, 5,
		51, 1, 0, 0, 0, 7, 55, 1, 0, 0, 0, 9, 57, 1, 0, 0, 0, 11, 59, 1, 0, 0,
		0, 13, 62, 1, 0, 0, 0, 15, 65, 1, 0, 0, 0, 17, 67, 1, 0, 0, 0, 19, 69,
		1, 0, 0, 0, 21, 71, 1, 0, 0, 0, 23, 73, 1, 0, 0, 0, 25, 75, 1, 0, 0, 0,
		27, 77, 1, 0, 0, 0, 29, 81, 1, 0, 0, 0, 31, 37, 5, 34, 0, 0, 32, 33, 5,
		92, 0, 0, 33, 36, 5, 34, 0, 0, 34, 36, 9, 0, 0, 0, 35, 32, 1, 0, 0, 0,
		35, 34, 1, 0, 0, 0, 36, 39, 1, 0, 0, 0, 37, 38, 1, 0, 0, 0, 37, 35, 1,
		0, 0, 0, 38, 40, 1, 0, 0, 0, 39, 37, 1, 0, 0, 0, 40, 41, 5, 34, 0, 0, 41,
		2, 1, 0, 0, 0, 42, 46, 3, 7, 3, 0, 43, 45, 3, 5, 2, 0, 44, 43, 1, 0, 0,
		0, 45, 48, 1, 0, 0, 0, 46, 44, 1, 0, 0, 0, 46, 47, 1, 0, 0, 0, 47, 4, 1,
		0, 0, 0, 48, 46, 1, 0, 0, 0, 49, 52, 3, 7, 3, 0, 50, 52, 3, 19, 9, 0, 51,
		49, 1, 0, 0, 0, 51, 50, 1, 0, 0, 0, 52, 6, 1, 0, 0, 0, 53, 56, 3, 9, 4,
		0, 54, 56, 3, 11, 5, 0, 55, 53, 1, 0, 0, 0, 55, 54, 1, 0, 0, 0, 56, 8,
		1, 0, 0, 0, 57, 58, 8, 0, 0, 0, 58, 10, 1, 0, 0, 0, 59, 60, 5, 92, 0, 0,
		60, 61, 7, 1, 0, 0, 61, 12, 1, 0, 0, 0, 62, 63, 5, 45, 0, 0, 63, 64, 7,
		2, 0, 0, 64, 14, 1, 0, 0, 0, 65, 66, 5, 43, 0, 0, 66, 16, 1, 0, 0, 0, 67,
		68, 5, 124, 0, 0, 68, 18, 1, 0, 0, 0, 69, 70, 5, 45, 0, 0, 70, 20, 1, 0,
		0, 0, 71, 72, 5, 40, 0, 0, 72, 22, 1, 0, 0, 0, 73, 74, 5, 41, 0, 0, 74,
		24, 1, 0, 0, 0, 75, 76, 5, 34, 0, 0, 76, 26, 1, 0, 0, 0, 77, 78, 7, 2,
		0, 0, 78, 79, 1, 0, 0, 0, 79, 80, 6, 13, 0, 0, 80, 28, 1, 0, 0, 0, 81,
		82, 9, 0, 0, 0, 82, 30, 1, 0, 0, 0, 6, 0, 35, 37, 46, 51, 55, 1, 6, 0,
		0,
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

// MexQueryGrammarLexerInit initializes any static state used to implement MexQueryGrammarLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewMexQueryGrammarLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func MexQueryGrammarLexerInit() {
	staticData := &mexquerygrammarlexerLexerStaticData
	staticData.once.Do(mexquerygrammarlexerLexerInit)
}

// NewMexQueryGrammarLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewMexQueryGrammarLexer(input antlr.CharStream) *MexQueryGrammarLexer {
	MexQueryGrammarLexerInit()
	l := new(MexQueryGrammarLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &mexquerygrammarlexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	l.channelNames = staticData.channelNames
	l.modeNames = staticData.modeNames
	l.RuleNames = staticData.ruleNames
	l.LiteralNames = staticData.literalNames
	l.SymbolicNames = staticData.symbolicNames
	l.GrammarFileName = "MexQueryGrammar.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// MexQueryGrammarLexer tokens.
const (
	MexQueryGrammarLexerQUOTED_TERM = 1
	MexQueryGrammarLexerTERM        = 2
	MexQueryGrammarLexerFREE_DASH   = 3
	MexQueryGrammarLexerAND         = 4
	MexQueryGrammarLexerOR          = 5
	MexQueryGrammarLexerNOT         = 6
	MexQueryGrammarLexerLPAR        = 7
	MexQueryGrammarLexerRPAR        = 8
	MexQueryGrammarLexerQUOTE       = 9
	MexQueryGrammarLexerWS          = 10
	MexQueryGrammarLexerLEFTOVER    = 11
)
