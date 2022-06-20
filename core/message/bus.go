package message

import (
	"context"
)

type Mailbox interface {
	Emit(ctx context.Context, message Message) error
	ReadC(ctx context.Context) (<-chan Message, context.CancelFunc)
}

type Bus interface {
	Subject(ctx context.Context, subject string) (Mailbox, error)
}
