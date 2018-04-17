package parse

import (
	"github.com/sgswtky/sqlparser"
)

// Builder Structure for SQL parsing.
type Builder struct {
	initialSQL  string
	changedSQL  string
	indentLevel int
}

// NewBuilder Get the parsed SQL builder.
func NewBuilder(initialSQL string) *Builder {
	return &Builder{
		initialSQL: initialSQL,
	}
}

// Parse Parses the SQL passed at builder creation to formatted SQL.
func (b *Builder) Parse() (string, error) {
	stmt, err := sqlparser.Parse(b.initialSQL)
	if err != nil {
		return "", err
	}
	return b.statementRoot(stmt), nil
}

func (b *Builder) statementRoot(statement sqlparser.Statement) string {
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
