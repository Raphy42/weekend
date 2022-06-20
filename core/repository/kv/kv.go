package kv

import (
	"context"
	"time"
)

type WriteApi interface {
	SetBytes(ctx context.Context, key string, value []byte) error
	SetString(ctx context.Context, key, value string) error
	SetAny(ctx context.Context, key string, value any) error
	SetTTL(ctx context.Context, key string, ttl time.Duration) error
}

type ReadApi interface {
	GetBytes(ctx context.Context, key string) ([]byte, error)
	GetString(ctx context.Context, key string) (string, error)
	GetAny(ctx context.Context, key string, valuePtr any) error
	GetTTL(ctx context.Context, key string) (*time.Duration, error)
}

type Keyable interface {
	Key(args ...string) string
}

type KV interface {
	WriteApi
	ReadApi
	Keyable
}
