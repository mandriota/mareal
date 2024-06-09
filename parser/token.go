package parser

import "fmt"

type TokenType uint8

func (t TokenType) Assert(tt TokenType) error {
	if t != tt {
		return fmt.Errorf("unexpected type %s: expected %s",
			TkStringifyTable[t],
			TkStringifyTable[tt])
	}
	return nil
}

type Token struct {
	Typ TokenType
	Val interface{}
}

const (
	TkIllegal TokenType = iota
	TkRoutine
	TkIdent
	TkNum
	TkStr
)

var TkStringifyTable = []string{
	TkIllegal: "ILLEGAL",
	TkRoutine: "ROUTINE",
	TkIdent:   "IDENT",
	TkNum:     "NUMBER",
	TkStr:     "STRING",
}
