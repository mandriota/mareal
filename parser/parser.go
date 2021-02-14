package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var (
	Scanner = bufio.NewScanner(os.Stdin)
	oWriter = bufio.NewWriter(os.Stdout)
	eWriter = bufio.NewWriter(os.Stderr)
)

type Parser struct {
	lets map[interface{}]*Node
}

func Parse(inner string) (err error) {
	defer func() {
		if v := recover(); v != nil {
			eWriter.WriteString(fmt.Sprintf("Error: %v", v))
		}
	}()

	p := new(Parser)
	p.init()

	p.parse(New(inner).Lex())

	return
}

func (self *Parser) init() {
	self.lets = map[interface{}]*Node {
		"nl": {Token: &Token{Typ: STR, Val: "\n"}},
	}
}

func (self *Parser) parse(inner *Node) *buff {
	ret := buff(make([]*Node, 0))

	switch inner.Typ {
	case NUM, STR, ARR:
		ret.Add(inner)
	case IDENT:
		if let, ok := self.lets[inner.Val]; ok {
			ret.Add(*self.parse(let)...)
		} else {
			panic(fmt.Errorf("undeclared variable \"%s\"", inner.str()))
		}
	case ROUTINE:
		if len(inner.Component) > 0 {
			switch cmd := inner.Component[0]; cmd.Val {
			case "_":
				for _, el := range inner.Component[1:] {
					ret.Add(*self.parse(el)...)
				}
			case "set":
				if len(inner.Component)%2 != 0 {
					for i := 1; i+1 < len(inner.Component); i += 2 {
						self.lets[inner.Component[i].Val] = inner.Component[i+1]
					}
				} else {
					panic(fmt.Errorf("missing assigned value"))
				}
			case "for":
				if len(inner.Component) > 2 && inner.Component[2].Typ == ROUTINE {
					v := new(Node)
					switch self.parse(inner.Component[1]).Sub(v); v.Typ {
					case NUM:
						n := int(v.Val.(float64))
						for i := 0; i < n; i++ {
							self.lets["x"] = &Node{Token: &Token{Typ: NUM, Val: float64(i)}}
							ret.Add(*self.parse(inner.Component[2])...)
						}
					case ARR:
						for _, el := range v.Component {
							self.parse(el).Sub(self.lets["x"])
							ret.Add(*self.parse(inner.Component[2])...)
						}
					default:
						panic(fmt.Errorf("wrong args for cmd \"%s\"", cmd.str()))
					}
				} else {
					panic(fmt.Errorf("wrong args for cmd \"%s\"", cmd.str()))
				}
			case "in":
				self.write(inner.Component[1:]...)

				Scanner.Scan()
				ret.Add(&Node{Token: &Token{Typ: STR, Val: Scanner.Text()}})
			case "out":
				self.write(inner.Component[1:]...)
			}
		}
	}

	return &ret
}

func (self *Parser) write(args ...*Node) {
	for _, el := range args {
		for _, el := range *self.parse(el) {
			oWriter.WriteString(el.str())
		}
	}

	oWriter.Flush()
}

func (self *Node) str() (v string) {
	switch self.Typ {
	case STR, IDENT:
		return self.Val.(string)
	case NUM:
		return strconv.FormatFloat(self.Val.(float64), 'f', -1, 64)
	case ARR:
		for _, el := range self.Component {
			v += el.str()
		}
	}

	return
}