[![Go Reference](https://pkg.go.dev/badge/github.com/Raphy42/weekend.svg)](https://pkg.go.dev/github.com/Raphy42/weekend)
[![Go Report Card](https://goreportcard.com/badge/github.com/Raphy42/weekend)](https://goreportcard.com/report/github.com/Raphy42/weekend)

# weekend
All included golang toolkit

## Status
Heavy WIP

## Default features
- latest go version for reduced boilerplate (generics)
- lightweight runtime DI system, including a way to define health-checks for injected services
- unique lexicographically sortable ids generated with `rs/xid` (supported by multiple languages)
- `context.Context` based scheduling system (lightly based on supervision trees)
- error management based on `palantir/stacktrace` thanks to bitmasks
- context aware structured logging with `uber/zap`
- bus system for event sourcing, available from within scheduled domain
- automatic tracing through opentelemetry, injected in asynchronous contexts, redis, gorm, and gin

## Gated features (use tag when compiling to enable)
- `ops.sentry`: sentry panic handler
- `gorm.postgres`: gorm postgres support
- `gorm.sqlite`: gorm sqlite support