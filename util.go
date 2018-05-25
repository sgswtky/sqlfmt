package main

import (
	"fmt"
	"strings"
	"runtime"
)

type tuple2String struct {
	str1 string
	str2 string
}

func linesIndent(lines string) string {
	ss := strings.Split(lines, "\n")
	result := make([]string, 0)
	for _, v := range ss {
		result = append(result, fmt.Sprintf("%s%s", strings.Repeat(" ", indentSpace), v))
	}
	return strings.Join(result, "\n")
}

func stringsContains(strs []string, str string) bool {
	for _, v := range strs {
		if v == str {
			return true
		}
	}
	return false
}

func isDual(tables []string) bool {
	return stringsContains(tables, dual)
}

func addIndent(str string) string {
	return strings.Replace(str, "\n", "\n  ", -1)
}

func unknownType(i interface{}) string {
	return unbleToContinueError(i, unknownTypeError)
}

func unsupportedType(i interface{}) string {
	return unbleToContinueError(i, unsportedTypeError)
}

func unbleToContinueError(i interface{}, sentence string) string {
	fmt.Println("----------------------")
	fmt.Println(fmt.Sprintf(sentence, i))
	if _, fileName, line, ok := runtime.Caller(2); ok {
		fmt.Println(fmt.Sprintf("%s:%d", fileName, line))
	}
	fmt.Println("----------------------")
	return ""
}
