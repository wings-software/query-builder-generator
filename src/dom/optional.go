package dom

type Optional struct {
	Query        *Query
	Name          string
	Filters 	[]Filter
}

func (optional Optional) Init() {
	for i := range optional.Filters {
		optional.Filters[i].Query = optional.Query
	}
}