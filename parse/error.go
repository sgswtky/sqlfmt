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
	if _, fileName, line, ok := runtime.Caller(1); ok {
		fmt.Println(fmt.Sprintf("%s:%d", fileName, line))
	}
	fmt.Println("----------------------")
}
