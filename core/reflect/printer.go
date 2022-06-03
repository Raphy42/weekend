package reflect

import (
	"fmt"
)

func SPrint(value interface{}) string {
	return Signature(value)
}

func Signature(value interface{}) string {
	return fmt.Sprintf("%T", value)
}
