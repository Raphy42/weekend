package di

import (
	"context"

	"github.com/palantir/stacktrace"
	"go.uber.org/dig"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/core/scheduler/schedulable"
	"github.com/Raphy42/weekend/pkg/reflect"
)

type Container struct {
	name    string
	modules []Module
	inner   *dig.Container
}

func NewContainer(name string) *Container {
	return &Container{
		name:    name,
		modules: make([]Module, 0),
		inner:   dig.New(),
	}
}

func (c *Container) Use(modules ...Module) *Container {
	c.modules = append(c.modules, modules...)
	return c
}

func (c *Container) startProviders(ctx context.Context, providers []interface{}) error {
	log := logger.FromContext(ctx)

	for _, provider := range providers {
		if err := c.inner.Provide(provider); err != nil {
			if dig.CanVisualizeError(err) {
				// error is related to dig's DAG construction
				// we assume it is a missing dependency but it's a bit hacky
				// todo: check if dig exposes it's error mechanism
				// todo: expose dig's DAG dotviz visualisation
				return stacktrace.PropagateWithCode(
					err,
					EDependencyMissing,
					"unable to start provider:\n%s", reflect.Signature(provider),
				)
			} else {
				// assume error is related to provider initialisation
				return stacktrace.PropagateWithCode(
					err,
					EDependencyInitFailed,
					"unable to start provider:\n%s", reflect.Signature(provider),
				)
			}
		}
		log.Debug("loaded provider", zap.String("di.provider.signature", reflect.Signature(provider)))
	}
	return nil
}

func (c *Container) startInvocations(ctx context.Context, invocations []interface{}) error {
	log := logger.FromContext(ctx)

	for _, invoc := range invocations {
		if err := c.inner.Invoke(invoc); err != nil {
			return stacktrace.PropagateWithCode(
				err,
				EInvocationFailed,
				"unable to run invocation:\n%s", reflect.Signature(invoc),
			)
		}
		log.Debug("ran invocation", zap.String("di.invocation.signature", reflect.Signature(invoc)))
	}
	return nil
}

func (c *Container) startChildren(ctx context.Context) error {
	log := logger.FromContext(ctx)

	exposed := make([]interface{}, 0)
	exports := make([]interface{}, 0)
	invocations := make([]interface{}, 0)
	log.Debug("loaded modules", zap.Int("di.module.count", len(c.modules)))

	for _, module := range c.modules {
		log.Debug("building module", zap.String("di.module.name", module.Name))
		exposed = append(exposed, module.Exposes...)
		exports = append(exports, module.Providers...)
		invocations = append(invocations, module.Invokes...)
	}
	log.Debug("building dependency tree",
		zap.Int("di.exposed", len(exposed)),
		zap.Int("di.exports", len(exports)),
		zap.Int("di.invocations", len(invocations)),
	)

	// inject parent context in dependency DAG
	exports = append(exports, func() context.Context {
		return ctx
	})

	// todo handle and inject lifecycles

	if err := c.startProviders(ctx, exports); err != nil {
		return stacktrace.Propagate(err, "cannot continue if a provider crashes")
	}

	if err := c.startInvocations(ctx, invocations); err != nil {
		return stacktrace.Propagate(err, "cannot continue if an invocation failed")
	}

	// todo install application lifecycle event handlers

	return nil
}

func (c *Container) start(ctx context.Context) error {
	return c.startChildren(ctx)
}

func (c *Container) Manifest() schedulable.Manifest {
	return schedulable.Of(
		schedulable.Name("container", c.name),
		c.start,
	)
}
