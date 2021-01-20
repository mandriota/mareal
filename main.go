package main

import (
	p "./parser"
)

func main() {
	test := `
(set name (get 'What is your name? '))
(put 'Hello, ' name '!' ln)

`

	p.Parse(test)
}
