package parser

type Token struct {
	Typ uint8
	Val interface{}
}

const (
	ILLEGAL uint8 = iota
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
