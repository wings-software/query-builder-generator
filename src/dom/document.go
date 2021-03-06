package dom

type Document struct {
	Package string
	Queries []Query
}

func (document Document) Init() {
	for i := range document.Queries {
		document.Queries[i].Init()
	}
}
