package dom

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