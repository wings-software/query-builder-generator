package parser

import (
	"github.com/query-builder-generator/src/dom"
	"github.com/stretchr/testify/assert"
	"testing"
)

//func TestSanity(t *testing.T) {
//	q := Parse("query DelegateTask for io.harness.beans.DelegateTasks {}")
//	assert.Equal(t, dom.Query{Name: "DelegateTask", Collection: "io.harness.beans.DelegateTasks"}, q)
//}

func TestSanityWithSingleFilter(t *testing.T) {
	document := Parse(`query DelegateTask for io.harness.beans.DelegateTasks
		{
			filter accountId as string ;
		}`)

	var query = dom.Query{
		Name: "DelegateTask",
		Collection: "io.harness.beans.DelegateTasks",
		Filters: []dom.Filter{
			{FieldType: "string", FieldName: "accountId", Operation: dom.Eq},
		},
	}
	query.Init()

	assert.Equal(t, query, document.Queries[0])
}

func TestSanityWithMultipleFilters(t *testing.T) {
	document := Parse(`query DelegateTask for io.harness.beans.DelegateTasks  
					{
						filter accountId as string ;
						filter delegateId as int ;
					}`)

	var query = dom.Query{
		Name: "DelegateTask",
		Collection: "io.harness.beans.DelegateTasks",
		Filters: []dom.Filter{
			{FieldType: "string", FieldName: "accountId", Operation: dom.Eq},
			{FieldType: "int", FieldName: "delegateId", Operation: dom.Eq},
		},
	}
	query.Init()

	assert.Equal(t, query, document.Queries[0])
}

func TestFiltersWithList(t *testing.T) {
	document := Parse(`query DelegateTask for io.harness.beans.DelegateTasks
					{
						filter accountId as io.harness.beans.Id from list ;
						filter delegateId as int from list ;
					}`)

	var query = dom.Query{
		Name: "DelegateTask",
		Collection: "io.harness.beans.DelegateTasks",
		Filters: []dom.Filter{
			{FieldType: "io.harness.beans.Id", FieldName: "accountId", Operation: dom.In},
			{FieldType: "int", FieldName: "delegateId", Operation: dom.In},
		},
	}
	query.Init()

	assert.Equal(t, query, document.Queries[0])
}

func TestFiltersWithFullPath(t *testing.T) {
	document := Parse(`query DelegateTask for io.harness.beans.DelegateTasks  
					{
						filter accountId as io.harness.beans.Id ;
					}`)

	var query = dom.Query{
		Name: "DelegateTask",
		Collection: "io.harness.beans.DelegateTasks",
		Filters: []dom.Filter{
			{FieldType: "io.harness.beans.Id", FieldName: "accountId", Operation: dom.Eq},
		},
	}
	query.Init()

	assert.Equal(t, query, document.Queries[0])
}

func TestWithNoProject(t *testing.T) {
	document := Parse(`query Select for io.harness.beans.DelegateTasks 
					{
						filter accountId as string ;
					}`)
	var query = dom.Query{
		Name: "Select",
		Collection: "io.harness.beans.DelegateTasks",
		Filters: []dom.Filter{
			{FieldType: "string", FieldName: "accountId", Operation: dom.Eq},
		},
	}
	query.Init()

	assert.Equal(t, query, document.Queries[0])
}

func TestWithProject(t *testing.T) {
	document := Parse(`query Select for io.harness.beans.DelegateTasks 
					{
						filter accountId as string ;
						project id ;
					}`)

	var query = dom.Query{
		Name: "Select",
		Collection: "io.harness.beans.DelegateTasks",
		Filters: []dom.Filter{
			{FieldType: "string", FieldName: "accountId", Operation: dom.Eq},
		},
		ProjectFields: []string{"id"},
	}
	query.Init()

	assert.Equal(t, query, document.Queries[0])
}

func TestWithMultipleProjects(t *testing.T) {
	document := Parse(`query Select for io.harness.beans.DelegateTasks 
					{
						filter accountId as string ;
						project id ;
						project foo_bar ;
					}`)
	var query = dom.Query{
		Name: "Select",
		Collection: "io.harness.beans.DelegateTasks",
		Filters: []dom.Filter{
			{FieldType: "string", FieldName: "accountId", Operation: dom.Eq},
		},
		ProjectFields: []string{"id", "foo_bar"},
	}
	query.Init()

	assert.Equal(t, query, document.Queries[0])
}
