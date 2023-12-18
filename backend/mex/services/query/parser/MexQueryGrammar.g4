grammar MexQueryGrammar;

// --- QUERY GRAMMAR ---

/*
A query is the most generic unit. It allows leftovers and includes the EOF tag to ensure that we always parse
to the end of the input. dangling_op must be included to handle the case where a part of the query is a single
dangling operator (does not match statement).
*/
query: (statement | dangling_op | leftover )* EOF;

/*
A statement is the most generic input containing at least one recognized, non-space item.
We use non-greedy matching (*?) to avoid consecutive bracketed statements being amalgated to a single
statement with dangling operators, e.g. that "(a) (b)" is interpreted as "(a b)" (with inner bracket discarded).
*/
statement: or_expr (or_expr)*?;

/*
The hierarchy of parse rules reflects the precendence of the operator: lowest precedence operator scome first.

The hierarchical naming might initially be confusing since e.g. an OR-expresion need not contain an OR-operator and
may, in fact, be an expression involving only an AND-operator.  Both AND- and OR-expressions, as defined below,
can be single expressions not involving the named opeator at the top level or at all.
*/

// A single AND-expression. or several of them joined with OR
or_expr: and_expr (OR and_expr)*;

// A single AND-operand, or several of them joined with AND
and_expr: operand_expr (AND operand_expr)*;

// Possible operand for a Boolean expression - an unary expression optionally padded with dangling operators
operand_expr: (dangling_op)*? unary_expr (dangling_op)*?;

// A single Boolean block
unary_expr
    : NOT operand_expr           # negated_block_expr
    | LPAR statement RPAR   # bracketed_statement
    | TERM                  # term
    | QUOTED_TERM           # phrase
    ;

// Operators are only matches here if they were not matched as part of a valid expression above
dangling_op
    : AND
    | OR
    | LPAR
    | RPAR
    | QUOTE
    | NOT
    | FREE_DASH
    ;

/*
Anything that does not match any of the above lexing rules - should not happen but is kept to ensure that we catch
it if it happens.
*/
leftover: LEFTOVER;


// --- LEXICAL GRAMMAR ---

/*
A quoted term is any string of characters except unescaped double quotes which is surrounded by double quotes
*/
QUOTED_TERM: '"' ('\\"'|.)*? '"';

/*
A term is any string which
1. does not start with the NOT-operator = '-'
2. does not contain any space characters (normal spaces, tabs, carriage returns, and newlines), even if escaped
3. does not contain any unescaped double quotes = '"'
4. does not contain any unescaped AND and OR operator symbols = '+' and '|'
5. does not contain any unescaped parentheses = '(' and ')'
*/
TERM : VALID_TERM_START_SYMBOL VALID_TERM_MIDDLE_SYMBOL*;

/*
Set of all characters that can occurs in a term after the first character.
The only difference is that the NOT operator symbol ('-') is not allowed
at the start of a term.
*/
fragment VALID_TERM_MIDDLE_SYMBOL: VALID_TERM_START_SYMBOL | NOT;
/*
Set of all characters that can occur as the first character in a term
*/
fragment VALID_TERM_START_SYMBOL: NORMAL_TERM_SYMBOL | ESCAPED_CONTROL_SYMBOL;
/*
Fragment for the characters which are allowed at ANY position in a term without escaping.
*/
fragment NORMAL_TERM_SYMBOL: ~[ \t\r\n\f"+|\-)(];
/*
Escaped versions of the MEx control symbols recognized by the parser. These are defined
separately as two-character blocks to ensure that they are also parsed as such. Otherwise,
smt. like 'hel\|lo' will be tokenized to 'hel\', '|', and 'lo', rather than the single
term 'hel\|lo'.
*/
fragment ESCAPED_CONTROL_SYMBOL: '\\' ["+|)(\-] ;

/*
MEx control symbols for Boolean and grouping operators, parsed as single symbols and not part of
terms or phrases.
*/
FREE_DASH: '-'[ \t\r\n\f];   // Match this before NOT to filter out dashes that cannot act as NOT operators
AND: '+';
OR: '|';
NOT: '-';
LPAR: '(';
RPAR: ')';
QUOTE: '"';

/*
Skip whitespace if unmatched until here (corner case with "disconnected dash" is handled by the FREE_DASH token)
*/
WS : [ \t\r\n\f] -> skip;

/*
Catch any input not matched by other rules. This should never happers but ensure that we catch issues
e.g. if the grammar is changed.
*/
LEFTOVER: .;

