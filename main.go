package main

import (
	p "./parser"
	"bufio"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("No input file specified...")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln("Cannot open input file...")
	}

	in, _ := bufio.NewReader(file).ReadString(0)
	p.Parse(in)
}
