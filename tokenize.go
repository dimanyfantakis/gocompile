package main

import (
	"fmt"
	"strconv"
	"unicode"
)

func tokenize(input string) []token {
	tokens := make([]token, 0, 512)
	i := 0
	for i < len(input) {
		r := rune(input[i])
		// Skip whitespaces.
		if unicode.IsSpace(r) {
			if r == '\n' {
				t := getToken(&i, input)
				logger.Printf("Token value: %s type: %d", strconv.Quote(t.contents), t.tokenType)
				tokens = append(tokens, *t)
			}
			pop(&i)
			continue
		}
		// Skip comments.
		if r == '#' {
			for i < len(input) {
				if input[i] == '\n' {
					break
				}
				pop(&i)
			}
			if i == len(input) {
				break
			}
		}

		t := getToken(&i, input)
		logger.Printf("Token value: %s type: %d", strconv.Quote(t.contents), t.tokenType)
		tokens = append(tokens, *t)
		pop(&i)
	}
	return tokens
}

func getToken(i *int, s string) *token {
	r := rune(s[*i])
	switch r {
	case '+':
		return &token{contents: string(r), tokenType: PLUS}
	case '-':
		return &token{contents: string(r), tokenType: MINUS}
	case '*':
		return &token{contents: string(r), tokenType: ASTERISK}
	case '/':
		return &token{contents: string(r), tokenType: SLASH}
	case '\n':
		return &token{contents: string(r), tokenType: NEWLINE}
	case '=':
		if peek(*i, s) == '=' {
			pop(i)
			return &token{contents: "==", tokenType: EQEQ}
		}
		return &token{contents: string(r), tokenType: EQ}
	case '>':
		if peek(*i, s) == '=' {
			pop(i)
			return &token{contents: ">=", tokenType: GTEQ}
		}
		return &token{contents: string(r), tokenType: GT}
	case '<':
		if peek(*i, s) == '=' {
			pop(i)
			return &token{contents: ">=", tokenType: LTEQ}
		}
		return &token{contents: string(r), tokenType: LT}
	case '!':
		if peek(*i, s) == '=' {
			pop(i)
			return &token{contents: "!=", tokenType: NOTEQ}
		} else {
			e := fmt.Sprintf("Unexpected character %c", r)
			panic(e)
		}

	case '"':
		pop(i)
		q := parseQuotes(i, s)
		return &token{contents: string(q), tokenType: STRING}
	default:
		if unicode.IsDigit(r) {
			n := parseNumber(i, s)
			return &token{contents: string(n), tokenType: NUMBER}
		} else if unicode.IsLetter(r) {
			w, t := parseWord(i, s)
			return &token{contents: w, tokenType: t}
		}
		e := fmt.Sprintf("Unexpected character %c", r)
		panic(e)
	}
}

// Small utility functions.

func parseWord(i *int, s string) (string, int32) {
	w := string(s[*i])
	for *i < len(s) {
		if unicode.IsLetter(peek(*i, s)) {
			w += string(peek(*i, s))
			pop(i)
		} else {
			break
		}
	}
	t := parseWordType(w)
	return w, t
}

// TODO float.
func parseNumber(i *int, s string) string {
	r := string(s[*i])
	for *i < len(s) {
		if unicode.IsDigit(peek(*i, s)) {
			r += string(peek(*i, s))
			pop(i)
		} else {
			break
		}
	}
	return r
}

func parseQuotes(i *int, s string) string {
	r := ""
	for *i < len(s) {
		if s[*i] != '"' {
			r += string(s[*i])
			pop(i)
		} else {
			break
		}
	}
	return r
}

type token struct {
	tokenType int32
	contents  string
}

type special = int32
type keyword = int32
type operator = int32
type comparator = int32

const (
	EOF special = iota - 1
	NEWLINE
	NUMBER
	IDENT
	STRING
)
const (
	LABEL keyword = iota + 4
	GOTO
	PRINT
	INPUT
	LET
	IF
	THEN
	ENDIF
	WHILE
	REPEAT
	ENDWHILE
)
const (
	PLUS operator = iota + 15
	MINUS
	ASTERISK
	SLASH
)
const (
	EQ comparator = iota + 19
	EQEQ
	NOTEQ
	LT
	LTEQ
	GT
	GTEQ
)

func parseWordType(s string) int32 {
	switch s {
	case "LABEL":
		return LABEL
	case "GOTO":
		return GOTO
	case "PRINT":
		return PRINT
	case "INPUT":
		return INPUT
	case "LET":
		return LET
	case "IF":
		return IF
	case "THEN":
		return THEN
	case "ENDIF":
		return ENDIF
	case "WHILE":
		return WHILE
	case "REPEAT":
		return REPEAT
	case "ENDWHILE":
		return ENDWHILE
	default:
		return IDENT
	}
}

func peek(i int, s string) rune {
	return rune(s[i+1])
}

func pop(i *int) {
	*i++
}
