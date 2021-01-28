package main

import (
	p "./parser"
	r "./reader"
	"log"
	"os"
)

func main() {
	defer func() {
		if v := recover(); v != nil {
			log.Fatalln(v)
		}
	}()

	if len(os.Args) < 2 {
		log.Fatalln("No input file specified...")
	}

	p.Parse(r.Read(os.Args[1]))
}
