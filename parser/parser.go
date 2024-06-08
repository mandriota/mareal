package parser

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strings"
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
	src *bufio.Reader
	col int
	row int
	cnt int
	cch byte
	eof bool
}

func New(src io.Reader) *Parser {
	return &Parser{src: bufio.NewReaderSize(src, 2048)}
}

func (p *Parser) Parse() (*Node, error) {
	root := &Node{
		Component: make([]*Node, 0),
		Token:     Token{Typ: TkRoutine},
	}

	p.readByte()
	for !p.eof {
		switch p.cch {
		case '#':
			p.skipLine()
			fallthrough
		case '\n':
			p.col = -1
			p.row++
		case ' ', '\t', '\r':
		case '(':
			node, err := p.Parse()
			if err != nil {
				return nil, err
			}
			root.Component.Add(node)
		case ')':
			return root, nil
		case '"':
			p.readByte()
			root.Component.Add(&Node{
				Token: Token{
					Typ: TkStr,
					Val: p.read(func(b byte) bool {
						return b != '"'
					}),
				},
			})
		default:
			switch {
			case isLetter(p.cch):
				root.Component.Add(&Node{
					Token: Token{
						Typ: TkIdent,
						Val: p.read(func(b byte) bool {
							return isLetter(b) || isDigit(b)
						}),
					},
				})
				continue
			case isDigit(p.cch) || p.cch == '+' || p.cch == '-':
				root.Component.Add(&Node{
					Token: Token{
						Typ: TkNum,
						Val: p.readNum(),
					},
				})
				continue
			default:
				return nil, fmt.Errorf("%d:%d: illegal character", p.row+1, p.col)
			}
		}
		p.readByte()
	}

	return root, nil
}

func (p *Parser) readByte() {
	cc, err := p.src.ReadByte()
	if err == io.EOF {
		p.eof = true
		return
	}
	if err != nil {
		panic(fmt.Errorf("failed to read file: %v", err))
	}

	p.col++
	p.cnt++
	p.cch = cc
}

func (p *Parser) read(while func(b byte) bool) string {
	sb := &strings.Builder{}

	for !p.eof && while(p.cch) {
		sb.WriteByte(p.cch)
		p.readByte()
	}

	return sb.String()
}

func (p *Parser) skipLine() {
	for !p.eof && p.cch != '\n' {
		p.readByte()
	}
}

func (p *Parser) readUint() (n uint) {
	for !p.eof && isDigit(p.cch) {
		n *= 10
		n += uint(p.cch - '0')
		p.readByte()
	}

	return n
}

func (p *Parser) readInt() (n int) {
	switch p.cch {
	case '-':
		p.readByte()
		return -int(p.readUint())
	case '+':
		p.readByte()
		fallthrough
	default:
		return int(p.readUint())
	}
}

func (p *Parser) readNum() float64 {
	i := p.readInt()
	d := 0
	dlen := -1
	m := 0

	if p.cch == '.' {
		p.cnt = 0
		d = int(p.readUint())
		dlen = p.cnt
	}

	if p.cch|0x20 == 'e' {
		p.readByte()
		m = p.readInt()
	}

	return (float64(i) + float64(d)/math.Pow10(dlen)) * math.Pow10(m)
}
