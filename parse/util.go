package parse

import (
	"strings"
	"fmt"
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
