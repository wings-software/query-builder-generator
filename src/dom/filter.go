package dom

import (
    "fmt"
    pluralize "github.com/gertd/go-pluralize"
    "github.com/query-builder-generator/src/lang/java"
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
    case Eq: fallthrough
    case Lt: fallthrough
    case Mod:
        field = filter.FieldName
    case In:
        var pluralize = pluralize.NewClient()
        field = pluralize.Plural(filter.FieldName)
    default:
        panic(fmt.Sprintf("Unknown filter operation %+v", filter.Operation))
    }

    field = strings.Title(field)
    return fmt.Sprintf("%sQuery%s", name , field)
}

func (_ Filter) ReturnFromThis() string {
    return "return this;"
}

func (filter Filter) identifier() string {
    var identifier string
    switch filter.Operation {
    case Mod: fallthrough
    case Eq: fallthrough
    case Lt:
        identifier = filter.FieldName
    case In:
        var pluralize = pluralize.NewClient()
        identifier = pluralize.Plural(filter.FieldName)
    default:
        panic(fmt.Sprintf("Unknown filter operation %+v", filter.Operation))
    }

    // make sure the first letter is lower case
    return identifier
}

func (filter Filter) MethodName() string {
    name := filter.FieldName
    switch filter.Operation {
    case Eq:
    case Lt:
        name += "LessThan"
    case In:
        name += "In"
    case Mod:
        name += "Module"
    default:
        panic(fmt.Sprintf("Unknown filter operation %+v", filter.Operation))
    }
    return name
}

func (filter Filter) MethodArguments() string {
    var pattern string
    switch filter.Operation {
    case Eq: fallthrough
    case Lt:
        pattern = "%s %s"
    case In:
        pattern = "Iterable<%s> %s"
    case Mod:
        return "long divisor, long remainder"
    default:
        panic(fmt.Sprintf("Unknown filter operation %+v", filter.Operation))
    }
    return fmt.Sprintf(pattern, filter.FieldType, filter.identifier())
}

func (filter Filter) MethodPrototype() string {
    return fmt.Sprintf("%s(%s)", filter.MethodName(), filter.MethodArguments())
}

const methodBodyTemplate = `{
      query.field(%sKeys.%s).%s(%s);%s
    }`

const methodBodyModuleTemplate = `{
      query.field(%sKeys.%s).mod(divisor, remainder);%s
    }`

func (filter Filter) MethodBody(returning java.Interface) string {
    returnInterface := returning.ReturnFromThis()
    if len(returnInterface) > 0 {
        returnInterface = "\n      " + returnInterface
    }

    var method string
    switch filter.Operation {
    case Eq:
        method = "equal"
    case In:
        method = "in"
    case Lt:
        method = "lessThan"
    case Mod:
        return fmt.Sprintf(methodBodyModuleTemplate, filter.Query.CollectionName(), filter.FieldName, returnInterface)
    default:
        panic(fmt.Sprintf("Unknown filter operation %+v", filter.Operation))
    }
    return fmt.Sprintf(methodBodyTemplate, filter.Query.CollectionName(), filter.FieldName, method, filter.identifier(), returnInterface)
}
