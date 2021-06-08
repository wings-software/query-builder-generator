package parser

import (
	"fmt"
	"github.com/query-builder-generator/src/dom"
	"text/scanner"
)

type Lexer struct {
	scanner.Scanner
	result dom.Query
}

func (l *Lexer) Lex(lval *DomSymType) int {
	token := l.Scan()
	lit := l.TokenText()
	tok := int(token)

	switch tok {
	case -1:
	case 46:
	case 123:
	case 125:
	default:
		switch lit {
		case "query":
			tok = QUERY
		case "for":
			tok = FOR
		case "project":
			tok = PROJECT
		case "filter":
			tok = FILTER
		case "as":
			tok = AS
		default:
			tok = IDENTIFIER
			lval.identifier = lit
		}
	}
	lval.token = Token{token: tok, literal: lit}
	fmt.Printf("Scanner: %+v, token: %+v, lit: %+v, tok: %+v \n", l.Scanner, token, lit, tok)
	return tok
}
func (l *Lexer) Error(e string) {
	panic(e)
}
