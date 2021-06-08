package dom

const (
    EQUAL = iota
    IN
)

type Filter struct {
    FieldType string
    FieldName string
    Operation int
}