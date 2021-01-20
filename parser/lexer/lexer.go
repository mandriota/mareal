package lexer

import (
	. "unicode"

	"./token"
)

type Lexer struct {
	in string
	pos int
}

func New(in string) *Lexer {
	return &Lexer{in: in}
}

type Node struct {
	Component []*Node
	*token.Token
}

func (l *Lexer) Lex() *Node {
	main := &Node{
		Component: make([]*Node, 0),
		Token: &token.Token{Typ: token.ROUTINE},
	}

loop:
	for; l.pos < len(l.in); l.pos++ {
		switch l.in[l.pos] {
		case ' ', '\t', '\n', '\r':
		case '(':
			l.pos++
			main.Component = append(main.Component, l.Lex())
		case ')':
			break loop
		case '\'':
			l.pos++
			main.Component = append(main.Component, &Node{Token: &token.Token{Typ: token.STR, Val: l.read(func(r rune) bool { return r == '\'' })}})
			l.pos++
		default:
			switch {
			case IsLetter([]rune(l.in)[l.pos]):
				main.Component = append(main.Component, &Node{Token: &token.Token{Typ: token.IDENT, Val: l.read(func(r rune) bool { return !IsLetter(r) && !IsDigit(r) && r != '_' })}})
			case IsDigit([]rune(l.in)[l.pos]):
				main.Component = append(main.Component, &Node{Token: &token.Token{Typ: token.NUM, Val: l.read(func(r rune) bool { return !IsDigit(r) })}})
			default:
				main.Component = append(main.Component, &Node{Token: &token.Token{Typ: token.ILLEGAL, Val: string([]rune(l.in)[l.pos]) }})
			}
		}
	}

	return main
}

func (l *Lexer) read(delimit func(rune) bool) string {
	beg := l.pos
	for; l.pos < len(l.in) && !delimit([]rune(l.in)[l.pos]); l.pos++ {}
	l.pos--

	return l.in[beg:l.pos+1]
}