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

type Parser struct {
	src string
	pos int
	row int
}

func New(src string) *Parser {
	return &Parser{src: src}
}

func (p *Parser) Parse() (*Node, error) {
	root := &Node{
		Component: make([]*Node, 0),
		Token:     Token{Typ: TkRoutine},
	}

	for p.pos < len(p.src) {
		switch p.src[p.pos] {
		case '\n':
			p.row++
		case ' ', '\t', '\r':
		case '(':
			p.pos++
			node, err := p.Parse()
			if err != nil {
				return nil, err
			}
			root.Component.Add(node)
		case ')':
			return root, nil
		case '"':
			p.pos++
			root.Component.Add(&Node{
				Token: Token{
					Typ: TkStr,
					Val: p.readString(),
				},
			})
		default:
			switch {
			case isLetter(p.src[p.pos]):
				root.Component.Add(&Node{
					Token: Token{
						Typ: TkIdent,
						Val: p.readIdent(),
					},
				})
				continue
			case isDigit(p.src[p.pos]) || p.src[p.pos] == '+' || p.src[p.pos] == '-':
				root.Component.Add(&Node{
					Token: Token{
						Typ: TkNum,
						Val: p.readNum(),
					},
				})
				continue
			default:
				return nil, fmt.Errorf("illegal character at line %d", p.row)
			}
		}
		p.pos++
	}

	return root, nil
}

func (p *Parser) readString() string {
	beg := p.pos

	r, n := utf8.DecodeRuneInString(p.src[p.pos:])
	for p.pos < len(p.src) && r != '"' {
		p.pos += n
		r, n = utf8.DecodeRuneInString(p.src[p.pos:])
	}

	return p.src[beg:p.pos]
}

func (p *Parser) readIdent() string {
	beg := p.pos
	for p.pos < len(p.src) && (isLetter(p.src[p.pos]) || isDigit(p.src[p.pos])) {
		p.pos++
	}

	return p.src[beg:p.pos]
}

func (p *Parser) readUint() (n uint) {
	for p.pos < len(p.src) && isDigit(p.src[p.pos]) {
		n *= 10
		n += uint(p.src[p.pos] - '0')
		p.pos++
	}

	return n
}

func (p *Parser) readInt() (n int) {
	switch p.src[p.pos] {
	case '-':
		p.pos++
		return -int(p.readUint())
	case '+':
		p.pos++
		fallthrough
	default:
		return int(p.readUint())
	}
}

func (p *Parser) readNum() float64 {
	i := p.readInt()
	d := 0
	dbeg := -1
	dend := -1
	m := 0

	if p.src[p.pos] == '.' {
		dbeg = p.pos
		d = int(p.readUint())
		dend = p.pos
	}

	if p.src[p.pos]|0x20 == 'e' {
		p.pos++
		m = p.readInt()
	}

	return (float64(i) + float64(d)/math.Pow10(dend-dbeg)) * math.Pow10(m)
}
