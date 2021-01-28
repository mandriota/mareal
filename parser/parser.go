package parser

import (
	"bufio"
	"os"
	"strconv"
)

var (
	Scanner = bufio.NewScanner(os.Stdin)
	oWriter = bufio.NewWriter(os.Stdout)
	eWriter = bufio.NewWriter(os.Stderr)
)

type Parser struct {
	lets map[Literal]*Node
}

func Parse(in string) {
	p := &Parser{lets: map[Literal]*Node{
		"nl": {Token: &Token{Typ: STR, Val: "\n"}}},
	}

	p.parse(New(in).Lex())
}

func (p *Parser) parse(node *Node) *buff {
	res := buff(make([]*Node, 0))

	switch node.Typ {
	case NUM, STR, ARR:
		res.Add(node)
	case IDENT:
		res.Add(*p.parse(p.lets[node.Val])...)
	case ROUTINE:
		if len(node.Component) > 0 {
			switch cmd := node.Component[0]; cmd.Val {
			case "_":
				for _, el := range node.Component[1:] {
					res.Add(*p.parse(el)...)
				}
			case "set":
				for i := 1; i < len(node.Component)-1; i += 2 {
					p.lets[node.Component[i].Val] = node.Component[i+1]
				}
			case "for":
				if len(node.Component) > 2 && node.Component[2].Typ == ROUTINE {
					v := new(Node)
					switch p.parse(node.Component[1]).Sub(v); v.Typ {
					case NUM:
						n := int(v.Val.(float64))
						for i := 0; i < n; i++ {
							p.lets["x"] = &Node{Token: &Token{Typ: NUM, Val: float64(i)}}
							res.Add(*p.parse(node.Component[2])...)
						}
					case ARR:
						for _, el := range v.Component {
							p.parse(el).Sub(p.lets["x"])
							res.Add(*p.parse(node.Component[2])...)
						}
					}
				}
			case "out":
				p.output(node.Component[1:]...)
			case "in":
				p.output(node.Component[1:]...)

				Scanner.Scan()
				res.Add(&Node{Token: &Token{Typ: STR, Val: Scanner.Text()}})
			}
		}
	}

	return &res
}

func (p *Parser) output(args ...*Node) {
	for _, el := range args {
		for _, el := range *p.parse(el) {
			oWriter.WriteString(el.str())
		}
	}
	oWriter.Flush()
}

func (n *Node) str() (res string) {
	switch n.Typ {
	case STR:
		return n.Val.(string)
	case NUM:
		return strconv.FormatFloat(n.Val.(float64), 'f', -1, 64)
	case ARR:
		for _, el := range n.Component {
			res += el.str()
		}
	}

	return
}