package app

import (
	"github.com/Raphy42/weekend/core/scheduler/async"
	"github.com/Raphy42/weekend/core/supervisor"
)

type Engine struct {
	supervisor *supervisor.Supervisor
	manifests  []async.Manifest
	errors     chan error
}

func (e *Engine) Manifest() async.Manifest {
	return e.supervisor.Manifest()
}
