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

func (filter Filter) identifier() string {
    var identifier string
    switch filter.Operation {
    case Eq:
        identifier = filter.FieldName
    case In:
        var pluralize = pluralize.NewClient()
        identifier = pluralize.Plural(filter.FieldName)
    }

    // make sure the first letter is lower case
    return identifier
}

func (filter Filter) MethodName() string {
    return filter.identifier()
}

func (filter Filter) MethodArguments() string {
    var pattern string
    switch filter.Operation {
    case Eq:
        pattern = "%s %s"
    case In:
        pattern = "Iterable<%s> %s"
    }
    return fmt.Sprintf(pattern, filter.FieldType, filter.identifier())
}

func (filter Filter) MethodPrototype() string {
    return fmt.Sprintf("%s(%s)", filter.MethodName(), filter.MethodArguments())
}

