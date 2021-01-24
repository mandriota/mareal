package parser

import (
	"bufio"
	"io"
	"os"
	"strconv"
)

var (
	Stdin io.ReadWriter = os.Stdin
	Stdout io.ReadWriter = os.Stdout
	Stderr io.ReadWriter = os.Stderr

	Scanner = bufio.NewScanner(Stdin)
	Writer = bufio.NewWriter(Stdout)
)

type Parser struct {
	lets map[Literal]*Node
	stack buff
}

func Parse(in string) {
	p := &Parser{lets: map[Literal]*Node{
		"stack": {Token: &Token{Typ: ARR}, Component: make([]*Node, 0, 1024)},
		"nl": {Token: &Token{Typ: STR, Val: "\n"}}},
	}

	p.parse(New(in).Lex())
}

func (p *Parser) parse(node *Node) *Node {
	switch node.Typ {
	case FLOAT, STR, ARR:
	case IDENT:
		return p.lets[node.Val]
	case ROUTINE:
		if len(node.Component) > 0 {
			switch node.Component[0].Val {
			case nil:
				for _, el := range node.Component {
					p.stack.Add(p.parse(el))
				}
				Writer.Flush()
			case "_":
				for _, el := range node.Component[1:] {
					p.stack.Add(p.parse(el))
				}
			case "set":
				for i := 1; i < len(node.Component)-1; i += 2 {
					p.lets[node.Component[i].Val] = node.Component[i+1]
				}
			case "et":
				if len(node.Component) > 2 && node.Component[2].Typ == ROUTINE {
					switch v := p.parse(node.Component[1]); v.Typ {
					case FLOAT, STR:
						n := int(v.To(FLOAT).Val.(float64))
						for i := 0; i < n; i++ {
							p.lets["x"] = &Node{Token: &Token{Typ: FLOAT, Val: float64(i)}}
							p.parse(node.Component[2])
						}
					case ARR:
						for _, el := range v.Component {
							p.lets["x"] = p.parse(el)
							p.parse(node.Component[2])
						}
					}
				}
			case "nova":
				p.stack.Add(&Node{Component: node.Component[1:], Token: &Token{Typ: ARR}})
			case "out":
				p.output(node.Component[1:])
			case "in":
				p.output(node.Component[1:])

				Scanner.Scan()
				p.stack.Add(&Node{Token: &Token{Typ: STR, Val: Scanner.Text()}})
			}
		}
	}

	return node
}

func (p *Parser) output(args []*Node) {
	for _, el := range args {
		entry := len(p.stack)
		p.parse(el)
		for _, el := range p.stack.Get(len(p.stack)-entry) {
			Writer.WriteString(el.To(STR).Val.(string))
		}
	}
	Writer.Flush()
}

func (node *Node) To(t Type) *Node {
	res := &Node{Token: &Token{Typ: t, Val: node.Val}}

	switch t {
	case FLOAT:
		if node.Typ == STR {
			res.Val, _ = strconv.ParseFloat(node.Val.(string), 64)
		}
	case STR:
		if node.Typ == FLOAT {
			res.Val = strconv.FormatFloat(node.Val.(float64), 'f', -1, 64)
		} else if node.Typ == ARR {
			res.Val = ""
			for _, el := range node.Component {
				res.Val = res.Val.(string) + el.To(STR).Val.(string)
			}
		}
	}

	return res
}
