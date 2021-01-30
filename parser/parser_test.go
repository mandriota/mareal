package parser

import (
	"testing"
)

func TestParser_Parse(t *testing.T) {
	if err := Parse(TEST); err != nil {
		t.Fatal(err)
	}
}
