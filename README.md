[![Go Reference](https://pkg.go.dev/badge/github.com/Raphy42/weekend.svg)](https://pkg.go.dev/github.com/Raphy42/weekend)

# weekend
All included golang toolkit (heavy WIP)

# TODO
## cmd
### node
- controller/worker
### dashboard
- tbd
## core
### bitmask
- better bitmask (u8 instead of i16 ?)
### logger
- context metadata extraction (tracing, request_id, auth)
- context decorator
- better production defaults
### errors
- diagnostic system
- better bit handling for stacktrace.ErrorCode (0x7fff instead of 0xfff0): DONE
- friendlier end-user API
### dependency_injection
- finish lifecycle sub-system
- inject lifecycle sub-system at provider stage
- move scheduling to own package
### reflect
- clone
### scheduler
- agent interface
- start/stop scheduling
- retry scheduling
- delay scheduling
- task idempotency
- binpack algorithm
### metrics
- prometheus integration
- sentry integration
- custom alerting solution
### templating
- wrap template/text or use custom template
### transport
#### http
- gin based
#### ws
- gorilla based
#### mail
- TBD probably SMTP via POP3
### TBD: distributed system
- discovery
- consensus
- election
- failover
## api
- stabilise protos
### grpc
- write wrapper around client & server