package ast

import (
	"fmt"
	"io"
	. "unicode"

	"./token"
)

type Parser struct {
	in string
	ch rune
	ps,np int

	Stdout, Stderr, Stdin io.ReadWriter
}

func New(in string, stdout, stderr, stdin io.ReadWriter) *Parser {
	l := &Parser{in: in, Stdout: stdout, Stderr: stderr, Stdin: stdin}
	return l
}

func (p *Parser) readCh() rune {
	if p.np < len(p.in) {
		p.ch = []rune(p.in)[p.np]
	} else {
		p.ch = 0
	}
	p.ps = p.np
	p.np++

	return p.ch
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

	for p.readCh();; p.readCh() {
		switch p.ch {
		case ' ', '\t', '\n', '\r':
		case '(':
			main.Component = append(main.Component, p.Parse())
		case ')', 0:
			p.parseRoutine(main)
			return main
		case '\'':
			p.readCh()
			main.Component = append(main.Component, &Node{Token: &token.Token{Typ: token.STR, Val: p.read(func(r rune) bool {return r == '\''})}})
			p.readCh()
		default:
			switch {
			case IsLetter(p.ch):
				main.Component = append(main.Component, &Node{Token: &token.Token{Typ: token.IDENT, Val: p.read(func(r rune) bool {return !IsLetter(r) && !IsDigit(r) && r != '_'})}})
			case IsDigit(p.ch):
				main.Component = append(main.Component, &Node{Token: &token.Token{Typ: token.NUM, Val: p.read(func(r rune) bool {return !IsDigit(r)})}})
			default:
				fmt.Fprintf(p.Stderr, "pos %d, Illegal char %c;\n", p.ps, p.ch)
				return main
			}
		}
	}
}

func (p *Parser) read(delimit func(rune) bool) (str string) {
	beg := p.ps
	p.np = p.ps
	for; p.np < len(p.in) && !delimit([]rune(p.in)[p.np]); p.np++ {}
	
	return p.in[beg:p.np]
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