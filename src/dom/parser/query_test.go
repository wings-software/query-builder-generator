//go:generate goyacc -o query.y.go -p Dom query.y

package parser

import (
	"github.com/query-builder-generator/src/dom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSanity(t *testing.T) {
	q := Parse("query foo for io.harness.beans.DelegateTasks {}")
	assert.Equal(t, dom.Query{Name: "foo", Collection:"io.harness.beans.DelegateTasks"}, q)
}

//func TestSanity1(t *testing.T) {
//	doc := Parse("query foo {}")
//	assert.Equal(t, doc.queries[0].name, "foo")
//}