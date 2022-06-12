package dep

import (
	"fmt"
	"strings"

	"github.com/Raphy42/weekend/pkg/reflect"
)

type Module struct {
	Name        string
	Factories   []interface{}
	Exposes     []interface{}
	SideEffects []interface{}
}

type ModuleOption func(module *Module)

func Factories(exports ...interface{}) ModuleOption {
	return func(module *Module) {
		module.Factories = append(module.Factories, exports...)
	}
}

func Expose(exposed interface{}) ModuleOption {
	return func(module *Module) {
		module.Exposes = append(module.Exposes, exposed)
	}
}

func SideEffects(invocation interface{}) ModuleOption {
	return func(module *Module) {
		module.SideEffects = append(module.SideEffects, invocation)
	}
}

func Declare(name string, options ...ModuleOption) Module {
	mod := Module{
		Name:        name,
		Exposes:     make([]interface{}, 0),
		Factories:   make([]interface{}, 0),
		SideEffects: make([]interface{}, 0),
	}
	for _, opt := range options {
		opt(&mod)
	}
	return mod
}

func sprint(items []interface{}) string {
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
Exposes: %d
%s
Side Effects: %d
%s
`,
		m.Name,
		len(m.Factories),
		sprint(m.Factories),
		len(m.Exposes),
		sprint(m.Exposes),
		len(m.SideEffects),
		sprint(m.SideEffects))
}
