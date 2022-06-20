package kv

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/pkg/concurrent_set"
)

type InMemory struct {
	pages concurrent_set.Set[string, []byte]
	ttl   concurrent_set.Set[string, time.Duration]
}

func NewInMemory() *InMemory {
	return &InMemory{
		pages: concurrent_set.New[string, []byte](),
		ttl:   concurrent_set.New[string, time.Duration](),
	}
}

func (i *InMemory) SetBytes(_ context.Context, key string, value []byte) error {
	i.pages.Insert(key, value)
	return nil
}

func (i *InMemory) SetString(_ context.Context, key, value string) error {
	i.pages.Insert(key, []byte(value))
	return nil
}

func (i *InMemory) SetAny(_ context.Context, key string, value any) error {
	buf, err := json.Marshal(value)
	if err != nil {
		return err
	}
	i.pages.Insert(key, buf)
	return nil
}

func (i *InMemory) SetTTL(ctx context.Context, key string, ttl time.Duration) error {
	i.ttl.Insert(key, ttl)

	time.AfterFunc(ttl, func() {
		i.pages.Delete(key)
		i.ttl.Delete(key)
	})
	return nil
}

func (i *InMemory) GetBytes(_ context.Context, key string) ([]byte, error) {
	bytes, ok := i.pages.Get(key)
	if !ok {
		return nil, stacktrace.NewErrorWithCode(EEntryNotFound, "no such key: '%s'", key)
	}
	return bytes, nil
}

func (i *InMemory) GetString(ctx context.Context, key string) (string, error) {
	bytes, err := i.GetBytes(ctx, key)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (i *InMemory) GetAny(ctx context.Context, key string, valuePtr any) error {
	bytes, err := i.GetBytes(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, valuePtr)
}

func (i *InMemory) GetTTL(ctx context.Context, key string) (*time.Duration, error) {
	duration, ok := i.ttl.Get(key)
	if !ok {
		return nil, nil
	}
	// cloned to prevent undesired ttl mutation
	return &*&duration, nil
}

func (i *InMemory) Key(args ...string) string {
	return strings.Join(args, ".")
}
