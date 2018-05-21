package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"strings"
)

func fmtSQL(sql string, w io.Writer, mode int) error {
	builder := NewBuilder(sql)
	result, err := builder.Parse()
	if err != nil {
		return err
	}

	if mode == modeDialog {
		// add last line
		result = fmt.Sprintf("-- >> formated sql\n%s", result)
		result += "\n"
		result += fmt.Sprintf("<< -- formated sql\n")
	}
	if mode == modeCommand || mode == modePipe {
		result = fmt.Sprintf("%s\n", result)
	}

	writeCount, err := fmt.Fprint(w, result)
	if writeCount != len([]byte(result)) {
		// TODO error message
		return errors.New("write result byte error")
	}
	return err
}

const (
	commentConst = "sqlfmt"
	variableName = "sql"
)

func parseAssignStmt(assignStmt *ast.AssignStmt) error {
	if len(assignStmt.Lhs) == 1 && len(assignStmt.Rhs) == 1 {
		basicLit := getBasicLit(assignStmt.Rhs)
		ident := getIdent(assignStmt.Lhs)

		if basicLit == nil || ident == nil {
			return errors.New("unknown error")
		}

		if err := replaceFormatedSQL(basicLit, ident); err != nil {
			return err
		}
	}
	return nil
}

func replaceFormatedSQL(basicLit *ast.BasicLit, ident *ast.Ident) error {
	if ident.Name == variableName {
		sqlRune := []rune(basicLit.Value)
		trimSQL := string(sqlRune[1: len(sqlRune)-1])
		sql, err := NewBuilder(trimSQL).Parse()
		basicLit.Value = fmt.Sprintf("`\n%s`", sql)
		return err
	}
	return nil
}

func getBasicLit(expr []ast.Expr) *ast.BasicLit {
	switch parsed := expr[0].(type) {
	case *ast.BasicLit:
		return parsed
	default:
		return nil
	}
}

func getIdent(expr []ast.Expr) *ast.Ident {
	switch parsed := expr[0].(type) {
	case *ast.Ident:
		return parsed
	default:
		return nil
	}
}

func fmtFile(astFilename string, fileReader io.Reader, fileWriter io.Writer) error {

	// Read source
	src, err := ioutil.ReadAll(fileReader)
	if err != nil {
		return err
	}

	// Parse ast
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, astFilename, src, parser.ParseComments)
	if err != nil {
		return err
	}

	// Parse comment map
	cmap := ast.NewCommentMap(fset, f, f.Comments)

	// Symbol search from comments and replace contents of variable
	for n, cgroups := range cmap {
		for _, cgroup := range cgroups {
			for _, c := range cgroup.List {
				if isFmtTargetComment(c.Text) {
					if assignStmt, isBasicLit := n.(*ast.AssignStmt); isBasicLit {
						if err := parseAssignStmt(assignStmt); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	f.Comments = cmap.Filter(f).Comments()
	return format.Node(fileWriter, fset, f)
}

func isFmtTargetComment(comment string) bool {
	comments := strings.Split(comment, "//")
	return len(comments) > 1 && strings.TrimSpace(comments[1]) == commentConst
}
