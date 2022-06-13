package task

import (
	"encoding/json"

	"github.com/vmihailenco/msgpack/v5"

	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/core/errors"
)

const (
	payloadTypeJSON    = "json"
	payloadTypeMsgPack = "msgpack"
)

type Manifest struct {
	Name        string `json:"name"`
	Payload     []byte `json:"payload"`
	PayloadType string `json:"payloadType"`
}

func (m *Manifest) Unmarshal(v any) error {
	switch m.PayloadType {
	case payloadTypeJSON:
		if err := json.Unmarshal(m.Payload, &v); err != nil {
			return stacktrace.Propagate(err, "JSON decoding failed")
		}
	case payloadTypeMsgPack:
		if err := msgpack.Unmarshal(m.Payload, &v); err != nil {
			return stacktrace.Propagate(err, "JSON decoding failed")
		}
	}

	return stacktrace.NewErrorWithCode(errors.EUnreachable, "no such encoding '%s' for task manifests", m.PayloadType)
}

func newManifest(name, encoding string, buf []byte) *Manifest {
	return &Manifest{
		Name:        name,
		Payload:     buf,
		PayloadType: encoding,
	}
}

func MarshallJSON(name string, args any) (*Manifest, error) {
	buf, err := json.Marshal(args)
	if err != nil {
		return nil, stacktrace.Propagate(err, "JSON encoding failed")
	}
	return newManifest(name, payloadTypeJSON, buf), nil
}

func MarshallMsgPack(name string, args any) (*Manifest, error) {
	buf, err := msgpack.Marshal(args)
	if err != nil {
		return nil, stacktrace.Propagate(err, "JSON encoding failed")
	}

	return newManifest(name, payloadTypeMsgPack, buf), nil
}
