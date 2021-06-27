package compiler

import (
	"fmt"
	pluralize "github.com/gertd/go-pluralize"
	"github.com/query-builder-generator/src/compiler/java"
	"github.com/query-builder-generator/src/dom"
	"strings"
)

type Compiler struct {
}

const createTemplate = `
  public static %s create(HPersistence persistence) {
    return new QueryImpl(persistence.createQuery(%s.class)%s);
  }`

const interfaceTemplate = `
  public interface %s {
    %s %s;
  }`

const interfaceFinalTemplate = `
  public interface %s {
    Query<%s> query();
  }`

const filterMethodTemplate = `
    @Override
    public %s %s {
      query.filter(%sKeys.%s, %s);
      return this;
    }`

const filterMethodOperatorTemplate = `
    @Override
    public %s %s {
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
    @Override
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
  @Override
  public List<String> queryCanonicalForms() {
    return ImmutableList.<String>builder()%s    .build();
  }`

const canonicalFormTemplate=`
      .add("collection(%s)"
         + "\n    .filter(%s)"%s)
`

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

	var pluralize = pluralize.NewClient()

	var name = query.Name
	var collectionName = compiler.collectionName(query.Collection)

	var projections strings.Builder
	if query.ProjectFields != nil {
		for _, field := range query.ProjectFields {
			projections.WriteString(fmt.Sprintf(
				"\n                                    .project(%sKeys.%s, true)", collectionName, field))
		}
	}

	var titleFieldName = strings.Title(query.Filters[0].FieldName)
	if query.Filters[0].Operation == dom.In {
		titleFieldName = pluralize.Plural(titleFieldName)
	}

	createMethod := fmt.Sprintf(createTemplate, query.Filters[0].InterfaceName(), collectionName, projections.String())

	var interfaces strings.Builder
	var interfaceNames strings.Builder

	// Generate #3
	var methods strings.Builder
	for i := range query.Filters {
		var nextInterface java.Interface

		if i == len(query.Filters)-1 {
			nextInterface = query
		} else {
			nextInterface = query.Filters[i+1]
		}

		var currentInterface java.Interface
		currentInterface = query.Filters[i]

		var currentMethod java.Method
		currentMethod = query.Filters[i]

		interfaceNames.WriteString(currentInterface.InterfaceName())
		interfaceNames.WriteString(", ")

		interfaces.WriteString(fmt.Sprintf(interfaceTemplate,
			currentInterface.InterfaceName(), nextInterface.InterfaceName(),
			currentMethod.MethodPrototype()))

		var currFieldName = query.Filters[i].FieldName
		var currOperationType = query.Filters[i].Operation
		switch currOperationType {
		case dom.Eq:
			methods.WriteString(fmt.Sprintf(filterMethodTemplate,
				nextInterface.InterfaceName(), currentMethod.MethodPrototype(),
				collectionName, currFieldName, currFieldName))
		case dom.In:
			var pluralCurrentFieldName = pluralize.Plural(currFieldName)
			methods.WriteString(fmt.Sprintf(filterMethodOperatorTemplate,
				nextInterface.InterfaceName(), currentMethod.MethodPrototype(),
				collectionName, currFieldName, "in", pluralCurrentFieldName))

		}
		methods.WriteString("\n")
	}

	interfaces.WriteString(fmt.Sprintf(interfaceFinalTemplate, query.InterfaceName(), collectionName))
	interfaceNames.WriteString(query.InterfaceName())

	var queryImpl = fmt.Sprintf(queryImplTemplate, interfaceNames.String(), collectionName, collectionName, methods.String(), collectionName)

	var imports = fmt.Sprintf(importsTemplate, query.Collection, query.Collection, collectionName)

	var canonicalExpression strings.Builder
	for _, filter := range query.Filters {
		if len(canonicalExpression.String()) != 0 {
			canonicalExpression.WriteString(", ")
		}
		var currFieldName = filter.FieldName
		var currOperationType = filter.Operation

		switch currOperationType {
		case dom.Eq:
			canonicalExpression.WriteString(currFieldName + " = <+>")
		case dom.In:
			canonicalExpression.WriteString(currFieldName + " in list<+>")
		}
	}

	var canonicalProjections strings.Builder
	if query.ProjectFields != nil && len(query.ProjectFields) !=0 {
		for _, field := range query.ProjectFields {
			if len(canonicalProjections.String()) != 0 {
				canonicalProjections.WriteString(", ")
			} else {
				canonicalProjections.WriteString("\n         + \"\\n    .project(")
			}
			canonicalProjections.WriteString(field)
		}
		canonicalProjections.WriteString(")\"")
	}


	var queryCanonicalForms = fmt.Sprintf(queryCanonicalFormsTemplate, fmt.Sprintf(canonicalFormTemplate, collectionName, canonicalExpression.String(), canonicalProjections.String()))

	return fmt.Sprintf(generatedFileTemplate, imports, collectionName, name, createMethod, interfaces.String(), queryImpl, queryCanonicalForms)
}
