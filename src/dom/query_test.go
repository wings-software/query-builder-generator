package dom

import (
	"testing"
	"fmt"
	"strings"
)

func TestSanity(t *testing.T) {
    query := Query{
                    Collection: "DelegateTask",
                    Filters: []Filter {
                        Filter {FieldType: "String", FieldName : "accountId"},
                        Filter {FieldType: "String", FieldName : "uuid"},
                    },
                 }

    fmt.Println("Generating Java File")

    var collectionName = query.Collection

    //generate #1
    var s1_template = `
  public static %sQueryAccountId create(HPersistence persistence) {
    return new QueryImpl(persistence.createQuery(%s.class));
  }`
    s1 := fmt.Sprintf(s1_template, collectionName, collectionName)
    fmt.Println(s1)


    // Generate #2

    var s2_template_1 = `
  public interface %sQuery%s {
    %sQuery%s %s(%s %s);
  }`
    var s2_template_2 = `
  public interface %sQueryFinal {
    Query<%s> query();
  }
    `
    var s2 strings.Builder
    var filtersCount = len(query.Filters)
    for i := 0; i < filtersCount; i++ {
        var nextFieldName = "";
        if i == filtersCount-1 {
            nextFieldName = "Final"
        } else {
            nextFieldName = query.Filters[i+1].FieldType
        }

        var currFieldType = query.Filters[i].FieldType
        var currFieldName = query.Filters[i].FieldName
        s2.WriteString(fmt.Sprintf(s2_template_1, collectionName, currFieldType, collectionName, nextFieldName, currFieldName, currFieldType, currFieldName))
    }


    s2.WriteString(fmt.Sprintf(s2_template_2, collectionName, collectionName))

    fmt.Println(s2.String())
}