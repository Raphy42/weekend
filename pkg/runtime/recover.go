package runtime

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/palantir/stacktrace"
)

func Recover(fn func()) error {
	errC := make(chan error)
	go func() {
		defer func() {
			if result := recover(); result != nil {
				switch v := result.(type) {
				case string:
					errC <- stacktrace.NewError(v)
				case error:
					errC <- stacktrace.Propagate(v, "recovered panic")
				default:
					errC <- stacktrace.NewError("recovered panic: %s", spew.Sdump(result))
				}
			} else {
				errC <- nil
			}
		}()
		fn()
	}()
	return <-errC
}
