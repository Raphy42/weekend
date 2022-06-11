package dependency

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestNewGraphBuilder(t *testing.T) {
	g := NewGraphBuilder()
	g.Factories(
		func(ctx context.Context) (string, error) {
			return "go.mod", nil
		},
		func(filename string) ([]byte, error) {
			return ioutil.ReadFile(filename)
		},
	)
	spew.Dump(g.Build())
}
