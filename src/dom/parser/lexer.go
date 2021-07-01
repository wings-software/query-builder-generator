package parser

import (
	"github.com/query-builder-generator/src/dom"
	"text/scanner"
)

type Lexer struct {
	scanner.Scanner
	result dom.Document
}

func (l *Lexer) Lex(lval *DomSymType) int {
	token := l.Scan()
	lit := l.TokenText()
	tok := int(token)

	switch tok {
	case -1:
	case 46:
	case 59:
	case 123:
	case 125:
	default:
		switch lit {
		case "equal":
			tok = EQUAL
		case "for":
			tok = FOR
		case "in":
			tok = IN
		case "is":
			tok = IS
		case "less":
			tok = LESS
		case "list":
			tok = LIST
		case "module":
			tok = MODULE
		case "of":
			tok = OF
		case "optional":
			tok = OPTIONAL
		case "project":
			tok = PROJECT
		case "query":
			tok = QUERY
		case "remainder":
			tok = REMAINDER
		case "than":
			tok = THAN

		default:
			tok = IDENTIFIER
			lval.identifier = lit
		}
	}
	lval.token = Token{token: tok, literal: lit}
	//fmt.Printf("Scanner: %+v, token: %+v, lit: %+v, tok: %+v \n", l.Scanner, token, lit, tok)
	return tok
}
func (l *Lexer) Error(e string) {
	panic(e)
}
