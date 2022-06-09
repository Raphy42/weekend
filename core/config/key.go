package config

import "strings"

func Key(elems ...string) string {
	return "." + strings.Join(elems, ".")
}
