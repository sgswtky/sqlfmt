package parse

import (
	"github.com/sgswtky/sqlparser"
	"fmt"
)

func valTypeFormat(valType sqlparser.ValType, body []byte) string {
	switch valType {
	case sqlparser.StrVal:
		return fmt.Sprintf("\"%s\"", string(body))
	case sqlparser.IntVal:
		fallthrough
	case sqlparser.FloatVal:
		fallthrough
	case sqlparser.HexNum:
		fallthrough
	case sqlparser.HexVal:
		fallthrough
	case sqlparser.ValArg:
		fallthrough
	case sqlparser.BitVal:
		return string(body)
	default:
		unknownType(valType)
	}
	return ""
}
