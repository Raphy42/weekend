package message

import (
	"context"
)

type Bus interface {
	Emit(ctx context.Context, message Message) error
	Read(ctx context.Context) (<-chan Message, context.CancelFunc)
}
