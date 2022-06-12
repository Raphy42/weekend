package schedulable

import (
	"strings"
)

func Name(parts ...string) string {
	return "async." + strings.Join(parts, ".")
}
