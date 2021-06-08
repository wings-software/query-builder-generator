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


    //generate #1
    var collectionName = query.Collection
    var s1_template = `
  public static %sQuery%s create(HPersistence persistence) {
    return new QueryImpl(persistence.createQuery(%s.class));
  }`
    s1 := fmt.Sprintf(s1_template, collectionName, strings.Title(query.Filters[0].FieldName), collectionName)


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
    var s3_1 strings.Builder //
    var s2_1 strings.Builder
    var filtersCount = len(query.Filters)
    for i := 0; i < filtersCount; i++ {
        var nextFieldName = "";
        if i == filtersCount-1 {
            nextFieldName = "Final"
        } else {
            nextFieldName = query.Filters[i+1].FieldName
        }

        var currFieldType = query.Filters[i].FieldType
        var currFieldName = query.Filters[i].FieldName
        s2_1.WriteString(fmt.Sprintf(s2_template_1, collectionName, strings.Title(currFieldName), collectionName, strings.Title(nextFieldName), currFieldName, currFieldType, currFieldName))
        s3_1.WriteString(fmt.Sprintf("%sQuery%s, ", collectionName, strings.Title(currFieldName)))
    }

    s2_1.WriteString(fmt.Sprintf(s2_template_2, collectionName, collectionName))
    s3_1.WriteString(fmt.Sprintf("%sQuery%s", collectionName, "Final"))
    var s2 = s2_1.String()


    // Generate #3
    var s3_template_1=`
    public %sQuery%s %s(%s %s) {
      query.filter(%sKeys.%s, %s);
      return this;
    }`
    var s3_2 strings.Builder
    for i := 0; i < filtersCount; i++ {
        var nextFieldName = "";
        if i == filtersCount-1 {
            nextFieldName = "Final"
        } else {
            nextFieldName = query.Filters[i+1].FieldName
        }

        var currFieldType = query.Filters[i].FieldType
        var currFieldName = query.Filters[i].FieldName
        s3_2.WriteString(fmt.Sprintf(s3_template_1, collectionName, strings.Title(nextFieldName), currFieldName, currFieldType, currFieldName, collectionName,
        currFieldName, currFieldName))
    }
    var s3_template_2=`
  private static class QueryImpl implements %s {
    Query<%s> query;
    private QueryImpl(Query<%s> query) {
      this.query = query;
    }%s
    public Query<%s> query() {
      return query;
    }
  }
    `
    var s3 = fmt.Sprintf(s3_template_2, s3_1.String(), collectionName, collectionName, s3_2.String(), collectionName)



    var generatedFile_template=`
package io.harness.beans;
import io.harness.beans.%s.%sKeys;
import io.harness.persistence.HPersistence;
import io.harness.query.PersistentQuery;
import org.mongodb.morphia.query.Query;

public class %sQuery implements PersistentQuery {%s%s
%s
}
    `
    fmt.Println(fmt.Sprintf(generatedFile_template, collectionName, collectionName, collectionName, s2, s1, s3))
}