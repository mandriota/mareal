package lexer

import (
	"github.com/MarkMandriota/Mareallang/token"
	"testing"
)

func TestLexer_NextToken(t *testing.T) {
	l := NewLexer("_ :with `text` (display (repeat +1.45e+5 (join :once `=^_^=`nl)))")

	for tok, val := l.NextToken(); tok != token.EOF; tok, val = l.NextToken() {
		t.Log(tok.String(), val)
	}
}

func BenchmarkLexer_NextToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := NewLexer("_ :with `text` (display (repeat +1.45e+5 (join :once `=^_^=`nl)))")
		for tok, _ := l.NextToken(); tok != token.EOF; tok, _ = l.NextToken() {}
	}
}