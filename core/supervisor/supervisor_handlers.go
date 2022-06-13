package supervisor

import (
	"context"

	"github.com/palantir/stacktrace"
	"go.uber.org/atomic"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/core/scheduler"
)

func (s *Supervisor) handleSchedule(ctx context.Context, payload *scheduler.ScheduleMessagePayload) error {
	if _, ok := s.restarts[payload.ManifestID]; !ok {
		s.restarts[payload.ManifestID] = atomic.NewInt32(0)
	}

	spec, ok := s.specLut[payload.ManifestID]
	if !ok {
		return stacktrace.NewError("no spec associated with manifest ManifestID '%s'", payload.ManifestID)
	}
	handle, err := s.scheduler.Schedule(ctx, spec.Manifest, spec.Args)
	if err != nil {
		return stacktrace.Propagate(err, "could not schedule manifest '%s'", spec.Manifest.Name)
	}
	s.children[handle.ID] = handle
	return nil
}

func (s *Supervisor) handleFailure(ctx context.Context, payload *scheduler.FailureMessagePayload) error {
	log := logger.FromContext(ctx)

	if payload.Error != nil {
		spec, ok := s.specLut[payload.ManifestID]
		if !ok {
			return stacktrace.NewError("no spec associated with manifest ManifestID '%s'", payload.ManifestID)
		}
		strategy := spec.Strategy
		switch strategy.Restart {
		case PermanentRestartStrategy:
			log.Debug("re-scheduling", zap.Stringer("wk.manifest.id", payload.ManifestID))
			return s.restart(ctx, payload.ManifestID, payload.HandleID, payload.Error)
		case TransientRestartStrategy:
			if errors.IsTransient(payload.Error) {
				log.Debug("re-scheduling", zap.Stringer("wk.manifest.id", payload.ManifestID))
				return s.restart(ctx, payload.ManifestID, payload.HandleID, payload.Error)
			}
			log.Debug("not re-scheduling due to restart strategy",
				zap.Stringer("wk.manifest.id", payload.ManifestID),
			)
			return stacktrace.Propagate(payload.Error,
				"error was not transient, terminating supervision",
			)
		case TemporaryRestartStrategy:
			// according to strategy we don't terminate the whole supervision tree on error
			log.Debug("not re-scheduling due to restart strategy",
				zap.Stringer("wk.manifest.id", payload.ManifestID),
			)
		default:
			panic(stacktrace.NewErrorWithCode(errors.EUnreachable, "invalid supervision strategy"))
		}
	}
	return nil
}
