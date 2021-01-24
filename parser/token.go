package parser

type Type uint16
type Literal interface{}

type Token struct {
	Typ Type
	Val Literal
	Row int32
}

const (
	ILLEGAL Type = iota
	ROUTINE
	IDENT
	FLOAT
	STR
	ARR
)

var Tokens = []string{
	ILLEGAL: "ILLEGAL",
	ROUTINE: "ROUTINE",
	IDENT:   "IDENT",
	FLOAT:   "FLOAT",
	STR:     "STRING",
	ARR:     "ARRAY",
}
