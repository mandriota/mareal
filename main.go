package main

import (
	"log"
	"os"

	e "github.com/mandriota/mareal/executor"
)

func main() {
	l := log.New(os.Stderr, "", 0)
	
	if len(os.Args) < 2 {
		l.Fatalln("no input file specified...")
	}

	fs, err := os.Open(os.Args[1])
	if err != nil {
		l.Fatalln("failed to read file")
	}
	defer fs.Close()

	// defer func() {
	// 	if err, ok := recover().(error); ok && err != nil {
	// 		l.Fatalln(err)
	// 	}
	// }()

	if err := e.Execute(fs); err != nil {
		l.Fatalln(err)
	}
}
