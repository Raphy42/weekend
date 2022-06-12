package app

import "github.com/Raphy42/weekend/core/scheduler/schedulable"

type EngineBuilder struct {
	background []schedulable.Manifest
}

func NewEngineBuilder() *EngineBuilder {
	return &EngineBuilder{
		background: make([]schedulable.Manifest, 0),
	}
}

func (e *EngineBuilder) Background(manifests ...schedulable.Manifest) *EngineBuilder {
	e.background = append(e.background, manifests...)
	return e
}

func (e *EngineBuilder) Build() (*Engine, error) {
	return &Engine{manifests: e.background}, nil
}