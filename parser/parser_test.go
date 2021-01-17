package ast

import (
	"./token"
	"fmt"
	"os"
	"strings"
	"testing"
)

const TEST = `
(put 'Hello, World!')
`

func TestLexer_Ast(t *testing.T) {
	p := New(TEST, os.Stdout, os.Stderr, os.Stdin)
	if v := p.Parse(); v != nil {
		t.Log(traversal("", v))
	}
}

func traversal(ind string, node *Node) string {
	var out string

	for i, el := range node.Component {
		out += fmt.Sprintf("\n%s%d) Type: %s; Literal: %s", ind, i, token.Tokens[el.Typ], el.Val)
		if el.Component != nil {
			out += traversal(strings.Repeat("\t", len(ind)+1), el)
		}
	}
	return out
}