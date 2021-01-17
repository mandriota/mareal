package token

type Type uint16
type Token struct {
	Typ Type
	Val string
}

const (
	ILLEGAL Type = iota
	BLOCK
	IDENT
	ERROR
	STR
	NUM
)

var Tokens = []string {
	ILLEGAL: "ILLEGAL",
	BLOCK: "BLOCK",
	IDENT: "IDENT",
	ERROR: "ERROR",
	STR: "STRING",
	NUM: "NUMBER",
}