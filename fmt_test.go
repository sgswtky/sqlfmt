package main

import (
	"bytes"
	"fmt"
	"testing"
	"go/ast"
)

func TestReplaceGoFile(t *testing.T) {
	gofile := `
package test_file
//aaaaaa
import "fmt"

func main() {
	// sqlfmt
	sql := "select * from user left join userb on user.user_id = userb.user_id"
	fmt.Println(sql)
}

//a/aaaa
`
	buff := new(bytes.Buffer)
	if err := fmtFile("gofile.go", bytes.NewBufferString(gofile), buff); err != nil {
		t.Fatal("expect not error, but occurred error.")
	}

	expect := fmt.Sprintf("%s%s%s", `package test_file

//aaaaaa
import "fmt"

func main() {
	// sqlfmt
	sql := `+ "`", `
SELECT
  *
FROM
  user
  LEFT JOIN userb
    ON user.user_id = userb.user_id`+ "`", `
	fmt.Println(sql)
}

//a/aaaa
`)
	if expect != buff.String() {
		fmt.Println(buff.String())
		t.Fatal("assert error.")
	}
}

func TestIsFmtTargetComment(t *testing.T) {
	expectComment := true
	resultComment := isFmtTargetComment("// " + commentConst)
	if resultComment != expectComment {
		t.Fatal(expectFmt(expectComment, resultComment))
	}

	expectNonComment := false
	resultNonComment := isFmtTargetComment("// aaaa")
	if resultNonComment != expectNonComment {
		t.Fatal(expectFmt(expectNonComment, resultNonComment))
	}
}

func TestGetBasicLit(t *testing.T) {
	basicLit := &ast.BasicLit{
		Value: "",
	}
	expectBasicLit := basicLit
	resultBasicLit := getBasicLit([]ast.Expr{basicLit})
	if resultBasicLit != expectBasicLit {
		t.Fatal(expectFmt(expectBasicLit, resultBasicLit))
	}

	nonBasicLit := &ast.FuncLit{}
	var expectNonBasicLit *ast.BasicLit = nil
	resultNonBasicLit := getBasicLit([]ast.Expr{nonBasicLit})
	if resultNonBasicLit != expectNonBasicLit {
		t.Fatal(expectFmt(expectNonBasicLit, resultNonBasicLit))
	}
}

func TestGetIdent(t *testing.T) {
	ident := &ast.Ident{
		Name: "",
	}
	expectIdent := ident
	resultIdent := getIdent([]ast.Expr{ident})
	if resultIdent != expectIdent {
		t.Fatal(expectFmt(expectIdent, resultIdent))
	}

	nonIdent := &ast.FuncLit{}
	var expectNonIdent *ast.Ident = nil
	resultNonIdent := getIdent([]ast.Expr{nonIdent})
	if resultNonIdent != expectNonIdent {
		t.Fatal(expectFmt(expectNonIdent, resultNonIdent))
	}
}

func TestParseAssignStmt(t *testing.T) {

	ident := &ast.Ident{
		NamePos: 55,
		Name:    "sql",
		Obj: &ast.Object{
			Kind: 4,
			Name: "sql",
			Decl: &ast.AssignStmt{},
			Data: nil,
			Type: nil,
		},
	}
	initialValue := "`select * from example`"
	basicLit := &ast.BasicLit{
		ValuePos: 62,
		Kind:     9,
		Value:    initialValue,
	}
	var expect error = nil
	result := replaceFormatedSQL(basicLit, ident)
	if result != expect {
		t.Fatal(expectFmt(expect, result))
	}
	// expect, changed basicLit Value
	if basicLit.Value == initialValue {
		t.Fatal("basicLit and initialValue is it should not be the same")
	}

	identNotTarget := ident
	ident.Name = "non target comment"
	basicLitNotTarget := &ast.BasicLit{
		ValuePos: 62,
		Kind:     9,
		Value:    initialValue,
	}
	var expectNotTarget error = nil
	resultNotTarget := replaceFormatedSQL(basicLitNotTarget, identNotTarget)
	if resultNotTarget != expectNotTarget {
		t.Fatal(expectFmt(expectNotTarget, resultNotTarget))
	}
	// expect, not changed basicLit Value
	if basicLitNotTarget.Value != initialValue {
		t.Fatal("basicLit to should not be write")
	}
}
