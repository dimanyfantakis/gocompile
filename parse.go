package main

import (
	"fmt"
)

var labels = make(map[string]bool)
var gotolabels = make(map[string]bool)
var idents = make(map[string]bool)

func parse(tokens []token) {
	logger.Print("PARSER")
	i := 0
	for i < len(tokens) {
		t := statement(&tokens[i], &i, tokens)
		// Skip newlines.
		newLine(t, tokens, &i)
	}
	// Validate GOTO labels.
	for l := range gotolabels {
		if !labels[l] {
			logger.Panicf("Label, %s hasn't been declared.", l)
		}
	}
	fmt.Print(&buf)
}

/*
statement ::=

	"PRINT" (expression | string) nl
	| "IF" comparison "THEN" nl {statement} "ENDIF" nl
	| "WHILE" comparison "REPEAT" nl {statement} "ENDWHILE" nl
	| "LABEL" ident nl
	| "GOTO" ident nl
	| "LET" ident "=" expression nl
	| "INPUT" ident nl
*/
func statement(t *token, i *int, tokens []token) *token {
	switch t.tokenType {
	// "PRINT" (expression | string) {nl}+.
	case PRINT:
		logger.Print("PRINT")
		*t = *next(i, tokens)
		if t.tokenType == STRING {
			*t = *next(i, tokens)
		} else {
			expression(t, i, tokens)
		}
	// "IF" comparison "THEN" nl {statement} "ENDIF" nl.
	case IF:
		logger.Print("IF")
		*t = *next(i, tokens)
		comparison(t, i, tokens)
		match(t, THEN)
		*t = *next(i, tokens)
		newLine(t, tokens, i)
		statement(t, i, tokens)
		newLine(t, tokens, i)
		match(t, ENDIF)
		logger.Print("ENDIF")
		*t = *next(i, tokens)
	// "WHILE" comparison "REPEAT" nl {statement nl} "ENDWHILE" nl.
	case WHILE:
		logger.Print("WHILE")
		*t = *next(i, tokens)
		comparison(t, i, tokens)
		match(t, REPEAT)
		*t = *next(i, tokens)
		newLine(t, tokens, i)
		for t.tokenType != ENDWHILE {
			statement(t, i, tokens)
			newLine(t, tokens, i)
		}
		// *t = *next(i, tokens)
		match(t, ENDWHILE)
		*t = *next(i, tokens)
	// "LABEL" ident nl.
	case LABEL:
		logger.Print("LABEL")
		*t = *next(i, tokens)
		match(t, IDENT)
		_, b := labels[t.contents]
		if b {
			logger.Panicf("Label, %s has already been declared.", t.contents)
		}
		labels[t.contents] = true
		*t = *next(i, tokens)
	// "GOTO" ident nl.
	case GOTO:
		logger.Print("GOTO")
		*t = *next(i, tokens)
		match(t, IDENT)
		gotolabels[t.contents] = true
		*t = *next(i, tokens)
	// "LET" ident "=" expression nl.
	case LET:
		logger.Print("LET")
		*t = *next(i, tokens)
		match(t, IDENT)
		idents[t.contents] = true
		*t = *next(i, tokens)
		match(t, EQ)
		*t = *next(i, tokens)
		expression(t, i, tokens)
	// "INPUT" ident nl
	case INPUT:
		logger.Print("INPUT")
		*t = *next(i, tokens)
		match(t, IDENT)
		idents[t.contents] = true
		*t = *next(i, tokens)
	}
	return t
}

// comparison ::= expression (("==" | "!=" | ">" | ">=" | "<" | "<=") expression)+
func comparison(t *token, i *int, tokens []token) {
	logger.Print("COMPARISON")
	expression(t, i, tokens)
	parseComparisonType(t)
	*t = *next(i, tokens)
	expression(t, i, tokens)
	if isComparisonType(t) {
		// TODO rest of comparison types.
		for *i < len(tokens) {
			parseComparisonType(t)
			*t = *next(i, tokens)
			expression(t, i, tokens)
		}
	}
}

// expression ::= term {( "-" | "+" ) term}
func expression(t *token, i *int, tokens []token) {
	logger.Print("EXPRESSION")
	term(t, i, tokens)
	if t.tokenType == PLUS || t.tokenType == MINUS {
		*t = *next(i, tokens)
		term(t, i, tokens)
	}
}

// term ::= unary {( "/" | "*" ) unary}
func term(t *token, i *int, tokens []token) {
	logger.Print("TERM")
	unary(t, i, tokens)
	if t.tokenType == SLASH || t.tokenType == ASTERISK {
		*t = *next(i, tokens)
		unary(t, i, tokens)
	}
}

// unary ::= ["+" | "-"] primary
func unary(t *token, i *int, tokens []token) {
	logger.Print("UNARY")
	if t.tokenType == PLUS || t.tokenType == MINUS {
		t = next(i, tokens)
	}
	primary(t, i, tokens)
}

// primary ::= number | ident
func primary(t *token, i *int, tokens []token) {
	logger.Printf("PRIMARY ( %s )", t.contents)
	if t.tokenType == NUMBER {
		logger.Printf("NUMBER ( %s )", t.contents)
		*t = *next(i, tokens)
	} else if t.tokenType == IDENT {
		// Has the identifier already been declared?
		if !idents[t.contents] {
			logger.Panicf("Identifier %s, hasn't been declared yet!", t.contents)
		}
		logger.Printf("IDENT ( %s )", t.contents)
		*t = *next(i, tokens)
	} else {
		logger.Panicf("Expected number or ident, but got %c, %s", t.tokenType, t.contents)
	}
}

// Small utility functions.

func match(t *token, exp int32) {
	if t.tokenType != exp {
		logger.Panicf("Expected type %c, but got %c", exp, t.tokenType)
	}
}

func newLine(t *token, tokens []token, i *int) {
	logger.Print("NEWLINE")
	match(t, NEWLINE)
	if *i == len(tokens)-1 {
		*i++
		return
	}
	*t = *next(i, tokens)
}

func parseComparisonType(t *token) comparator {
	switch t.tokenType {
	case EQ:
		return EQ
	case EQEQ:
		return EQEQ
	case NOTEQ:
		return NOTEQ
	case LT:
		return LT
	case LTEQ:
		return LTEQ
	case GT:
		return GT
	case GTEQ:
		return GTEQ
	}
	logger.Panicf("Expected comparison operator, but got %s", t.contents)
	return -1
}

func isComparisonType(t *token) bool {
	switch t.tokenType {
	case EQ:
	case EQEQ:
	case NOTEQ:
	case LT:
	case LTEQ:
	case GT:
	case GTEQ:
		return true
	}
	return false
}

func next(i *int, tokens []token) *token {
	*i++
	return &tokens[*i]
}
