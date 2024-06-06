package reader

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Read(pass string) string {
	file, err := os.Open(pass)
	if err != nil {
		panic(fmt.Errorf("error opening file %s", pass))
	}

	scanner := bufio.NewScanner(file)

	var inner string
	for i := 0; scanner.Scan(); i++ {
		inner += build(strings.TrimSpace(scanner.Text()))
	}

	return inner
}