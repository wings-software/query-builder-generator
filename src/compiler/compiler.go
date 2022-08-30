package compiler

import (
	"fmt"
	"math"
	"strings"

	pluralize "github.com/gertd/go-pluralize"
	"github.com/query-builder-generator/src/dom"
	"github.com/query-builder-generator/src/lang/java"
)

type Compiler struct {
}

const create1Template = `
  public static %s create(HPersistence persistence) {
    return new QueryImpl(persistence.createQuery(%s.class)%s);
  }`

const create2Template = `
  public static %s create(HPersistence persistence, Set<QueryChecks> queryChecks) {
    return new QueryImpl(persistence.createQuery(%s.class, queryChecks)%s);
  }`

const interfaceTemplate = `
  public interface %s {
    %s %s;
  }`

const interfaceFinalTemplate = `
  public interface %s {
    Query<%s> query();%s
  }`

const optionMethod = `
    default %s %s()  {
      return (%s) this;
    }`

const methodTemplate = `
    @Override
    public %s %s %s`

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
import io.harness.persistence.HQuery.QueryChecks;
import io.harness.query.PersistentQuery;

import com.google.common.collect.ImmutableList;
import java.util.List;
import java.util.Set;
import dev.morphia.query.Query;`

const queryCanonicalFormsTemplate = `
  @Override
  public List<String> queryCanonicalForms() {
    return ImmutableList.<String>builder()%s
      .build();
  }`

const canonicalFormTemplate = `
      .add("collection(%s)"
         + "\n    .filter(%s)"%s)`

const generatedFileTemplate = `package %s;

%s

public class %s%sQuery implements PersistentQuery {%s
%s
%s
%s
}
`

func powInt(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}

func (compiler *Compiler) Generate(document *dom.Document) string {
	fmt.Println("Generating Java File")

	var pluralize = pluralize.NewClient()

	var query = document.Queries[0]

	var name = query.Name
	var collectionName = query.CollectionName()

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

	createMethod := fmt.Sprintf(create1Template, query.Filters[0].InterfaceName(), collectionName, projections.String())
	createMethod += fmt.Sprintf(create2Template, query.Filters[0].InterfaceName(), collectionName, projections.String())

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

		methods.WriteString(fmt.Sprintf(methodTemplate,
			nextInterface.InterfaceName(),
			currentMethod.MethodPrototype(),
			currentMethod.MethodBody(nextInterface)))
		methods.WriteString("\n")
	}

	for _, optinal := range query.Optionals {
		for i := range optinal.Filters {
			var nextInterface java.Interface

			if i == len(optinal.Filters)-1 {
				nextInterface = java.Void{}
			} else {
				nextInterface = optinal.Filters[i+1]
			}

			var currentInterface java.Interface
			currentInterface = optinal.Filters[i]

			var currentMethod java.Method
			currentMethod = optinal.Filters[i]

			interfaceNames.WriteString(currentInterface.InterfaceName())
			interfaceNames.WriteString(", ")

			interfaces.WriteString(fmt.Sprintf(interfaceTemplate,
				currentInterface.InterfaceName(), nextInterface.InterfaceName(),
				currentMethod.MethodPrototype()))

			methods.WriteString(fmt.Sprintf(methodTemplate,
				nextInterface.InterfaceName(),
				currentMethod.MethodPrototype(),
				currentMethod.MethodBody(nextInterface)))
			methods.WriteString("\n")
		}
	}

	var options strings.Builder

	for _, optinal := range query.Optionals {
		interfaceName := optinal.Filters[0].InterfaceName()
		options.WriteString(fmt.Sprintf(optionMethod, interfaceName, optinal.Name, interfaceName))
	}

	interfaces.WriteString(fmt.Sprintf(interfaceFinalTemplate, query.InterfaceName(), collectionName, options.String()))
	interfaceNames.WriteString(query.InterfaceName())

	var queryImpl = fmt.Sprintf(queryImplTemplate, interfaceNames.String(), collectionName, collectionName, methods.String(), collectionName)

	var imports = fmt.Sprintf(importsTemplate, query.Collection, query.Collection, collectionName)

	var canonicalProjections strings.Builder
	if query.ProjectFields != nil && len(query.ProjectFields) != 0 {
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

	var canonicalExpressions strings.Builder
	for i := 0; i < powInt(2, len(query.Optionals)); i++ {
		var canonicalExpression strings.Builder
		filters := query.Filters
		flag := i
		for _, optional := range query.Optionals {
			if flag%2 == 1 {
				filters = append(filters, optional.Filters...)
			}
			flag /= 2
		}

		for _, filter := range filters {
			if len(canonicalExpression.String()) != 0 {
				canonicalExpression.WriteString(", ")
			}
			var currFieldName = filter.FieldName
			var currOperationType = filter.Operation

			switch currOperationType {
			case dom.Eq:
				canonicalExpression.WriteString(currFieldName + " == @")
			case dom.Lt:
				canonicalExpression.WriteString(currFieldName + " < @")
			case dom.In:
				canonicalExpression.WriteString(currFieldName + " in [@]")
			case dom.Mod:
				canonicalExpression.WriteString(currFieldName + " % @divisor == @remainder")
			default:
				panic(fmt.Sprintf("Unknown filter operation %+v", filter.Operation))
			}
		}

		canonicalExpressions.WriteString(fmt.Sprintf(canonicalFormTemplate, collectionName, canonicalExpression.String(), canonicalProjections.String()))
	}

	var queryCanonicalForms = fmt.Sprintf(queryCanonicalFormsTemplate, canonicalExpressions.String())

	return fmt.Sprintf(generatedFileTemplate,
		document.Package, imports,
		collectionName, name, createMethod, interfaces.String(), queryImpl, queryCanonicalForms)
}
