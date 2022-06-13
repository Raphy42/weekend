package dep

import (
	"fmt"
	"strings"

	"github.com/Raphy42/weekend/pkg/reflect"
)

type Module struct {
	Name        string
	Factories   []any
	SideEffects []any
}

type ModuleOption func(module *Module)

func Factories(exports ...any) ModuleOption {
	return func(module *Module) {
		module.Factories = append(module.Factories, exports...)
	}
}

func SideEffects(invocation ...any) ModuleOption {
	return func(module *Module) {
		module.SideEffects = append(module.SideEffects, invocation...)
	}
}

func Declare(name string, options ...ModuleOption) Module {
	mod := Module{
		Name:        name,
		Factories:   make([]any, 0),
		SideEffects: make([]any, 0),
	}
	for _, opt := range options {
		opt(&mod)
	}
	return mod
}

func sprint(items []any) string {
	out := make([]string, len(items))
	for idx, item := range items {
		out[idx] = "/t" + reflect.SPrint(item)
	}
	return strings.Join(out, "\n")
}

func (m Module) Print() string {
	return fmt.Sprintf(`%s
Factories: %d
%s
Side Effects: %d
%s
`,
		m.Name,
		len(m.Factories),
		sprint(m.Factories),
		len(m.SideEffects),
		sprint(m.SideEffects))
}
