package loop

import "strings"

func buildEnum(values []string) (s string) {
	switch len(values) {
	case 0:
		return ""
	case 1:
		return values[0]
	}

	return strings.Join(values[:len(values)-1], ", ") +
		" and " + values[len(values)-1]
}
