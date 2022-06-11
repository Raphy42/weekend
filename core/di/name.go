package di

import "strings"

func Name(parts ...string) string {
	return strings.Join(parts, ".")
}
