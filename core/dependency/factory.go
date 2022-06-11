package dependency

import (
	"github.com/rs/xid"
)

type Factory struct {
	id      xid.ID
	inputs  []xid.ID
	outputs []xid.ID
}

//
//func executeFactory(factory Factory, registry *TypeRegistry) (interface{}, error) {
//	factoryT, ok := registry.Get(factory.id)
//	if !ok {
//		return nil, stacktrace.NewError("trying to execute a factory which is not registered in the type registry")
//	}
//
//}
