package gorm

import (
	"context"

	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/core/config"
)

func gormNewFactory(name string) func(ctx context.Context, cfg *config.Config) (*DB, error) {
	return func(ctx context.Context, cfg *config.Config) (*DB, error) {
		dsn, err := cfg.URL(ctx, config.Key("db", name))
		if err != nil {
			return nil, stacktrace.Propagate(err, "DSN not found for db '%s'", name)
		}

		dialect := dsn.Scheme
		dialector, ok := globalDriver.Dialect(*dsn)
		if !ok {
			return nil, stacktrace.NewError(
				"no registered driver for sql dialect '%s', maybe it's behind the build tag 'gorm.%s' or your dsn scheme is invalid",
				dialect, dialect,
			)
		}
		return newDB(dialector)
	}
}
