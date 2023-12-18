// Code generated from MexQueryGrammar.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // MexQueryGrammar

import "github.com/antlr/antlr4/runtime/Go/antlr"

// MexQueryGrammarListener is a complete listener for a parse tree produced by MexQueryGrammarParser.
type MexQueryGrammarListener interface {
	antlr.ParseTreeListener

	// EnterQuery is called when entering the query production.
	EnterQuery(c *QueryContext)

	// EnterStatement is called when entering the statement production.
	EnterStatement(c *StatementContext)

	// EnterOr_expr is called when entering the or_expr production.
	EnterOr_expr(c *Or_exprContext)

	// EnterAnd_expr is called when entering the and_expr production.
	EnterAnd_expr(c *And_exprContext)

	// EnterOperand_expr is called when entering the operand_expr production.
	EnterOperand_expr(c *Operand_exprContext)

	// EnterNegated_block_expr is called when entering the negated_block_expr production.
	EnterNegated_block_expr(c *Negated_block_exprContext)

	// EnterBracketed_statement is called when entering the bracketed_statement production.
	EnterBracketed_statement(c *Bracketed_statementContext)

	// EnterTerm is called when entering the term production.
	EnterTerm(c *TermContext)

	// EnterPhrase is called when entering the phrase production.
	EnterPhrase(c *PhraseContext)

	// EnterDangling_op is called when entering the dangling_op production.
	EnterDangling_op(c *Dangling_opContext)

	// EnterLeftover is called when entering the leftover production.
	EnterLeftover(c *LeftoverContext)

	// ExitQuery is called when exiting the query production.
	ExitQuery(c *QueryContext)

	// ExitStatement is called when exiting the statement production.
	ExitStatement(c *StatementContext)

	// ExitOr_expr is called when exiting the or_expr production.
	ExitOr_expr(c *Or_exprContext)

	// ExitAnd_expr is called when exiting the and_expr production.
	ExitAnd_expr(c *And_exprContext)

	// ExitOperand_expr is called when exiting the operand_expr production.
	ExitOperand_expr(c *Operand_exprContext)

	// ExitNegated_block_expr is called when exiting the negated_block_expr production.
	ExitNegated_block_expr(c *Negated_block_exprContext)

	// ExitBracketed_statement is called when exiting the bracketed_statement production.
	ExitBracketed_statement(c *Bracketed_statementContext)

	// ExitTerm is called when exiting the term production.
	ExitTerm(c *TermContext)

	// ExitPhrase is called when exiting the phrase production.
	ExitPhrase(c *PhraseContext)

	// ExitDangling_op is called when exiting the dangling_op production.
	ExitDangling_op(c *Dangling_opContext)

	// ExitLeftover is called when exiting the leftover production.
	ExitLeftover(c *LeftoverContext)
}
