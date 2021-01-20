package ast

import (
	l "./lexer"
	t "./lexer/token"
	"fmt"
	"io"
	"os"
)

var (
	Stdin io.ReadWriter = os.Stdin
	Stdout io.ReadWriter = os.Stdout
	Stderr io.ReadWriter = os.Stderr
)

type Parser struct {
	lets map[string]*l.Node
}

func Parse(in string) {
	p := &Parser{lets: map[string]*l.Node{"ln": &l.Node{Token: &t.Token{Typ: t.STR, Val: "\n"}}}}

	p.parse(l.New(in).Lex())
}

func (p *Parser) parse(node *l.Node) *l.Node {
	put := func() {
		for _, el := range node.Component[1:] {
			if node := p.parse(el); node != nil {
				Stdout.Write([]byte(node.Val))
			}
		}
	}

	switch node.Typ {
	case t.NUM, t.STR, t.ARR:
	case t.IDENT:
		return p.lets[node.Val]
	case t.ROUTINE:
		if len(node.Component) > 0 {
			switch node.Component[0].Val {
			case "":
				for _, el := range node.Component {
					p.parse(el)
				}
			case "set":
				for i := 1; i < len(node.Component)-1; i += 2 {
					p.lets[node.Component[i].Val] = p.parse(node.Component[i+1])
				}
			case "put": put()
			case "get":
				put()
				node := &l.Node{Token: &t.Token{Typ: t.STR}}
				fmt.Fscanf(Stdin, "%s", &node.Val)
				return node
			}
		}
	}

	return node
}