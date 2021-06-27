package dom

import (
    "fmt"
    pluralize "github.com/gertd/go-pluralize"
    "strings"
)

type Filter struct {
    Query *Query
    FieldType string
    FieldName string
    Operation OperationType
}

func (filter Filter) InterfaceName() string {
    var name = filter.Query.Name
    var field string
    switch filter.Operation {
    case Eq:
        field = filter.FieldName
    case In:
        var pluralize = pluralize.NewClient()
        field = pluralize.Plural(filter.FieldName)
    }

    field = strings.Title(field)
    return fmt.Sprintf("%sQuery%s", name , field)
}