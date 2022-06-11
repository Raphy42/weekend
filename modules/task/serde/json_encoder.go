//go:build task.encoder.json

package serde

import "encoding/json"

func NewEncoding() *Encoding {
	return &Encoding{
		Serializer:   json.Marshal,
		Deserializer: json.Unmarshal,
	}
}
