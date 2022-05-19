package reflect

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

func Print(value interface{}) string {
	return spew.Sdump(value)
}

func Signature(value interface{}) string {
	return fmt.Sprintf("%T", value)
}
