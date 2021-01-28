package parser

type Type uint16
type Literal interface{}

type Token struct {
	Typ Type
	Val Literal
}

const (
	ILLEGAL Type = iota
	ROUTINE
	IDENT
	NUM
	STR
	ARR
)

var Tokens = []string{
	ILLEGAL: "ILLEGAL",
	ROUTINE: "ROUTINE",
	IDENT: "IDENT",
	NUM: "NUMBER",
	STR: "STRING",
	ARR: "ARRAY",
}
