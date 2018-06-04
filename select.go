package main

import (
	"github.com/sgswtky/sqlparser"
	"strings"
)

func (b *BuilderStruct) stmtSelect(stmt *sqlparser.Select) string {

	selectColumns := make([]string, 0)
	froms := make([]string, 0)
	wheres := ""
	groupBy := ""
	having := ""
	orderBy := ""
	limit := ""
	lock := ""

	// Columns
	selectColumns = b.columns(stmt.SelectExprs)

	// Table
	for _, v := range stmt.From {
		froms = append(froms, b.tableExpr(v))
	}

	//Where
	if stmt.Where != nil {
		wheres = b.where(stmt.Where)
	}

	//TODO?? Comments, Hints

	//GroupBy
	if stmt.GroupBy != nil {
		groupBy = b.groupBy(stmt.GroupBy)
	}
	//Having
	if stmt.Having != nil {
		having = b.having(stmt.Having)
	}
	//OrderBy
	if stmt.OrderBy != nil {
		orderBy = b.orderBy(stmt.OrderBy)
	}
	//Limit
	if stmt.Limit != nil {
		limit = b.limit(stmt.Limit)
	}
	//Lock
	if stmt.Lock != "" {
		lock = b.lock(stmt.Lock)
	}

	return formatSelect(stmt.Cache, selectColumns, stmt.Distinct, froms, wheres, groupBy, having, orderBy, limit, lock)
}

func (b *BuilderStruct) columns(exprs sqlparser.SelectExprs) []string {
	selectColumns := make([]string, 0)
	for _, column := range exprs {
		switch parsedColumn := column.(type) {
		case *sqlparser.AliasedExpr:
			exprStr := b.expr(parsedColumn.Expr)
			if parsedColumn.As.IsEmpty() {
				selectColumns = append(selectColumns, exprStr)
			} else {
				selectColumns = append(selectColumns, formatAS(exprStr, parsedColumn.As.String()))
			}
		default:
			selectColumns = append(selectColumns, "*")
		}
	}
	return selectColumns
}

func (b *BuilderStruct) where(wheres *sqlparser.Where) string {
	return formatWhere(b.expr(wheres.Expr))
}

func (b *BuilderStruct) groupBy(g sqlparser.GroupBy) string {
	result := make([]string, 0)
	for _, v := range g {
		result = append(result, b.expr(v))
	}
	return formatGroupBy(result)
}

func (b *BuilderStruct) having(h *sqlparser.Where) string {
	return formatHaving(b.expr(h.Expr))
}

func (b *BuilderStruct) orderBy(o sqlparser.OrderBy) string {
	tuples := make([]*tuple2String, 0)
	for _, v := range o {
		tuples = append(tuples,
			&tuple2String{
				str1: v.Direction,
				str2: b.expr(v.Expr),
			})
	}
	return formatOrderBy(tuples)
}

func (b *BuilderStruct) limit(l *sqlparser.Limit) string {
	limit := ""
	offset := ""
	if l.Rowcount != nil {
		limit = b.expr(l.Rowcount)
	}
	if l.Offset != nil {
		offset = b.expr(l.Offset)
	}
	return formatLimit(limit, offset)
}

func (b *BuilderStruct) lock(l string) string {
	return strings.ToUpper(strings.TrimSpace(l))
}

func (b *BuilderStruct) selectStatement(selectStmt sqlparser.SelectStatement) string {
	switch parsedSelectStatement := selectStmt.(type) {
	case *sqlparser.Select:
		return b.stmtSelect(parsedSelectStatement)
	case *sqlparser.Union:
		return formatUnion(b.selectStatement(parsedSelectStatement.Left), b.selectStatement(parsedSelectStatement.Right), parsedSelectStatement.Type)
	case *sqlparser.ParenSelect:
		//TODO cording
		//return fmt.Sprintf("(%s)", b.selectStatement(parsedSelectStatement.Select))
	default:
		unknownType(parsedSelectStatement)
	}
	return ""
}

func (b *BuilderStruct) asterOption(str string) string {
	if str == "" {
		return "*"
	}
	return str
}
