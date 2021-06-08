package compiler

import (
	"github.com/query-builder-generator/src/dom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSanity(t *testing.T) {
	query := dom.Query{
		Name: "Select",
		Collection: "io.harness.beans.DelegateTask",
		Filters: []dom.Filter{
			{FieldType: "String", FieldName: "accountId"},
			{FieldType: "String", FieldName: "uuid"},
		},
	}

	compiler := Compiler{}

	expected := `package io.harness.beans;

import io.harness.beans.DelegateTask;
import io.harness.beans.DelegateTask.DelegateTaskKeys;
import io.harness.persistence.HPersistence;
import io.harness.query.PersistentQuery;
import org.mongodb.morphia.query.Query;

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

    public SelectQueryUuid accountId(String accountId) {
      query.filter(DelegateTaskKeys.accountId, accountId);
      return this;
    }

    public SelectQueryFinal uuid(String uuid) {
      query.filter(DelegateTaskKeys.uuid, uuid);
      return this;
    }

    public Query<DelegateTask> query() {
      return query;
    }
  }
}
`

	result := compiler.Generate(&query)
	assert.Equal(t, expected, result)


	query1 := dom.Query{
		Name: "Select",
		Collection: "io.harness.beans.DelegateTask",
		Filters: []dom.Filter{
			{FieldType: "String", FieldName: "accountId"},
			{FieldType: "String", FieldName: "uuid", Operation: dom.In},
		},
	}
	expected1 := `package io.harness.beans;

import io.harness.beans.DelegateTask;
import io.harness.beans.DelegateTask.DelegateTaskKeys;
import io.harness.persistence.HPersistence;
import io.harness.query.PersistentQuery;
import org.mongodb.morphia.query.Query;

public class DelegateTaskSelectQuery implements PersistentQuery {
  public static SelectQueryAccountId create(HPersistence persistence) {
    return new QueryImpl(persistence.createQuery(DelegateTask.class));
  }

  public interface SelectQueryAccountId {
    SelectQueryUuid accountId(String accountId);
  }
  public interface SelectQueryUuids {
    SelectQueryFinal uuids(Iterable<String> uuids);
  }
  public interface SelectQueryFinal {
    Query<DelegateTask> query();
  }

  private static class QueryImpl implements SelectQueryAccountId, SelectQueryUuids, SelectQueryFinal {
    Query<DelegateTask> query;

    private QueryImpl(Query<DelegateTask> query) {
      this.query = query;
    }

    public SelectQueryUuid accountId(String accountId) {
      query.filter(DelegateTaskKeys.accountId, accountId);
      return this;
    }

    public SelectQueryFinal uuids(Iterable<String> uuids) {
      query.field(DelegateTaskKeys.uuid).in(uuids);
      return this;
    }

    public Query<DelegateTask> query() {
      return query;
    }
  }
}
`
	result = compiler.Generate(&query1)
	assert.Equal(t, expected1, result)
}
