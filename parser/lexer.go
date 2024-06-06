package parser

import (
	"fmt"
	"math"
	"unicode/utf8"
)

func isLetter(r byte) bool {
	return r|0x20 >= 'a' && r|0x20 <= 'z' || r == '_' ||
		r == '+' || r == '-' || r == '*' || r == '/' ||
		r == '%' || r == '^'
}

func isDigit(r byte) bool {
	return r >= '0' && r <= '9'
}

type Lexer struct {
	src string
	pos int
	row int
}

func New(src string) *Lexer {
	return &Lexer{src: src}
}

func (l *Lexer) Lex() (*Node, error) {
	root := &Node{
		Component: make([]*Node, 0),
		Token:     Token{Typ: TkRoutine},
	}

	for l.pos < len(l.src) {
		switch l.src[l.pos] {
		case '\n':
			l.row++
		case ' ', '\t', '\r':
		case '(':
			l.pos++
			node, err := l.Lex()
			if err != nil {
				return nil, err
			}
			root.Component.Add(node)
		case ')':
			return root, nil
		case '"':
			l.pos++
			root.Component.Add(&Node{
				Token: Token{
					Typ: TkStr,
					Val: l.readString(),
				},
			})
		default:
			switch {
			case isLetter(l.src[l.pos]):
				root.Component.Add(&Node{
					Token: Token{
						Typ: TkIdent,
						Val: l.readIdent(),
					},
				})
				continue
			case isDigit(l.src[l.pos]) || l.src[l.pos] == '+' || l.src[l.pos] == '-':
				root.Component.Add(&Node{
					Token: Token{
						Typ: TkNum,
						Val: l.readNum(),
					},
				})
				continue
			default:
				return nil, fmt.Errorf("illegal character at line %d", l.row)
			}
		}
		l.pos++
	}

	return root, nil
}

func (l *Lexer) readString() string {
	beg := l.pos

	r, n := utf8.DecodeRuneInString(l.src[l.pos:])
	for l.pos < len(l.src) && r != '"' {
		l.pos += n
		r, n = utf8.DecodeRuneInString(l.src[l.pos:])
	}

	return l.src[beg:l.pos]
}

func (l *Lexer) readIdent() string {
	beg := l.pos
	for l.pos < len(l.src) && (isLetter(l.src[l.pos]) || isDigit(l.src[l.pos])) {
		l.pos++
	}

	return l.src[beg:l.pos]
}

func (l *Lexer) readUint() (n uint) {
	for l.pos < len(l.src) && isDigit(l.src[l.pos]) {
		n *= 10
		n += uint(l.src[l.pos] - '0')
		l.pos++
	}

	return n
}

func (l *Lexer) readInt() (n int) {
	switch l.src[l.pos] {
	case '-':
		l.pos++
		return -int(l.readUint())
	case '+':
		l.pos++
		fallthrough
	default:
		return int(l.readUint())
	}
}

func (l *Lexer) readNum() float64 {
	i := l.readInt()
	d := 0
	dbeg := -1
	dend := -1
	m := 0

	if l.src[l.pos] == '.' {
		dbeg = l.pos
		d = int(l.readUint())
		dend = l.pos
	}

	if l.src[l.pos]|0x20 == 'e' {
		l.pos++
		m = l.readInt()
	}

	return (float64(i) + float64(d)/math.Pow10(dend-dbeg)) * math.Pow10(m)
}
