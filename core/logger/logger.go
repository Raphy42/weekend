package logger

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/palantir/stacktrace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	logmode := os.Getenv("WEEKEND_LOG_MODE")
	if logmode == "" {
		logmode = "DEV"
	}
	var logger *zap.Logger
	var err error

	switch logmode {
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

	globalLoggerInstance.Lock()
	defer globalLoggerInstance.Unlock()

	globalLoggerInstance.logger = logger
}

var globalLoggerInstance = &globalLogger{}

func New(opts ...LoggerOption) *zap.Logger {
	options := newLoggerOptions()
	options.apply(opts...)

	callerPtr, _, _, _ := runtime.Caller(options.SkipCallFrame)
	caller := runtime.FuncForPC(callerPtr)

	globalLoggerInstance.RLock()
	defer globalLoggerInstance.RUnlock()

	name := caller.Name()
	if options.Name != "" {
		name = options.Name
	}

	return globalLoggerInstance.logger.Named(name)
}

func ctxDecorator(ctx context.Context) []LoggerOption {
	opts := make([]LoggerOption, 0)
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

func FromContext(ctx context.Context, opts ...LoggerOption) *zap.Logger {
	//todo context aware metadata retrieval
	// - http
	// - business
	// - domain specific
	opts = append(opts, ctxDecorator(ctx)...)
	return New(opts...)
}
