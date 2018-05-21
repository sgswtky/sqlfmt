package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"go/ast"
)

const (
	splitAreaString = "--------------------------------------------------\n"
)

func TestSelectSQLFmt(t *testing.T) {
	filepath.Walk("test_select_file/", func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)
		// not work if directory
		if info.IsDir() {
			return nil
		}

		b, err := ioutil.ReadFile(path)
		strs := strings.Split(string(b), splitAreaString)
		if len(strs) < 2 {
			t.Error("require source and answer")
		}
		if err != nil {
			t.Error(err)
		}
		source := strs[0]
		answer := strs[1]

		builder := NewBuilder(source)
		sql, err := builder.Parse()

		if err != nil {
			t.Error(err)
		}

		if sql != answer {
			errMsg := fmt.Sprintf(">>>> file answer >>>>>>>>\n%s\n--------------------------------------------------\n%s\n<<<< formatted SQL <<<<<<<<",
				answer,
				sql)
			fmt.Println(errMsg)
			t.Error("failure, not expected string.")
		}
		return nil
	})

}

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
