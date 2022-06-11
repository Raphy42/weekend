package core

import (
	"github.com/Raphy42/weekend/core/message"
	"github.com/Raphy42/weekend/core/scheduler"
	"github.com/Raphy42/weekend/core/scheduler/schedulable"
)

type Engine struct {
	bus        message.Bus
	background *scheduler.Pipeline
}

type EngineBuilder struct {
	background []schedulable.Manifest
}

func newEngineBuilder() *EngineBuilder {
	return &EngineBuilder{
		background: make([]schedulable.Manifest, 0),
	}
}

func (e *EngineBuilder) Background(manifests ...schedulable.Manifest) *EngineBuilder {
	e.background = append(e.background, manifests...)
	return e
}

func (e *EngineBuilder) Build() (*Engine, error) {
	return nil, nil
}
