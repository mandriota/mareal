package parser

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

var (
	Stdin io.ReadWriter = os.Stdin
	Stdout io.ReadWriter = os.Stdout
	Stderr io.ReadWriter = os.Stderr
)

type Parser struct {
	lets map[string]*Node
}

func Parse(in string) {
	p := &Parser{lets: map[string]*Node{"ln": {Token: &Token{Typ: STR, Val: "\n"}}}}

	p.parse(New(in).Lex())
}

func (p *Parser) parse(node *Node) *Node {
	switch node.Typ {
	case NUM, STR:
	case ARR:
		for _, el := range node.Component {
			if el.Typ != STR && el.Typ != NUM {
				el = p.parse(el)
			}
		}
	case IDENT:
		return p.lets[node.Val]
	case ROUTINE:
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
			case "put": p.put(node.Component[1:])
			case "get":
				p.put(node.Component[1:])
				node := &Node{Token: &Token{Typ: STR}}
				fmt.Fscanf(Stdin, "%s", &node.Val)
				return node
			case "for":
				if len(node.Component) > 2 && node.Component[2].Typ == ROUTINE {
					switch v := p.parse(node.Component[1]); v.Typ {
					case ARR:
						for _, el := range v.Component {
							p.lets["x"] = p.parse(el)
							p.parse(node.Component[2])
						}
					case NUM:
						n, _ := strconv.Atoi(v.Val)
						for i := 0; i < n; i++ {
							p.lets["x"] = &Node{Token: &Token{Typ: NUM, Val: strconv.Itoa(i)}}
							p.parse(node.Component[2])
						}
					}
				}
			case "new":
				return &Node{Component: node.Component[1:], Token: &Token{Typ: ARR}}
			}
		}
	}

	return node
}

func (p *Parser) put(args []*Node) {
	for _, el := range args {
		if node := p.parse(el); node != nil {
			Stdout.Write([]byte(node.Val))
		}
	}
}