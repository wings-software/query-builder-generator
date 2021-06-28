package dom

import (
	"fmt"
	"strings"
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

func (query Query) CollectionName() string {
	ss := strings.Split(query.Collection, ".")
	return ss[len(ss)-1]
}

func (query Query) InterfaceName() string {
	var name = query.Name
	return fmt.Sprintf("%sQueryFinal", name )
}