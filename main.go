package main

import (
	"log"
	"os"

	e "github.com/mandriota/mareal/executor"
	p "github.com/mandriota/mareal/preprocessor"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("No input file specified...")
	}

	if err := e.Execute(p.Read(os.Args[1])); err != nil {
		log.Fatalln(err)
	}
}
