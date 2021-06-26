package dom

import (
    "fmt"
    "strings"
)

type Filter struct {
    Query *Query
    FieldType string
    FieldName string
    Operation OperationType
}

func (f Filter) InterfaceName() string {
    var name = f.Query.Name
    var field = strings.Title(f.FieldName)
    return fmt.Sprintf("%sQuery%s", name , field)
}