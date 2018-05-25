package main

import (
	"testing"
	"github.com/sgswtky/sqlparser"
)

// The method attached to select confirms only that the expected value is returned from the structure.
// Essential tests should be done in format if possible.

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
}
