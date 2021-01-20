package ast

import (
	"testing"
)

const TEST = `
(set name (get 'What is your name? '))
(put 'Hello, ' name '!' ln)
(put 'Hello, ' name '!' ln)
`

func TestParser_Parse(t *testing.T) {
	Parse(TEST)
}
