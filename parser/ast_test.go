package parser

import (
	"testing"
)

func TestBuff_Add(t *testing.T) {
	buff := &Buff{}
	buff.Add(&Node{Token: Token{Typ: TkStr, Val: "Mareal"}})

	n := new(Node)
	buff.Sub(n)
	t.Log(n)
}
