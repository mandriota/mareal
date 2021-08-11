package token

import (
	"testing"
)

func TestToken(t *testing.T) {
	// (put :msg 666 "Hello, Mareal!" nl)
	text := []Token{
		LPAREN,
		WORD,
		LABEL,
		STR,
		NUM,
		RPAREN,
		ILLEGAL,
		EOF,
	}

	for _, token := range text {
		t.Log(token.String())
	}
}