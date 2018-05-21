package parse

import (
	"fmt"
	"runtime"
)

func unknownType(i interface{}) {
	fmt.Println("----------------------")
	fmt.Println("There is a possibility of SQL including an unsupported implementation.")
	fmt.Println("Please describe SQL and create an github issue or contact me. twitter: @sgswtky")
	fmt.Println(fmt.Sprintf("unknown value: %+v", i))
	printRuntime()
	fmt.Println("----------------------")
}

func unsupportedType(i interface{}) {
	fmt.Println("----------------------")
	fmt.Println("Support for select statement only.")
	fmt.Println("If in the case of select statement looked this error, there is a possibility of a bug.")
	fmt.Println("There in the case Please describe SQL and create an github issue or contact me. twitter: @sgswtky")
	fmt.Println(fmt.Sprintf("unknown value: %+v", i))
	printRuntime()
	fmt.Println("----------------------")
}

func printRuntime() {
	if _, fileName, line, ok := runtime.Caller(2); ok {
		fmt.Println(fmt.Sprintf("%s:%d", fileName, line))
	}
}
