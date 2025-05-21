package strings

import (
	"strings"
)

func TrimSuffix(sep string, in string) string {
	if pos := strings.LastIndex(in, sep); pos >= 0 {
		return in[:pos]
	}

	return in
}

func BaseName(sep string, in string) string {
	if pos := strings.LastIndex(in, sep); pos >= 0 {
		return in[pos+len(sep):]
	}

	return in
}

func Replace(from, to string, in string) string {
	return strings.ReplaceAll(in, from, to)
}

func Split(sep string, s string) []string {
	if out := strings.Split(s, sep); len(out) > 0 {
		return out
	}

	return []string{}
}

func Join(sep string, s ...string) string {
	return strings.Join(s, sep)
}
