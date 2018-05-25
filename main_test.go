package main

import (
	"testing"
	"reflect"
	"runtime"
)

func TestParseMode(t *testing.T) {
	isTerminal = func(fd int) bool { return true }

	// file mode
	str := "aaa"
	none := ""
	fileOpt = &str
	resultModeFile := parseMode()
	fileOpt = &none
	if resultModeFile != modeFile {
		t.Fatal(expectFmt(modeFile, resultModeFile))
	}

	// sql mode
	sqlOpt = &str
	resultModeCommand := parseMode()
	sqlOpt = &none
	if resultModeCommand != modeCommand {
		t.Fatal(expectFmt(modeCommand, resultModeCommand))
	}

	// interactive mode
	boolean := true
	interactiveOpt = &boolean
	resultModeInteractive := parseMode()
	interactiveOpt = nil
	if resultModeInteractive != modeDialog {
		t.Fatal(expectFmt(modeDialog, resultModeInteractive))
	}

	// pipe mode
	isTerminal = func(fd int) bool { return false }
	resultModePipe := parseMode()
	if resultModePipe != modePipe {
		t.Fatal(expectFmt(modePipe, resultModePipe))
	}
}

func getFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func TestExec(t *testing.T) {
	unknownFunc := exec(modeUnknown)
	resultUnknownFunc := getFuncName(unknownFunc)
	expectUnknownFunc := getFuncName(usageMode)
	if resultUnknownFunc != expectUnknownFunc {
		t.Fatal(expectFmt(expectUnknownFunc, resultUnknownFunc))
	}

	dialogFunc := exec(modeDialog)
	resultDialogFunc := getFuncName(dialogFunc)
	expectDialogFunc := getFuncName(dialogMode)
	if resultDialogFunc != expectDialogFunc {
		t.Fatal(expectDialogFunc, resultDialogFunc)
	}

	commandFunc := exec(modeCommand)
	resultCommandFunc := getFuncName(commandFunc)
	expectCommandFunc := getFuncName(commandMode)
	if resultCommandFunc != expectCommandFunc {
		t.Fatal(expectCommandFunc, resultCommandFunc)
	}

	fileFunc := exec(modeFile)
	resultFileFunc := getFuncName(fileFunc)
	expectFileFunc := getFuncName(fileMode)
	if resultFileFunc != expectFileFunc {
		t.Fatal(expectFileFunc, resultFileFunc)
	}

	pipeFunc := exec(modePipe)
	resultPipeFunc := getFuncName(pipeFunc)
	expectPipeFunc := getFuncName(pipeMode)
	if resultPipeFunc != expectPipeFunc {
		t.Fatal(expectPipeFunc, resultPipeFunc)
	}
}
