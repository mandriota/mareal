// Copyright 2021 Mark Mandriota. All right reserved.

package token

const (
	EOF = iota
	ILLEGAL
	COMMENT
	WORD
	NUM
	STR

	LABEL

	LPAREN
	RPAREN
)

type Token uint8

var Stringify = [...]string {
	EOF: "END OF FILE",
	ILLEGAL: "ILLEGAL",
	COMMENT: "COMMENT",
	WORD: "WORD",
	NUM: "NUMBER",
	STR: "STRING",
	LABEL: "LABEL",
	LPAREN: "LEFT PAREN",
	RPAREN: "RIGHT PAREN",
}

func (t *Token) String() string {
	return Stringify[*t]
}