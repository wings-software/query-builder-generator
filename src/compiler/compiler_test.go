package compiler

import (
	"github.com/query-builder-generator/src/dom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSanity1(t *testing.T) {
	query := dom.Query{
		Name:       "Select",
		Collection: "io.harness.beans.DelegateTask",
		Filters: []dom.Filter{
			{FieldType: "String", FieldName: "accountId", Operation: dom.Eq},
			{FieldType: "String", FieldName: "uuid", Operation: dom.Eq},
		},
	}

	compiler := Compiler{}

	expected := `package io.harness.beans;

import io.harness.beans.DelegateTask;
import io.harness.beans.DelegateTask.DelegateTaskKeys;
import io.harness.persistence.HPersistence;
import io.harness.query.PersistentQuery;
import org.mongodb.morphia.query.Query;
import com.google.common.collect.ImmutableList;
import java.util.List;

public class DelegateTaskSelectQuery implements PersistentQuery {
  public static SelectQueryAccountId create(HPersistence persistence) {
    return new QueryImpl(persistence.createQuery(DelegateTask.class));
  }

  public interface SelectQueryAccountId {
    SelectQueryUuid accountId(String accountId);
  }
  public interface SelectQueryUuid {
    SelectQueryFinal uuid(String uuid);
  }
  public interface SelectQueryFinal {
    Query<DelegateTask> query();
  }

  private static class QueryImpl implements SelectQueryAccountId, SelectQueryUuid, SelectQueryFinal {
    Query<DelegateTask> query;

    private QueryImpl(Query<DelegateTask> query) {
      this.query = query;
    }

    @Override
    public SelectQueryUuid accountId(String accountId) {
      query.filter(DelegateTaskKeys.accountId, accountId);
      return this;
    }

    @Override
    public SelectQueryFinal uuid(String uuid) {
      query.filter(DelegateTaskKeys.uuid, uuid);
      return this;
    }

    @Override
    public Query<DelegateTask> query() {
      return query;
    }
  }

  @Override
  public List<String> queryCanonicalForms() {
    return ImmutableList.<String>builder()
      .add("collection(DelegateTask)"
         + "\n    .filter(accountId = <+>, uuid = <+>)")
    .build();
  }
}
`
	result := compiler.Generate(&query)
	assert.Equal(t, expected, result)
}

func TestSanity2(t *testing.T) {
	query := dom.Query{
		Name: "Select",
		Collection: "io.harness.beans.DelegateTask",
		Filters: []dom.Filter{
			{FieldType: "String", FieldName: "orange", Operation: dom.In},
			{FieldType: "String", FieldName: "worm", Operation: dom.Eq},
			{FieldType: "String", FieldName: "apple", Operation: dom.In},
			{FieldType: "String", FieldName: "banana", Operation: dom.In},
		},
		ProjectFields: []string{"foo", "bar"},
	}
	expected := `package io.harness.beans;

import io.harness.beans.DelegateTask;
import io.harness.beans.DelegateTask.DelegateTaskKeys;
import io.harness.persistence.HPersistence;
import io.harness.query.PersistentQuery;
import org.mongodb.morphia.query.Query;
import com.google.common.collect.ImmutableList;
import java.util.List;

public class DelegateTaskSelectQuery implements PersistentQuery {
  public static SelectQueryOranges create(HPersistence persistence) {
    return new QueryImpl(persistence.createQuery(DelegateTask.class)
                                    .project(DelegateTaskKeys.foo, true)
                                    .project(DelegateTaskKeys.bar, true));
  }

  public interface SelectQueryOranges {
    SelectQueryWorm oranges(Iterable<String> oranges);
  }
  public interface SelectQueryWorm {
    SelectQueryApples worm(String worm);
  }
  public interface SelectQueryApples {
    SelectQueryBananas apples(Iterable<String> apples);
  }
  public interface SelectQueryBananas {
    SelectQueryFinal bananas(Iterable<String> bananas);
  }
  public interface SelectQueryFinal {
    Query<DelegateTask> query();
  }

  private static class QueryImpl implements SelectQueryOranges, SelectQueryWorm, SelectQueryApples, SelectQueryBananas, SelectQueryFinal {
    Query<DelegateTask> query;

    private QueryImpl(Query<DelegateTask> query) {
      this.query = query;
    }

    @Override
    public SelectQueryWorm oranges(Iterable<String> oranges) {
      query.field(DelegateTaskKeys.orange).in(oranges);
      return this;
    }

    @Override
    public SelectQueryApples worm(String worm) {
      query.filter(DelegateTaskKeys.worm, worm);
      return this;
    }

    @Override
    public SelectQueryBananas apples(Iterable<String> apples) {
      query.field(DelegateTaskKeys.apple).in(apples);
      return this;
    }

    @Override
    public SelectQueryFinal bananas(Iterable<String> bananas) {
      query.field(DelegateTaskKeys.banana).in(bananas);
      return this;
    }

    @Override
    public Query<DelegateTask> query() {
      return query;
    }
  }

  @Override
  public List<String> queryCanonicalForms() {
    return ImmutableList.<String>builder()
      .add("collection(DelegateTask)"
         + "\n    .filter(orange in list<+>, worm = <+>, apple in list<+>, banana in list<+>)"
         + "\n    .project(foo, bar)")
    .build();
  }
}
`
	compiler := Compiler{}

	var result = compiler.Generate(&query)
	assert.Equal(t, expected, result)
}