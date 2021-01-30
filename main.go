package main

import (
	p "./parser"
	r "./reader"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("No input file specified...")
	}

	if err := p.Parse(r.Read(os.Args[1])); err != nil {
		log.Fatalln(err)
	}
}