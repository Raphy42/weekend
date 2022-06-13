package reflect

import (
	"fmt"
)

func SPrint(value any) string {
	return Signature(value)
}

func Signature(value any) string {
	return fmt.Sprintf("%T", value)
}
