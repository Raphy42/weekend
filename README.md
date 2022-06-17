[![Go Reference](https://pkg.go.dev/badge/github.com/Raphy42/weekend.svg)](https://pkg.go.dev/github.com/Raphy42/weekend)
[![Go Report Card](https://goreportcard.com/badge/github.com/Raphy42/weekend)](https://goreportcard.com/report/github.com/Raphy42/weekend)

# weekend

All included golang toolkit

## Status

Heavy WIP

## Current features

- latest go version for reduced boilerplate (generics)
- unique lexicographically sortable ids generated with `rs/xid` (supported by multiple languages)
- `context.Context` based scheduling system (lightly based on supervision trees)
- error management based on `palantir/stacktrace` thanks to bitmasks
- context aware structured logging with `uber/zap`
- bus system for event sourcing, available from within scheduled domain

## Coming soon

- stabilisation of `platform` module (environment, secrets)
- stabilisation of `scheduler.Supervisor` (policies, lifecycle)
- `task` module for distributed job scheduling

## Planned

- monitoring gRPC API
- driver interface for message brokers (redis, kafka)
- database interface for RDBMS drivers (postgres, redis)
- resilient plugin system to reduce dependency bloat (for external drivers)