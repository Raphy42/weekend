package di

import (
	"fmt"
	"strings"

	"github.com/Raphy42/weekend/pkg/reflect"
)

type Module struct {
	Name      string
	Providers []interface{}
	Exposes   []interface{}
	Invokes   []interface{}
}

type ModuleOption func(module *Module)

func Providers(exports ...interface{}) ModuleOption {
	return func(module *Module) {
		module.Providers = append(module.Providers, exports...)
	}
}

func Expose(exposed interface{}) ModuleOption {
	return func(module *Module) {
		module.Exposes = append(module.Exposes, exposed)
	}
}

func Invoke(invocation interface{}) ModuleOption {
	return func(module *Module) {
		module.Invokes = append(module.Invokes, invocation)
	}
}

func Declare(name string, options ...ModuleOption) Module {
	mod := Module{
		Name:      name,
		Exposes:   make([]interface{}, 0),
		Providers: make([]interface{}, 0),
		Invokes:   make([]interface{}, 0),
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
Providers: %d
%s
Exposes: %d
%s
Invokes: %d
%s
`,
		m.Name,
		len(m.Providers),
		sprint(m.Providers),
		len(m.Exposes),
		sprint(m.Exposes),
		len(m.Invokes),
		sprint(m.Invokes))
}
