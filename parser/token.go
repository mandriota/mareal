package parser

type Type uint16
type Token struct {
	Typ Type
	Val string
}

const (
	ILLEGAL Type = iota
	ROUTINE
	IDENT
	STR
	NUM
	ARR

	ADD
	SUB
)

var Tokens = []string {
	ILLEGAL: "ILLEGAL",
	ROUTINE: "ROUTINE",
	IDENT: "IDENT",
	STR: "STRING",
	NUM: "NUMBER",
	ARR: "ARRAY",
}