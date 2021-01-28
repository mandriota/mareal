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
	inner := make([]byte, 0, 8192)

	for i := 0; scanner.Scan(); i++ {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "+") {
			words := strings.Fields(line[1:])
			if len(words) > 0 {
				switch words[0] {
				case "insert":
					pass := strings.Join(words[1:], " ")
					if !strings.HasSuffix(pass, ".mr") {
						panic(fmt.Sprintf("at line %d: not mareal file %s", i, pass))
					}

					inner = append(inner, Read(pass)...)
				}
			}

			continue
		}

		inner = append(inner, line...)
	}

	return string(inner)
}