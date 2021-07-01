package dom

import (
	"fmt"
	"strings"
)

type Query struct {
	Name            string
	Collection      string
	Filters 	  []Filter
	Optionals 	  []Optional
	ProjectFields []string
}

func (query Query) Init() {
	for i := range query.Filters {
		query.Filters[i].Query = &query
	}
	if query.Optionals != nil {
		for i := range query.Optionals {
			query.Optionals[i].Query = &query
			query.Optionals[i].Init()
		}
	}
}

func (query Query) CollectionName() string {
	ss := strings.Split(query.Collection, ".")
	return ss[len(ss)-1]
}

func (query Query) InterfaceName() string {
	var name = query.Name
	if len(query.Optionals) == 0 {
		return fmt.Sprintf("%sQueryFinal", name)
	} else {
		return fmt.Sprintf("%sQueryOptions", name)
	}
}

func (query Query) ReturnFromThis() string {
	return "return this;"
}