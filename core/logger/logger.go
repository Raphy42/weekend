package logger

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/palantir/stacktrace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Raphy42/weekend/pkg/runtime"
)

const (
	KFile    = "file"
	KLine    = "line"
	LDebug   = zapcore.DebugLevel
	LInfo    = zapcore.InfoLevel
	LWarning = zapcore.WarnLevel
	LError   = zapcore.ErrorLevel
	LFatal   = zapcore.FatalLevel
)

type globalLogger struct {
	sync.RWMutex
	logger *zap.Logger
}

func init() {
	logMode := os.Getenv("WEEKEND_LOG_MODE")
	if logMode == "" {
		logMode = "DEV"
	}
	var logger *zap.Logger
	var err error

	switch logMode {
	case "DEV":
		logger, err = zap.NewDevelopment()
	case "PROD":
		logger, err = zap.NewProduction()
	default:
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(stacktrace.Propagate(err, "unable to initialise root logger"))
	}

	zap.ReplaceGlobals(logger)
}

func New(opts ...Option) *zap.Logger {
	options := newLoggerOptions()
	options.apply(opts...)

	caller := runtime.CallerName(options.SkipCallFrame)

	name := caller
	if options.Name != "" {
		name = options.Name
	}

	return zap.L().Named(name)
}

func ctxDecorator(ctx context.Context) []Option {
	opts := make([]Option, 0)
	deadline, ok := ctx.Deadline()
	if ok {
		opts = append(opts,
			Decorate(func(s string) string {
				return fmt.Sprintf("%s!", s)
			}),
			Fields(zap.Time("context.deadline", deadline)),
		)
	}

	opts = append(opts, SkipCallFrame(2))
	return opts
}

func FromContext(ctx context.Context, opts ...Option) *zap.Logger {
	//todo context aware metadata retrieval
	// - http
	// - business
	// - domain specific
	opts = append(opts, ctxDecorator(ctx)...)
	return New(opts...)
}
