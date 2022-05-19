package scheduler

import "context"

type SchedulingContext struct {
	context.Context
	cancel    context.CancelFunc
	metadatas map[string]interface{}
}

func (s *SchedulingContext) WithMetadata(key string, value interface{}) *SchedulingContext {
	s.metadatas[key] = value
	return s
}

func newSchedulingContext(parent ...context.Context) *SchedulingContext {
	parentCtx := context.Background()
	if len(parent) != 0 {
		parentCtx = parent[0]
	}
	ctx, cancel := context.WithCancel(parentCtx)
	return &SchedulingContext{
		Context:   ctx,
		cancel:    cancel,
		metadatas: make(map[string]interface{}),
	}
}

func Context(parent ...context.Context) *SchedulingContext {
	return newSchedulingContext(parent...)
}
