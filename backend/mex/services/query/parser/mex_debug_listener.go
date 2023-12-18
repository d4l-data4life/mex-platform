package parser

import (
	"encoding/json"
	"fmt"
)

type SubTreeQueryObject map[string]interface{}

// ObjectStack provides a simple stack implementation
type ObjectStack struct {
	elem []SubTreeQueryObject
}

func (s *ObjectStack) isEmpty() bool {
	return len(s.elem) == 0
}

func (s *ObjectStack) push(e SubTreeQueryObject) {
	s.elem = append(s.elem, e)
}

func (s *ObjectStack) pop() (SubTreeQueryObject, error) {
	if s.isEmpty() {
		return SubTreeQueryObject{}, fmt.Errorf("conversionStack is empty")
	}
	lastIndex := len(s.elem) - 1
	e := s.elem[lastIndex]
	s.elem = s.elem[:lastIndex]
	return e, nil
}

type DebugListener struct {
	*BaseMexQueryGrammarListener

	showFull bool // If true, all nodes will be included

	conversionStack ObjectStack
	Depth           int
	errs            []error
}

func NewDebugListener(showFull bool) *DebugListener {
	return &DebugListener{
		showFull: showFull,
	}
}

// PopMultiple returns multiple objects from stack (last popped first)
func (listener *DebugListener) PopMultiple(count int) ([]map[string]interface{}, error) {
	terms := make([]map[string]interface{}, count)
	// We fill the slice from the end to recover the original term order
	for i := count - 1; i >= 0; i-- {
		q, err := listener.conversionStack.pop()
		if err != nil {
			return []map[string]interface{}{}, err
		}
		terms[i] = q
	}
	return terms, nil
}

func (listener *DebugListener) GetResult() (map[string]interface{}, error) {
	// After successful conversion, the listener should contain the final Solr query as the only element
	res, err := listener.PopMultiple(1)
	if err != nil {
		return nil, fmt.Errorf("stack does not contain parse tree object")
	} else if !listener.conversionStack.isEmpty() {
		topElem, _ := listener.PopMultiple(1)
		topObj, _ := json.Marshal(topElem[0])
		return nil, fmt.Errorf("listener stack contained multiple elements at end of parsing. "+
			"result popped: %s\n, top element left on stack: %s", res[0], topObj)
	}
	return res[0], nil
}

func (listener *DebugListener) handleUnaryNode(nodeName string) {
	term, err := listener.PopMultiple(1)
	if err != nil {
		listener.errs = append(listener.errs, err)
		return
	}
	res := map[string]interface{}{
		nodeName: term[0],
	}
	listener.conversionStack.push(res)
}

func (listener *DebugListener) handleMultiTermNode(childCount int, nodeName string, ignoreUnaryNodes bool) {
	if ignoreUnaryNodes && childCount == 1 && !listener.showFull {
		return
	}
	var res map[string]interface{}
	if childCount > 0 {
		terms, err := listener.PopMultiple(childCount)
		if err != nil {
			listener.errs = append(listener.errs, err)
			return
		}

		res = map[string]interface{}{
			nodeName: terms,
		}
	}
	listener.conversionStack.push(res)
}

// ExitQuery is called when production query is exited.
func (listener *DebugListener) ExitQuery(ctx *QueryContext) {
	listener.handleMultiTermNode(len(ctx.AllStatement()), "query", true)
}

// ExitStatement is called when production statement is exited.
func (listener *DebugListener) ExitStatement(ctx *StatementContext) {
	listener.handleMultiTermNode(len(ctx.AllOr_expr()), "statement", true)
}

// ExitOr_expr is called when production or_expr is exited.
func (listener *DebugListener) ExitOr_expr(ctx *Or_exprContext) {
	listener.handleMultiTermNode(len(ctx.AllAnd_expr()), "or_expr", true)
}

// ExitAnd_expr is called when production and_expr is exited.
func (listener *DebugListener) ExitAnd_expr(ctx *And_exprContext) {
	listener.handleMultiTermNode(len(ctx.AllOperand_expr()), "and_expr", true)
}

// ExitOperand is called when production operand is exited.
func (listener *DebugListener) ExitOperand_expr(_ *Operand_exprContext) {
	if !listener.showFull {
		return
	}
	listener.handleUnaryNode("operand")
}

// ExitUnary_expr is called when production unary_expr is exited.
func (listener *DebugListener) ExitUnary_expr(_ *Unary_exprContext) {
	if !listener.showFull {
		return
	}
	listener.handleUnaryNode("unary_expr")
}

// ExitNegated_block_expr is called when production negated_block_expr is exited.
func (listener *DebugListener) ExitNegated_block_expr(_ *Negated_block_exprContext) {
	listener.handleUnaryNode("negated_block_expr")
}

// ExitTerm is called when production term is exited.
func (listener *DebugListener) ExitTerm(ctx *TermContext) {
	res := map[string]interface{}{
		"term": ctx.GetText(),
	}
	listener.conversionStack.push(res)
}

// ExitPhrase is called when production phrase is exited.
func (listener *DebugListener) ExitPhrase(ctx *PhraseContext) {
	res := map[string]interface{}{
		"phrase": ctx.GetText(),
	}
	listener.conversionStack.push(res)
}

// ExitBracketed_statement is called when production bracketed_query is exited.
func (listener *DebugListener) ExitBracketed_statement(_ *Bracketed_statementContext) {
	listener.handleUnaryNode("bracketed_query")
}

// ExitDangling_op is called when exiting the dangling_op production.
func (listener *DebugListener) ExitDangling_op(ctx *Dangling_opContext) {
	res := map[string]interface{}{
		"dangling_operator": ctx.GetText(),
	}
	listener.conversionStack.push(res)
}

// ExitLeftover is called when production leftover is exited.
func (listener *DebugListener) ExitLeftover(ctx *LeftoverContext) {
	res := map[string]interface{}{
		"leftover": ctx.GetText(),
	}
	listener.conversionStack.push(res)
}
