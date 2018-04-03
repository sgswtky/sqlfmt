package parse

import (
	"github.com/sgswtky/sqlparser"
)

type builder struct {
	initialSQL  string
	changedSQL  string
	indentLevel int
}

func NewBuilder(initialSQL string) *builder {
	return &builder{
		initialSQL: initialSQL,
	}
}

func (b *builder) Parse() (string, error) {
	stmt, err := sqlparser.Parse(b.initialSQL)
	if err != nil {
		return "", err
	}
	return b.statementRoot(stmt), nil
}

func (b *builder) statementRoot(statement sqlparser.Statement) string {
	switch parsedStmt := statement.(type) {
	case *sqlparser.Select:
		return b.stmtSelect(parsedStmt)
	case *sqlparser.Union:
		return b.selectStatement(parsedStmt)
	default:
		unknownType(parsedStmt)
	}
	return ""
}
