package kernel

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        uuid.UUID
	Kind      string
	CreatedAt time.Time
	Payload   interface{}
}

func NewEvent(kind string, payload interface{}) *Event {
	return &Event{
		ID:        uuid.New(),
		Kind:      kind,
		CreatedAt: time.Now(),
		Payload:   payload,
	}
}
