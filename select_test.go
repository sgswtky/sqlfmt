package main

import "testing"

func TestAsterOption(t *testing.T) {
	builder := NewBuilder("")

	expect := "*"
	result := builder.asterOption("")
	if result != expect {
		t.Fatal(expectFmt(expect, result))
	}

	expectNonAster := "default"
	resultNonAster := builder.asterOption("default")
	if resultNonAster != expectNonAster {
		t.Fatal(expectFmt(expectNonAster, resultNonAster))
	}
}
