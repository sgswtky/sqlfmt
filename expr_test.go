package main

import (
	"testing"
	"github.com/sgswtky/sqlparser"
	"fmt"
	"strings"
)

func TestGetConvertTypeQualifier(t *testing.T) {

	normalConvertQualifier := []string{convertTypeDate, convertTypeDatetime, convertTypeTime, convertTypeSigned, convertTypeUnsigned}
	for _, v := range normalConvertQualifier {
		expect := formatConvertTypeQualifierFormat(v)
		convertType := &sqlparser.ConvertType{
			Type: v,
		}
		result := getConvertTypeQualifier(convertType)
		if result != expect {
			t.Fatal(expectFmt(expect, result))
		}
	}

	binaryAndCharQualifier := []string{convertTypeBinary, convertTypeChar}
	for _, v := range binaryAndCharQualifier {
		// not nil length
		lengthConvertType := &sqlparser.ConvertType{
			Type: v,
			Length: &sqlparser.SQLVal{
				Type: sqlparser.StrVal,
				Val:  []byte("EXAMPLE EXAMPLE"),
			},
		}
		lengthConvertTypeContent := valTypeFormat(lengthConvertType.Length.Type, lengthConvertType.Length.Val)
		lengthExpect := formatConvertTypeQualifierFormat(fmt.Sprintf("%s(%s)", v, lengthConvertTypeContent))
		lengthResult := getConvertTypeQualifier(lengthConvertType)
		if lengthResult != lengthExpect {
			t.Fatal(expectFmt(lengthExpect, lengthResult))
		}

		// nil length
		nilLengthConvertType := &sqlparser.ConvertType{
			Type:   v,
			Length: nil,
		}
		nilLengthExpect := formatConvertTypeQualifierFormat(v)
		nilLengthResult := getConvertTypeQualifier(nilLengthConvertType)
		if nilLengthResult != nilLengthExpect {
			t.Fatal(expectFmt(nilLengthExpect, nilLengthResult))
		}
	}
	// convertTypeDecimal

	// not zero + exists scale
	notZeroScaleConvertType := &sqlparser.ConvertType{
		Type: convertTypeDecimal,
		Length: &sqlparser.SQLVal{
			Type: sqlparser.StrVal,
			Val:  []byte("EXAMPLE EXAMPLE"),
		},
		Scale: &sqlparser.SQLVal{
			Type: sqlparser.StrVal,
			Val:  []byte("EXAMPLE EXAMPLE"),
		},
	}
	notZeroScaleFormat := "DECIMAL(" +
		valTypeFormat(notZeroScaleConvertType.Length.Type, notZeroScaleConvertType.Length.Val) +
		", " +
		valTypeFormat(notZeroScaleConvertType.Scale.Type, notZeroScaleConvertType.Scale.Val) + ")"
	notZeroScaleFormatExpect := formatConvertTypeQualifierFormat(notZeroScaleFormat)
	notZeroScaleFormatResult := getConvertTypeQualifier(notZeroScaleConvertType)
	if notZeroScaleFormatResult != notZeroScaleFormatExpect {
		t.Fatal(expectFmt(notZeroScaleFormatExpect, notZeroScaleFormatResult))
	}

	// not zero
	notZeroConvertType := &sqlparser.ConvertType{
		Type: convertTypeDecimal,
		Length: &sqlparser.SQLVal{
			Type: sqlparser.StrVal,
			Val:  []byte("EXAMPLE EXAMPLE"),
		},
	}
	notZeroFormat := "DECIMAL(" +
		valTypeFormat(notZeroScaleConvertType.Length.Type, notZeroScaleConvertType.Length.Val) + ")"
	notZeroFormatExpect := formatConvertTypeQualifierFormat(notZeroFormat)
	notZeroFormatResult := getConvertTypeQualifier(notZeroConvertType)
	if notZeroFormatResult != notZeroFormatExpect {
		t.Fatal(expectFmt(notZeroFormatExpect, notZeroFormatResult))
	}

	// zero
	zeroConvertType := &sqlparser.ConvertType{
		Type: convertTypeDecimal,
	}
	expect := formatConvertTypeQualifierFormat(convertTypeDecimal)
	result := getConvertTypeQualifier(zeroConvertType)
	if result != expect {
		t.Fatal(expectFmt(expect, result))
	}
}

func TestGetFuncExpr(t *testing.T) {
	var builder Builder = &BuilderStruct{}
	constFunc := []string{
		avg, bitAnd, bitOr, bitXor, count, countDistinct,
		groupConcat, max, min, std, stdDev, stdDevPop,
		stdDevSamp, varPop, varSamp, variance, now, concat, ifnull, round, sum,
	}

	for _, v := range constFunc {
		// single expr
		funcExpr := &sqlparser.FuncExpr{
			Name:     sqlparser.NewColIdent(v),
			Distinct: false,
			Exprs: []sqlparser.SelectExpr{
				&sqlparser.StarExpr{
					TableName: sqlparser.TableName{
						Name: sqlparser.NewTableIdent("table1"),
					},
				},
			},
		}
		expect := formatFuncs(strings.ToUpper(v), fmt.Sprintf("(%s)", strings.Join([]string{"table1"}, "")), funcExpr.Distinct)
		result := builder.getFuncExpr(funcExpr)

		if expect != result {
			t.Fatal(expectFmt(expect, result))
		}

		// multiple exprs
		funcExprs := &sqlparser.FuncExpr{
			Name:     sqlparser.NewColIdent(v),
			Distinct: false,
			Exprs: []sqlparser.SelectExpr{
				&sqlparser.StarExpr{
					TableName: sqlparser.TableName{
						Name: sqlparser.NewTableIdent("table1"),
					},
				},
				&sqlparser.StarExpr{
					TableName: sqlparser.TableName{
						Name: sqlparser.NewTableIdent("table2"),
					},
				},
			},
		}
		expectExprs := formatFuncs(
			strings.ToUpper(funcExprs.Name.String()),
			formatSimpleArray([]string{"table1", "table2"}),
			funcExpr.Distinct,
		)
		resultExprs := builder.getFuncExpr(funcExprs)
		if expectExprs != resultExprs {
			t.Fatal(expectFmt(expect, result))
		}
	}
}
