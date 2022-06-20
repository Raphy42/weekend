package config

import (
	"context"
	"sync"

	"github.com/itchyny/gojq"
	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/pkg/set"
)

type InMemoryConfig struct {
	sync.RWMutex
	Values    map[any]any
	queryable map[string]any
}

func NewInMemoryConfig(values ...map[any]any) *InMemoryConfig {
	if len(values) > 0 {
		return &InMemoryConfig{Values: values[0]}
	}
	return &InMemoryConfig{
		Values:    make(map[any]any),
		queryable: make(map[string]any),
	}
}

func (i *InMemoryConfig) Refresh(_ context.Context) error {
	i.Lock()
	defer i.Unlock()

	return nil
}

func (i *InMemoryConfig) Merge(ctx context.Context, configurable Configurable) (Configurable, error) {
	i.RLock()
	defer i.RUnlock()

	allValues, err := configurable.Get(ctx, confAll)
	if err != nil {
		return nil, stacktrace.Propagate(err, "unable to get config values from input configurable")
	}
	mappable, ok := allValues.(map[string]any)
	if !ok {
		return nil, stacktrace.Propagate(err, "underlying config implementation is not mergeable")
	}

	return NewInMemoryConfig(set.Merge(set.AsMapInterfaceInterface(mappable), i.Values)), nil
}

func (i *InMemoryConfig) Get(ctx context.Context, key string) (any, error) {
	query, err := gojq.Parse(key)
	if err != nil {
		return nil, stacktrace.NewError("not a valid key: (invalid jq query)")
	}

	i.Lock()
	defer i.Unlock()

	if i.queryable == nil {
		i.queryable = set.AsMapStringInterface(i.Values)
	}

	iter := query.RunWithContext(ctx, i.queryable)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if _, ok := v.(error); ok {
			return nil, stacktrace.NewError("jq execution error")
		}
		return v, nil
	}
	return nil, stacktrace.NewErrorWithCode(errors.EUnreachable, "unreachable part of code reached")
}
