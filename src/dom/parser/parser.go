package parser

import (
	"github.com/query-builder-generator/src/dom"
	"strings"
)

func Parse(exp string) dom.Document {
	l := new(Lexer)
	l.Init(strings.NewReader(exp))
	DomParse(l)
	return l.result
}
