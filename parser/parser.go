package ast

import (
	"fmt"
	"io"
	"os"
	. "unicode"

	"./token"
)

var (
	Stdin io.ReadWriter = os.Stdin
	Stdout io.ReadWriter = os.Stdout
	Stderr io.ReadWriter = os.Stderr
)

type Parser struct {
	in string
	pos int

	lets map[string]*Node
}

func New(in string) *Parser {
	l := &Parser{in: in, lets: make(map[string]*Node)}
	return l
}

type Node struct {
	Component []*Node
	*token.Token
}

func (p *Parser) Parse() *Node {
	main := &Node{
		Component: make([]*Node, 0),
		Token: &token.Token{Typ: token.ROUTINE},
	}

	for; p.pos < len(p.in); p.pos++ {
		switch p.in[p.pos] {
		case ' ', '\t', '\n', '\r':
		case '(':
			p.pos++
			main.Component = append(main.Component, p.Parse())
		case ')':
			goto end
		case '\'':
			p.pos++
			main.Component = append(main.Component, &Node{Token: &token.Token{Typ: token.STR, Val: p.read(&p.pos, func(r rune) bool {return r == '\''})}})
			p.pos++
		default:
			switch {
			case IsLetter([]rune(p.in)[p.pos]):
				main.Component = append(main.Component, &Node{Token: &token.Token{Typ: token.IDENT, Val: p.read(&p.pos, func(r rune) bool {return !IsLetter(r) && !IsDigit(r) && r != '_'})}})
			case IsDigit([]rune(p.in)[p.pos]):
				main.Component = append(main.Component, &Node{Token: &token.Token{Typ: token.NUM, Val: p.read(&p.pos, func(r rune) bool {return !IsDigit(r)})}})
			default:
				fmt.Fprintf(Stderr, "pos %d, Illegal char %c\n", p.pos, p.in[p.pos])
				return main
			}
		}
	}
	end:
		main.Val = p.parse(main).Val
	return main
}

func (p *Parser) read(pos *int, delimit func(rune) bool) (str string) {
	beg := *pos
	for; *pos < len(p.in) && !delimit([]rune(p.in)[*pos]); *pos++ {}

	return p.in[beg:*pos]
}

func (p *Parser) parse(routine *Node) *Node {
	switch routine.Typ {
	case token.NUM, token.STR:
		return routine
	case token.IDENT:
		return p.lets[routine.Val]
	case token.ROUTINE:
		if len(routine.Component) > 0 {
			switch routine.Component[0].Val {
			case "set":
				for i := 1; i < len(routine.Component)-1; i += 2 {
					p.lets[routine.Component[i].Val] = routine.Component[i+1]
				}
			case "put":
				for _, el := range routine.Component[1:] {
					if node := p.parse(el); node != nil {
						Stdout.Write([]byte(node.Val))
					}
				}
				return &Node{Token: &token.Token{Typ: token.NUM}}
			}
		}
	}
	return &Node{Token: &token.Token{}}
}
