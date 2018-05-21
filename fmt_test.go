package main

import (
	"bytes"
	"fmt"
	"github.com/sgswtky/sqlfmt/parse"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	splitAreaString = "--------------------------------------------------\n"
)

func expectFmt(expect, result interface{}) string {
	return fmt.Sprintf("expect: `%v` but result: `%v`", expect, result)
}

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

		builder := parse.NewBuilder(source)
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
