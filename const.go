package main

const (
	as  = "AS"
	and = "AND"
	or  = "OR"

	indentSpace = 2
)

const (
	null = "NULL"
)

const (
	dual = "dual"
)

const (
	convertTypeBinary   = "BINARY"
	convertTypeChar     = "CHAR"
	convertTypeDate     = "DATE"
	convertTypeDatetime = "DATETIME"
	convertTypeDecimal  = "DECIMAL"
	convertTypeSigned   = "SIGNED"
	convertTypeUnsigned = "UNSIGNED"
	convertTypeTime     = "TIME"
)

const (
	avg           = "AVG"
	bitAnd        = "BIT_AND"
	bitOr         = "BIT_OR"
	bitXor        = "BIT_XOR"
	count         = "COUNT"
	countDistinct = "COUNT(DISTINCT)" //TODO
	groupConcat   = "GROUP_CONCAT"
	max           = "MAX"
	min           = "MIN"
	std           = "STD"
	stdDev        = "STDDEV"
	stdDevPop     = "STDDEV_POP"
	stdDevSamp    = "STDDEV_SAMP"
	sum           = "SUM"
	varPop        = "VAR_POP"
	varSamp       = "VAR_SAMP"
	variance      = "VARIANCE"
	now           = "NOW"
	concat        = "CONCAT"
	ifnull        = "IFNULL"
	round         = "ROUND"
)

const (
	asc  = "ASC"
	desc = "DESC"
)

const (
	unknownTypeError = `
There is a possibility of SQL including an unsupported implementation.
Please describe SQL and create an github issue or contact me. twitter: @sgswtky
unknown value: %+v
`
	unsportedTypeError = `
Support for select statement only.
If in the case of select statement looked this error, there is a possibility of a bug.
There in the case Please describe SQL and create an github issue or contact me. twitter: @sgswtky
unknown value: %+v
`
)
