package config

import (
	"context"
	"sync"

	"github.com/itchyny/gojq"
	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/pkg/std"
)

type InMemoryConfig struct {
	sync.RWMutex
	Values map[interface{}]interface{}
}

func NewInMemoryConfig(values ...map[interface{}]interface{}) *InMemoryConfig {
	if len(values) > 0 {
		return &InMemoryConfig{Values: values[0]}
	}
	return &InMemoryConfig{Values: make(map[interface{}]interface{})}
}

func (i *InMemoryConfig) Refresh(ctx context.Context) error {
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
	mappable, ok := allValues.(map[interface{}]interface{})
	if !ok {
		return nil, stacktrace.Propagate(err, "underlying config implementation is not mergeable")
	}

	return NewInMemoryConfig(std.MergeMap(i.Values, mappable)), nil
}

func (i *InMemoryConfig) Get(ctx context.Context, key string) (interface{}, error) {
	i.RLock()
	defer i.RUnlock()

	query, err := gojq.Parse(key)
	if err != nil {
		return nil, stacktrace.Propagate(err, "not a valid key: (invalid jq query)")
	}

	iter := query.RunWithContext(ctx, i.Values)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return nil, stacktrace.Propagate(err, "jq execution error")
		}
		return v, nil
	}
	return nil, stacktrace.NewErrorWithCode(errors.EUnreachable, "unreachable part of code reached")
}
