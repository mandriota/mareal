package parser

import (
	"errors"
	"regexp"
	"strconv"
	u "unicode"
)

var (
	number = regexp.MustCompile(`[-+]?([0-9]*[.])?[0-9]+([eE][-+]?\d+)?`)
)

type Lexer struct {
	inner string
	pos int
}

func New(inner string) *Lexer {
	return &Lexer{inner: inner}
}

func (self *Lexer) Lex() *Node {
	main := &Node{
		Component: make([]*Node, 0),
		Token: &Token{Typ: ROUTINE},
	}

loop:
	for self.pos < len(self.inner) {
		switch self.inner[self.pos] {
		case ' ', '\t', '\r', '\n':
		case '(':
			self.pos++
			main.Component.Add(self.Lex())
		case ')':
			break loop
		case '"':
			self.pos++
			main.Component.Add(&Node{Token: &Token{Typ: STR, Val: self.read(func(r rune) bool { return r == '"' })}})
		default:
			switch {
			case u.IsLetter([]rune(self.inner)[self.pos]) || []rune(self.inner)[self.pos] == '_':
				main.Component.Add(&Node{Token: &Token{Typ: IDENT, Val: self.read(func(r rune) bool { return !u.IsLetter(r) && !u.IsDigit(r) && r != '_' })}})
				continue
			case u.IsDigit([]rune(self.inner)[self.pos]) || []rune(self.inner)[self.pos] == '+' || []rune(self.inner)[self.pos] == '-':
				main.Component.Add(&Node{Token: &Token{Typ: NUM, Val: self.readNum()}})
				continue
			default:
				panic(errors.New("Illegal character at line "))
			}
		}
		self.pos++
	}

	return main
}

func (self *Lexer) read(delimit func(rune) bool) string {
	beg := self.pos
	for ; self.pos < len(self.inner) && !delimit([]rune(self.inner)[self.pos]); self.pos++ {}

	return self.inner[beg:self.pos]
}

func (self *Lexer) readNum() float64 {
	str := number.FindString(self.inner[self.pos:])
	self.pos += len(str)

	res, _ := strconv.ParseFloat(str, 64)
	return res
}
