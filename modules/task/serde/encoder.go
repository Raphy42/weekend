package serde

import (
	"sync"

	"github.com/palantir/stacktrace"
)

type Serializer func(value interface{}) ([]byte, error)

type Deserializer func(buf []byte, value interface{}) error

type Encoding struct {
	Serializer
	Deserializer
}

type Encodings struct {
	sync.RWMutex
	defaultEncoder *Encoding
	registry       map[string]Encoding
}

func (e *Encodings) RegisterDefault(name string, encoding Encoding) {
	e.Lock()
	defer e.Unlock()

	e.defaultEncoder = &encoding
	e.registry[name] = encoding
}

func (e *Encodings) Register(name string, encoding Encoding) {
	e.Lock()
	defer e.Unlock()

	if e.defaultEncoder == nil {
		e.defaultEncoder = &encoding
	}
	e.registry[name] = encoding
}

func (e *Encodings) Serialize(kind string, value interface{}) ([]byte, error) {
	e.RLock()
	defer e.RUnlock()

	encoding, ok := e.registry[kind]
	if !ok {
		return nil, stacktrace.NewErrorWithCode(
			EEncodingNotFound,
			"no registered encoding for '%s', did you forget to enable a particular build tag ?",
			encoding,
		)
	}
	return encoding.Serializer(value)
}

func (e *Encodings) Deserialize(kind string, buf []byte, value interface{}) error {
	e.RLock()
	defer e.RUnlock()

	encoding, ok := e.registry[kind]
	if !ok {
		return stacktrace.NewErrorWithCode(
			EEncodingNotFound,
			"no registered encoding for '%s', did you forget to enable a particular build tag ?",
			encoding,
		)
	}
	return encoding.Deserializer(buf, value)
}
