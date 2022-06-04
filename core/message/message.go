package message

import (
	"time"

	"github.com/rs/xid"
)

type Message struct {
	ID      xid.ID
	Kind    string
	Payload interface{}
}

func New(kind string, payload interface{}) Message {
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
