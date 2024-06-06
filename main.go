package main

import (
	p "github.com/mandriota/mareal/parser"
	pp "github.com/mandriota/mareal/preprocessor"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("No input file specified...")
	}

	if err := p.Execute(pp.Read(os.Args[1])); err != nil {
		log.Fatalln(err)
	}
}
