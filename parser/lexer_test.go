package parser

import (
	"fmt"
	"strings"
	"testing"
)

const TestProgram =
`_
# drawCats declaration
(def drawCats (lambda (quote n) (put (rep i n (_ i ": =^_^=" nl)))))

# call drawCats with arguments n=10
(drawCats 10)
`

func TestLexer_Ast(t *testing.T) {
	l := New(TestProgram)
	rt, err := l.Lex()
	if err != nil {
		t.Fatal(err)
	}

	if rt == nil {
		t.Fatal("result must be not nil")
	}
	
	t.Log(traversal("", rt))
}

func traversal(ind string, node *Node) string {
	var out string
	for i, el := range node.Component {
		out += fmt.Sprintf("\n%s%d) Type: %s; Literal: %v", ind, i, TkStringifyTable[el.Typ], el.Val)
		if el.Component != nil {
			out += traversal(strings.Repeat("\t", len(ind)+1), el)
		}
	}
	return out
}
