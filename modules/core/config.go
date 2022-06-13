package core

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/palantir/stacktrace"
	"gopkg.in/yaml.v3"

	"github.com/Raphy42/weekend/core/config"
)

func configFromFilenames(ctx context.Context, filenames ...string) (config.Configurable, error) {
	cfg := config.NewInMemoryConfig()
	for _, filename := range filenames {
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, stacktrace.Propagate(err, "unable to read config file: '%s'", filename)
		}
		switch v := filepath.Ext(filename); v {
		case ".yml", ".yaml":
			var values map[any]any
			if err = yaml.Unmarshal(content, &values); err != nil {
				return nil, stacktrace.Propagate(err, "YAML deserialization failed")
			}
			newest := config.NewInMemoryConfig(values)
			merged, err := cfg.Merge(ctx, newest)
			if err != nil {
				return nil, stacktrace.Propagate(err, "unable to merge configurations")
			}
			cfg = merged.(*config.InMemoryConfig)
		case ".json":
			var values map[any]any
			if err = json.Unmarshal(content, &values); err != nil {
				return nil, stacktrace.Propagate(err, "JSON deserialization failed")
			}
			newest := config.NewInMemoryConfig(values)
			merged, err := cfg.Merge(ctx, newest)
			if err != nil {
				return nil, stacktrace.Propagate(err, "unable to merge configurations")
			}
			cfg = merged.(*config.InMemoryConfig)
		default:
			return nil, stacktrace.Propagate(err, "unhandled config format extension: '%s'", v)
		}
	}
	return cfg, nil
}
