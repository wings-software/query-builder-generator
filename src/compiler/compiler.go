package compiler

import (
	"fmt"
	"github.com/query-builder-generator/src/dom"
	"strings"
)

type Compiler struct {
}

const createTemplate = `
  public static %sQuery%s create(HPersistence persistence) {
    return new QueryImpl(persistence.createQuery(%s.class));
  }`

const interfaceTemplate = `
  public interface %sQuery%s {
    %sQuery%s %s(%s %s);
  }`

const interfaceFinalTemplate = `
  public interface %sQueryFinal {
    Query<%s> query();
  }`

const filterMethodTemplate = `
    public %sQuery%s %s(%s %s) {
      query.filter(%sKeys.%s, %s);
      return this;
    }`

const queryImplTemplate = `
  private static class QueryImpl implements %s {
    Query<%s> query;

    private QueryImpl(Query<%s> query) {
      this.query = query;
    }
%s
    public Query<%s> query() {
      return query;
    }
  }`

const importsTemplate = `import %s;
import %s.%sKeys;
import io.harness.persistence.HPersistence;
import io.harness.query.PersistentQuery;
import org.mongodb.morphia.query.Query;`

const queryCanonicalFormsTemplate = `
  List<String> queryCanonicalForms() {
    return ImmutableList.<String>builder()%s.build();
  }`

const generatedFileTemplate = `package io.harness.beans;

%s

public class %s%sQuery implements PersistentQuery {%s
%s
%s
%s
}
`

func (compiler *Compiler) collectionName(collection string) string {
	ss := strings.Split(collection, ".")
	return ss[len(ss)-1]
}

func (compiler *Compiler) Generate(query *dom.Query) string {
	fmt.Println("Generating Java File")

	var name = query.Name
	var collectionName = compiler.collectionName(query.Collection)

	createMethod := fmt.Sprintf(createTemplate, name, strings.Title(query.Filters[0].FieldName), collectionName)

	var interfaces strings.Builder
	var interfaceNames strings.Builder

	var filtersCount = len(query.Filters)
	for i := 0; i < filtersCount; i++ {
		var nextFieldName = ""
		if i == filtersCount-1 {
			nextFieldName = "Final"
		} else {
			nextFieldName = query.Filters[i+1].FieldName
		}

		var currFieldType = query.Filters[i].FieldType
		var currFieldName = query.Filters[i].FieldName
		var currFieldNameTitle = strings.Title(currFieldName)
		interfaces.WriteString(fmt.Sprintf(interfaceTemplate, name, currFieldNameTitle, name, strings.Title(nextFieldName),
			currFieldName, currFieldType, currFieldName))
		interfaceNames.WriteString(fmt.Sprintf("%sQuery%s, ", name, currFieldNameTitle))
	}

	interfaces.WriteString(fmt.Sprintf(interfaceFinalTemplate, name, collectionName))
	interfaceNames.WriteString(fmt.Sprintf("%sQuery%s", name, "Final"))

	// Generate #3
	var methods strings.Builder
	for i := 0; i < filtersCount; i++ {
		var nextFieldName = ""
		if i == filtersCount-1 {
			nextFieldName = "Final"
		} else {
			nextFieldName = query.Filters[i+1].FieldName
		}

		var currFieldType = query.Filters[i].FieldType
		var currFieldName = query.Filters[i].FieldName
		methods.WriteString(fmt.Sprintf(filterMethodTemplate, name, strings.Title(nextFieldName), currFieldName, currFieldType, currFieldName,
			collectionName, currFieldName, currFieldName))
		methods.WriteString("\n")
	}

	var queryImpl = fmt.Sprintf(queryImplTemplate, interfaceNames.String(), collectionName, collectionName, methods.String(), collectionName)

	var imports = fmt.Sprintf(importsTemplate, query.Collection, query.Collection, collectionName)

	var queryCanonicalForms = fmt.Sprintf(queryCanonicalFormsTemplate, "");

	return fmt.Sprintf(generatedFileTemplate, imports, collectionName, name, createMethod, interfaces.String(), queryImpl, queryCanonicalForms)
}
