package main

import (
	p "./parser"
)

func main() {
	test := `
(set arr (new 0 1 2))
(for arr (put x ln))
`

	p.Parse(test)
}
