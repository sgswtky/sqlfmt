package parse

import (
	"github.com/sgswtky/sqlparser"
	"fmt"
	"strings"
)

func (b *builder) expr(e sqlparser.Expr) string {
	switch parsedExpr := e.(type) {
	case *sqlparser.ComparisonExpr:
		return fmt.Sprintf("%s %s %s", b.expr(parsedExpr.Left), parsedExpr.Operator, b.expr(parsedExpr.Right))
		// TODO cording for escape consideration
	case *sqlparser.AndExpr:
		return formatAND(b.expr(parsedExpr.Left), b.expr(parsedExpr.Right))
	case *sqlparser.OrExpr:
		return formatOR(b.expr(parsedExpr.Left), b.expr(parsedExpr.Right))
	case *sqlparser.NotExpr:
		fmt.Println("*sqlparser.NotExpr")
	case *sqlparser.ParenExpr:
		return formatParenthesis(b.expr(parsedExpr.Expr))
	case *sqlparser.RangeCond:
		return formatBetween(b.expr(parsedExpr.Left), b.expr(parsedExpr.From), b.expr(parsedExpr.To), parsedExpr.Operator)
	case *sqlparser.IsExpr:
		return formatIS(b.expr(parsedExpr.Expr), parsedExpr.Operator)
	case *sqlparser.ExistsExpr:
		return formatExists(b.statementRoot(parsedExpr.Subquery.Select))
	case *sqlparser.SQLVal:
		return valTypeFormat(parsedExpr.Type, parsedExpr.Val)
	case *sqlparser.NullVal:
		return "NULL"
	case sqlparser.BoolVal:
		return fmt.Sprintf("%v", parsedExpr)
	case *sqlparser.ColName:
		if b.simpleTableExpr(parsedExpr.Qualifier) == "" {
			return parsedExpr.Name.String()
		}
		return formatTableColumn(b.simpleTableExpr(parsedExpr.Qualifier), parsedExpr.Name.String())
	case sqlparser.ValTuple:
		values := make([]string, 0)
		for _, v := range parsedExpr {
			values = append(values, b.expr(v))
		}
		return formatTuple(values...)
	case *sqlparser.Subquery:
		return formatSubquery(b.statementRoot(parsedExpr.Select))
		//case *sqlparser.ListArg:
		//	// TODO cording
		//	fmt.Println("*sqlparser.ListArg")
		//case *sqlparser.BinaryExpr:
		//	// TODO cording
		//	fmt.Println("*sqlparser.BinaryExpr")
	case *sqlparser.UnaryExpr:
		return formatUnary(parsedExpr.Operator, b.expr(parsedExpr.Expr))
		//case *sqlparser.IntervalExpr:
		//	// TODO cording
		//	fmt.Println("*sqlparser.IntervalExpr")
	case *sqlparser.CollateExpr:
		value := b.expr(parsedExpr.Expr)
		charset := parsedExpr.Charset
		return formatCollate(value, charset)
	case *sqlparser.FuncExpr:
		return b.getFuncExpr(parsedExpr)
	case *sqlparser.CaseExpr:
		caseValue := b.expr(parsedExpr.Expr)
		elseValue := b.expr(parsedExpr.Else)
		whens := make([]string, 0)
		for _, v := range parsedExpr.Whens {
			cond := b.expr(v.Cond)
			val := b.expr(v.Val)
			whens = append(whens, formatWhen(cond, val))
		}
		return formatCase(caseValue, elseValue, whens)
		//case *sqlparser.ValuesFuncExpr:
		//	// TODO cording
		//	fmt.Println("*sqlparser.ValuesFuncExpr")
	case *sqlparser.ConvertExpr:
		// CAST and CONVERT
		format := getConvertTypeQualifier(parsedExpr.Type)
		return fmt.Sprintf(format, b.expr(parsedExpr.Expr))
	case *sqlparser.ConvertUsingExpr:
		value := b.expr(parsedExpr.Expr)
		return formatConvertUsing(value, parsedExpr.Type)
	case *sqlparser.MatchExpr:
		columns := b.selectExprs(parsedExpr.Columns)
		expr := b.expr(parsedExpr.Expr)
		option := parsedExpr.Option
		return formatMatch(columns, expr, option)
	case *sqlparser.GroupConcatExpr:
		distinct := parsedExpr.Distinct
		values := b.selectExprs(parsedExpr.Exprs)
		orderBy := ""
		if parsedExpr.OrderBy != nil {
			orderBy = b.orderBy(parsedExpr.OrderBy)
		}
		separator := parsedExpr.Separator
		return formatGroupConcat(distinct, values, orderBy, separator)
		//case *sqlparser.Default:
		//	// TODO
		//	fmt.Println("*sqlparser.Default")
	default:
		unknownType(parsedExpr)
	}
	return ""
}

func getConvertTypeQualifier(t *sqlparser.ConvertType) string {
	qualifier := strings.ToUpper(t.Type)
	switch qualifier {
	case ConvertTypeBinary:
		fallthrough
	case ConvertTypeChar:
		// only length
		if t.Length != nil {
			v := valTypeFormat(t.Length.Type, t.Length.Val)
			return formatConvertTypeQualifierFormat(fmt.Sprintf("%s(%s)", qualifier, v))
		}
		return formatConvertTypeQualifierFormat(qualifier)
	case ConvertTypeDate:
		fallthrough
	case ConvertTypeDatetime:
		fallthrough
	case ConvertTypeTime:
		fallthrough
	case ConvertTypeSigned:
		fallthrough
	case ConvertTypeUnsigned:
		// only qualifier
		return formatConvertTypeQualifierFormat(qualifier)
	case ConvertTypeDecimal:
		// only decimal
		length := ""
		if t.Length != nil {
			length = valTypeFormat(t.Length.Type, t.Length.Val)
			if t.Scale != nil {
				scale := valTypeFormat(t.Scale.Type, t.Scale.Val)
				format := fmt.Sprintf("%s(%s, %s)", qualifier, length, scale)
				return formatConvertTypeQualifierFormat(format)
			}
			format := fmt.Sprintf("%s(%s)", qualifier, length)
			return formatConvertTypeQualifierFormat(format)
		}
		return formatConvertTypeQualifierFormat(qualifier)
	default:
		unknownType(qualifier)
	}
	return ""
}

func (b *builder) getFuncExpr(funcExpr *sqlparser.FuncExpr) string {
	funcType := funcExpr.Name.String()
	strExprs := make([]string, 0)
	for _, expr := range funcExpr.Exprs {
		strExprs = append(strExprs, b.selectExpr(expr))
	}

	switch strings.ToUpper(funcType) {
	case Avg:
		fallthrough
	case BitAnd:
		fallthrough
	case BitOr:
		fallthrough
	case BitXor:
		fallthrough
	case Count:
		fallthrough
	case CountDistinct:
		fallthrough
	case GroupConcat: // TODO not tried
		fallthrough
	case Max:
		fallthrough
	case Min:
		fallthrough
	case Std:
		fallthrough
	case StdDev:
		fallthrough
	case StdDevPop:
		fallthrough
	case StdDevSamp:
		fallthrough
	case VarPop:
		fallthrough
	case VarSamp:
		fallthrough
	case Variance:
		fallthrough
	case Now:
		fallthrough
	case Concat:
		fallthrough
	case Sum:
		if len(strExprs) > 1 {
			return formatFuncs(strings.ToUpper(funcType), formatSimpleArray(strExprs), funcExpr.Distinct)
		}
		return formatFuncs(strings.ToUpper(funcType), fmt.Sprintf("(%s)", strings.Join(strExprs, "")), funcExpr.Distinct)
	default:
		unknownType(funcType)
	}
	return ""
}
