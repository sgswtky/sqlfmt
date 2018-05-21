package main

import (
	"testing"
)

func TestFormatAs(t *testing.T) {
	const expect = "user_id AS uid"
	left := "user_id"
	right := "uid"
	result := formatAS(left, right)
	if result != expect {
		t.Fatalf(expectFmt(expect, result))
	}
}

func TestFormatAND(t *testing.T) {
	const expect = "user_id IS NOT NULL\nAND user_id != 0"
	left := "user_id IS NOT NULL"
	right := "user_id != 0"
	result := formatAND(left, right)
	if result != expect {
		t.Fatalf(expectFmt(expect, result))
	}
}

func TestFormatOR(t *testing.T) {
	const expect = "user_id IS NOT NULL\n  OR user_id != 0"
	left := "user_id IS NOT NULL"
	right := "user_id != 0"
	result := formatOR(left, right)
	if result != expect {
		t.Fatalf(expectFmt(expect, result))
	}
}

func TestFormatSelect(t *testing.T) {
	cache := "sql_cache"
	columns := []string{
		"user_id",
	}
	distinct := "distinct"
	table := []string{
		"users",
	}
	// normal select
	expectNormalSelect := "SELECT\n  user_id\nFROM\n  users"
	resultNormalSelect := formatSelect("", columns, "", table, "", "", "", "", "", "")
	if resultNormalSelect != expectNormalSelect {
		t.Fatalf(expectFmt(expectNormalSelect, resultNormalSelect))
	}

	// normal select + cache
	expectCacheSelect := "SELECT SQL_CACHE\n  user_id\nFROM\n  users"
	resultCacheSelect := formatSelect(cache, columns, "", table, "", "", "", "", "", "")
	if resultCacheSelect != expectCacheSelect {
		t.Fatalf(expectFmt(expectCacheSelect, resultCacheSelect))
	}

	// normal select + distinct
	expectDistinctSelect := "SELECT\n  DISTINCT\n    user_id\nFROM\n  users"
	resultDistinctSelect := formatSelect("", columns, distinct, table, "", "", "", "", "", "")
	if resultDistinctSelect != expectDistinctSelect {
		t.Fatal(expectFmt(expectDistinctSelect, resultDistinctSelect))
	}

	// 'select 1'
	expectDualSelect := "SELECT\n  1"
	resultDualSelect := formatSelect("", []string{"1"}, "", []string{"dual"}, "", "", "", "", "", "")
	if resultDualSelect != expectDualSelect {
		t.Fatal(expectFmt(expectDualSelect, resultDualSelect))
	}

	// normal select + where
	where := "WHERE\n  user_id = 1"
	expectWhereSelect := "SELECT\n  user_id\nFROM\n  users\nWHERE\n  user_id = 1"
	resultWhereSelect := formatSelect("", columns, "", table, where, "", "", "", "", "")
	if resultWhereSelect != expectWhereSelect {
		t.Fatal(expectFmt(expectWhereSelect, resultWhereSelect))
	}

	// normal select + group by
	groupBy := "GROUP BY user_id, user_name"
	expectGroupBySelect := "SELECT\n  user_id\nFROM\n  users\nGROUP BY user_id, user_name"
	resultGroupBySelect := formatSelect("", columns, "", table, "", groupBy, "", "", "", "")
	if resultGroupBySelect != expectGroupBySelect {
		t.Fatal(expectFmt(expectGroupBySelect, resultGroupBySelect))
	}

	// normal select + group by + having
	having := "HAVING\n  is_deleted = false"
	expectHavingSelect := "SELECT\n  user_id\nFROM\n  users\nGROUP BY user_id, user_name\nHAVING\n  is_deleted = false"
	resultHavingSelect := formatSelect("", columns, "", table, "", groupBy, having, "", "", "")
	if resultHavingSelect != expectHavingSelect {
		t.Fatal(expectFmt(expectHavingSelect, resultHavingSelect))
	}

	// normal select + order by
	orderBy := "ORDER BY user_name, is_deleted"
	expectOrderBySelect := "SELECT\n  user_id\nFROM\n  users\nORDER BY user_name, is_deleted"
	resultOrderBySelect := formatSelect("", columns, "", table, "", "", "", orderBy, "", "")
	if resultOrderBySelect != expectOrderBySelect {
		t.Fatal(expectFmt(expectOrderBySelect, resultOrderBySelect))
	}

	// normal select + limit
	limit := "LIMIT 1"
	expectLimitSelect := "SELECT\n  user_id\nFROM\n  users\nLIMIT 1"
	resultLimitSelect := formatSelect("", columns, "", table, "", "", "", "", limit, "")
	if resultLimitSelect != expectLimitSelect {
		t.Fatal(expectFmt(expectLimitSelect, resultLimitSelect))
	}

	// normal select + lock
	lock := "LOCK IN SHARE MODE"
	expectLockSelect := "SELECT\n  user_id\nFROM\n  users\nLOCK IN SHARE MODE"
	resultLockSelect := formatSelect("", columns, "", table, "", "", "", "", "", lock)
	if resultLockSelect != expectLockSelect {
		t.Fatal(expectFmt(expectLockSelect, resultLockSelect))
	}
}

func TestFormatWhere(t *testing.T) {
	const expect = "WHERE\n  user_id = 1\n  AND is_deleted = false"
	whereContent := "user_id = 1\nAND is_deleted = false"
	result := formatWhere(whereContent)
	if result != expect {
		t.Fatal(expectFmt(expect, result))
	}
}

func TestFormatParenthesis(t *testing.T) {
	const expect = "(\n  paren1,\n  paren2\n)"
	result := formatParenthesis("  paren1,\n  paren2")
	if result != expect {
		t.Fatal(expectFmt(expect, result))
	}
}

func TestFormatIS(t *testing.T) {
	const expect = "user_name IS NOT NULL"
	result := formatIS("user_name", "IS NOT NULL")
	if result != expect {
		t.Fatal(expectFmt(expect, result))
	}
}

func TestFormatBetween(t *testing.T) {
	const expect = "col BETWEEN 1 AND 100"
	result := formatBetween("col", "1", "100", "between")
	if result != expect {
		t.Fatal(expectFmt(expect, result))
	}
}

func TestFormatTuple(t *testing.T) {
	const expect = "(value1, value2)"
	result := formatTuple("value1", "value2")
	if result != expect {
		t.Fatal(expectFmt(expect, result))
	}
}

func TestFormatExists(t *testing.T) {
	const expect = "EXISTS (\nSELECT\n  *\nFROM\n  users\n)"
	result := formatExists("SELECT\n  *\nFROM\n  users")
	if result != expect {
		t.Fatal(expectFmt(expect, result))
	}
}

func TestFormatSubQuery(t *testing.T) {
	const expect = "(\n  SELECT\n    *\n  FROM\n    users\n)"
	result := formatSubquery("SELECT\n  *\nFROM\n  users")
	if result != expect {
		t.Fatal(expectFmt(expect, result))
	}
}

func TestFormatConvertTypeQualifierFormat(t *testing.T) {
	const expect = "CONVERT(%s, DECIMAL)"
	result := formatConvertTypeQualifierFormat("DECIMAL")
	if result != expect {
		t.Fatal(expectFmt(expect, result))
	}
}

func TestFormatUnary(t *testing.T) {
	const expect = "OP(val)"
	result := formatUnary("op", "val")
	if result != expect {
		t.Fatal(expectFmt(expect, result))
	}
}

func TestFormatWhen(t *testing.T) {
	const expect = "WHEN user_type = 1 THEN ADMIN"
	result := formatWhen("user_type = 1", "ADMIN")
	if result != expect {
		t.Fatal(expectFmt(expect, result))
	}
}

func TestFormatCase(t *testing.T) {
	const expect = "CASE \n  WHEN user_type = 1 THEN ADMIN\n  ELSE NORMAL\n  END"
	result := formatCase("", "NORMAL", []string{"WHEN user_type = 1 THEN ADMIN"})
	if result != expect {
		t.Fatal(expectFmt(expect, result))
	}
}

func TestFormatJoin(t *testing.T) {
	const expect = "user AS u\nINNER JOIN user_rank AS ur\n  ON u.user_id = ur.user_id"
	result := formatJoin("INNER JOIN", "user AS u", "user_rank AS ur", "u.user_id = ur.user_id")
	if result != expect {
		t.Fatal(expectFmt(expect, result))
	}
}

func TestFormatGroupBy(t *testing.T) {
	const expect = "GROUP BY col1, col2, col3"
	result := formatGroupBy([]string{"col1", "col2", "col3"})
	if result != expect {
		t.Fatal(expectFmt(expect, result))
	}
}

func TestFormatHaving(t *testing.T) {
	const expect = "HAVING\n  user_id = 1"
	result := formatHaving("user_id = 1")
	if result != expect {
		t.Fatal(expectFmt(expect, result))
	}
}

func TestFormatLimit(t *testing.T) {
	const expectLimitAndOffset = "LIMIT 10 OFFSET 3"
	resultLimitAndOffset := formatLimit("10", "3")
	if resultLimitAndOffset != expectLimitAndOffset {
		t.Fatal(expectFmt(resultLimitAndOffset, expectLimitAndOffset))
	}

	const expectLimit = "LIMIT 10"
	resultLimit := formatLimit("10", "")
	if resultLimit != expectLimit {
		t.Fatal(expectFmt(resultLimit, expectLimit))
	}
}

func TestFormatOrderBy(t *testing.T) {
	const expect = "ORDER BY col1, col2, col3 DESC"
	orderBys := []*tuple2String{
		{
			asc,
			"col1",
		},
		{
			asc,
			"col2",
		},
		{
			desc,
			"col3",
		},
	}
	result := formatOrderBy(orderBys)
	if result != expect {
		t.Fatal(expectFmt(result, expect))
	}
}

func TestFormatTableColumn(t *testing.T) {
	const expect = "users_a.user_id"
	result := formatTableColumn("users_a", "user_id")
	if result != expect {
		t.Fatal(expectFmt(result, expect))
	}
}

func TestFormatUnion(t *testing.T) {
	const expect = "SELECT\n  *\nFROM\n  users_a\nUNION ALL\nSELECT\n  *\nFROM\n  users_b"
	result := formatUnion("SELECT\n  *\nFROM\n  users_a", "SELECT\n  *\nFROM\n  users_b", "UNION ALL")
	if result != expect {
		t.Fatal(expectFmt(result, expect))
	}
}

func TestFormatSimpleArray(t *testing.T) {
	const expect = "(\n  content1,\n  content2\n)"
	result := formatSimpleArray([]string{"content1", "content2"})
	if result != expect {
		t.Fatal(expectFmt(result, expect))
	}
}

func TestFormatFuncs(t *testing.T) {
	// Distinct
	const expectDistinct = "SUM(DISTINCT login_count)"
	resultDistinct := formatFuncs(sum, "login_count", true)
	if resultDistinct != expectDistinct {
		t.Fatal(expectFmt(resultDistinct, expectDistinct))
	}

	// not Distinct
	const expect = "SUM(login_count)"
	result := formatFuncs(sum, "(login_count)", false)
	if result != expect {
		t.Fatal(expectFmt(result, expect))
	}
}

func TestFormatGroupConcat(t *testing.T) {
	const expect = "GROUP_CONCAT(col1,\ncol2)"
	result := formatGroupConcat("", []string{"col1", "col2"}, "", "")
	if result != expect {
		t.Fatal(expectFmt(result, expect))
	}
	// Distinct
	const expectDistinct = "GROUP_CONCAT(DISTINCT col1,\ncol2)"
	resultDistinct := formatGroupConcat("distinct", []string{"col1", "col2"}, "", "")
	if resultDistinct != expectDistinct {
		t.Fatal(expectFmt(resultDistinct, expectDistinct))
	}

	// OrderBy
	const expectOrderBy = "GROUP_CONCAT(\n  col1,\n  col2 ORDER BY col1\n)"
	resultOrderBy := formatGroupConcat("", []string{"col1", "col2"}, "ORDER BY col1", "")
	if resultOrderBy != expectOrderBy {
		t.Fatal(expectFmt(resultOrderBy, expectOrderBy))
	}

	// Separator
	const expectSeparator = "GROUP_CONCAT(\n  col1,\n  col2 SEPARATOR ' '\n)"
	resultSeparator := formatGroupConcat("", []string{"col1", "col2"}, "", "separator ' '")
	if resultSeparator != expectSeparator {
		t.Fatal(expectFmt(resultSeparator, expectOrderBy))
	}

	// all
	const expectAll = "GROUP_CONCAT(\n  DISTINCT col1,\n  col2 ORDER BY col1 SEPARATOR ' '\n)"
	resultAll := formatGroupConcat("distinct", []string{"col1", "col2"}, "ORDER BY col1", "separator ' '")
	if resultAll != expectAll {
		t.Fatal(expectFmt(resultAll, expectAll))
	}
}

func TestFormatMatch(t *testing.T) {
	const expect = "MATCH (\n  user_comment_a,\n  user_comment_b\n) AGAINST(xxx in boolean mode)"
	result := formatMatch([]string{"user_comment_a", "user_comment_b"}, "xxx", "in boolean mode")
	if result != expect {
		t.Fatal(expectFmt(result, expect))
	}

	const expectSingle = "MATCH (user_comment) AGAINST(xxx in boolean mode)"
	resultSingle := formatMatch([]string{"user_comment"}, "xxx", "in boolean mode")
	if resultSingle != expectSingle {
		t.Fatal(expectFmt(resultSingle, expectSingle))
	}
}

func TestFormatCollate(t *testing.T) {
	const expect = "k COLLATE latin1_german2_ci"
	result := formatCollate("k", "latin1_german2_ci")
	if result != expect {
		t.Fatal(expectFmt(result, expect))
	}
}

func TestFormatConvertUsing(t *testing.T) {
	const expect = "CONVERT(abc USING utf8)"
	result := formatConvertUsing("abc", "utf8")
	if result != expect {
		t.Fatal(expectFmt(result, expect))
	}
}

func TestFormatTable(t *testing.T) {
	const expect = "users USE INDEX (UID_IDX,UID_CIDX)"
	result := formatTable("users", "USE", []string{"UID_IDX", "UID_CIDX"})
	if result != expect {
		t.Fatal(expectFmt(result, expect))
	}

	const expectNoHints = "users"
	resultNoHints := formatTable("users", "", []string{})
	if resultNoHints != expectNoHints {
		t.Fatal(expectFmt(resultNoHints, expectNoHints))
	}
}

func TestFormatAsTable(t *testing.T) {
	const expect = "users AS u USE INDEX (UID_IDX)"
	result := formatAsTable("users", "u", "USE", []string{"UID_IDX"})
	if result != expect {
		t.Fatal(expectFmt(result, expect))
	}
}

func TestFormatDBTable(t *testing.T) {
	const expect = "db_name.table_name"
	result := formatDBTable("db_name", "table_name")
	if result != expect {
		t.Fatal(expectFmt(result, expect))
	}

	const expectNoDB = "table_name"
	resultNoDB := formatDBTable("", "table_name")
	if resultNoDB != expectNoDB {
		t.Fatal(expectFmt(resultNoDB, expectNoDB))
	}
}
