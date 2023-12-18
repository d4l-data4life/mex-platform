// Code generated from MexQueryGrammar.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // MexQueryGrammar

import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseMexQueryGrammarListener is a complete listener for a parse tree produced by MexQueryGrammarParser.
type BaseMexQueryGrammarListener struct{}

var _ MexQueryGrammarListener = &BaseMexQueryGrammarListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseMexQueryGrammarListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseMexQueryGrammarListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseMexQueryGrammarListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseMexQueryGrammarListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterQuery is called when production query is entered.
func (s *BaseMexQueryGrammarListener) EnterQuery(ctx *QueryContext) {}

// ExitQuery is called when production query is exited.
func (s *BaseMexQueryGrammarListener) ExitQuery(ctx *QueryContext) {}

// EnterStatement is called when production statement is entered.
func (s *BaseMexQueryGrammarListener) EnterStatement(ctx *StatementContext) {}

// ExitStatement is called when production statement is exited.
func (s *BaseMexQueryGrammarListener) ExitStatement(ctx *StatementContext) {}

// EnterOr_expr is called when production or_expr is entered.
func (s *BaseMexQueryGrammarListener) EnterOr_expr(ctx *Or_exprContext) {}

// ExitOr_expr is called when production or_expr is exited.
func (s *BaseMexQueryGrammarListener) ExitOr_expr(ctx *Or_exprContext) {}

// EnterAnd_expr is called when production and_expr is entered.
func (s *BaseMexQueryGrammarListener) EnterAnd_expr(ctx *And_exprContext) {}

// ExitAnd_expr is called when production and_expr is exited.
func (s *BaseMexQueryGrammarListener) ExitAnd_expr(ctx *And_exprContext) {}

// EnterOperand_expr is called when production operand_expr is entered.
func (s *BaseMexQueryGrammarListener) EnterOperand_expr(ctx *Operand_exprContext) {}

// ExitOperand_expr is called when production operand_expr is exited.
func (s *BaseMexQueryGrammarListener) ExitOperand_expr(ctx *Operand_exprContext) {}

// EnterNegated_block_expr is called when production negated_block_expr is entered.
func (s *BaseMexQueryGrammarListener) EnterNegated_block_expr(ctx *Negated_block_exprContext) {}

// ExitNegated_block_expr is called when production negated_block_expr is exited.
func (s *BaseMexQueryGrammarListener) ExitNegated_block_expr(ctx *Negated_block_exprContext) {}

// EnterBracketed_statement is called when production bracketed_statement is entered.
func (s *BaseMexQueryGrammarListener) EnterBracketed_statement(ctx *Bracketed_statementContext) {}

// ExitBracketed_statement is called when production bracketed_statement is exited.
func (s *BaseMexQueryGrammarListener) ExitBracketed_statement(ctx *Bracketed_statementContext) {}

// EnterTerm is called when production term is entered.
func (s *BaseMexQueryGrammarListener) EnterTerm(ctx *TermContext) {}

// ExitTerm is called when production term is exited.
func (s *BaseMexQueryGrammarListener) ExitTerm(ctx *TermContext) {}

// EnterPhrase is called when production phrase is entered.
func (s *BaseMexQueryGrammarListener) EnterPhrase(ctx *PhraseContext) {}

// ExitPhrase is called when production phrase is exited.
func (s *BaseMexQueryGrammarListener) ExitPhrase(ctx *PhraseContext) {}

// EnterDangling_op is called when production dangling_op is entered.
func (s *BaseMexQueryGrammarListener) EnterDangling_op(ctx *Dangling_opContext) {}

// ExitDangling_op is called when production dangling_op is exited.
func (s *BaseMexQueryGrammarListener) ExitDangling_op(ctx *Dangling_opContext) {}

// EnterLeftover is called when production leftover is entered.
func (s *BaseMexQueryGrammarListener) EnterLeftover(ctx *LeftoverContext) {}

// ExitLeftover is called when production leftover is exited.
func (s *BaseMexQueryGrammarListener) ExitLeftover(ctx *LeftoverContext) {}
