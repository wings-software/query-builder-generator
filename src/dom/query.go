package dom

import (
	"fmt"
)

type Query struct {
	Name          string
	Collection    string
	Filters 	[]Filter
	ProjectFields []string
}

func (query Query) Init() {
	for i := range query.Filters {
		query.Filters[i].Query = &query
	}
}

func (query Query) InterfaceName() string {
	var name = query.Name
	return fmt.Sprintf("%sQueryFinal", name )
}