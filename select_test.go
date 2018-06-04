package main

import (
	"testing"
	"github.com/sgswtky/sqlparser"
)

// The method attached to select confirms only that the expected value is returned from the structure.
// Essential tests should be done in format if possible.

func TestColumnes(t *testing.T) {
	columns := sqlparser.SelectExprs{
		&sqlparser.AliasedExpr{
			Expr: &sqlparser.ColName{
				Metadata: nil,
				Name:     sqlparser.NewColIdent("user_id"),
				Qualifier: sqlparser.TableName{
					Name:      sqlparser.NewTableIdent(""),
					Qualifier: sqlparser.NewTableIdent(""),
				},
			},
			As: sqlparser.NewColIdent(""),
		},
		&sqlparser.AliasedExpr{
			Expr: &sqlparser.ColName{
				Metadata: nil,
				Name:     sqlparser.NewColIdent("group_user_name"),
				Qualifier: sqlparser.TableName{
					Name:      sqlparser.NewTableIdent(""),
					Qualifier: sqlparser.NewTableIdent(""),
				},
			},
			As: sqlparser.NewColIdent("user_name"),
		},
		&sqlparser.AliasedExpr{
			Expr: &sqlparser.FuncExpr{
				Qualifier: sqlparser.NewTableIdent(""),
				Name:      sqlparser.NewColIdent("count"),
				Distinct:  false,
				Exprs: sqlparser.SelectExprs{
					&sqlparser.StarExpr{
						TableName: sqlparser.TableName{
							Name:      sqlparser.NewTableIdent(""),
							Qualifier: sqlparser.NewTableIdent(""),
						},
					},
				},
			},
			As: sqlparser.NewColIdent(""),
		},
		&sqlparser.StarExpr{
			TableName: sqlparser.TableName{
				Name:      sqlparser.NewTableIdent(""),
				Qualifier: sqlparser.NewTableIdent(""),
			},
		},
	}
	result := NewBuilder("").columns(columns)
	// user_id group_user_name AS user_name COUNT(*) *
	resultCount := len(result)
	expectCount := 4
	if resultCount != expectCount {
		t.Fatal(expectFmt(expectCount, resultCount))
	}

	expect0 := "user_id"
	if result[0] != expect0 {
		t.Fatal(expectFmt(expect0, result[0]))
	}

	expect1 := "group_user_name AS user_name"
	if result[1] != expect1 {
		t.Fatal(expectFmt(expect1, result[1]))
	}

	expect2 := "COUNT(*)"
	if result[2] != expect2 {
		t.Fatal(expectFmt(expect2, result[2]))
	}

	expect3 := "*"
	if result[3] != expect3 {
		t.Fatal(expectFmt(expect3, result[3]))
	}
}

func TestAsterOption(t *testing.T) {
	builder := NewBuilder("")

	expect := "*"
	result := builder.asterOption("")
	if result != expect {
		t.Fatal(expectFmt(expect, result))
	}

	expectNonAster := "default"
	resultNonAster := builder.asterOption("default")
	if resultNonAster != expectNonAster {
		t.Fatal(expectFmt(expectNonAster, resultNonAster))
	}
}

func TestWhere(t *testing.T) {
	sampleWhere := &sqlparser.Where{
		Type: "where",
		Expr: &sqlparser.ComparisonExpr{
			Operator: "=",
			Left: &sqlparser.ColName{
				Metadata: nil,
				Name:     sqlparser.NewColIdent("user_id"),
				Qualifier: sqlparser.TableName{
					Name:      sqlparser.NewTableIdent("user"),
					Qualifier: sqlparser.NewTableIdent(""),
				},
			},
			Right: &sqlparser.SQLVal{
				Type: 1,
				Val: []uint8{
					0x31,
				},
			},
			Escape: nil,
		},
	}
	result := NewBuilder("").where(sampleWhere)
	expect := `WHERE
  user.user_id = 1`
	if result != expect {
		t.Fatal(expectFmt(expect, result))
	}
}

func TestGroupBy(t *testing.T) {
	groupBy := sqlparser.GroupBy{
		&sqlparser.ColName{
			Metadata: nil,
			Name:     sqlparser.NewColIdent("user_id"),
			Qualifier: sqlparser.TableName{
				Name:      sqlparser.NewTableIdent(""),
				Qualifier: sqlparser.NewTableIdent(""),
			},
		},
		&sqlparser.ColName{
			Metadata: nil,
			Name:     sqlparser.NewColIdent("user_name"),
			Qualifier: sqlparser.TableName{
				Name:      sqlparser.NewTableIdent(""),
				Qualifier: sqlparser.NewTableIdent(""),
			},
		},
		&sqlparser.ColName{
			Metadata: nil,
			Name:     sqlparser.NewColIdent("user_status"),
			Qualifier: sqlparser.TableName{
				Name:      sqlparser.NewTableIdent(""),
				Qualifier: sqlparser.NewTableIdent(""),
			},
		},
	}
	result := NewBuilder("").groupBy(groupBy)
	expect := "GROUP BY user_id, user_name, user_status"
	if result != expect {
		t.Fatal(expectFmt(expect, result))
	}
}

func TestHaving(t *testing.T) {
	having := &sqlparser.Where{
		Type: "having",
		Expr: &sqlparser.ComparisonExpr{
			Operator: ">",
			Left: &sqlparser.ColName{
				Metadata: nil,
				Name:     sqlparser.NewColIdent("lc"),
				Qualifier: sqlparser.TableName{
					Name:      sqlparser.NewTableIdent(""),
					Qualifier: sqlparser.NewTableIdent(""),
				},
			},
			Right: &sqlparser.SQLVal{
				Type: 1,
				Val: []uint8{
					0x30,
				},
			},
			Escape: nil,
		},
	}
	result := NewBuilder("").having(having)
	expect := `HAVING
  lc > 0`
	if result != expect {
		t.Fatal(expectFmt(expect, result))
	}
}

func TestOrderBy(t *testing.T) {
	orderBy := sqlparser.OrderBy{
		&sqlparser.Order{
			Expr: &sqlparser.ColName{
				Metadata: nil,
				Name:     sqlparser.NewColIdent("login_count"),
				Qualifier: sqlparser.TableName{
					Name:      sqlparser.NewTableIdent(""),
					Qualifier: sqlparser.NewTableIdent(""),
				},
			},
			Direction: "asc",
		},
		&sqlparser.Order{
			Expr: &sqlparser.ColName{
				Metadata: nil,
				Name:     sqlparser.NewColIdent("user_id"),
				Qualifier: sqlparser.TableName{
					Name:      sqlparser.NewTableIdent(""),
					Qualifier: sqlparser.NewTableIdent(""),
				},
			},
			Direction: "asc",
		},
	}
	result := NewBuilder("").orderBy(orderBy)
	expect := "ORDER BY login_count, user_id"
	if result != expect {
		t.Fatal(expectFmt(expect, result))
	}

	orderByDesc := sqlparser.OrderBy{
		&sqlparser.Order{
			Expr: &sqlparser.ColName{
				Metadata: nil,
				Name:     sqlparser.NewColIdent("login_count"),
				Qualifier: sqlparser.TableName{
					Name:      sqlparser.NewTableIdent(""),
					Qualifier: sqlparser.NewTableIdent(""),
				},
			},
			Direction: "desc",
		},
		&sqlparser.Order{
			Expr: &sqlparser.ColName{
				Metadata: nil,
				Name:     sqlparser.NewColIdent("user_id"),
				Qualifier: sqlparser.TableName{
					Name:      sqlparser.NewTableIdent(""),
					Qualifier: sqlparser.NewTableIdent(""),
				},
			},
			Direction: "desc",
		},
	}
	orderByDescResult := NewBuilder("").orderBy(orderByDesc)
	orderByDescExpect := "ORDER BY login_count DESC, user_id DESC"
	if orderByDescResult != orderByDescExpect {
		t.Fatal(expectFmt(orderByDescExpect, orderByDescResult))
	}
}

func TestLimit(t *testing.T) {
	limit := &sqlparser.Limit{
		Offset: &sqlparser.SQLVal{
			Type: 1,
			Val: []uint8{
				0x32,
			},
		},
		Rowcount: &sqlparser.SQLVal{
			Type: 1,
			Val: []uint8{
				0x31,
			},
		},
	}
	result := NewBuilder("").limit(limit)
	expect := "LIMIT 1 OFFSET 2"
	if result != expect {
		t.Fatal(expectFmt(expect, result))
	}

	onlyLimit := &sqlparser.Limit{
		Rowcount: &sqlparser.SQLVal{
			Type: 1,
			Val: []uint8{
				0x31,
			},
		},
	}
	onlyLimitResult := NewBuilder("").limit(onlyLimit)
	onlyLimitExpect := "LIMIT 1"
	if onlyLimitResult != onlyLimitExpect {
		t.Fatal(expectFmt(onlyLimitExpect, onlyLimitResult))
	}
}

func TestLock(t *testing.T) {
	lock := "    lock in share mode   "
	result := NewBuilder("").lock(lock)
	expect := "LOCK IN SHARE MODE"
	if result != expect {
		t.Fatal(expectFmt(expect, result))
	}
}

func TestSelect(t *testing.T) {
	result, err := NewBuilder(`select count(*) as c from users_a, users_b
where user_type = 3 
group by user_product 
having c < 10 
order by department 
limit 100 
lock in share mode`).
		Parse()
	if err != nil {
		t.Fatal(expectFmt("not error.", err))
	}

	cache := ""
	columns := []string{"COUNT(*) AS c"}
	distinct := ""
	froms := []string{"users_a", "users_b"}
	wheres := "WHERE\n  user_type = 3"
	groupBy := "GROUP BY user_product"
	having := "HAVING\n  c < 10"
	orderBy := "ORDER BY department"
	limit := "LIMIT 100"
	lock := "LOCK IN SHARE MODE"
	expect := formatSelect(cache, columns, distinct, froms, wheres, groupBy, having, orderBy, limit, lock)

	if result != expect {
		t.Fatal(expectFmt(expect, result))
	}
}
