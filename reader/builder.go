package reader

import "strings"

func build(line string) string {
	inner := strings.TrimSpace(line)

	if strings.HasPrefix(inner, "#") {
		return ""
	}

	if strings.HasPrefix(inner, "+") {
		args := strings.Fields(inner[1:])

		if len(args) > 0 {
			switch args[0] {
			case "insert":
				pass := strings.Join(args[1:], " ")
				return Read(pass)
			}
		}
	}

	return line
}
