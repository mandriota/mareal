package reader

import (
	"bufio"
	"os"
	"strings"
)

const (
	pre = "#"
)

func Read(pass string) (string, error) {
	file, err := os.OpenFile(pass, os.O_RDONLY, 0777)
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(file)

	var result string

	for scanner.Scan() {
		inner := scanner.Text()

		if str := strings.TrimSpace(inner); len(inner) > 0 && strings.HasPrefix(str, pre) {
			switch args := strings.Fields(str[len(pre):]); args[0] {
			case "inline":
				inner, err = Read(strings.Join(args[1:], " "))
				if err != nil {
					return "", err
				}
			}
		}

		result += inner
	}

	return result, nil
}
