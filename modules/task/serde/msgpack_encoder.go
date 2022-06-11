//go:build task.encoder.msgpack

package serde

import "github.com/vmihailenco/msgpack/v5"

func NewEncoding() *Encoding {
	return &Encoding{
		Serializer:   msgpack.Marshal,
		Deserializer: msgpack.Unmarshal,
	}
}
