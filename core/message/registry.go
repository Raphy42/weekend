package message

import (
	"github.com/palantir/stacktrace"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/Raphy42/weekend/pkg/concurrent_set"
)

func NewDecoder[T any]() DecoderFunc {
	return func(buf []byte) (any, error) {
		value := new(T)
		if err := msgpack.Unmarshal(buf, value); err != nil {
			return nil, err
		}
		return value, nil
	}
}

type DecoderFunc func(buf []byte) (any, error)

type DecoderRegistry struct {
	decoders concurrent_set.Set[string, DecoderFunc]
}

func NewRegistry() *DecoderRegistry {
	return &DecoderRegistry{decoders: concurrent_set.New[string, DecoderFunc]()}
}

func (r *DecoderRegistry) Register(kind string, decoderFunc DecoderFunc) *DecoderRegistry {
	r.decoders.Insert(kind, decoderFunc)
	return r
}

func (r *DecoderRegistry) Decode(kind string, buf []byte) (any, error) {
	decoder, ok := r.decoders.Get(kind)
	if !ok {
		return nil, stacktrace.NewError(
			"unknown payload type for '%s', add it with message.GlobalDecoderRegistry.Register(%s, func(...))",
			kind, kind,
		)
	}
	return decoder(buf)
}

var GlobalDecoderRegistry = NewRegistry()
