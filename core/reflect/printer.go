package reflect

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

func SPrint(value interface{}) string {
	return spew.Sdump(value)
}

func Signature(value interface{}) string {
	return fmt.Sprintf("%T", value)
}
