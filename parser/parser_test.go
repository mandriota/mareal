package parser

import (
	"testing"
)

func TestParser_Parse(t *testing.T) {
	if err := Execute(TestProgram); err != nil {
		t.Fatal(err)
	}
}
