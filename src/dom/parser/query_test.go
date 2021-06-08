//go:generate goyacc -o query.y.go -p Dom query.y

package parser

import (
	"github.com/query-builder-generator/src/dom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSanity(t *testing.T) {
	q := Parse("query DelegateTask for io.harness.beans.DelegateTasks {}")
	assert.Equal(t, dom.Query{Name: "DelegateTask", Collection: "io.harness.beans.DelegateTasks"}, q)
}

func TestSanityWithStatements(t *testing.T) {
	q := Parse("query DelegateTask for io.harness.beans.DelegateTasks " +
		"{" +
		"filter accountId as string " +
		"project {  uuid  }" +
		" }")
	assert.Equal(t, dom.Query{Name: "DelegateTask", Collection: "io.harness.beans.DelegateTasks", QueryStmtType: "{filter string accountId}", ProjectFields:"uuid"},  q)
}
