package main

import (
	"github.com/sgswtky/sqlparser"
)

// Builder Structure for SQL parsing.
type BuilderStruct struct {
	initialSQL  string
	changedSQL  string
	indentLevel int
}

type Builder interface {
	Parse() (string, error)
	statementRoot(statement sqlparser.Statement) string
	expr(e sqlparser.Expr) string
	getFuncExpr(funcExpr *sqlparser.FuncExpr) string
	stmtSelect(stmt *sqlparser.Select) string
	columns(exprs sqlparser.SelectExprs) []string
	where(wheres *sqlparser.Where) string
	groupBy(g sqlparser.GroupBy) string
	having(h *sqlparser.Where) string
	orderBy(o sqlparser.OrderBy) string
	limit(l *sqlparser.Limit) string
	lock(l string) string
	tableExpr(expr sqlparser.TableExpr) string
	simpleTableExpr(expr sqlparser.SimpleTableExpr) string
	selectStatement(selectStmt sqlparser.SelectStatement) string
	selectExprs(selectExprs sqlparser.SelectExprs) []string
	selectExpr(selectExpr sqlparser.SelectExpr) string
	asterOption(str string) string
}

// NewBuilder Get the parsed SQL builder.
func NewBuilder(initialSQL string) Builder {
	return &BuilderStruct{
		initialSQL: initialSQL,
	}
}

// Parse Parses the SQL passed at builder creation to formatted SQL.
func (b *BuilderStruct) Parse() (string, error) {
	stmt, err := sqlparser.Parse(b.initialSQL)
	if err != nil {
		return "", err
	}
	return b.statementRoot(stmt), nil
}

func (b *BuilderStruct) statementRoot(statement sqlparser.Statement) string {
	switch parsedStmt := statement.(type) {
	case *sqlparser.Select:
		return b.stmtSelect(parsedStmt)
	case *sqlparser.Union:
		return b.selectStatement(parsedStmt)
	default:
		return unsupportedType(parsedStmt)
	}
}
