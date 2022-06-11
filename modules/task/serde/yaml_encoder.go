//go:build task.encoder.yaml

package serde

import (
	"gopkg.in/yaml.v2"
)

func NewEncoding() *Encoding {
	return &Encoding{
		Serializer:   yaml.Marshal,
		Deserializer: yaml.Unmarshal,
	}
}
