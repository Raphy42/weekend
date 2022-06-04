package message

import "context"

type Handler func(ctx context.Context, message Message) (Message, error)
