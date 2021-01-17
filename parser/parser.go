package ast

import (
	"fmt"
	"io"
	. "unicode"

	"./token"
)

type Parser struct {
	in string
	pos int

	Stdout, Stderr, Stdin io.ReadWriter
}

func New(in string, stdout, stderr, stdin io.ReadWriter) *Parser {
	l := &Parser{in: in, Stdout: stdout, Stderr: stderr, Stdin: stdin}
	return l
}

type Node struct {
	Component []*Node
	*token.Token
}

func (p *Parser) Parse() *Node {
	main := &Node{
		Component: make([]*Node, 0),
		Token: &token.Token{Typ: token.BLOCK},
	}

	for; p.pos < len(p.in); p.pos++ {
		switch p.in[p.pos] {
		case ' ', '\t', '\n', '\r':
		case '(':
			p.pos++
			main.Component = append(main.Component, p.Parse())
		case ')':
			p.parseRoutine(main)
			return main
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
				fmt.Fprintf(p.Stderr, "pos %d, Illegal char %c\n", p.pos, p.in[p.pos])
				return main
			}
		}
	}
	return main
}

func (p *Parser) read(pos *int, delimit func(rune) bool) (str string) {
	beg := *pos
	for; *pos < len(p.in) && !delimit([]rune(p.in)[*pos]); *pos++ {}
	
	return p.in[beg:*pos]
}

func (p *Parser) parseRoutine(routine *Node) {
	if routine.Typ == token.BLOCK && len(routine.Component) > 0 {
		switch routine.Component[0].Val {
		case "put":
			for _, arg := range routine.Component[1:] {
				fmt.Fprintf(p.Stdout, arg.Val)
			}
		}
	}
}
