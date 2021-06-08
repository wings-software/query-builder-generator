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

func TestSanityWithSingleFilter(t *testing.T) {
	q := Parse("query DelegateTask for io.harness.beans.DelegateTasks " +
		"{" +
		"filter accountId as string ;" +
		" }")
	assert.Equal(t, dom.Query{Name: "DelegateTask",
		Collection: "io.harness.beans.DelegateTasks",
		Filters:    []dom.Filter{dom.Filter{FieldType: "accountId", FieldName: "string"}}}, q)
}

func TestSanityWithMultipleFilters(t *testing.T) {
	q := Parse(`query DelegateTask for io.harness.beans.DelegateTasks  
					{
						filter accountId as string ;
						filter delegateId as int;
					}`)
	assert.Equal(t, dom.Query{Name: "DelegateTask",
		Collection: "io.harness.beans.DelegateTasks",
		Filters: []dom.Filter{dom.Filter{FieldType: "accountId", FieldName: "string"},
			dom.Filter{FieldType: "delegateId", FieldName: "int"}}}, q)
}

func TestSanityWithMultipleFiltersWithList(t *testing.T) {
	q := Parse(`query DelegateTask for io.harness.beans.DelegateTasks  
					{
						filter accountId as string from list ;
						filter delegateId as int from list;
					}`)
	assert.Equal(t, dom.Query{Name: "DelegateTask",
		Collection: "io.harness.beans.DelegateTasks",
		Filters: []dom.Filter{dom.Filter{FieldType: "accountId", FieldName: "string", Operation: dom.IN},
			dom.Filter{FieldType: "delegateId", FieldName: "int", Operation: dom.IN}}}, q)
}

