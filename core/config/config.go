package config

import (
	"context"
	"net/url"

	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/pkg/std/slice"
)

const (
	confAll = "."
)

type Configurable interface {
	Refresh(ctx context.Context) error
	Merge(ctx context.Context, configurable Configurable) (Configurable, error)
	Get(ctx context.Context, key string) (any, error)
}

//type Facade interface {
//	Bool(ctx context.Context, key string, defaultTo ...bool) (bool, error)
//	String(ctx context.Context, key string, defaultTo ...string) (string, error)
//	Strings(ctx context.Context, key string, defaultTo ...[]string) ([]string, error)
//	Bytes(ctx context.Context, key string, defaultTo ...[]byte) ([]byte, error)
//	Float(ctx context.Context, key string, defaultTo ...float64) (float64, error)
//	URL(ctx context.Context, key string, defaultTo ...*url.URL) (*url.URL, error)
//}

type Config struct {
	Configurable
}

func New(configurable Configurable) *Config {
	return &Config{configurable}
}

func (c Config) Bool(ctx context.Context, key string, defaultTo ...bool) (bool, error) {
	v, err := c.Get(ctx, key)
	if err != nil {
		return false, err
	}
	if v == nil {
		if len(defaultTo) > 0 {
			return defaultTo[0], nil
		}
		return false, stacktrace.NewError("entry not found in config")
	}
	value, ok := v.(bool)
	if !ok {
		return false, errors.InvalidCast(true, v)
	}
	return value, nil
}

func (c Config) String(ctx context.Context, key string, defaultTo ...string) (string, error) {
	v, err := c.Get(ctx, key)
	if err != nil {
		return "", err
	}
	if v == nil {
		if len(defaultTo) > 0 {
			return defaultTo[0], nil
		}
		return "", stacktrace.NewError("entry not found in config")
	}
	value, ok := v.(string)
	if !ok {
		return "", errors.InvalidCast("", v)
	}
	return value, nil
}

func (c Config) Strings(ctx context.Context, key string, defaultTo ...string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (c Config) Bytes(ctx context.Context, key string, defaultTo ...byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (c Config) Float(ctx context.Context, key string, defaultTo ...float64) (float64, error) {
	v, err := c.Get(ctx, key)
	if err != nil {
		return 0.0, err
	}
	if v == nil {
		if len(defaultTo) > 0 {
			return defaultTo[0], nil
		}
		return 0.0, stacktrace.NewError("entry not found in config")
	}
	value, ok := v.(float64)
	if !ok {
		return 0.0, errors.InvalidCast(0.0, v)
	}
	return value, nil
}

func (c Config) Int(ctx context.Context, key string, defaultTo ...int) (int, error) {
	v, err := c.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	if v == nil {
		if len(defaultTo) > 0 {
			return defaultTo[0], nil
		}
		return 0, stacktrace.NewError("entry not found in config")
	}
	value, ok := v.(int)
	if !ok {
		return 0, errors.InvalidCast(0, v)
	}
	return value, nil
}

func (c Config) URL(ctx context.Context, key string, defaultTo ...*url.URL) (*url.URL, error) {
	v, err := c.String(ctx, key, slice.Map(defaultTo, func(_ int, in *url.URL) string {
		return in.String()
	})...)
	if err != nil {
		return nil, stacktrace.Propagate(err, "entry not found in config")
	}
	uri, err := url.Parse(v)
	if err != nil {
		return nil, stacktrace.Propagate(err, "invalid url format")
	}
	return uri, nil
}
