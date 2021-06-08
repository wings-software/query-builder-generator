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
	q := Parse(`query DelegateTask for io.harness.beans.DelegateTasks
		{
			filter accountId as string;
		}`)
	assert.Equal(t, dom.Query{Name: "DelegateTask",
		Collection: "io.harness.beans.DelegateTasks",
		Filters: []dom.Filter{
			{FieldType: "string", FieldName: "accountId"},
		},
	}, q)
}

func TestSanityWithMultipleFilters(t *testing.T) {
	q := Parse(`query DelegateTask for io.harness.beans.DelegateTasks  
					{
						filter accountId as string ;
						filter delegateId as int;
					}`)
	assert.Equal(t, dom.Query{Name: "DelegateTask",
		Collection: "io.harness.beans.DelegateTasks",
		Filters: []dom.Filter{
			{FieldType: "string", FieldName: "accountId"},
			{FieldType: "int", FieldName: "delegateId"},
		},
	}, q)
}

func TestFiltersWithList(t *testing.T) {
	q := Parse(`query DelegateTask for io.harness.beans.DelegateTasks  
					{
						filter accountId as string from list ;
						filter delegateId as int from list;
					}`)
	assert.Equal(t, dom.Query{Name: "DelegateTask",
		Collection: "io.harness.beans.DelegateTasks",
		Filters: []dom.Filter{
			{FieldType: "string", FieldName: "accountId", Operation: dom.IN},
			{FieldType: "int", FieldName: "delegateId", Operation: dom.IN},
		},
	}, q)
}

func TestFiltersWithFullPath(t *testing.T) {
	q := Parse(`query DelegateTask for io.harness.beans.DelegateTasks  
					{
						filter accountId as io.harness.beans.Id;
					}`)
	assert.Equal(t, dom.Query{Name: "DelegateTask",
		Collection: "io.harness.beans.DelegateTasks",
		Filters: []dom.Filter{
			{FieldType: "io.harness.beans.Id", FieldName: "accountId", Operation: dom.EQUAL},
		},
	}, q)
}

func TestSanityWithSingleFilterAndProject(t *testing.T) {
	q := Parse(`query Select for io.harness.beans.DelegateTasks 
					{
						filter accountId as string ;
						project id ;
					}`)
	assert.Equal(t, dom.Query{Name: "Select",
		Collection: "io.harness.beans.DelegateTasks",
		Filters:    []dom.Filter{dom.Filter{FieldType: "accountId", FieldName: "string"}}, ProjectFields: []string{"id"}}, q)
}
