package parser

import (
	"regexp"
	"strconv"
	u "unicode"
)

var (
	number = regexp.MustCompile(`[-+]?([0-9]*[.])?[0-9]+([eE][-+]?\d+)?`)
)

type Lexer struct {
	in string
	rp int
}

func New(in string) *Lexer {
	return &Lexer{in: in}
}

func (l *Lexer) Lex() *Node {
	main := &Node{
		Component: make([]*Node, 0),
		Token: &Token{Typ: ROUTINE},
	}

loop:
	for l.rp < len(l.in) {
		switch l.in[l.rp] {
		case ' ', '\n', '\t', '\r':
		case '(':
			l.rp++
			main.Component.Add(l.Lex())
		case ')':
			break loop
		case '"':
			l.rp++
			main.Component.Add(&Node{Token: &Token{Typ: STR, Val: l.read(func(r rune) bool { return r == '"' })}})
		default:
			switch {
			case u.IsLetter([]rune(l.in)[l.rp]) || []rune(l.in)[l.rp] == '_':
				main.Component.Add(&Node{Token: &Token{Typ: IDENT, Val: l.read(func(r rune) bool { return !u.IsLetter(r) && !u.IsDigit(r) && r != '_' })}})
				continue
			case u.IsDigit([]rune(l.in)[l.rp]) || []rune(l.in)[l.rp] == '+' || []rune(l.in)[l.rp] == '-':
				main.Component.Add(&Node{Token: &Token{Typ: NUM, Val: l.readNum()}})
				continue
			default:
				main.Component.Add(&Node{Token: &Token{Typ: ILLEGAL, Val: string([]rune(l.in)[l.rp])}})
				break loop
			}
		}
		l.rp++
	}

	return main
}

func (l *Lexer) read(delimit func(rune) bool) string {
	beg := l.rp
	for ; l.rp < len(l.in) && !delimit([]rune(l.in)[l.rp]); l.rp++ {}

	return l.in[beg:l.rp]
}

func (l *Lexer) readNum() float64 {
	str := number.FindString(l.in[l.rp:])
	l.rp += len(str)

	res, _ := strconv.ParseFloat(str, 64)
	return res
}
