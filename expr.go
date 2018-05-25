package main

import (
	"fmt"
	"github.com/sgswtky/sqlparser"
	"strings"
)

func (b *BuilderStruct) expr(e sqlparser.Expr) string {
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
		return null
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
	case *sqlparser.BinaryExpr:
		return formatBinaly(parsedExpr.Operator, b.expr(parsedExpr.Left), b.expr(parsedExpr.Right))
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
	case nil:
		return ""
	default:
		unknownType(parsedExpr)
	}
	return ""
}

func getConvertTypeQualifier(t *sqlparser.ConvertType) string {
	qualifier := strings.ToUpper(t.Type)
	switch qualifier {
	case convertTypeBinary:
		fallthrough
	case convertTypeChar:
		// only length
		if t.Length != nil {
			v := valTypeFormat(t.Length.Type, t.Length.Val)
			return formatConvertTypeQualifierFormat(fmt.Sprintf("%s(%s)", qualifier, v))
		}
		return formatConvertTypeQualifierFormat(qualifier)
	case convertTypeDate:
		fallthrough
	case convertTypeDatetime:
		fallthrough
	case convertTypeTime:
		fallthrough
	case convertTypeSigned:
		fallthrough
	case convertTypeUnsigned:
		// only qualifier
		return formatConvertTypeQualifierFormat(qualifier)
	case convertTypeDecimal:
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

func (b *BuilderStruct) getFuncExpr(funcExpr *sqlparser.FuncExpr) string {
	funcType := funcExpr.Name.String()
	strExprs := b.selectExprs(funcExpr.Exprs)

	switch strings.ToUpper(funcType) {
	case avg:
		fallthrough
	case bitAnd:
		fallthrough
	case bitOr:
		fallthrough
	case bitXor:
		fallthrough
	case count:
		fallthrough
	case countDistinct:
		fallthrough
	case groupConcat: // TODO not tried
		fallthrough
	case max:
		fallthrough
	case min:
		fallthrough
	case std:
		fallthrough
	case stdDev:
		fallthrough
	case stdDevPop:
		fallthrough
	case stdDevSamp:
		fallthrough
	case varPop:
		fallthrough
	case varSamp:
		fallthrough
	case variance:
		fallthrough
	case now:
		fallthrough
	case concat:
		fallthrough
	case ifnull:
		fallthrough
	case round:
		fallthrough
	case sum:
		if len(strExprs) > 1 {
			return formatFuncs(strings.ToUpper(funcType), formatSimpleArray(strExprs), funcExpr.Distinct)
		}
		return formatFuncs(strings.ToUpper(funcType), fmt.Sprintf("(%s)", strings.Join(strExprs, "")), funcExpr.Distinct)
	default:
		unknownType(funcType)
	}
	return ""
}

func valTypeFormat(valType sqlparser.ValType, body []byte) string {
	switch valType {
	case sqlparser.StrVal:
		return fmt.Sprintf("\"%s\"", string(body))
	case sqlparser.IntVal:
		fallthrough
	case sqlparser.FloatVal:
		fallthrough
	case sqlparser.HexNum:
		fallthrough
	case sqlparser.HexVal:
		fallthrough
	case sqlparser.ValArg:
		fallthrough
	case sqlparser.BitVal:
		return string(body)
	default:
		unknownType(valType)
	}
	return ""
}
