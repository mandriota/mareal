package parser

import (
	"bufio"
	"fmt"
	"io"
	"math/big"
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
	cch byte
	eof bool
	
	sb strings.Builder
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
	p.cch = cc
}

func (p *Parser) read(while func(b byte) bool) string {
	p.sb.Reset()
	
	for !p.eof && while(p.cch) {
		p.sb.WriteByte(p.cch)
		p.readByte()
	}

	return p.sb.String()
}

func (p *Parser) skipLine() {
	for !p.eof && p.cch != '\n' {
		p.readByte()
	}
}

func (p *Parser) readUint() {
	for !p.eof && isDigit(p.cch) {
		p.sb.WriteByte(p.cch)
		p.readByte()
	}
}

func (p *Parser) readInt() {
	switch p.cch {
	case '-', '+':
		p.sb.WriteByte(p.cch)
		p.readByte()
		fallthrough
	default:
		p.readUint()
	}
}

func (p *Parser) readNum() *big.Float {
	p.sb.Reset()
	p.readInt()

	if p.cch == '.' {
		p.sb.WriteByte(p.cch)
		p.readByte()
		p.readUint()
	}

	if p.cch|0x20 == 'e' {
		p.sb.WriteByte(p.cch)
		p.readByte()
		p.readInt()
	}

	n, _ := big.NewFloat(0).SetString(p.sb.String())
	return n
}
