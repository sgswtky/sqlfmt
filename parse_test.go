package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	splitAreaString = "--------------------------------------------------\n"
)

func TestSelectSQLFmt(t *testing.T) {
	filepath.Walk("test_select_file/", func(path string, info os.FileInfo, err error) error {
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
