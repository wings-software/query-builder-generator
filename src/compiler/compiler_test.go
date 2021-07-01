package compiler

import (
	"github.com/query-builder-generator/src/dom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSanity1(t *testing.T) {
	document := dom.Document{
		Package: "io.harness.beans",
		Queries: []dom.Query{
			{
				Name:       "Select",
				Collection: "io.harness.beans.DelegateTask",
				Filters: []dom.Filter{
					{FieldType: "String", FieldName: "accountId", Operation: dom.Eq},
					{FieldType: "String", FieldName: "uuid", Operation: dom.Eq},
				},
			},
		},
	}
	document.Init()

	compiler := Compiler{}

	expected := `package io.harness.beans;

import io.harness.beans.DelegateTask;
import io.harness.beans.DelegateTask.DelegateTaskKeys;

import io.harness.persistence.HPersistence;
import io.harness.persistence.HQuery.QueryChecks;
import io.harness.query.PersistentQuery;

import com.google.common.collect.ImmutableList;
import java.util.List;
import java.util.Set;
import org.mongodb.morphia.query.Query;

public class DelegateTaskSelectQuery implements PersistentQuery {
  public static SelectQueryAccountId create(HPersistence persistence) {
    return new QueryImpl(persistence.createQuery(DelegateTask.class));
  }
  public static SelectQueryAccountId create(HPersistence persistence, Set<QueryChecks> queryChecks) {
    return new QueryImpl(persistence.createQuery(DelegateTask.class, queryChecks));
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
      query.field(DelegateTaskKeys.accountId).equal(accountId);
      return this;
    }

    @Override
    public SelectQueryFinal uuid(String uuid) {
      query.field(DelegateTaskKeys.uuid).equal(uuid);
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
         + "\n    .filter(accountId == @, uuid == @)")
      .build();
  }
}
`
	result := compiler.Generate(&document)
	assert.Equal(t, expected, result)
}

func TestOperators(t *testing.T) {
	document := dom.Document{
		Package: "io.harness.qbg",
		Queries: []dom.Query{
			{
				Name:       "Select",
				Collection: "io.harness.beans.DelegateTask",
				Filters: []dom.Filter{
					{FieldType: "String", FieldName: "equal", Operation: dom.Eq},
					{FieldType: "String", FieldName: "in", Operation: dom.In},
					{FieldType: "String", FieldName: "lessThan", Operation: dom.Lt},
					{FieldType: "String", FieldName: "mod", Operation: dom.Mod},
				},
			},
		},
	}
	document.Init()

	expected := `package io.harness.qbg;

import io.harness.beans.DelegateTask;
import io.harness.beans.DelegateTask.DelegateTaskKeys;

import io.harness.persistence.HPersistence;
import io.harness.persistence.HQuery.QueryChecks;
import io.harness.query.PersistentQuery;

import com.google.common.collect.ImmutableList;
import java.util.List;
import java.util.Set;
import org.mongodb.morphia.query.Query;

public class DelegateTaskSelectQuery implements PersistentQuery {
  public static SelectQueryEqual create(HPersistence persistence) {
    return new QueryImpl(persistence.createQuery(DelegateTask.class));
  }
  public static SelectQueryEqual create(HPersistence persistence, Set<QueryChecks> queryChecks) {
    return new QueryImpl(persistence.createQuery(DelegateTask.class, queryChecks));
  }

  public interface SelectQueryEqual {
    SelectQueryIns equal(String equal);
  }
  public interface SelectQueryIns {
    SelectQueryLessThan inIn(Iterable<String> ins);
  }
  public interface SelectQueryLessThan {
    SelectQueryMod lessThanLessThan(String lessThan);
  }
  public interface SelectQueryMod {
    SelectQueryFinal modModule(long divisor, long remainder);
  }
  public interface SelectQueryFinal {
    Query<DelegateTask> query();
  }

  private static class QueryImpl implements SelectQueryEqual, SelectQueryIns, SelectQueryLessThan, SelectQueryMod, SelectQueryFinal {
    Query<DelegateTask> query;

    private QueryImpl(Query<DelegateTask> query) {
      this.query = query;
    }

    @Override
    public SelectQueryIns equal(String equal) {
      query.field(DelegateTaskKeys.equal).equal(equal);
      return this;
    }

    @Override
    public SelectQueryLessThan inIn(Iterable<String> ins) {
      query.field(DelegateTaskKeys.in).in(ins);
      return this;
    }

    @Override
    public SelectQueryMod lessThanLessThan(String lessThan) {
      query.field(DelegateTaskKeys.lessThan).lessThan(lessThan);
      return this;
    }

    @Override
    public SelectQueryFinal modModule(long divisor, long remainder) {
      query.field(DelegateTaskKeys.mod).mod(divisor, remainder);
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
         + "\n    .filter(equal == @, in in [@], lessThan < @, mod % @divisor == @remainder)")
      .build();
  }
}
`
	compiler := Compiler{}

	var result = compiler.Generate(&document)
	assert.Equal(t, expected, result)
}

func TestPlurals(t *testing.T) {
	document := dom.Document{
		Package: "io.harness.qbg",
		Queries: []dom.Query{
			{
				Name:       "Select",
				Collection: "io.harness.beans.DelegateTask",
				Filters: []dom.Filter{
					{FieldType: "String", FieldName: "orange", Operation: dom.In},
					{FieldType: "String", FieldName: "worm", Operation: dom.Eq},
					{FieldType: "String", FieldName: "apple", Operation: dom.In},
					{FieldType: "String", FieldName: "banana", Operation: dom.In},
				},
				ProjectFields: []string{"foo", "bar"},
			},
		},
	}
	document.Init()

	expected := `package io.harness.qbg;

import io.harness.beans.DelegateTask;
import io.harness.beans.DelegateTask.DelegateTaskKeys;

import io.harness.persistence.HPersistence;
import io.harness.persistence.HQuery.QueryChecks;
import io.harness.query.PersistentQuery;

import com.google.common.collect.ImmutableList;
import java.util.List;
import java.util.Set;
import org.mongodb.morphia.query.Query;

public class DelegateTaskSelectQuery implements PersistentQuery {
  public static SelectQueryOranges create(HPersistence persistence) {
    return new QueryImpl(persistence.createQuery(DelegateTask.class)
                                    .project(DelegateTaskKeys.foo, true)
                                    .project(DelegateTaskKeys.bar, true));
  }
  public static SelectQueryOranges create(HPersistence persistence, Set<QueryChecks> queryChecks) {
    return new QueryImpl(persistence.createQuery(DelegateTask.class, queryChecks)
                                    .project(DelegateTaskKeys.foo, true)
                                    .project(DelegateTaskKeys.bar, true));
  }

  public interface SelectQueryOranges {
    SelectQueryWorm orangeIn(Iterable<String> oranges);
  }
  public interface SelectQueryWorm {
    SelectQueryApples worm(String worm);
  }
  public interface SelectQueryApples {
    SelectQueryBananas appleIn(Iterable<String> apples);
  }
  public interface SelectQueryBananas {
    SelectQueryFinal bananaIn(Iterable<String> bananas);
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
    public SelectQueryWorm orangeIn(Iterable<String> oranges) {
      query.field(DelegateTaskKeys.orange).in(oranges);
      return this;
    }

    @Override
    public SelectQueryApples worm(String worm) {
      query.field(DelegateTaskKeys.worm).equal(worm);
      return this;
    }

    @Override
    public SelectQueryBananas appleIn(Iterable<String> apples) {
      query.field(DelegateTaskKeys.apple).in(apples);
      return this;
    }

    @Override
    public SelectQueryFinal bananaIn(Iterable<String> bananas) {
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
         + "\n    .filter(orange in [@], worm == @, apple in [@], banana in [@])"
         + "\n    .project(foo, bar)")
      .build();
  }
}
`
	compiler := Compiler{}

	var result = compiler.Generate(&document)
	assert.Equal(t, expected, result)
}

func TestOptional(t *testing.T) {
	document := dom.Document{
		Package: "io.harness.qbg",
		Queries: []dom.Query{
			{
				Name:       "Timeout",
				Collection: "io.harness.beans.DelegateTask",
				Filters: []dom.Filter{
					{FieldType: "String", FieldName: "orange", Operation: dom.Eq},
				},
				Optionals: []dom.Optional{
					{Name: "parasite", Filters: []dom.Filter{
						{FieldType: "String", FieldName: "worm", Operation: dom.Eq},
					}},
				},
			},
		},
	}
	document.Init()

	expected := `package io.harness.qbg;

import io.harness.beans.DelegateTask;
import io.harness.beans.DelegateTask.DelegateTaskKeys;

import io.harness.persistence.HPersistence;
import io.harness.persistence.HQuery.QueryChecks;
import io.harness.query.PersistentQuery;

import com.google.common.collect.ImmutableList;
import java.util.List;
import java.util.Set;
import org.mongodb.morphia.query.Query;

public class DelegateTaskTimeoutQuery implements PersistentQuery {
  public static TimeoutQueryOrange create(HPersistence persistence) {
    return new QueryImpl(persistence.createQuery(DelegateTask.class));
  }
  public static TimeoutQueryOrange create(HPersistence persistence, Set<QueryChecks> queryChecks) {
    return new QueryImpl(persistence.createQuery(DelegateTask.class, queryChecks));
  }

  public interface TimeoutQueryOrange {
    TimeoutQueryOptions orange(String orange);
  }
  public interface TimeoutQueryWorm {
    void worm(String worm);
  }
  public interface TimeoutQueryOptions {
    Query<DelegateTask> query();
    default TimeoutQueryWorm parasite()  {
      return (TimeoutQueryWorm) this;
    }
  }

  private static class QueryImpl implements TimeoutQueryOrange, TimeoutQueryWorm, TimeoutQueryOptions {
    Query<DelegateTask> query;

    private QueryImpl(Query<DelegateTask> query) {
      this.query = query;
    }

    @Override
    public TimeoutQueryOptions orange(String orange) {
      query.field(DelegateTaskKeys.orange).equal(orange);
      return this;
    }

    @Override
    public void worm(String worm) {
      query.field(DelegateTaskKeys.worm).equal(worm);
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
         + "\n    .filter(orange == @)")
      .add("collection(DelegateTask)"
         + "\n    .filter(orange == @, worm == @)")
      .build();
  }
}
`
	compiler := Compiler{}

	var result = compiler.Generate(&document)
	assert.Equal(t, expected, result)
}
