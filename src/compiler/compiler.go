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
    return new QueryImpl(persistence.createQuery(%s.class)%s);
  }`

const interfaceTemplate = `
  public interface %sQuery%s {
    %sQuery%s %s(%s %s);
  }`

const interfaceOperationTemplate = `
  public interface %sQuery%s {
    %sQuery%s %s(Iterable<%s> %s);
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

const filterMethodOperatorTemplate = `
    public %sQuery%s %s(Iterable<%s> %s) {
      query.field(%sKeys.%s).%s(%s);
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
import org.mongodb.morphia.query.Query;
import com.google.common.collect.ImmutableList;
import java.util.List;`

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

	var projections strings.Builder
	if query.ProjectFields != nil {
		for _, field := range query.ProjectFields {
			projections.WriteString(fmt.Sprintf("\n                                    .project(%sKeys.%s, true)", collectionName, field))
		}
	}

	createMethod := fmt.Sprintf(createTemplate, name, strings.Title(query.Filters[0].FieldName), collectionName, projections.String())

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
		var currOperationType = query.Filters[i].Operation
		switch currOperationType {
		case dom.Eq:
			interfaces.WriteString(fmt.Sprintf(interfaceTemplate, name, currFieldNameTitle, name, strings.Title(nextFieldName),
				currFieldName, currFieldType, currFieldName))
			interfaceNames.WriteString(fmt.Sprintf("%sQuery%s, ", name, currFieldNameTitle))
		case dom.In:
			var pluralCurrentFieldName = currFieldName + "s"
			interfaces.WriteString(fmt.Sprintf(interfaceOperationTemplate, name, currFieldNameTitle+"s", name, strings.Title(nextFieldName),
				pluralCurrentFieldName, currFieldType, pluralCurrentFieldName))
			interfaceNames.WriteString(fmt.Sprintf("%sQuery%s, ", name, currFieldNameTitle+"s"))
		}

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
		var currOperationType = query.Filters[i].Operation
		switch currOperationType {
		case dom.None:
			methods.WriteString(fmt.Sprintf(filterMethodTemplate, name, strings.Title(nextFieldName), currFieldName, currFieldType, currFieldName,
				collectionName, currFieldName, currFieldName))
		case dom.In:
			var pluralCurrentFieldName = currFieldName + "s"
			methods.WriteString(fmt.Sprintf(filterMethodOperatorTemplate, name, strings.Title(nextFieldName), pluralCurrentFieldName, currFieldType,
				pluralCurrentFieldName, collectionName, currFieldName, "in", pluralCurrentFieldName))

		}
		methods.WriteString("\n")
	}

	var queryImpl = fmt.Sprintf(queryImplTemplate, interfaceNames.String(), collectionName, collectionName, methods.String(), collectionName)

	var imports = fmt.Sprintf(importsTemplate, query.Collection, query.Collection, collectionName)

	var queryCanonicalForms = fmt.Sprintf(queryCanonicalFormsTemplate, "");

	return fmt.Sprintf(generatedFileTemplate, imports, collectionName, name, createMethod, interfaces.String(), queryImpl, queryCanonicalForms)
}
