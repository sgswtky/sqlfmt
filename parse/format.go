package parse

import (
	"fmt"
	"strings"
)

func formatAS(left, right string) string {
	return fmt.Sprintf("%s %s %s", left, as, right)
}

func formatAND(left, right string) string {
	return fmt.Sprintf("%s\n%s %s", left, and, right)
}

func formatOR(left, right string) string {
	return fmt.Sprintf("%s\n  %s %s", left, or, right)
}

func formatSelect(cache string, columns []string, distinct string, table []string, wheres, groupBy, having, orderBy, limit, lock string) (selectSQL string) {

	selectStr := ""
	if cache == "" {
		selectStr = "SELECT"
	} else {
		selectStr = fmt.Sprintf("SELECT %s", strings.ToUpper(strings.TrimSpace(cache)))
	}

	selectColumns := ""
	if distinct == "" {
		selectColumns = fmt.Sprintf("%s\n  %s", selectStr, addIndent(strings.Join(columns, ",\n")))
	} else {
		selectColumns = fmt.Sprintf("%s\n  %s\n    %s", selectStr, strings.ToUpper(distinct), addIndent(addIndent(strings.Join(columns, ",\n"))))
	}

	// sqlparser parses it as dual if there is no table specification. dual is not displayed in sqlfmt.
	if !isDual(table) {
		if len(table) > 0 {
			selectFroms := fmt.Sprintf("FROM\n  %s", addIndent(strings.Join(table, ",\n")))
			selectSQL = fmt.Sprintf("%s\n%s", selectColumns, selectFroms)
		} else {
			selectSQL = fmt.Sprintf("%s", selectColumns)
		}
	} else {
		selectSQL = selectColumns
	}

	if wheres != "" {
		selectSQL += "\n" + wheres
	}

	if groupBy != "" {
		selectSQL += "\n" + groupBy
	}

	if having != "" {
		selectSQL += "\n" + having
	}

	if orderBy != "" {
		selectSQL += "\n" + orderBy
	}

	if limit != "" {
		selectSQL += "\n" + limit
	}

	if lock != "" {
		selectSQL += "\n" + lock
	}
	return selectSQL
}

func formatWhere(where string) string {
	return fmt.Sprintf(`WHERE
%s`, linesIndent(where))
}

func formatParenthesis(s string) string {
	// TODO to only indent of content
	return fmt.Sprintf("(\n%s\n)", s)
}

func formatIS(target, operate string) string {
	return fmt.Sprintf("%s %s", target, strings.ToUpper(operate))
}

func formatBetween(left, from, to, operate string) string {
	return fmt.Sprintf("%s %s %s AND %s",
		left,
		strings.ToUpper(operate),
		from,
		to)
}

func formatTuple(values ...string) string {
	return fmt.Sprintf("(%s)", strings.Join(values, ", "))
}

func formatExists(subQuery string) string {
	return fmt.Sprintf("EXISTS (\n%s\n)", subQuery)
}

func formatSubquery(subQuery string) string {
	return fmt.Sprintf("(\n%s\n)", linesIndent(subQuery))
}

func formatConvertTypeQualifierFormat(s string) string {
	return "CONVERT(%s, " + s + ")"
}

func formatUnary(operater, value string) string {
	o := strings.TrimSpace(operater)
	v := strings.TrimSpace(value)
	return fmt.Sprintf("%s(%s)", strings.ToUpper(o), v)
}

func formatWhen(cond, val string) string {
	return fmt.Sprintf("WHEN %s THEN %s", cond, val)
}

func formatCase(caseValue, elseValue string, whens []string) string {
	whensStr := linesIndent(strings.Join(whens, "\n"))
	return fmt.Sprintf("CASE %s\n%s\n%s\n%s", caseValue, whensStr, linesIndent(fmt.Sprintf("ELSE %s", elseValue)), linesIndent("END"))
}

func formatJoin(format, leftExpr, rightStr, condition string) string {
	if condition != "" {
		condition = linesIndent("ON " + condition)
	}
	return fmt.Sprintf("%s\n%s %s\n%s", leftExpr, strings.ToUpper(format), rightStr, condition)
}

func formatGroupBy(groupBys []string) string {
	return fmt.Sprintf("GROUP BY %s", strings.Join(groupBys, ", "))
}

func formatHaving(wheres string) string {
	return fmt.Sprintf("HAVING\n  %s", wheres)
}

func formatLimit(limit, offset string) string {
	// only limit
	if limit != "" && offset == "" {
		return fmt.Sprintf("LIMIT %s", limit)
	}
	return fmt.Sprintf("LIMIT %s OFFSET %s", limit, offset)
}

func formatOrderBy(tuples []*tuple2String) string {
	result := make([]string, 0)
	for _, v := range tuples {
		if strings.ToUpper(v.str1) == asc {
			result = append(result, v.str2)
		} else {
			result = append(result, fmt.Sprintf("%s %s", v.str2, desc))
		}
	}
	return fmt.Sprintf("ORDER BY %s", strings.Join(result, ", "))
}

func formatTableColumn(table, column string) string {
	return fmt.Sprintf("%s.%s", table, column)
}

func formatUnion(left, right, union string) string {
	return fmt.Sprintf("%s\n%s\n%s", left, strings.ToUpper(union), right)
}

func formatTables(tables []string) string {
	return formatSimpleArray(tables)
}

func formatSimpleArray(array []string) string {
	return fmt.Sprintf("(\n%s\n)", linesIndent(strings.Join(array, ",\n")))
}

func formatFuncs(funcs, val string, isDistinct bool) string {
	if isDistinct {
		return fmt.Sprintf("%s(DISTINCT %s)", funcs, val)
	}
	return fmt.Sprintf("%s%s", funcs, val)
}

func formatGroupConcat(distinct string, values []string, orderBy, separator string) string {
	distinct = strings.TrimSpace(distinct)
	separator = strings.TrimSpace(separator)
	result := make([]string, 0)
	if distinct != "" {
		result = append(result, strings.ToUpper(distinct))
	}
	result = append(result, fmt.Sprintf("%s", strings.Join(values, ",\n")))
	if orderBy != "" {
		result = append(result, fmt.Sprintf("%s", orderBy))
	}
	if separator != "" {
		result = append(result, strings.ToUpper(separator))
	}
	if len(result) == 1 || (len(result) == 2 && distinct != "") {
		// this case is dont create new line
		// case "GROUP_CONCAT(column_name) or GROUP_CONCAT(DISTINCT column_name)"
		return fmt.Sprintf("GROUP_CONCAT(%s)", strings.Join(result, " "))
	}
	return fmt.Sprintf("GROUP_CONCAT(\n%s\n)", linesIndent(strings.Join(result, " ")))
}

func formatMatch(columns []string, expr, option string) string {
	fieldColumns := ""
	if len(columns) > 1 {
		fieldColumns = strings.Join(columns, ",\n")
		return fmt.Sprintf("MATCH (\n%s\n) AGAINST(%s %s)", linesIndent(fieldColumns), expr, option)
	}
	fieldColumns = strings.Join(columns, ",")
	return fmt.Sprintf("MATCH (%s) AGAINST(%s %s)", fieldColumns, expr, option)
}

func formatCollate(value, charset string) string {
	return fmt.Sprintf("%s COLLATE %s", value, charset)
}

func formatConvertUsing(value, convertType string) string {
	return fmt.Sprintf("CONVERT(%s USING %s)", value, convertType)
}

func formatTable(tableName, hintType string, hints []string) string {
	if len(hints) > 0 {
		hintString := strings.Join(hints, ",")
		return fmt.Sprintf("%s %s INDEX (%s)", tableName, strings.ToUpper(strings.TrimSpace(hintType)), hintString)
	}
	return fmt.Sprintf("%s", tableName)
}

func formatAsTable(tableName, asString, hintType string, hints []string) string {
	tblName := formatAS(tableName, asString)
	return formatTable(tblName, hintType, hints)
}

func formatDBTable(DBName, tableName string) string {
	if DBName == "" {
		return tableName
	}
	return fmt.Sprintf("%s.%s", DBName, tableName)
}

func formatBinaly(operator, left, right string) string {
	return fmt.Sprintf("%s %s %s", left, strings.ToUpper(operator), right)
}
