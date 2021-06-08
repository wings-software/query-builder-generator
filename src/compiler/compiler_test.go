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
    return ImmutableList.<String>builder().build();
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
			{FieldType: "String", FieldName: "uuid", Operation: dom.In},
			{FieldType: "String", FieldName: "accountId", Operation: dom.Eq},
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
  public static SelectQueryUuids create(HPersistence persistence) {
    return new QueryImpl(persistence.createQuery(DelegateTask.class)
                                    .project(DelegateTaskKeys.foo, true)
                                    .project(DelegateTaskKeys.bar, true));
  }

  public interface SelectQueryUuids {
    SelectQueryAccountId uuids(Iterable<String> uuids);
  }
  public interface SelectQueryAccountId {
    SelectQueryFinal accountId(String accountId);
  }
  public interface SelectQueryFinal {
    Query<DelegateTask> query();
  }

  private static class QueryImpl implements SelectQueryUuids, SelectQueryAccountId, SelectQueryFinal {
    Query<DelegateTask> query;

    private QueryImpl(Query<DelegateTask> query) {
      this.query = query;
    }

    @Override
    public SelectQueryAccountId uuids(Iterable<String> uuids) {
      query.field(DelegateTaskKeys.uuid).in(uuids);
      return this;
    }

    @Override
    public SelectQueryFinal accountId(String accountId) {
      query.filter(DelegateTaskKeys.accountId, accountId);
      return this;
    }

    @Override
    public Query<DelegateTask> query() {
      return query;
    }
  }

  @Override
  public List<String> queryCanonicalForms() {
    return ImmutableList.<String>builder().build();
  }
}
`
	compiler := Compiler{}

	var result = compiler.Generate(&query)
	assert.Equal(t, expected, result)
}
