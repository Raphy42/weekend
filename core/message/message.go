package message

import (
	"encoding/json"
	"time"

	"github.com/rs/xid"
	"github.com/vmihailenco/msgpack/v5"
)

type Message struct {
	ID      xid.ID `json:"id"`
	Kind    string `json:"kind"`
	Payload any    `json:"-"`
}

func New(kind string, payload any) Message {
	return Message{
		ID:      xid.New(),
		Kind:    kind,
		Payload: payload,
	}
}

func (m Message) CreatedAt() time.Time {
	return m.ID.Time()
}

func (m Message) Node() []byte {
	return m.ID.Machine()
}

type Envelope struct {
	Message Message `json:"message"`
	Payload []byte  `json:"payload"`
}

func (m Message) Marshall() ([]byte, error) {
	envelope := Envelope{Message: m}
	buf, err := msgpack.Marshal(m.Payload)
	if err != nil {
		return nil, err
	}
	envelope.Payload = buf
	return json.Marshal(envelope)
}

func Unmarshall(buf []byte) (*Message, error) {
	var e Envelope
	if err := json.Unmarshal(buf, &e); err != nil {
		return nil, err
	}
	payload, err := GlobalDecoderRegistry.Decode(e.Message.Kind, e.Payload)
	if err != nil {
		return nil, err
	}
	return &Message{
		ID:      e.Message.ID,
		Kind:    e.Message.Kind,
		Payload: payload,
	}, nil
}
