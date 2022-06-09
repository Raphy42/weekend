package config

import "fmt"

func MissingKeyMessage(key string) func(format string, args ...interface{}) string {
	return func(format string, args ...interface{}) string {
		return fmt.Sprintf("%s: %s",
			fmt.Sprintf(format, args),
			fmt.Sprintf("missing required key: '%s'", key),
		)
	}
}
