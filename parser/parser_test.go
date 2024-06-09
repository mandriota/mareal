package parser

import (
	"fmt"
	"strings"
	"testing"
)

const TestProgram =
`_

(def fib-helper (lambda (quote n x y)
		 (if n (fib-helper (- n 1) (+ x y) x) x)))

(def fib (lambda (quote n)
		 (fib-helper n 0 1)))

(put (_ "fibonacci: "
		 		(fib (num
						 (get "enter number: ")))) nl)
`

func TestLexer_Ast(t *testing.T) {
	l := New(strings.NewReader(TestProgram))
	rt, err := l.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if rt == nil {
		t.Fatal("result must be not nil")
	}

	t.Log("parsed successfully")
	
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
