package main

import (
	"testing"
	"fmt"
)

func expectFmt(expect, result interface{}) string {
	return fmt.Sprintf("expect: `%v` but result: `%v`", expect, result)
}

func TestLinesIndent(t *testing.T) {
	// 1 indents
	const expectSingle = "  test1\n  test2\n  test3"
	resultSingle := linesIndent("test1\ntest2\ntest3")
	if resultSingle != expectSingle {
		t.Fatal(expectFmt(expectSingle, resultSingle))
	}
	// 2 indents
	const expectDouble = "    test1\n    test2\n    test3"
	resultDouble := linesIndent("  test1\n  test2\n  test3")
	if resultDouble != expectDouble {
		t.Fatal(expectFmt(expectDouble, resultDouble))
	}
}

func TestIsDual(t *testing.T) {
	result := isDual([]string{dual})
	expectTrue := true
	if result != expectTrue {
		t.Fatal(expectFmt(expectTrue, result))
	}

	resultNonDual := isDual([]string{"not_dual"})
	expectFalse := false
	if resultNonDual != expectFalse {
		t.Fatal(expectFmt(expectFalse, result))
	}
}

func TestStringsContains(t *testing.T) {
	// false
	const expectFalse = false
	resultFalse := stringsContains([]string{"val1", "val2", "val3"}, "val4")
	if resultFalse != expectFalse {
		t.Fatal(expectFmt(expectFalse, resultFalse))
	}

	// true
	const expectTrue = true
	resultTrue := stringsContains([]string{"val1", "val2", "val3"}, "val1")
	if resultTrue != expectTrue {
		t.Fatal(expectFmt(expectTrue, resultTrue))
	}
}
