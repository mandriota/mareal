package parser

import (
	"testing"
)

func TestBuff_Add(t *testing.T) {
	buff := &buff{}
	buff.Add(&Node{Token: &Token{Typ: STR, Val: "Mareal"}})

	n := new(Node)
	buff.Sub(n)
	t.Log(n)
}
