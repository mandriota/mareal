// Copyright 2021 Mark Mandriota. All right reserved.

package lexer

import . "github.com/MarkMandriota/Mareallang/token"

func NewLexer(f string) *Lexer {
	return &Lexer{fi: f}
}

type Lexer struct {
	fi       string
	beg, end uint32
	cc       byte
}

func (l *Lexer) NextByte() byte {
	if l.end < uint32(len(l.fi)) {
		l.cc = l.fi[l.end]
	} else {
		l.cc = 0
	}

	l.end++
	return l.cc
}

func (l *Lexer) NextToken() (t Token, v string) {
	for l.NextByte(); l.isSpace(); {
		l.NextByte()
	}
	l.beg = l.end

	switch l.cc {
	case '(':
		return LPAREN, ""
	case ')':
		return RPAREN, ""
	case ':':
		return LABEL, ""
	case '*':
		for l.NextByte() != 0 && l.cc != '*' {
		}
		return COMMENT, l.fi[l.beg:l.end-1]
	case '`':
		for l.NextByte() != 0 && l.cc != '`' {
		}
		return STR, l.fi[l.beg:l.end-1]
	case 0x0:
		return EOF, ""
	default:
		switch {
		case l.isDigit(), l.isSign():
			var dotC, sigC, expC byte

			for ;; l.NextByte() {
				for l.isDigit() {
					l.NextByte()
				}

				switch {
				case l.cc == '.' && dotC == 0 && expC == 0:
					dotC++
				case l.isSign() && sigC == 0:
					sigC++
				case l.cc|0x20 == 'e' && expC == 0:
					expC++
					sigC = 0
				default:
					return NUM, l.fi[l.beg-1 : l.end-1]
				}
			}
		case l.isLetter():
			for l.isLetter() || l.isDigit() {
				l.NextByte()
			}
			return WORD, l.fi[l.beg-1:l.end-1]
		}
	}

	return ILLEGAL, l.fi[l.beg-1:l.end]
}

func (l *Lexer) isSpace() bool {
	return l.cc < '!' && l.cc != 0
}

func (l *Lexer) isDigit() bool {
	return l.cc >= '0' && l.cc <= '9'
}

func (l *Lexer) isSign() bool {
	return l.cc == '+' || l.cc == '-'
}

func (l *Lexer) isLetter() bool {
	cc := l.cc|0x20
	return cc >= 'a' && cc <= 'z'
}
